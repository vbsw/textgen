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
	"strconv"
	"math/rand"
	"time"
//	"strings"
//	"sync"
	"runtime"
)

const (
	wordMIN_LEN = 2
	wordMAX_LEN = 30
	newLinePROBABILITY = 0.1
	wordsMAX_PER_LINE = 20
)

type tParameters struct {
	help           *osargs.Result
	version        *osargs.Result
	example        *osargs.Result
	copyright      *osargs.Result
	size        *osargs.Result
	threads        *osargs.Result
	system         *osargs.Result
	buffer         *osargs.Result
	output         *osargs.Result
	infoParams []*osargs.Result
	cmdParams  []*osargs.Result
}

func main() {
	params := new(tParameters)
	err := params.initFromOSArgs()
	if err == nil {
		if params.infoAvailable() {
			printInfo(params)
		} else {
			var size, threads, buffer int
			lineSeparator := interpretLineSeparator(params)
			size, err = interpretSize(params, err)
			threads, err = interpretThreads(params, err)
			buffer, err = interpretBuffer(params, len(lineSeparator) + 1, err)
			if err == nil {
				if threads == 1 {
					err = generate(params, size, buffer, lineSeparator)
				} else {
					err = generateGo(params, size, buffer, threads, lineSeparator)
				}
			}
		}
	}
	if err != nil {
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
		params.buffer = args.ParsePairs(delimiter, "-b", "--buffer", "-buffer", "buffer")
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
	params.cmdParams = make([]*osargs.Result, 5)
	params.cmdParams[0] = params.size
	params.cmdParams[1] = params.threads
	params.cmdParams[2] = params.system
	params.cmdParams[3] = params.buffer
	params.cmdParams[4] = params.output
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

func interpretSize(params *tParameters, err error) (int, error) {
	if err == nil {
		var size int
		size, err = parseBytes(params.size.Values[0])
		if err == nil || size <= 0 {
			return size, nil
		}
		return 0, errors.New("can't parse output file size")
	}
	return 0, err
}

func interpretThreads(params *tParameters, err error) (int, error) {
	if err == nil {
		threads, err := strconv.Atoi(params.threads.Values[0])
		if err == nil {
			if threads > 0 {
				return threads, nil
			}
			return 1, nil
		}
		return 0, errors.New("can't parse number of threads")
	}
	return 0, err
}

func interpretBuffer(params *tParameters, sizeMin int, err error) (int, error) {
	if err == nil {
		if params.buffer.Available() {
			bytes, err := parseBytes(params.buffer.Values[0])
			if err == nil {
				if bytes > 0 {
					if bytes >= sizeMin {
						return bytes, nil
					}
					return sizeMin, nil
				}
			} else {
				return 0, errors.New("can't parse size of buffer")
			}
		}
		return 1024 * 1024 * 8, nil
	}
	return 0, err
}

func interpretLineSeparator(params *tParameters) []byte {
	var lineSeparator []byte
	if params.system.Values[0] == "win" || params.system.Values[0] == "windows" {
		lineSeparator = make([]byte, 2)
		lineSeparator[0] = '\r'
		lineSeparator[1] = '\n'
	} else {
		lineSeparator = make([]byte, 1)
		lineSeparator[0] = '\n'
	}
	return lineSeparator
}

func generate(params *tParameters, size, buffer int, lineSeparator []byte) error {
	pathOut := params.output.Values[0]
	out, err := os.OpenFile(pathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer out.Close()
		var copiedTotal, copied int
		bytes := make([]byte, buffer)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		for copiedTotal < size && err == nil {
			if size - copiedTotal >= buffer {
				generateText(bytes, random, lineSeparator)
			} else {
				bytes = bytes[:size - copiedTotal]
				generateText(bytes, random, lineSeparator)
			}
			copied, err = out.Write(bytes)
			copiedTotal += copied
		}
	}
	return err
}

func generateGo(params *tParameters, size, buffer, threads int, lineSeparator []byte) error {
	pathOut := params.output.Values[0]
	out, err := os.OpenFile(pathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer out.Close()
		var copiedTotal, copied int
		bytes := make([]byte, buffer)
		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		for copiedTotal < size && err == nil {
			if size - copiedTotal >= buffer {
				generateText(bytes, random, lineSeparator)
			} else {
				bytes = bytes[:size - copiedTotal]
				generateText(bytes, random, lineSeparator)
			}
			copied, err = out.Write(bytes)
			copiedTotal += copied
		}
	}
	return err
}

func generateText(bytes []byte, random *rand.Rand, lineSeparator []byte) {
	var writtenTotal, words int
	limit, lengthSep := len(bytes), len(lineSeparator)
	for writtenTotal < limit {
		written := randWordLength(random)
		newLine := randNewLine(random)
		if words >= wordsMAX_PER_LINE {
			newLine = true
		}
		if newLine {
			words = 0
			writtenMax := limit - lengthSep - writtenTotal
			if written > writtenMax {
				written = writtenMax
			}
			randFill(bytes[writtenTotal:writtenTotal+written], random)
			writtenTotal += written
			for i, b := range lineSeparator {
				bytes[writtenTotal+i] = b
			}
			writtenTotal += lengthSep
		} else {
			words++
			writtenMax := limit - 1 - writtenTotal
			if written > writtenMax {
				written = writtenMax
			}
			randFill(bytes[writtenTotal:writtenTotal+written], random)
			writtenTotal += written
			bytes[writtenTotal] = ' '
			writtenTotal++
		}
	}
}

func randWordLength(random *rand.Rand) int {
	randomFloat := random.Float32()
	numberFloat := randomFloat * float32(wordMAX_LEN - wordMIN_LEN + 1)
	return int(numberFloat) + wordMIN_LEN
}

func randNewLine(random *rand.Rand) bool {
	randomFloat := random.Float32()
	if randomFloat > newLinePROBABILITY {
		return false
	}
	return true
}

func randFill(bytes []byte, random *rand.Rand) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32((90 - 65) * 2 + 2)
		letter := byte(numberFloat)
		if letter > 90 - 65 {
			bytes[i] = letter + 65 + 6
		} else {
			bytes[i] = letter + 65
		}
	}
}

