/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"testing"
)

func TestParseNumber(t *testing.T) {
	if parseNumber("100") == false {
		t.Error()
	}
	if parseNumber("100K") == false {
		t.Error()
	}
	if parseNumber("100k") == false {
		t.Error()
	}
	if parseNumber("100M") == false {
		t.Error()
	}
	if parseNumber("100m") == false {
		t.Error()
	}
	if parseNumber("100G") == false {
		t.Error()
	}
	if parseNumber("100g") == false {
		t.Error()
	}
	if parseNumber("100g12") == true {
		t.Error()
	}
	if parseNumber("K100") == true {
		t.Error()
	}
	if parseNumber("10k785") == true {
		t.Error()
	}
	if parseNumber("1asdf") == true {
		t.Error()
	}
}
