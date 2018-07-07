package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	keyName    = "aws_access_key_id"
	secretName = "aws_secret_access_key"
)

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
	fp := getFile()
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 256)

	credentials := []map[string]string{}
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
		} else {
			pair := strings.Split(line, "=")
			if strings.Trim(pair[0], " ") == keyName {
				credentials[current][keyName] = strings.Trim(pair[1], " ")
			} else if strings.Trim(pair[0], " ") == secretName {
				credentials[current][secretName] = strings.Trim(pair[1], " ")
			}
		}
	}

	var defaultKey string
	for _, v := range credentials {
		if v["profile"] == "default" {
			defaultKey = v[keyName]
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

func main() {
	c := cli.NewCLI("awsswitch", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"ls": func() (cli.Command, error) {
			return &LsCommand{}, nil
		},
	}

	exitCode, err := c.Run()
	if err != nil {
		fmt.Printf("Failed to execute: %s\n", err.Error())
	}
	os.Exit(exitCode)
}

func getFile() *os.File {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	credentials := filepath.Join(u.HomeDir, ".aws", "credentials")
	fp, err := os.Open(credentials)
	if err != nil {
		log.Fatal(err)
	}
	return fp
}