func parseBytes(bytesStr string) (int, error) {
	if len(bytesStr) > 0 {
		lastByte := bytesStr[len(bytesStr)-1]
		if lastByte == 'k' || lastByte == 'K' || lastByte == 'm' || lastByte == 'M' || lastByte == 'g' || lastByte == 'G' {
			bytesStr = bytesStr[:len(bytesStr)-1]
		} else {
			lastByte = 0
		}
		bytes, err := strconv.Atoi(bytesStr)
		if err == nil {
			switch lastByte {
			case 'k':
				bytes = bytes * 1024
			case 'K':
				bytes = bytes * 1000
			case 'm':
				bytes = bytes * 1024 * 1024
			case 'M':
				bytes = bytes * 1000 * 1000
			case 'g':
				bytes = bytes * 1024 * 1024 * 1024
			case 'G':
				bytes = bytes * 1000 * 1000 * 1000
			}
		}
		return bytes, err
	}
	return 0, nil
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
	message += "  fsplit ( INFO | SIZE OUTPUT-FILE {OPTION} )\n\n"
	message += "INFO\n"
	message += "  -h, --help    print this help\n"
	message += "  -v, --version print version\n"
	message += "  --copyright   print copyright\n\n"
	message += "SIZE\n"
	message += "  -s=N[U]       size of file, U = unit (k/K, m/M or g/G)\n\n"
	message += "OPTION\n"
	message += "  -t=N          maximum number of threads (default 1)\n"
	message += "  -y=Y          operating system (e.g. -y=windows, for CRLF)\n"
	message += "  -b=N[U]       buffer size per thread, U = unit (k/K, m/M or g/G)"
	fmt.Println(message)
}

func printVersion() {
	fmt.Println("0.2.0")
}

func printExample() {
	message := "\nEXAMPLES\n"
	message += "   ... not available"
	fmt.Println(message)
}

func printCopyright() {
	message := "Copyright 2021, 2022, Vitali Baumtrok (vbsw@mailbox.org).\n"
	message += "Distributed under the Boost Software License, Version 1.0."
	fmt.Println(message)
}

func printError(err error) {
	fmt.Println("error:", err.Error())
}