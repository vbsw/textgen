/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import(
	"fmt"
)

type tStdWriter struct {
}

func newStdWriter() *tStdWriter {
	stdWriter := new(tStdWriter)
	return stdWriter
}

func (this *tStdWriter) Write(bytes []byte) (int, error) {
	fmt.Printf("%s",bytes)
	return len(bytes), nil
}
