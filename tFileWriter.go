/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


import (
	"os"
)

type tFileWriter struct {
	fileName string
	file *os.File
}

func newFileWriter() *tFileWriter {
	fileWriter := new(tFileWriter)
	return fileWriter
}

func (this *tFileWriter) setOutputFileName(fileName string) {
	this.fileName = fileName
}

func (this *tFileWriter) open() {
	var err error
	this.file, err = os.Create(this.fileName)

	if err != nil {
		panic(err)
	}
}

func (this *tFileWriter) close() {
	if err := this.file.Close(); err != nil {
		panic(err)
	}
}

func (this *tFileWriter) Write(bytes []byte) (int, error) {
	n, err := this.file.Write(bytes)
	return n, err
}
