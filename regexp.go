// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

package flagarize

import (
	"regexp"
)

type Regexp struct {
	*regexp.Regexp
}

// FlagarizeSetValue registers Regexp flag.
func (r *Regexp) FlagarizeSetValue(v string) (err error) {
	rg, err := regexp.Compile(v)
	if err != nil {
		return err
	}
	r.Regexp = rg
	return nil
}

type AnchoredRegexp struct {
	*regexp.Regexp
}

// FlagarizeSetValue registers anchored Regexp flag.
func (r *AnchoredRegexp) FlagarizeSetValue(v string) (err error) {
	rg, err := regexp.Compile("^(?:" + v + ")$")
	if err != nil {
		return err
	}
	r.Regexp = rg
	return nil
}
