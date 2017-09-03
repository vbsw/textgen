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

type tInfoWriter struct {
	clArgsInterpreter *tCLArgsInterpreter
}

func newInfoWriter() *tInfoWriter {
	infoWriter := new(tInfoWriter)
	return infoWriter
}

func (this *tInfoWriter) setInterpreter(clArgsInterpreter *tCLArgsInterpreter) {
	this.clArgsInterpreter = clArgsInterpreter
}

func (this *tInfoWriter) writeInfo() {
	if this.clArgsInterpreter.isHelp() {
		this.printHelp()
	} else if this.clArgsInterpreter.isVersion() {
		this.printVersion()
	} else if this.clArgsInterpreter.isCopyright() {
		this.printCopyright()
	} else {
		this.printUnforeseenState()
	}
}

func (this *tInfoWriter) printHelp() {
	fmt.Println()
	fmt.Println("USAGE")
	fmt.Println("  textgen [INFO-OPTION]")
	fmt.Println("  textgen SIZE OUTPUT {GENERATOR-OPTION}")
	fmt.Println("INFO-OPTION")
	fmt.Println("  -help       prints this help")
	fmt.Println("  -version    prints version numer of textgen")
	fmt.Println("  -copyright  prints copyright of textgen")
	fmt.Println("GENERATOR-OPTION")
	fmt.Println("  -cN         sets the number of logical CPUs to N")
	fmt.Println("  -tN         sets the number of threads to N")
	fmt.Println("OPTION")
	fmt.Println("  std         prints to standard output (i.e. terminal)")
	fmt.Println("  <file name> prints to file")
	fmt.Println("Examples")
	fmt.Println("  textgen 100 std        prints 100 bytes of text to standard output")
	fmt.Println("  textgen 100K test.txt  prints 100 kilobytes of text to file test.txt")
}

func (this *tInfoWriter) printVersion() {
	version := Version()
	fmt.Println(version)
}

func (this *tInfoWriter) printCopyright() {
	fmt.Println("Copyright 2017, Vitali Baumtrok (vbsw@mailbox.org).")
	fmt.Println("Text Generator is distributed under the Boost Software License, version 1.0.")
	fmt.Println("(See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)")
}

func (this *tInfoWriter) printUnforeseenState() {
	fmt.Println("error: unforseen state in printing info")
}
