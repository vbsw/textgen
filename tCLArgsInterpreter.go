/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main


type tCLArgsInterpreter struct {
	clArgs *tCLArgs
}

func newCLArgsInterpreter() *tCLArgsInterpreter {
	clArgsInterpreter := new(tCLArgsInterpreter)
	return clArgsInterpreter
}

func (this *tCLArgsInterpreter) setCLArgs(clArgs *tCLArgs) {
	this.clArgs = clArgs
}

func (this *tCLArgsInterpreter) isValidInfo() bool {
	for _, result := range this.clArgs.results {
		if result == result_UNKNOWN_ARGUMENT || result == result_WRONG_COMBINATION {
			return false
		}
	}
	return len(this.clArgs.results) == 1 || isResultInfo(this.clArgs.results[1])
}

func (this *tCLArgsInterpreter) isValidGenerate() bool {
	if len(this.clArgs.results) > 2 {
		if this.clArgs.results[1] == result_GENERATE_SIZE {
			if this.clArgs.results[2] == result_OUTPUT_STD || this.clArgs.results[2] == result_OUTPUT_FILE {
				if len(this.clArgs.results) > 3 {
					for _, result := range this.clArgs.results[3:] {
						if result == result_UNKNOWN_ARGUMENT || result == result_WRONG_COMBINATION {
							return false
						}
					}
				}
				return true
			}
		}
	}
	return false
}

func (this *tCLArgsInterpreter) isOutputStd() bool {
	if len(this.clArgs.results) > 2 {
		for _, result := range this.clArgs.results[2:] {
			if result == result_OUTPUT_STD {
				return true
			}
		}
	}
	return false
}

func (this *tCLArgsInterpreter) outputFileName() string {
	return this.clArgs.outputFileName
}

func (this *tCLArgsInterpreter) isHelp() bool {
	if len(this.clArgs.results) == 1 {
		return true
	} else if len(this.clArgs.results) == 2 {
		return this.clArgs.results[1] == result_HELP
	} else {
		return false
	}
}

func (this *tCLArgsInterpreter) isVersion() bool {
	if len(this.clArgs.results) == 2 {
		return this.clArgs.results[1] == result_VERSION
	} else {
		return false
	}
}

func (this *tCLArgsInterpreter) isCopyright() bool {
	if len(this.clArgs.results) == 2 {
		return this.clArgs.results[1] == result_COPYRIGHT
	} else {
		return false
	}
}
