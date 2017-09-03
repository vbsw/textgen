/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"io"
)

type tGenerator struct {
	writer io.Writer
}

func newGenerator() *tGenerator {
	generator := new(tGenerator)
	return generator
}

func (this *tGenerator) setWriter(writer io.Writer) {
	this.writer = writer
}

func (this *tGenerator) generate() {
}
