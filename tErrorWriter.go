/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


import (
	"fmt"
)

type tErrorWriter struct {
	clArgsInterpreter *tCLArgsInterpreter
}

func newErrorWriter() *tErrorWriter {
	errorWriter := new(tErrorWriter)
	return errorWriter
}

func (this *tErrorWriter) setInterpreter(clArgsInterpreter *tCLArgsInterpreter) {
	this.clArgsInterpreter = clArgsInterpreter
}

func (this *tErrorWriter) writeError() {
	fmt.Println("error: unknown")
}
