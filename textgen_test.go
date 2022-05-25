/*
 *      Copyright 2021, 2022 Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"github.com/vbsw/golib/osargs"
	"testing"
)

func TestParseOSArgsA(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)
	err := params.initFromArgs(args)
	if err != nil {
		t.Error(err.Error())
	}

	args.Values = []string{"--help", "--version"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"--help", "-t=2"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}
}

func TestParseOSArgsB(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{"100k"}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)

	err := params.initFromArgs(args)
	if err == nil {
		t.Error("unspecified out file not recognized")
	}

	args.Values = []string{"100k"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("unspecified input directory not recognized")
	}

	args.Values = []string{"100k", "./"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("existing file not recognized")
	}
}

func TestParseOSArgsC(t *testing.T) {
	args := new(osargs.Arguments)
	args.Values = []string{"100k", "--help"}
	args.Parsed = make([]bool, len(args.Values))
	params := new(tParameters)

	err := params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"-t=2", "-t=5"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"--help", "100k", "./a.txt"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err == nil {
		t.Error("incompatible parameters not recognized")
	}

	args.Values = []string{"-h"}
	args.Parsed = make([]bool, len(args.Values))
	params = new(tParameters)
	err = params.initFromArgs(args)
	if err != nil {
		t.Error("valid parameter not recognized: " + err.Error())
	}
}
