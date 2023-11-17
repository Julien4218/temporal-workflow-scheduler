package main

import (
	"log"

	flag "github.com/spf13/pflag"
)

func main() {
	if err := Execute(); err != nil {
		if err != flag.ErrHelp {
			log.Fatal(err)
		}
	}
}
