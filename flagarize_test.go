// Copyright (c) Bartłomiej Płotka @bwplotka
// Licensed under the Apache License 2.0.

package flagarize

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/bwplotka/flagarize/testutil"
	"github.com/pkg/errors"
)

func TestParseTag(t *testing.T) {
	t.Run("default separator", func(t *testing.T) {
		testParseTag(t, struct {
			noTag         bool
			wrongNoHelp1  bool `flagarize:""`
			wrongFormat1  bool `flagarize:"wrong"`
			wrongFormat2  bool `flagarize:"|"`
			wrongFormat3  bool `flagarize:"help=|"`
			wrongNoHelp3  bool `flagarize:"help="`
			noName        bool `flagarize:"help=help"`
			noName2       bool `flagarize:"name=|help=help"`
			simple        bool `flagarize:"name=case1|help=help"`
			simple_       string
			complexHelp   bool `flagarize:"name=case2a"`
			complexHelp_  string
			complexHelp1  bool `flagarize:"name=case2b|help="`
			complexHelp1_ string
			wrongFormat4  bool `flagarize:"name=...|help=help|nonexistingfield="`
			wrongFormat5  bool `flagarize:"name=...|help=help|wrongformat"`
			hidden        bool `flagarize:"name=case3|help=help|hidden=true"`
			required      bool `flagarize:"name=case4|help=help|required=true"`
			wrongNoHelp6  bool `flagarize:"name=...|hidden=true|required=true"`
			defaultValue  bool `flagarize:"name=case5|help=help|default=some default value ms 2213"`
			envVarWrong   bool `flagarize:"name=case6|help=help|envvar=lowerCASEnotallowed"`
			envVar        bool `flagarize:"name=case6|help=help|envvar=SOME_ENVVAR"`
			shortWrong    bool `flagarize:"name=...|help=help|short=tooLong"`
			short         bool `flagarize:"name=case7|help=help|short=l"`
			placeHolder   bool `flagarize:"name=case8|help=help|placeholder=<something>"`
			all           bool `flagarize:"name=case9|help=help|hidden=true|required=true|default=some|envvar=LOL|short=z|placeholder=<something2>"`
		}{
			// Most of those fields have default value. This could be skipped but is needed for (unused) lint.
			noTag:         false,
			wrongNoHelp1:  false,
			wrongFormat1:  false,
			wrongFormat2:  false,
			wrongFormat3:  false,
			wrongNoHelp3:  false,
			noName:        false,
			noName2:       false,
			simple:        false,
			simple_:       "yolo; this help should not be used, it was already specified.",
			complexHelp:   false,
			complexHelp_:  fmt.Sprintf("Some runtime evaluated help in %s.", flagTagName),
			complexHelp1:  false,
			complexHelp1_: fmt.Sprintf("Some runtime evaluated help2 in %s.", flagTagName),
			wrongFormat4:  false,
			wrongFormat5:  false,
			hidden:        false,
			required:      false,
			wrongNoHelp6:  false,
			defaultValue:  false,
			envVarWrong:   false,
			envVar:        false,
			shortWrong:    false,
			short:         false,
			placeHolder:   false,
			all:           false,
		}, "|")
	})
	t.Run("comma separator", func(t *testing.T) {
		testParseTag(t, struct {
			noTag         bool
			wrongNoHelp1  bool `flagarize:""`
			wrongFormat1  bool `flagarize:"wrong"`
			wrongFormat2  bool `flagarize:","`
			wrongFormat3  bool `flagarize:"help=,"`
			wrongNoHelp3  bool `flagarize:"help="`
			noName        bool `flagarize:"help=help"`
			noName2       bool `flagarize:"name=,help=help"`
			simple        bool `flagarize:"name=case1,help=help"`
			simple_       string
			complexHelp   bool `flagarize:"name=case2a"`
			complexHelp_  string
			complexHelp1  bool `flagarize:"name=case2b,help="`
			complexHelp1_ string
			wrongFormat4  bool `flagarize:"name=...,help=help,nonexistingfield="`
			wrongFormat5  bool `flagarize:"name=...,help=help,wrongformat"`
			hidden        bool `flagarize:"name=case3,help=help,hidden=true"`
			required      bool `flagarize:"name=case4,help=help,required=true"`
			wrongNoHelp6  bool `flagarize:"name=...,hidden=true,required=true"`
			defaultValue  bool `flagarize:"name=case5,help=help,default=some default value ms 2213"`
			envVarWrong   bool `flagarize:"name=case6,help=help,envvar=lowerCASEnotallowed"`
			envVar        bool `flagarize:"name=case6,help=help,envvar=SOME_ENVVAR"`
			shortWrong    bool `flagarize:"name=...,help=help,short=tooLong"`
			short         bool `flagarize:"name=case7,help=help,short=l"`
			placeHolder   bool `flagarize:"name=case8,help=help,placeholder=<something>"`
			all           bool `flagarize:"name=case9,help=help,hidden=true,required=true,default=some,envvar=LOL,short=z,placeholder=<something2>"`
		}{
			// Most of those fields have default value. This could be skipped but is needed for (unused) lint.
			noTag:         false,
			wrongNoHelp1:  false,
			wrongFormat1:  false,
			wrongFormat2:  false,
			wrongFormat3:  false,
			wrongNoHelp3:  false,
			noName:        false,
			noName2:       false,
			simple:        false,
			simple_:       "yolo; this help should not be used, it was already specified.",
			complexHelp:   false,
			complexHelp_:  fmt.Sprintf("Some runtime evaluated help in %s.", flagTagName),
			complexHelp1:  false,
			complexHelp1_: fmt.Sprintf("Some runtime evaluated help2 in %s.", flagTagName),
			wrongFormat4:  false,
			wrongFormat5:  false,
			hidden:        false,
			required:      false,
			wrongNoHelp6:  false,
			defaultValue:  false,
			envVarWrong:   false,
			envVar:        false,
			shortWrong:    false,
			short:         false,
			placeHolder:   false,
			all:           false,
		}, ",")
	})
}
func testParseTag(t *testing.T, input interface{}, sep string) {
	expected := []struct {
		tag *Tag
		err error
	}{
		{},
		{err: errors.New("flagarize: no help=<help> in struct Tag for field \"wrongNoHelp1\" and no help var; help=<help> in struct Tag or \"wrongNoHelp1_\" is required for help/usage of the flag; be helpful! :)")},
		{err: errors.New("flagarize: expected map-like Tag elements (e.g hidden=true), found non supported format \"wrong\" for field \"wrongFormat1\"")},
		{err: errors.New("flagarize: expected map-like Tag elements (e.g hidden=true), found non supported format \"\" for field \"wrongFormat2\"")},
		{err: errors.New("flagarize: expected map-like Tag elements (e.g hidden=true), found non supported format \"\" for field \"wrongFormat3\"")},
		{err: errors.New("flagarize: no help=<help> in struct Tag for field \"wrongNoHelp3\" and no help var; help=<help> in struct Tag or \"wrongNoHelp3_\" is required for help/usage of the flag; be helpful! :)")},
		{tag: &Tag{Name: "no_name", Help: "help"}},
		{tag: &Tag{Name: "no_name2", Help: "help"}},
		{tag: &Tag{Name: "case1", Help: "help"}},
		{},
		{tag: &Tag{Name: "case2a", Help: "Some runtime evaluated help in flagarize."}},
		{},
		{tag: &Tag{Name: "case2b", Help: "Some runtime evaluated help2 in flagarize."}},
		{},
		{err: errors.Errorf("flagarize: expected map-like Tag elements (e.g hidden=true) separated with %s, found but no supported key found \"nonexistingfield\" for field \"wrongFormat4\"; only [name help hidden required default envvar short placeholder] are supported", sep)},
		{err: errors.New("flagarize: expected map-like Tag elements (e.g hidden=true), found non supported format \"wrongformat\" for field \"wrongFormat5\"")},
		{tag: &Tag{Name: "case3", Help: "help", Hidden: true}},
		{tag: &Tag{Name: "case4", Help: "help", Required: true}},
		{err: errors.New("flagarize: no help=<help> in struct Tag for field \"wrongNoHelp6\" and no help var; help=<help> in struct Tag or \"wrongNoHelp6_\" is required for help/usage of the flag; be helpful! :)")},
		{tag: &Tag{Name: "case5", Help: "help", DefaultValue: "some default value ms 2213"}},
		{err: errors.New("flagarize: environment variable name has to be upper case, but it's not \"lowerCASEnotallowed\" for field \"envVarWrong\"")},
		{tag: &Tag{Name: "case6", Help: "help", EnvName: "SOME_ENVVAR"}},
		{err: errors.New("flagarize: short cannot be longer than one character got \"tooLong\" for field \"shortWrong\"")},
		{tag: &Tag{Name: "case7", Help: "help", Short: 'l'}},
		{tag: &Tag{Name: "case8", Help: "help", PlaceHolder: "<something>"}},
		{tag: &Tag{Name: "case9", Help: "help", Required: true, Hidden: true, Short: 'z', EnvName: "LOL", DefaultValue: "some", PlaceHolder: "<something2>"}},
	}

	val := reflect.ValueOf(input)

	helpVars := parseHelpVars(val)
	if val.NumField() != len(expected) {
		t.Fatalf("Different number of fields than elements in expected slice")
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)

		t.Run(field.Name, func(t *testing.T) {
			tag, err := parseTag(field, helpVars[field.Name], sep)
			if expected[i].err != nil {
				testutil.NotOk(t, err)
				testutil.Equals(t, expected[i].err.Error(), err.Error())
				return
			}
			testutil.Ok(t, err)
			testutil.Equals(t, expected[i].tag, tag)
		})
	}
}
