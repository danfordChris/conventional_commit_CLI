package main

import (
	"os"

	"github.com/user/convcommit/cmd/convcommit"
)

func main() {
	os.Exit(convcommit.Execute())
}