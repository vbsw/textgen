/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


type tStdWriter struct {
}

func newStdWriter() *tStdWriter {
	stdWriter := new(tStdWriter)
	return stdWriter
}

func (this *tStdWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
