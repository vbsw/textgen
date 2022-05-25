/*
 *          Copyright 2022, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"errors"
	"fmt"
	"github.com/vbsw/golib/osargs"
//	"io"
	"os"
	"path/filepath"
//	"strconv"
//	"strings"
//	"sync"
	"runtime"
)

type tParameters struct {
	help           *osargs.Result
	version        *osargs.Result
	example        *osargs.Result
	copyright      *osargs.Result
	size        *osargs.Result
	threads        *osargs.Result
	system         *osargs.Result
	output         *osargs.Result
	infoParams []*osargs.Result
	cmdParams  []*osargs.Result
}

func main() {
	var params tParameters
	err := params.initFromOSArgs()
	if err == nil {
		if params.infoAvailable() {
			printInfo(&params)
		} else {
		}
	} else {
		printError(err)
	}
}

func (params *tParameters) initFromOSArgs() error {
	args := osargs.New()
	err := params.initFromArgs(args)
	return err
}

// initFromArgs is for test purposes.
func (params *tParameters) initFromArgs(args *osargs.Arguments) error {
	var err error
	if len(args.Values) > 0 {
		delimiter := osargs.NewDelimiter(true, true, "=")
		params.help = args.Parse("-h", "--help", "-help", "help")
		params.version = args.Parse("-v", "--version", "-version", "version")
		params.example = args.Parse("-e", "--example", "-example", "example")
		params.copyright = args.Parse("-c", "--copyright", "-copyright", "copyright")
		params.size = args.ParsePairs(delimiter, "-s", "--size", "-size", "size")
		params.threads = args.ParsePairs(delimiter, "-t", "--threads", "-threads", "threads")
		params.system = args.ParsePairs(delimiter, "-y", "--system", "-system", "system")
		params.output = new(osargs.Result)
		params.poolInfoParams()
		params.poolCmdParams()

		unparsedArgs := args.UnparsedArgs()
		unparsedArgs = params.parseSize(unparsedArgs)
		unparsedArgs = params.parseOutput(unparsedArgs)

		err = params.validateParameters(unparsedArgs)
		if err == nil && !params.infoAvailable() {
			params.ensureThreads()
			params.ensureSystem()
		}
	}
	return err
}

func (params *tParameters) parseSize(unparsedArgs []string) []string {
	// just accept the first unparsed argument
	if !params.size.Available() && len(unparsedArgs) > 0 {
		params.size.Values = append(params.size.Values, unparsedArgs[0])
		return unparsedArgs[1:]
	}
	return unparsedArgs
}

func (params *tParameters) parseOutput(unparsedArgs []string) []string {
	if params.output.Available() {
		outputPath, err := filepath.Abs(params.output.Values[0])
		if err == nil {
			params.output.Values[0] = outputPath
		} else {
			panic(err.Error())
		}
	} else if len(unparsedArgs) > 0 {
		// just accept the first unparsed argument
		outputPath, err := filepath.Abs(unparsedArgs[0])
		if err == nil {
			params.output.Values = append(params.output.Values, outputPath)
		} else {
			panic(err.Error())
		}
		return unparsedArgs[1:]
	}
	return unparsedArgs
}

func (params *tParameters) validateParameters(unparsedArgs []string) error {
	var err error
	if len(unparsedArgs) > 0 {
		unknownArg := unparsedArgs[0]
		err = errors.New("unknown argument \"" + unknownArg + "\"")
	} else {
		if params.isCompatible() {
			if !params.infoAvailable() {
				if !params.size.Available() {
					err = errors.New("file size not specified")
				} else {
					err = params.validateIODirectories()
				}
			}
		} else {
			err = errors.New("wrong argument usage")
		}
	}
	return err
}

func (params *tParameters) ensureThreads() {
	if !params.threads.Available() {
		params.threads.Values = append(params.threads.Values, "1")
	}
}

func (params *tParameters) ensureSystem() {
	if !params.system.Available() {
		if runtime.GOOS == "windows" {
			params.system.Values = append(params.system.Values, "windows")
		} else {
			params.system.Values = append(params.system.Values, "unix")
		}
	}
}

func (params *tParameters) poolInfoParams() {
	params.infoParams = make([]*osargs.Result, 4)
	params.infoParams[0] = params.help
	params.infoParams[1] = params.version
	params.infoParams[2] = params.example
	params.infoParams[3] = params.copyright
}

func (params *tParameters) poolCmdParams() {
	params.cmdParams = make([]*osargs.Result, 4)
	params.cmdParams[0] = params.size
	params.cmdParams[1] = params.threads
	params.cmdParams[2] = params.system
	params.cmdParams[3] = params.output
}

func (params *tParameters) infoAvailable() bool {
	return anyAvailable(params.infoParams)
}

func (params *tParameters) validateIODirectories() error {
	var err error
	if !params.output.Available() {
		err = errors.New("output file is not specified")
	} else {
		_, errInfo := os.Stat(params.output.Values[0])
		if errInfo == nil || !os.IsNotExist(errInfo) {
			err = errors.New("output file exists already")
		}
	}
	return err
}

func (params *tParameters) isCompatible() bool {
	// same parameter must not be multiple
	if isMultiple(params.infoParams) || isMultiple(params.cmdParams) {
		return false
	}
	// either info or command
	if anyAvailable(params.infoParams) && anyAvailable(params.cmdParams) {
		return false
	}
	// no mixed info parameters
	if isMixed(params.infoParams...) {
		return false
	}
	return true
}

func anyAvailable(results []*osargs.Result) bool {
	for _, result := range results {
		if result.Available() {
			return true
		}
	}
	return false
}

func isMultiple(paramsMult []*osargs.Result) bool {
	for _, param := range paramsMult {
		if param.Count() > 1 {
			return true
		}
	}
	return false
}

func isMixed(params ...*osargs.Result) bool {
	for i, paramA := range params {
		if paramA.Available() {
			for _, paramB := range params[i+1:] {
				if paramB.Available() {
					return true
				}
			}
			break
		}
	}
	return false
}













func printInfo(params *tParameters) {
	if params.help == nil {
		printShortInfo()
	} else if params.help.Available() {
		printHelp()
	} else if params.version.Available() {
		printVersion()
	} else if params.example.Available() {
		printExample()
	} else if params.copyright.Available() {
		printCopyright()
	} else {
		printShortInfo()
	}
}

func printShortInfo() {
	fmt.Println("Run 'fsplit --help' for usage.")
}

func printHelp() {
	message := "\nUSAGE\n"
	message += "  fsplit ( INFO | SPLIT/CONCATENATE )\n\n"
	message += "INFO\n"
	message += "  -h, --help    print this help\n"
	message += "  -v, --version print version\n"
	message += "  --copyright   print copyright\n\n"
	message += "SPLIT/CONCATENATE\n"
	message += "  [COMMAND] INPUT-FILE [OUTPUT-FILE]\n\n"
	message += "COMMAND\n"
	message += "  -p=N          split file into N parts (chunks)\n"
	message += "  -b=N[U]       split file into N bytes per chunk, U = unit (k/K, m/M or g/G)\n"
	message += "  -l=N          split file into N lines per chunk\n"
	message += "  -c            concatenate files (INPUT-FILE is only one file, the first one)"
	fmt.Println(message)
}

func printVersion() {
	fmt.Println("1.0.3")
}

func printExample() {
	message := "\nEXAMPLES\n"
	message += "   ... not available"
	fmt.Println(message)
}

func printCopyright() {
	message := "Copyright 2019 - 2022, Vitali Baumtrok (vbsw@mailbox.org).\n"
	message += "Distributed under the Boost Software License, Version 1.0."
	fmt.Println(message)
}

func printError(err error) {
	fmt.Println("error:", err.Error())
}