// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bwplotka/flagarize"
	"gopkg.in/alecthomas/kingpin.v2"
)

const prefix = "AlwaysAddingPrefix"

type YourCustomType string

func (cst *YourCustomType) Set(s string) error {
	*cst = YourCustomType(fmt.Sprintf("%s%s", prefix, s))
	return nil
}

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "<Your CLI description>")

	type ConfigForCLI struct {
		Field1              string         `flagarize:"name=flag1|help=Some help.|default=Some Value|envvar=FLAG1|short=f|placeholder=<put some string here>"`
		Field2              YourCustomType `flagarize:"name=custom-type-flag|required=true|short=c|placeholder=<put some my custom type here>"`
		Field2FlagarizeHelp string
	}

	// Flagarize your config! (Register flags from config).
	cfg := &ConfigForCLI{
		Field2FlagarizeHelp: fmt.Sprintf("Custom Type always add prefix %q to the given string value.", prefix),
	}
	if err := flagarize.Flagarize(a, cfg); err != nil {
		log.Fatal(err)
	}
	if _, err := a.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Config Value after flagarizing & flag parsing: %+v\n", cfg)
}
