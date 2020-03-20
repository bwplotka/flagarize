// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

package flagarize_test

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os"
	"testing"
	"time"
	"unsafe"

	"github.com/alecthomas/units"
	"github.com/bwplotka/flagarize"
	"github.com/bwplotka/flagarize/extflag/pathorcontent"
	"github.com/bwplotka/flagarize/extflag/timeduration"
	"github.com/bwplotka/flagarize/testutil"
	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
	"gopkg.in/alecthomas/kingpin.v2"
)

type customDuration model.Duration

func (d *customDuration) Flagarize(r flagarize.FlagRegisterer, tag *flagarize.Tag, ptr unsafe.Pointer) error {
	if tag == nil {
		return nil
	}
	m := (*model.Duration)(ptr)
	tag.Flag(r).DurationVar((*time.Duration)(m))
	return nil
}

type nonReceiveFlagarizeType model.Duration

func (t nonReceiveFlagarizeType) Flagarize(r flagarize.FlagRegisterer, tag *flagarize.Tag, ptr unsafe.Pointer) error {
	if tag == nil {
		return nil
	}
	m := (*model.Duration)(ptr)
	tag.Flag(r).DurationVar((*time.Duration)(m))
	return nil
}

type failFlagarizeType struct{}

func (*failFlagarizeType) Flagarize(flagarize.FlagRegisterer, *flagarize.Tag, unsafe.Pointer) error {
	return errors.New("fail")
}

func newTestKingpin(t *testing.T) *kingpin.Application {
	app := kingpin.New("test", "test")
	app.Terminate(func(code int) { t.Fatal("kingping terminates with code", code) })
	return app
}

