/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


const (
	result_PROG_NAME int = iota
	result_HELP
	result_VERSION
	result_COPYRIGHT
	result_GENERATE_SIZE
	result_OUTPUT_STD
	result_OUTPUT_FILE
	result_THREADS_COUNT
	result_CORES_COUNT
	result_UNKNOWN_ARGUMENT
	result_WRONG_COMBINATION
)

type tCLArgs struct {
	results []int
	outputFileName string
	unknowns []string
	progName string
}

func newCLArgs() *tCLArgs {
	clArgs := new(tCLArgs)
	return clArgs
}

func (this *tCLArgs) parse(args []string) {
	this.results = make([]int, len(args))
	this.results[0] = result_PROG_NAME

	for i:=1; i<len(args); i+=1 {
		arg := args[i]

		if parseHelp(arg) {
			this.results[i] = result_HELP
		} else if parseVersion(arg) {
			this.results[i] = result_VERSION
		} else if parseCopyright(arg) {
			this.results[i] = result_COPYRIGHT
		} else {
			this.results[i] = result_UNKNOWN_ARGUMENT
		}
	}
	this.flagDoubles(result_HELP)
	this.flagDoubles(result_VERSION)
	this.flagDoubles(result_COPYRIGHT)
	this.flagDoubles(result_OUTPUT_STD)
	this.flagDoubles(result_OUTPUT_FILE)
	this.flagDoubles(result_THREADS_COUNT)
	this.flagDoubles(result_CORES_COUNT)
	this.flagInfoAndGenerateIncompatibility()
}

func (this *tCLArgs) flagDoubles(result int) {
	alreadyExists := false

	for i, resultCurr := range this.results {
		if resultCurr == result {
			if alreadyExists {
				this.results[i] = result_WRONG_COMBINATION
			} else {
				alreadyExists = true
			}
		}
	}
}

func (this *tCLArgs) flagInfoAndGenerateIncompatibility() {
	if len(this.results) > 2 {
		if isResultInfo(this.results[1]) {
			for i, result := range this.results[2:] {
				if result != result_UNKNOWN_ARGUMENT {
					this.results[i] = result_WRONG_COMBINATION
				}
			}
		} else {
			for i, result := range this.results[2:] {
				if isResultInfo(result) {
					this.results[i] = result_WRONG_COMBINATION
				}
			}
		}
	} else if len(this.results) > 1 {
		if !isResultInfo(this.results[1]) {
			this.results[1] = result_WRONG_COMBINATION
		}
	}
}

func parseHelp(argument string) bool {
	return argument == "--help" || argument == "-help" || argument == "-h" || argument == "--usage" || argument == "-usage"
}

func parseVersion(argument string) bool {
	return argument == "--version" || argument == "-version" || argument == "-v"
}

func parseCopyright(argument string) bool {
	return argument == "--copyright" || argument == "-copyright" || argument == "-c"
}

func isResultInfo(result int) bool {
	return result == result_HELP || result == result_VERSION || result == result_COPYRIGHT
}
