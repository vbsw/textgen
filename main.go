/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"os"
	"io"
	"github.com/vbsw/semver"
)

func Version() semver.Version {
	version := semver.New(0,1,0)
	return version
}

func main() {
	clArgs := newCLArgs()
	clArgsInterpreter := newCLArgsInterpreter()
	clArgs.parse(os.Args)
	clArgsInterpreter.setCLArgs(clArgs)

	if clArgsInterpreter.isValidInfo() {
		infoWriter := newInfoWriter()
		infoWriter.setInterpreter(clArgsInterpreter)
		infoWriter.writeInfo()

	} else if clArgsInterpreter.isValidGenerate() {
		writer := createWriterFromInterpreter(clArgsInterpreter)
		generator := newGenerator()
		generator.setWriter(writer)
		generator.generate()

	} else {
		errorWriter := newErrorWriter()
		errorWriter.setInterpreter(clArgsInterpreter)
		errorWriter.writeError()
	}
}

func createWriterFromInterpreter(clArgsInterpreter *tCLArgsInterpreter) io.Writer {
		if clArgsInterpreter.isOutputStd() {
			stdWriter := newStdWriter()
			return io.Writer(stdWriter)

		} else {
			fileWriter := newFileWriter()
			fileWriter.setOutputFileName(clArgsInterpreter.outputFileName())
			return io.Writer(fileWriter)
		}
}
