package main

import (
	"os"

	"github.com/okamos/go-aws-profile-switch"
)

func main() {
	os.Exit(awsswitch.Run())
}
