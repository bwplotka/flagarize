// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

// Taken from Thanos project.
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package flagarize_test

import (
	"testing"
	"time"

	"github.com/bwplotka/flagarize"
	"github.com/bwplotka/flagarize/internal/timestamp"
	"github.com/bwplotka/flagarize/testutil"
)

func TestTimeOrDuration(t *testing.T) {
	minTime := &flagarize.TimeOrDuration{}
	testutil.Ok(t, minTime.Set("10s"))
	maxTime := &flagarize.TimeOrDuration{}
	testutil.Ok(t, maxTime.Set("9999-12-31T23:59:59Z"))

	testutil.Equals(t, "10s", minTime.String())
	testutil.Equals(t, "9999-12-31 23:59:59 +0000 UTC", maxTime.String())

	prevTime := timestamp.FromTime(time.Now())
	afterTime := timestamp.FromTime(time.Now().Add(15 * time.Second))

	testutil.Assert(t, minTime.PrometheusTimestamp() > prevTime, "minTime prometheus timestamp is less than time now.")
	testutil.Assert(t, minTime.PrometheusTimestamp() < afterTime, "minTime prometheus timestamp is more than time now + 15s")

	testutil.Assert(t, maxTime.PrometheusTimestamp() == 253402300799000, "maxTime is not equal to 253402300799000")
}
