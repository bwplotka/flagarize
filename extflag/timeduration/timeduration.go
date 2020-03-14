// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

// Taken from Thanos project.
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package timeduration

import (
	"bytes"
	"fmt"
	"time"
	"unsafe"

	"github.com/bwplotka/flagarize"
	"github.com/bwplotka/flagarize/extflag/internal/timestamp"
	"github.com/prometheus/common/model"
)

// multiError is a slice of errors implementing the error interface. It is used
// by a Gatherer to report multiple errors during MetricFamily gathering.
type multiError []error

func (errs multiError) Error() string {
	if len(errs) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d error(s) occurred:", len(errs))
	for _, err := range errs {
		fmt.Fprintf(buf, "\n* %s", err)
	}
	return buf.String()
}

// Append appends the provided error if it is not nil.
func (errs *multiError) Append(err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
}

// Value is a custom kingping parser for time in RFC3339
// or duration in Go's duration format, such as "300ms", "-1.5h" or "2h45m".
// Only one will be set.
type Value struct {
	Time *time.Time
	Dur  *model.Duration
}

// Set converts string to Value.
func (tdv *Value) Set(s string) error {
	var merr multiError
	t, err := time.Parse(time.RFC3339, s)
	if err == nil {
		tdv.Time = &t
		return nil
	}
	merr.Append(err)

	// error parsing time, let's try duration.
	var minus bool
	if s[0] == '-' {
		minus = true
		s = s[1:]
	}
	dur, err := model.ParseDuration(s)
	if err != nil {
		merr.Append(err)
		return merr
	}

	if minus {
		dur = dur * -1
	}
	tdv.Dur = &dur
	return nil
}

// String returns either time or duration.
func (tdv *Value) String() string {
	switch {
	case tdv.Time != nil:
		return tdv.Time.String()
	case tdv.Dur != nil:
		return tdv.Dur.String()
	}

	return "nil"
}

// PrometheusTimestamp returns Value converted to PrometheusTimestamp
// if duration is set now+duration is converted to Timestamp.
func (tdv *Value) PrometheusTimestamp() int64 {
	switch {
	case tdv.Time != nil:
		return timestamp.FromTime(*tdv.Time)
	case tdv.Dur != nil:
		return timestamp.FromTime(time.Now().Add(time.Duration(*tdv.Dur)))
	}

	return 0
}

// Flagarize registers PathOrContent flag.
func (tdv *Value) Flagarize(r flagarize.FlagRegisterer, tag *flagarize.Tag, _ unsafe.Pointer) error {
	if tag == nil {
		return nil
	}
	tag.Flag(r).SetValue(tdv)
	return nil
}
