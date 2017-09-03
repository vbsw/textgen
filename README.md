# Text Generator

## About
Text Generator is a program to generate random text. It is written in Go and published on <https://github.com/vbsw/textgen>.

## Copyright
Copyright 2017, Vitali Baumtrok (vbsw@mailbox.org).

Text Generator is distributed under the Boost Software License, version 1.0. (See accompanying file LICENSE or copy at http://www.boost.org/LICENSE_1_0.txt)

Text Generator is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the Boost Software License for more details.

## Example
Input:

	$ textgen 10 std

Output:

	asdf d ewq

## Using Go
Get this project:

	$ go get github.com/vbsw/textgen

Update a local copy:

	$ go get -u github.com/vbsw/textgen

Compile:

	$ go install github.com/vbsw/textgen

Run tests:

	$ go test github.com/vbsw/textgen

## Using Git
Get the master branch and all refs of this project:

	$ git clone https://github.com/vbsw/textgen.git

See all tags:

	$ git tag -l

See local and remote branches:

	$ git branch -a

Checkout other branches than master, for example the development branch:

	$ git branch development origin/development
	$ git checkout development

See tracked remote branches:

	$ git branch -vv

Update all tracked branches and all refs:

	$ git fetch

## References
- <https://golang.org/doc/install>
- <https://git-scm.com/book/en/v2/Getting-Started-Installing-Git>
