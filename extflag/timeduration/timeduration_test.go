// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

// Taken from Thanos project.
//
// Copyright (c) The Thanos Authors.
// Licensed under the Apache License 2.0.

package timeduration_test

//func TestTimeOrDurationValue(t *testing.T) {
//	cmd := kingpin.New("test", "test")
//
//	minTime := timeduration.TimeOrDuration(cmd.Flag("min-time", "Start of time range limit to serve"))
//
//	maxTime := timeduration.TimeOrDuration(cmd.Flag("max-time", "End of time range limit to serve").
//		Default("9999-12-31T23:59:59Z"))
//
//	_, err := cmd.Parse([]string{"--min-time", "10s"})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	testutil.Equals(t, "10s", minTime.String())
//	testutil.Equals(t, "9999-12-31 23:59:59 +0000 UTC", maxTime.String())
//
//	prevTime := timestamp.FromTime(time.Now())
//	afterTime := timestamp.FromTime(time.Now().Add(15 * time.Second))
//
//	testutil.Assert(t, minTime.PrometheusTimestamp() > prevTime, "minTime prometheus timestamp is less than time now.")
//	testutil.Assert(t, minTime.PrometheusTimestamp() < afterTime, "minTime prometheus timestamp is more than time now + 15s")
//
//	testutil.Assert(t, maxTime.PrometheusTimestamp() == 253402300799000, "maxTime is not equal to 253402300799000")
//}
