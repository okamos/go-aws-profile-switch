package awsswitch

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/okamos/go-aws-profile-switch/ui"
)

const (
	keyName    = "aws_access_key_id"
	secretName = "aws_secret_access_key"
	regionName = "region"
	outputName = "output"
	swUsage    = "Usage: awsswitch ws -p [profile name]"
)

var credentialsPath string

// LsCommand defines list AWS profiles
type LsCommand struct{}

// Synopsis about ls
func (c *LsCommand) Synopsis() string {
	return "Lists available AWS profiles"
}

// Help about ls
func (c *LsCommand) Help() string {
	return "Usage: awsswitch ls"
}

// Run ls command
func (c *LsCommand) Run(args []string) int {
	credentials, _ := LoadCredentials()

	var defaultKey string
	for _, v := range credentials {
		if v["profile"] == "default" {
			defaultKey = v[keyName]
			break
		}
	}
	fmt.Println("Available profiles")
	for _, v := range credentials {
		if v["profile"] != "default" && v[keyName] == defaultKey {
			fmt.Println("* " + v["profile"])
		} else {
			fmt.Println("  " + v["profile"])
		}
	}

	return 0
}

// SwCommand defines switch default aws profile
type SwCommand struct{}

// Synopsis about sw
func (c *SwCommand) Synopsis() string {
	return "Switches default your AWS profile"
}

// Help about sw
func (c *SwCommand) Help() string {
	return swUsage
}

// Run sw command
func (c *SwCommand) Run(args []string) int {
	var (
		profile string
		err     error
	)

	flags := flag.NewFlagSet("sw", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(swUsage) }
	flags.StringVar(&profile, "profile", "", "A profile name from your credentials")
	flags.StringVar(&profile, "p", "", "A profile name from your credentials")

	if err := flags.Parse(args); err != nil {
		return 1
	}
	credentials, comments := LoadCredentials()

	if profile == "" {
		profile, err = dynamicSwitch(credentials)
		if err != nil {
			log.Print(err)
			return 1
		}
	}

	var defaultKey string
	candidate := map[string]string{}
	for _, v := range credentials {
		if v["profile"] == "default" {
			defaultKey = v[keyName]
		} else if v["profile"] == profile {
			candidate["profile"] = profile
			candidate[keyName] = v[keyName]
			candidate[secretName] = v[secretName]
			candidate[regionName] = v[regionName]
			candidate[outputName] = v[outputName]
		}
	}

	if candidate["profile"] == "" {
		fmt.Printf("The profile %s is not found\n", profile)
		return 1
	}
	if defaultKey == candidate[keyName] {
		fmt.Printf("Your default profile set %s already\n", profile)
		return 0
	}

	isBackup := defaultKey == ""
	for _, v := range credentials {
		if v[keyName] == defaultKey {
			isBackup = true
			break
		}
	}
	if !isBackup {
		s := "Your default profile will be overwrite `" + profile + "`\n"
		fmt.Printf("%s You should be backup default profile.", s)
		return 1
	}

	var buf []byte
	buf = append(buf, ("[default]\n")...)
	if candidate[keyName] != "" {
		buf = append(buf, (keyName + "=" + candidate[keyName] + "\n")...)
	}
	if candidate[secretName] != "" {
		buf = append(buf, (secretName + "=" + candidate[secretName] + "\n")...)
	}
	if candidate[regionName] != "" {
		buf = append(buf, (regionName + "=" + candidate[regionName] + "\n")...)
	}
	if candidate[outputName] != "" {
		buf = append(buf, (outputName + "=" + candidate[outputName] + "\n")...)
	}
	buf = append(buf, "\n"...)

	for _, v := range credentials {
		if v["profile"] == "default" {
			continue
		}
		buf = append(buf, ("[" + v["profile"] + "]\n")...)
		if v[keyName] != "" {
			buf = append(buf, (keyName + "=" + v[keyName] + "\n")...)
		}
		if v[secretName] != "" {
			buf = append(buf, (secretName + "=" + v[secretName] + "\n")...)
		}
		if v[regionName] != "" {
			buf = append(buf, (regionName + "=" + v[regionName] + "\n")...)
		}
		if v[outputName] != "" {
			buf = append(buf, (outputName + "=" + v[outputName] + "\n")...)
		}
		buf = append(buf, "\n"...)
	}
	buf = append(buf, comments...)

	file, err := os.OpenFile(credentialsPath, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return 1
	}
	defer file.Close()
	file.Write(buf)
	fmt.Println("Your default profile is overwrote:", profile)

	return 0
}

// Run the awsswitch
func Run() int {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	credentialsPath = filepath.Join(u.HomeDir, ".aws", "credentials")
	c := cli.NewCLI("awsswitch", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"ls": func() (cli.Command, error) {
			return &LsCommand{}, nil
		},
		"sw": func() (cli.Command, error) {
			return &SwCommand{}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("Failed to execute: %s\n", err.Error())
	}
	return exitCode
}

// LoadCredentials parses AWS credentials file
func LoadCredentials() ([]map[string]string, []byte) {
	fp := getFile()
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 256)
	credentials := []map[string]string{}
	var comments []byte
	current := -1
	for {
		bufLine, _, err := reader.ReadLine()
		line := strings.Trim(string(bufLine), " ")
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(line) == 0 {
			continue
		}
		if line[0] == byte('[') {
			credentials = append(credentials, map[string]string{
				"profile": strings.Trim(line[1:len(line)-1], " "),
			})
			current++
		} else if line[0] == byte('#') {
			comments = append(comments, (line + "\n")...)
		} else if strings.Contains(line, "=") {
			pair := strings.Split(line, "=")
			val := strings.Trim(pair[1], " ")
			switch strings.Trim(pair[0], " ") {
			case keyName:
				credentials[current][keyName] = val
			case secretName:
				credentials[current][secretName] = val
			case regionName:
				credentials[current][regionName] = val
			case outputName:
				credentials[current][outputName] = val
			}
		}
	}
	return credentials, comments
}

func getFile() *os.File {
	fp, err := os.Open(credentialsPath)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}

func dynamicSwitch(credentials []map[string]string) (string, error) {
	values := make([]string, 0)
	for _, v := range credentials {
		if v["profile"] != "default" {
			values = append(values, v["profile"])
		}
	}
	ui.InitSelector(values)
	return ui.Select()
}