func TestFlagarize(t *testing.T) {
	t.Run("flagarize on private field", func(t *testing.T) {
		type wrong struct {
			f map[string]int `flagarize:"help=help"`
		}
		w := &wrong{f: nil}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on private field \"f\"; it has to be exported", err.Error())
	})
	t.Run("flagarize on not supported field: map", func(t *testing.T) {
		type wrong struct {
			F map[string]int `flagarize:"help=help"`
		}
		w := &wrong{F: nil}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on not supported type map map[string]int for field \"F\"", err.Error())
	})
	t.Run("flagarize on pointer for standard type", func(t *testing.T) {
		type wrong struct {
			F *string `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on not supported type ptr *string for field \"F\"", err.Error())
	})
	t.Run("flagarize on custom struct that does not have flagarizer method", func(t *testing.T) {
		type wrong struct {
			F struct{} `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on not supported type struct struct {} for field \"F\"", err.Error())
	})
	t.Run("flagarize on interface{}", func(t *testing.T) {
		type wrong struct {
			F interface{} `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on not supported type interface <nil> for field \"F\"", err.Error())
	})
	t.Run("flagarize on *interface{}", func(t *testing.T) {
		type wrong struct {
			F *interface{} `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize struct Tag found on not supported type ptr *interface {} for field \"F\"", err.Error())
	})
	t.Run("custom failing flagarizer", func(t *testing.T) {
		type wrong struct {
			F *failFlagarizeType `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: custom Flagarizer for field F: fail", err.Error())
	})
	t.Run("custom non receiver pointer flagarizer", func(t *testing.T) {
		type wrong struct {
			F *nonReceiveFlagarizeType `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize field \"F\" custom Flagarizer is non receiver pointer", err.Error())
	})
	t.Run("custom non pointer flagarizer", func(t *testing.T) {
		type wrong struct {
			F customDuration `flagarize:"help=help"`
		}
		w := &wrong{}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, w)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: flagarize field \"F\" is not a pointer, but custom Flagarizer was used", err.Error())
	})

	type testConfig struct {
		Ignore    int
		F1        bool     `flagarize:"help=1"`
		F1Slice   []bool   `flagarize:"help=1slice"`
		F2        string   `flagarize:"help=2"`
		F2Slice   []string `flagarize:"help=2slice"`
		F4        int      `flagarize:"help=3"`
		F4Slice   []int    `flagarize:"help=3slice"`
		F5        int8     `flagarize:"help=4"`
		F5Slice   []int8   `flagarize:"help=4slice"`
		F6        int16    `flagarize:"help=5"`
		F6Slice   []int16  `flagarize:"help=5slice"`
		F7        int32    `flagarize:"help=6"`
		F7Slice   []int32  `flagarize:"help=6slice"`
		F8        int64    `flagarize:"help=7"`
		F8Slice   []int64  `flagarize:"help=7slice"`
		F9        uint     `flagarize:"help=8"`
		F9Slice   []uint   `flagarize:"help=8slice"`
		F10       uint8    `flagarize:"help=9"`
		F10Slice  []uint8  `flagarize:"help=9slice"`
		F11       uint16   `flagarize:"help=10"`
		F11Slice  []uint16 `flagarize:"help=10slice"`
		F12       uint32   `flagarize:""`
		F12_      string
		F12Slice  []uint32 `flagarize:""`
		F12Slice_ string
		F13       uint64            `flagarize:"help=11"`
		F13Slice  []uint64          `flagarize:"help=11slice"`
		F14       float32           `flagarize:"help=12"`
		F14Slice  []float32         `flagarize:"help=12slice"`
		F15       float64           `flagarize:"help=13"`
		F15Slice  []float64         `flagarize:"help=13slice"`
		F17       map[string]string `flagarize:"help=14"`
		F18       *net.TCPAddr      `flagarize:"help=15"`
		F18Slice  []*net.TCPAddr    `flagarize:"help=15slice"`
		F19       *url.URL          `flagarize:"help=16"`
		F19Slice  []*url.URL        `flagarize:"help=16slice"`
		F20       *os.File          `flagarize:"help=17"`
		F21       time.Duration     `flagarize:"help=18"`
		F21Slice  []time.Duration   `flagarize:"help=18slice"`
		F22       net.IP            `flagarize:"help=19"`
		F22Slice  []net.IP          `flagarize:"help=19slice"`
		F23       units.Base2Bytes  `flagarize:"help=20"`

		Cf1 *customDuration              `flagarize:"help=21"`
		Cf2 *timeduration.Value          `flagarize:"help=22"`
		Cf3 *pathorcontent.PathOrContent `flagarize:"help=23"`
	}
	const expectedUsage = `usage: test [<flags>]

test

Flags:
  --help                     Show context-sensitive help (also try --help-long
                             and --help-man).
  --f1                       1
  --f1_slice=F1_SLICE ...    1slice
  --f2=F2                    2
  --f2_slice=F2_SLICE ...    2slice
  --f4=F4                    3
  --f4_slice=F4_SLICE ...    3slice
  --f5=F5                    4
  --f5_slice=F5_SLICE ...    4slice
  --f6=F6                    5
  --f6_slice=F6_SLICE ...    5slice
  --f7=F7                    6
  --f7_slice=F7_SLICE ...    6slice
  --f8=F8                    7
  --f8_slice=F8_SLICE ...    7slice
  --f9=F9                    8
  --f9_slice=F9_SLICE ...    8slice
  --f10=F10                  9
  --f10_slice=F10_SLICE ...  9slice
  --f11=F11                  10
  --f11_slice=F11_SLICE ...  10slice
  --f12=F12                  12
  --f12_slice=F12_SLICE ...  12Slice
  --f13=F13                  11
  --f13_slice=F13_SLICE ...  11slice
  --f14=F14                  12
  --f14_slice=F14_SLICE ...  12slice
  --f15=F15                  13
  --f15_slice=F15_SLICE ...  13slice
  --f17=F17 ...              14
  --f18=F18                  15
  --f18_slice=F18_SLICE ...  15slice
  --f19=F19                  16
  --f19_slice=F19_SLICE ...  16slice
  --f20=F20                  17
  --f21=F21                  18
  --f21_slice=F21_SLICE ...  18slice
  --f22=F22                  19
  --f22_slice=F22_SLICE ...  19slice
  --f23=F23                  20
  --cf1=CF1                  21
  --cf2=CF2                  22
  --cf3-file=<file-path>     Path to 23
  --cf3=<content>            Alternative to 'cf3-file' flag (lower priority).
                             Content of 23

`

	type testEmbeddedConfig testConfig
	type testParseConfig struct {
		testEmbeddedConfig

		// Private so ignored.
		config2 testConfig
		// Slice without flagarize or customer Flagarizer so ignored.
		Configs []testConfig
	}
	t.Run("nil pointer struct parsed", func(t *testing.T) {
		var c *testParseConfig

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, c)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: object cannot be nil", err.Error())
	})
	t.Run("non pointer struct parsed", func(t *testing.T) {
		c := testParseConfig{testEmbeddedConfig: testEmbeddedConfig{F12_: "12", F12Slice_: "12Slice"}, config2: testConfig{}, Configs: nil}

		app := newTestKingpin(t)
		err := flagarize.Flagarize(app, c)
		testutil.NotOk(t, err)
		testutil.Equals(t, "flagarize: object must be a pointer to struct or interface", err.Error())
	})

	t.Run("expected help message", func(t *testing.T) {
		c := &testParseConfig{testEmbeddedConfig: testEmbeddedConfig{F12_: "12", F12Slice_: "12Slice"}}
		app := newTestKingpin(t)
		b := bytes.Buffer{}
		app.UsageWriter(&b)

		var terminates bool
		app.Terminate(func(code int) { terminates = true })

		testutil.Ok(t, flagarize.Flagarize(app, c))
		_, err := app.Parse([]string{"--help"})
		testutil.Ok(t, err)
		testutil.Assert(t, terminates, "parse did not terminate")
		testutil.Equals(t, expectedUsage, b.String())
	})

	t.Run("expected help message", func(t *testing.T) {
		type testParseConfig2 struct {
			Config1 testConfig
		}
		c := &testParseConfig2{Config1: testConfig{F12_: "12", F12Slice_: "12Slice"}}

		app := newTestKingpin(t)
		b := bytes.Buffer{}
		app.UsageWriter(&b)

		var terminates bool
		app.Terminate(func(code int) { terminates = true })

		testutil.Ok(t, flagarize.Flagarize(app, c))
		_, err := app.Parse([]string{"--help"})
		testutil.Ok(t, err)
		testutil.Assert(t, terminates, "parse did not terminate")
		testutil.Equals(t, expectedUsage, b.String())
	})

	twoTwoFourDuration := customDuration(244 * time.Hour)
	var someString string
	fileLICENSEPath := "LICENSE"
	for _, tcase := range []struct {
		input    []string
		expected *testParseConfig
	}{
		{
			input: []string{},
			expected: &testParseConfig{testEmbeddedConfig: testEmbeddedConfig{
				F12_:      "12",
				F12Slice_: "12Slice",
				F17:       map[string]string{},
				Cf1:       new(customDuration),
				Cf2:       &timeduration.Value{},
				Cf3:       pathorcontent.New("cf3", false, &someString, &someString),
			}},
		},
		{
			input: []string{
				"--f1",
				"--f1_slice", "true", "--f1_slice", "FALSE", "--f1_slice", "true",
				"--f2", "testString",
				"--f2_slice", "a", "--f2_slice", "b",
				"--f4=-1234",
				"--f4_slice=-1", "--f4_slice=-424", "--f4_slice", "3",
				"--f5=-122",
				"--f5_slice=-12", "--f5_slice=-42", "--f5_slice", "32",
				"--f6=-12343",
				"--f6_slice=-13", "--f6_slice=-4243", "--f6_slice", "33",
				"--f7=-12344",
				"--f7_slice=-14", "--f7_slice=-4244", "--f7_slice", "34",
				"--f8=-12345",
				"--f8_slice=-15", "--f8_slice=-4245", "--f8_slice", "35",
				"--f9", "1234",
				"--f9_slice", "1", "--f9_slice", "424", "--f9_slice", "3",
				"--f10", "122",
				"--f10_slice", "12", "--f10_slice", "42", "--f10_slice", "32",
				"--f11", "12343",
				"--f11_slice", "13", "--f11_slice", "4243", "--f11_slice", "33",
				"--f12", "12344",
				"--f12_slice", "14", "--f12_slice", "4244", "--f12_slice", "34",
				"--f13", "12345",
				"--f13_slice", "15", "--f13_slice", "4245", "--f13_slice", "35",
				"--f14", "12.4",
				"--f14_slice=-2.1", "--f14_slice", "1.12345", "--f14_slice", "0",
				"--f15=-12.43265",
				"--f15_slice", "2.1", "--f15_slice=-1.12345", "--f15_slice", "0",
				"--f17", `a=b`,
				"--f18", "1.2.3.4:34",
				"--f18_slice", "1.2.3.5:35", "--f18_slice", "1.2.3.6:36", "--f18_slice", "1.2.3.7:37",
				"--f19", "http://example.com/path",
				"--f19_slice", "http://example.com30/path", "--f19_slice", "http://example.com10/path2", "--f19_slice", "http://example.com20/path3",
				"--f20", fileLICENSEPath,
				"--f21", "204s",
				"--f21_slice", "12000ns", "--f21_slice", "204234s", "--f21_slice", "2s",
				"--f22", "1.2.3.4",
				"--f22_slice", "1.2.3.5", "--f22_slice", "1.2.3.6", "--f22_slice", "1.2.3.7",
				"--f23", "232MB",
				"--cf1", "244h",
				"--cf2", "2020-03-18T12:01:33Z",
				"--cf3-file", fileLICENSEPath,
			},
			expected: &testParseConfig{testEmbeddedConfig: testEmbeddedConfig{
				F1:        true,
				F1Slice:   []bool{true, false, true},
				F2:        "testString",
				F2Slice:   []string{"a", "b"},
				F4:        -1234,
				F4Slice:   []int{-1, -424, 3},
				F5:        -122,
				F5Slice:   []int8{-12, -42, 32},
				F6:        -12343,
				F6Slice:   []int16{-13, -4243, 33},
				F7:        -12344,
				F7Slice:   []int32{-14, -4244, 34},
				F8:        -12345,
				F8Slice:   []int64{-15, -4245, 35},
				F9:        1234,
				F9Slice:   []uint{1, 424, 3},
				F10:       122,
				F10Slice:  []uint8{12, 42, 32},
				F11:       12343,
				F11Slice:  []uint16{13, 4243, 33},
				F12:       12344,
				F12_:      "12",
				F12Slice:  []uint32{14, 4244, 34},
				F12Slice_: "12Slice",
				F13:       12345,
				F13Slice:  []uint64{15, 4245, 35},
				F14:       12.4,
				F14Slice:  []float32{-2.1, 1.12345, 0},
				F15:       -12.43265,
				F15Slice:  []float64{2.1, -1.12345, 0},
				F17:       map[string]string{"a": "b"},
				F18:       &net.TCPAddr{IP: net.IPv4(0x1, 0x2, 0x3, 0x4), Port: 34},
				F18Slice:  []*net.TCPAddr{{IP: net.IPv4(0x1, 0x2, 0x3, 0x5), Port: 35}, {IP: net.IPv4(0x1, 0x2, 0x3, 0x6), Port: 36}, {IP: net.IPv4(0x1, 0x2, 0x3, 0x7), Port: 37}},
				F19:       &url.URL{Host: "example.com", Path: "/path", Scheme: "http"},
				F19Slice:  []*url.URL{{Host: "example.com30", Path: "/path", Scheme: "http"}, {Host: "example.com10", Path: "/path2", Scheme: "http"}, {Host: "example.com20", Path: "/path3", Scheme: "http"}},
				F20: func() *os.File {
					f, _ := os.OpenFile(fileLICENSEPath, os.O_RDONLY, 0)
					return f
				}(),
				F21:      204 * time.Second,
				F21Slice: []time.Duration{12000 * time.Nanosecond, 204234 * time.Second, 2 * time.Second},
				F22:      net.IPv4(0x1, 0x2, 0x3, 0x4),
				F22Slice: []net.IP{net.IPv4(0x1, 0x2, 0x3, 0x5), net.IPv4(0x1, 0x2, 0x3, 0x6), net.IPv4(0x1, 0x2, 0x3, 0x7)},
				F23:      units.Base2Bytes(232 * 1024 * 1024),
				Cf1:      &twoTwoFourDuration,
				Cf2: &timeduration.Value{
					Time: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2020-03-18T12:01:33Z"); return &t }(),
				},
				Cf3: pathorcontent.New("cf3", false, &fileLICENSEPath, &someString),
			}},
		},
	} {
		t.Run(fmt.Sprintf("%v", tcase.input), func(t *testing.T) {
			c := &testParseConfig{testEmbeddedConfig: testEmbeddedConfig{F12_: "12", F12Slice_: "12Slice"}}

			app := newTestKingpin(t)
			testutil.Ok(t, flagarize.Flagarize(app, c, flagarize.WithElemSep(",")))

			_, err := app.Parse(tcase.input)
			testutil.Ok(t, err)

			if tcase.expected.F20 != nil && c.F20 != nil {
				// Treat file specially. Equals does not compare properly.
				testutil.Equals(t, tcase.expected.F20.Name(), c.F20.Name())
				c.F20 = nil
				tcase.expected.F20 = nil
			}
			testutil.Equals(t, tcase.expected, c)
		})
	}
}
