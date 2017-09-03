/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


type tFileWriter struct {
	fileName string
}

func newFileWriter() *tFileWriter {
	fileWriter := new(tFileWriter)
	return fileWriter
}

func (this *tFileWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (this *tFileWriter) setOutputFileName(fileName string) {
	this.fileName = fileName
}
