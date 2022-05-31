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
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	wordLEN_MIN        = 2
	wordLEN_MAX        = 30
	newLinePROBABILITY = 0.1
	wordsPerLineMAX    = 20
)

type tParameters struct {
	help       *osargs.Result
	version    *osargs.Result
	example    *osargs.Result
	copyright  *osargs.Result
	alpha      *osargs.Result
	lower      *osargs.Result
	upper      *osargs.Result
	size       *osargs.Result
	threads    *osargs.Result
	system     *osargs.Result
	buffer     *osargs.Result
	output     *osargs.Result
	infoParams []*osargs.Result
	cmdParams  []*osargs.Result
}

type tGenerator struct {
	bytes      []byte
	random     *rand.Rand
	randomFill func(*rand.Rand, []byte)
}

type tThreads struct {
	chnl           chan *tGenerator
	counter        int
	maxThreadsUsed int
	seedAdd        int64
	randomFill     func(*rand.Rand, []byte)
}

func main() {
	params := new(tParameters)
	err := params.initFromOSArgs()
	if err == nil {
		if params.infoAvailable() {
			printInfo(params)
		} else {
			var sizeFile, maxThreads, sizeBuffer int
			newLine := interpretNewLine(params)
			sizeFile, err = interpretSize(params, err)
			maxThreads, err = interpretThreads(params, err)
			sizeBuffer, err = interpretBuffer(params, len(newLine)+1, err)
			if err == nil {
				if maxThreads == 1 {
					if params.outputToFile() {
						err = generateFile(params, sizeFile, sizeBuffer, newLine)
					} else {
						err = generateStd(params, sizeFile, sizeBuffer, newLine)
					}
				} else {
					if params.outputToFile() {
						err = generateFileGo(params, sizeFile, sizeBuffer, maxThreads, newLine)
					} else {
						err = generateStdGo(params, sizeFile, sizeBuffer, maxThreads, newLine)
					}
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
		params.alpha = args.Parse("-a", "--alpha", "-alpha", "alpha")
		params.lower = args.Parse("-l", "--lower", "-lower", "lower")
		params.upper = args.Parse("-u", "--upper", "-upper", "upper")
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
		outputLowerCase := strings.ToLower(params.output.Values[0])
		if outputLowerCase == "std" {
			params.output.Values[0] = outputLowerCase
		} else {
			outputPath, err := filepath.Abs(params.output.Values[0])
			if err == nil {
				params.output.Values[0] = outputPath
			} else {
				panic(err.Error())
			}
		}
	} else if len(unparsedArgs) > 0 {
		outputLowerCase := strings.ToLower(unparsedArgs[0])
		if outputLowerCase == "std" {
			params.output.Values = append(params.output.Values, outputLowerCase)
		} else {
			// just accept the first unparsed argument
			outputPath, err := filepath.Abs(unparsedArgs[0])
			if err == nil {
				params.output.Values = append(params.output.Values, outputPath)
			} else {
				panic(err.Error())
			}
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
	params.cmdParams = make([]*osargs.Result, 8)
	params.cmdParams[0] = params.size
	params.cmdParams[1] = params.threads
	params.cmdParams[2] = params.system
	params.cmdParams[3] = params.buffer
	params.cmdParams[4] = params.output
	params.cmdParams[5] = params.alpha
	params.cmdParams[6] = params.lower
	params.cmdParams[7] = params.upper
}

func (params *tParameters) infoAvailable() bool {
	return anyAvailable(params.infoParams)
}

func (params *tParameters) outputToFile() bool {
	return params.output.Values[0] != "std"
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

func interpretNewLine(params *tParameters) []byte {
	var newLine []byte
	if params.system.Values[0] == "win" || params.system.Values[0] == "windows" {
		newLine = make([]byte, 2)
		newLine[0] = '\r'
		newLine[1] = '\n'
	} else {
		newLine = make([]byte, 1)
		newLine[0] = '\n'
	}
	return newLine
}

func interpretSize(params *tParameters, err error) (int, error) {
	if err == nil {
		var sizeFile int
		sizeFile, err = parseBytes(params.size.Values[0])
		if err == nil || sizeFile <= 0 {
			return sizeFile, nil
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

func interpretBuffer(params *tParameters, sizeFile int, err error) (int, error) {
	if err == nil {
		if params.buffer.Available() {
			bytes, err := parseBytes(params.buffer.Values[0])
			if err == nil {
				if bytes > 0 {
					if bytes >= sizeFile {
						return bytes, nil
					}
					return sizeFile, nil
				}
			} else {
				return 0, errors.New("can't parse size of buffer")
			}
		}
		return 1024 * 1024 * 8, nil
	}
	return 0, err
}

func generateFile(params *tParameters, sizeFile, sizeBuffer int, newLine []byte) error {
	pathOut := params.output.Values[0]
	out, err := os.OpenFile(pathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer out.Close()
		timeStart := time.Now().UnixNano()
		randomFill := randomFillFunc(params)
		generator := newGenerator(sizeBuffer, 0, randomFill)
		for sizeTotal, sizeAdd := 0, 0; sizeTotal < sizeFile && err == nil; sizeTotal += sizeAdd {
			sizeAdd = generator.adjustBuffer(sizeFile - sizeTotal)
			generator.generateText(newLine)
			err = generator.writeFile(out)
		}
		timeEnd := time.Now().UnixNano()
		printThreadsUsed(1)
		printTime(timeEnd - timeStart)
	}
	return err
}

func generateStd(params *tParameters, sizeFile, sizeBuffer int, newLine []byte) error {
	randomFill := randomFillFunc(params)
	generator := newGenerator(sizeBuffer, 0, randomFill)
	for sizeTotal, sizeAdd := 0, 0; sizeTotal < sizeFile; sizeTotal += sizeAdd {
		sizeAdd = generator.adjustBuffer(sizeFile - sizeTotal)
		generator.generateText(newLine)
		generator.writeStd()
	}
	return nil
}

func generateFileGo(params *tParameters, sizeFile, sizeBuffer int, maxThreads int, newLine []byte) error {
	pathOut := params.output.Values[0]
	out, err := os.OpenFile(pathOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		defer out.Close()
		timeStart := time.Now().UnixNano()
		randomFill := randomFillFunc(params)
		threads := newThreads(maxThreads, randomFill)
		for sizeTotal, sizeAdd := 0, 0; sizeTotal < sizeFile && err == nil; sizeTotal += sizeAdd {
			generator, result := threads.nextGenerator(sizeBuffer)
			if result {
				sizeAdd = 0
				err = generator.writeFile(out)
			} else {
				sizeAdd = generator.adjustBuffer(sizeFile - sizeTotal)
				go threads.generateTextGo(generator, newLine)
			}
		}
		for threads.counter > 0 && err == nil {
			generator := threads.nextGeneratorResult(sizeBuffer)
			err = generator.writeFile(out)
		}
		timeEnd := time.Now().UnixNano()
		printThreadsUsed(threads.maxThreadsUsed)
		printTime(timeEnd - timeStart)
	}
	return err
}

func generateStdGo(params *tParameters, sizeFile, sizeBuffer int, maxThreads int, newLine []byte) error {
	randomFill := randomFillFunc(params)
	threads := newThreads(maxThreads, randomFill)
	for sizeTotal, sizeAdd := 0, 0; sizeTotal < sizeFile; sizeTotal += sizeAdd {
		generator, result := threads.nextGenerator(sizeBuffer)
		if result {
			sizeAdd = 0
			generator.writeStd()
		} else {
			sizeAdd = generator.adjustBuffer(sizeFile - sizeTotal)
			go threads.generateTextGo(generator, newLine)
		}
	}
	for threads.counter > 0 {
		generator := threads.nextGeneratorResult(sizeBuffer)
		generator.writeStd()
	}
	return nil
}

func newThreads(maxThreads int, randomFill func(*rand.Rand, []byte)) *tThreads {
	threads := new(tThreads)
	threads.chnl = make(chan *tGenerator, maxThreads)
	threads.randomFill = randomFill
	return threads
}

func (threads *tThreads) nextGenerator(sizeBuffer int) (*tGenerator, bool) {
	if threads.counter < cap(threads.chnl) {
		select {
		case generator := <-threads.chnl:
			threads.counter--
			return generator, true
		default:
			threads.counter++
			threads.seedAdd++
			if threads.maxThreadsUsed < threads.counter {
				threads.maxThreadsUsed = threads.counter
			}
		}
		return newGenerator(sizeBuffer, threads.seedAdd, threads.randomFill), false
	}
	return threads.nextGeneratorResult(sizeBuffer), true
}

func (threads *tThreads) nextGeneratorResult(sizeBuffer int) *tGenerator {
	generator := <-threads.chnl
	threads.counter--
	return generator
}

func (threads *tThreads) generateTextGo(generator *tGenerator, newLine []byte) {
	generator.generateText(newLine)
	threads.chnl <- generator
}

func newGenerator(sizeBuffer int, seedAdd int64, randomFill func(*rand.Rand, []byte)) *tGenerator {
	generator := new(tGenerator)
	generator.bytes = make([]byte, sizeBuffer)
	generator.random = rand.New(rand.NewSource(time.Now().UnixNano() + seedAdd))
	generator.randomFill = randomFill
	return generator
}

func (generator *tGenerator) adjustBuffer(sizeRemaining int) int {
	if sizeRemaining < len(generator.bytes) {
		generator.bytes = generator.bytes[:sizeRemaining]
	}
	return len(generator.bytes)
}

func (generator *tGenerator) writeFile(out *os.File) error {
	_, err := out.Write(generator.bytes)
	return err
}

func (generator *tGenerator) writeStd() {
	fmt.Printf("%s", generator.bytes)
}

func (generator *tGenerator) generateText(newLine []byte) {
	var writtenTotal, words int
	writtenLimit := generator.writeLimit(newLine)
	for writtenTotal < writtenLimit {
		lineBreak := generator.randLineBreak(words)
		if lineBreak {
			lengthWord := generator.randWordLength(len(generator.bytes) - writtenTotal - len(newLine))
			words = 0
			generator.randomFill(generator.random, generator.bytes[writtenTotal:writtenTotal+lengthWord])
			writtenTotal += lengthWord
			for i, b := range newLine {
				generator.bytes[writtenTotal+i] = b
			}
			writtenTotal += len(newLine)
		} else {
			lengthWord := generator.randWordLength(len(generator.bytes) - writtenTotal - len(newLine))
			words++
			generator.randomFill(generator.random, generator.bytes[writtenTotal:writtenTotal+lengthWord])
			writtenTotal += lengthWord
			generator.bytes[writtenTotal] = ' '
			writtenTotal++
		}
	}
	generator.randomFill(generator.random, generator.bytes[writtenTotal:])
}

func (generator *tGenerator) writeLimit(newLine []byte) int {
	limit := len(generator.bytes) - wordLEN_MAX - 1
	if len(newLine) > 0 {
		limit -= len(newLine) - 1
	}
	if limit > 0 {
		return limit
	}
	return 0
}

func (generator *tGenerator) randWordLength(lengthMax int) int {
	randomFloat := generator.random.Float32()
	numberFloat := randomFloat * float32(wordLEN_MAX-wordLEN_MIN+1)
	lengthWord := int(numberFloat) + wordLEN_MIN
	if lengthWord < lengthMax {
		return lengthWord
	}
	return lengthMax
}

func (generator *tGenerator) randLineBreak(words int) bool {
	if words < wordsPerLineMAX {
		randomFloat := generator.random.Float32()
		if randomFloat > newLinePROBABILITY {
			return false
		}
	}
	return true
}

func randomFillFunc(params *tParameters) func(*rand.Rand, []byte) {
	if params.alpha.Available() {
		if params.lower.Available() {
			return randomFillL
		} else if params.upper.Available() {
			return randomFillU
		} else {
			return randomFillA
		}
	} else if params.lower.Available() {
		if params.upper.Available() {
			return randomFillA
		} else {
			return randomFillL
		}
	} else if params.upper.Available() {
		if params.lower.Available() {
			return randomFillA
		} else {
			return randomFillU
		}
	}
	return randomFillZ
}

func randomFillA(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32((90-65)*2+2)
		letter := byte(numberFloat)
		if letter > 90-65 {
			bytes[i] = letter + 65 + 6
		} else {
			bytes[i] = letter + 65
		}
	}
}

func randomFillL(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(122-97+1)
		letter := byte(numberFloat)
		bytes[i] = letter + 97
	}
}

func randomFillU(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(90-65+1)
		letter := byte(numberFloat)
		bytes[i] = letter + 65
	}
}

func randomFillZ(random *rand.Rand, bytes []byte) {
	for i := range bytes {
		randomFloat := random.Float32()
		numberFloat := randomFloat * float32(126-33+1-5)
		letter := byte(numberFloat)
		if letter > 95-33+1-5 {
			bytes[i] = letter + 33 + 5
		} else if letter > 93-33+1-4 {
			bytes[i] = letter + 33 + 4
		} else if letter > 91-33+1-3 {
			bytes[i] = letter + 33 + 3
		} else if letter > 38-33+1-2 {
			bytes[i] = letter + 33 + 2
		} else if letter > 33-33+1-1 {
			bytes[i] = letter + 33 + 1
		} else {
			bytes[i] = letter + 33
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
	fmt.Println("Run 'textgen --help' for usage.")
}

func printHelp() {
	message := "\nUSAGE\n"
	message += "  textgen ( INFO | SIZE OUTPUT {OPTION} )\n\n"
	message += "INFO\n"
	message += "  -h, --help       print this help\n"
	message += "  -v, --version    print version\n"
	message += "  -c, --copyright  print copyright\n"
	message += "SIZE\n"
	message += "  -s=N[U]          size of file, U = unit (k/K, m/M or g/G)\n"
	message += "OUTPUT\n"
	message += "  <path>           write output to file <path>\n"
	message += "  std              write output to standard output (e.g. console)\n"
	message += "OPTION\n"
	message += "  -t=N             maximum number of threads (default 1)\n"
	message += "  -y=Y             operating system (e.g. -y=windows, for CRLF)\n"
	message += "  -b=N[U]          buffer size per thread, U = unit (k/K, m/M or g/G)\n"
	message += "  -a               output letters, only\n"
	message += "  -l               output lower case letters, only\n"
	message += "  -u               output upper case letters, only"
	fmt.Println(message)
}

func printVersion() {
	fmt.Println("0.3.0")
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

func printThreadsUsed(maxThreadsUsed int) {
	fmt.Println("threads:", maxThreadsUsed)
}

func printTime(nanos int64) {
	micro := nanos / 1000
	millis := micro / 1000
	seconds := float64(millis) / 1000.0
	fmt.Println("seconds:", seconds)
}

func printError(err error) {
	fmt.Println("error:", err.Error())
}
