# flagarize

[![CI](https://github.com/bwplotka/flagarize/workflows/test/badge.svg)](https://github.com/bwplotka/flagarize/actions?query=workflow%3Atest)

Flagarize your Go struct to initialize your even complex CLI config struct from flags!

## Goals

* Allow flag parsing for any struct field using Go struct tags.
* Minimal dependencies: Only `"gopkg.in/alecthomas/kingpin.v2"`.
* Extensible with [custom types](#custom-type-parsing) and [custom flagarizing](#custom-flags).
* Native supports for all [kingpin](https://github.com/alecthomas/kingpin) flag types and more like [`regexp`](./regexp.go) , [`pathorcontent`](./pathorcontent.go), [`timeorduration`](./timeorduration.go).

## Requirements:

* Go 1.3+
* `gopkg.in/alecthomas/kingpin.v2`

## Usage

See below example for usage:

```go

type ComponentAOptions struct {
	Field1 string
}

func main() {
	// Create new kingpin app as usual.
	a := kingpin.New(filepath.Base(os.Args[0]), "<Your CLI description>")

	// Define you own config.
	type ConfigForCLI struct {
		Field1 string                   `flagarize:"name=config.file|help=Prometheus configuration file path.|default=prometheus.yml"`
		Field2 []string                 `flagarize:"name=web.external-url|help=The URL under which Prometheus is externally reachable (for example, if Prometheus is served via a reverse proxy). Used for generating relative and absolute links back to Prometheus itself. If the URL has a path portion, it will be used to prefix all HTTP endpoints served by Prometheus. If omitted, relevant URL components will be derived automatically.|placeholder=<URL>"`
		Field3 int                      `flagarize:"name=storage.tsdb.path|help=Base path for metrics storage.|default=data/"`
		Field4 flagarize.TimeOrDuration `flagarize:"name=storage.remote.flush-deadline|help=How long to wait flushing sample on shutdown or config reload.|default=1m|placeholder=<duration>"`

		NotFromFlags int `flagarize:"name=query.max-concurrency|help=Maximum number of queries executed concurrently.|default=20"`

		ComponentA ComponentAOptions
	}

	// Create new config.
	cfg := &ConfigForCLI{}

	// Flagarize it! (Register flags from config).
	if err := flagarize.Flagarize(a, &cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// You can define some fields as usual as well.
	var notInConfigField time.Duration
	a.Flag("some-field10", "...").
		DurationVar(&notInConfigField)

	// Parse flags as usual.
	if _, err := a.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	// Config is filled with flags from value!
	... = cfg.Field1
}
```

## Production Examples

To see production example see:

 * [Thanos](todo)
 * [Prometheus](todo)

## But Bartek, such Go projects already exists!

Yes, but not as simple, not focused on kingpin, and they does not allow custom flagarazing!  🤗

## But Bartek, normal flag registration is enough, don't overengineer!

Well, depends. It might get quite weird. Here is how it could look with and without flagarize:

### Without

```

```

### With flagarize

```

```

## Custom Type Parsing

Flagarize allows parsing of native types like int, string, etc (all that kingpin supports). For custom 
types it's enough if your type implements part of `kingping.Value` interface as follows:

```go
// ValueFlagarizer is the simplest way to extend flagarize to parse your custom type.
// If any field has `flagarize:` struct tag and it implements the ValueFlagarizer, this will be
// used by kingping to parse the flag value.
//
// For an example see: `./timeduration.go` or `./regexp.go`
type ValueFlagarizer interface {
	// FlagarizeSetValue is invoked on kinpgin.Parse with the flag value passed as string.
	// It is expected from this method to parse the string to the underlying type.
	// This method has to be a pointer receiver for the method to take effect.
	// Flagarize will return error otherwise.
	Set(s string) error
}
```

## Custom Flags

Sometimes custom parsing is not enough. Sometimes you need to register more flags than just one from
single flagarize definition. To do so the type has to implement following interface:

```go
// Flagarizer is more advanced way to extend flagarize to parse a type. It allows to register
// more than one flag or register them in a custom way. It's ok for a method to register nothing.
// If any field implements `Flagarizer` this method will be invoked even if field does not
// have `flagarize:` struct tag.
//
// If the field implements both ValueFlagarizer and Flagarizer, only Flagarizer will be used.
//
// For an example usage see: `./pathorcontent.go`
type Flagarizer interface {
	// Flagarize is invoked on Flagarize. If field type does not implement custom Flagarizer
	// default one will be used.
	// Tag argument is nil if no `flagarize` struct tag was specified. Otherwise it has parsed
	//`flagarize` struct tag.
	// The ptr argument is an address of the already allocated type, that can be used
	// by FlagRegisterer kingping *Var methods.
	Flagarize(r FlagRegisterer, tag *Tag, ptr unsafe.Pointer) error
}
```