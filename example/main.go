// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/bwplotka/flagarize"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	a := kingpin.New(filepath.Base(os.Args[0]), "<Your CLI description>")

	type ComponentAOptions struct {
		Field1 []string `flagarize:"name=a.flag1|help=Help for field 1 in nested struct for component A."`
	}
	type ConfigForCLI struct {
		Field1     string                   `flagarize:"name=flag1|help=Help for field 1.|default=something"`
		Field2     *url.URL                 `flagarize:"name=flag2|help=Help for field 2.|placeholder=<URL>"`
		Field3     int                      `flagarize:"name=flag3|help=Help for field 3.|default=2144"`
		Field4     flagarize.TimeOrDuration `flagarize:"name=flag4|help=Help for field 4.|default=1m|placeholder=<time or duration>"`
		ComponentA ComponentAOptions
	}
	var notInConfigField time.Duration
	a.Flag("some-field10", "Help for some help which is defined outside of ConfigForCLI struct.").
		DurationVar(&notInConfigField)

	// Flagarize your config! (Register flags from config).
	cfg := &ConfigForCLI{}
	if err := flagarize.Flagarize(a, cfg); err != nil {
		log.Fatal(err)
	}
	if _, err := a.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	// Config is filled with values from flags.
	fmt.Printf("Config Value after flagarizing & flag parsing: %+v\n", cfg)
}
