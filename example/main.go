package main

import (
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/bwplotka/flagarize"
	"gopkg.in/alecthomas/kingpin.v2"
)

type ComponentAOptions struct {
	Field1 []string `flagarize:"name=a.flag1|help=Help for field 1 in nested struct for component A."`
}

func main() {
	// Create new kingpin app as usual.
	a := kingpin.New(filepath.Base(os.Args[0]), "<Your CLI description>")

	// Define you own config.
	type ConfigForCLI struct {
		Field1 string                   `flagarize:"name=flag1|help=Help for field 1.|default=something"`
		Field2 *url.URL                 `flagarize:"name=flag2|help=Help for field 2.|placeholder=<URL>"`
		Field3 int                      `flagarize:"name=flag3|help=Help for field 3.|default=2144"`
		Field4 flagarize.TimeOrDuration `flagarize:"name=flag4|help=Help for field 4. for field 1.p4|default=1m|placeholder=<time or duration>"`

		NotFromFlags int

		ComponentA ComponentAOptions
	}

	// You can define some fields as usual as well.
	var notInConfigField time.Duration
	a.Flag("some-field10", "Help for some help which is defined outside of ConfigForCLI struct.").
		DurationVar(&notInConfigField)

	// Create new config.
	cfg := &ConfigForCLI{}

	// Flagarize your config! (Register flags from config). Parse the flags afterwards.
	if err := flagarize.Flagarize(a, cfg); err != nil {
		log.Fatal(err)
	}
	if _, err := a.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
	// Config is filled with flags from value!
	_ = cfg.Field1

	// Run your command...
}
