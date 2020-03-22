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

type YourCustomType string

func (cst *YourCustomType) Set(s string) error {
	*cst = YourCustomType(fmt.Sprintf("AlwaysAddingPrefix%s", s))
	return nil
}

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "<Your CLI description>")

	type ConfigForCLI struct {
		Field1 string         `flagarize:"name=flag1|help=Some help.|default=Some Value|envvar=FLAG1|short=f|placeholder=<put some string here>"`
		Field2 YourCustomType `flagarize:"name=custom-type-flag|help=Custom Type always add prefix 'AlwaysAddingPrefix' to the given string value.|required=true|short=c|placeholder=<put some my custom type here>"`
	}

	// Flagarize your config! (Register flags from config).
	cfg := &ConfigForCLI{}
	if err := flagarize.Flagarize(a, cfg); err != nil {
		log.Fatal(err)
	}
	if _, err := a.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
