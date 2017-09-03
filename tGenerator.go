/*
 *          Copyright 2017, Vitali Baumtrok.
 * Distributed under the Boost Software License, Version 1.0.
 *     (See accompanying file LICENSE or copy at
 *        http://www.boost.org/LICENSE_1_0.txt)
 */

package main

import (
	"io"
	"time"
	"math/rand"
)

const (
	const_BUFFER_SIZE   = 1024*1024
	const_MAX_WORD_SIZE = 20
)

type tGenerator struct {
	clArgsInterpreter *tCLArgsInterpreter
	writer io.Writer
}

func newGenerator() *tGenerator {
	generator := new(tGenerator)
	return generator
}

func (this *tGenerator) setInterpreter(clArgsInterpreter *tCLArgsInterpreter) {
	this.clArgsInterpreter = clArgsInterpreter
}

func (this *tGenerator) setWriter(writer io.Writer) {
	this.writer = writer
}

func (this *tGenerator) generate() {
	bytesWritten := 0
	bytes := make([]byte,const_BUFFER_SIZE)
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	for bytesWritten + const_BUFFER_SIZE <= this.clArgsInterpreter.clArgs.size {
		bytesWritten += const_BUFFER_SIZE
		fillBytes(bytes,random)
		this.writer.Write(bytes)
	}
	if bytesWritten < this.clArgsInterpreter.clArgs.size {
		bytesLeftToWrite := this.clArgsInterpreter.clArgs.size - bytesWritten
		bytes = bytes[:bytesLeftToWrite]
		fillBytes(bytes,random)
		this.writer.Write(bytes)
	}
}

func fillBytes(bytes []byte, random *rand.Rand) {
	wordLength := 0

	for i := range bytes {
		if i > 0 {
			if wordLength > 0 {
				bytes[i] = 'a'
				wordLength -= 1
			} else {
				bytes[i] = newRandomSeparator(random)
				wordLength = newRandomNumber(random,2,const_MAX_WORD_SIZE)
			}
		} else {
			bytes[i] = 'a'
			wordLength = newRandomNumber(random,2,const_MAX_WORD_SIZE) - 1
		}
	}
}

func newRandomNumber(random *rand.Rand, from, to int) int {
	randomFloat := random.Float32()
	number := int(randomFloat * float32(to - from)) + from
	return number
}

func newRandomSeparator(random *rand.Rand) byte {
	randomFloat := random.Float32()

	if randomFloat > 0.1 {
		return ' '
	} else {
		return '\n'
	}
}
