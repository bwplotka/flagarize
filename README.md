# flagarize

Flagarize your Go struct to initialize your even complex struct from flags!

## Goals

* Allow flag parsing for any struct field using Go struct tags.
* Minimal dependencies: Only `"gopkg.in/alecthomas/kingpin.v2"`.
* Extensible with [custom types](#custom-type-parsing) and [custom flagarizing](#custom-flags).
* Native supports for all [kingpin](https://github.com/alecthomas/kingpin) flag types and more like [`regexp`](./regexp.go) , [`pathorcontent`](./pathorcontent.go), [`timeorduration`](./timeorduration.go).

## Usage

## Custom Type Parsing



## Custom Flags

## Production Examples

To see production example see:

 * [Thanos](todo).
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