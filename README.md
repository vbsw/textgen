# textgen

[![Go Reference](https://pkg.go.dev/badge/github.com/vbsw/textgen.svg)](https://pkg.go.dev/github.com/vbsw/textgen) [![Go Report Card](https://goreportcard.com/badge/github.com/vbsw/textgen)](https://goreportcard.com/report/github.com/vbsw/textgen) [![Stability: Experimental](https://masterminds.github.io/stability/experimental.svg)](https://masterminds.github.io/stability/experimental.html)

## About
textgen creates a file with random text. textgen is published on <https://github.com/vbsw/textgen> and <https://gitlab.com/vbsw/textgen>.

## Copyright
Copyright 2021, 2022, Vitali Baumtrok (vbsw@mailbox.org).

textgen is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

textgen is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Usage

	textgen ( INFO | SIZE OUTPUT {OPTION} )

	INFO
		-h, --help       print this help
		-v, --version    print version
		-c, --copyright  print copyright
	SIZE
		-s=N[U]          size of file, U = unit (k/K, m/M or g/G)
	OUTPUT
		<path>           write output to file <path>
		std              write output to standard output (e.g. console)
	OPTION
		-t=N             maximum number of threads (default 1)
		-y=Y             operating system (e.g. -y=windows, for CRLF)
		-b=N[U]          buffer size per thread, U = unit (k/K, m/M or g/G)
		-a               output letters, only
		-l               output lower case letters, only
		-u               output upper case letters, only

## Example
Create a new file, named test.txt, in working directory, with 100 kilobytes of random text.

	$ textgen 100K test.txt

## References
- https://golang.org/doc/install
- https://git-scm.com/book/en/v2/Getting-Started-Installing-Git
