# goIJVM
[![Go Report Card](https://goreportcard.com/badge/github.com/TheTipo01/goIJVM)](https://goreportcard.com/report/github.com/TheTipo01/goIJVM)
[![Build Status](https://travis-ci.com/TheTipo01/goIJVM.svg?branch=master)](https://travis-ci.com/TheTipo01/goIJVM)

An [IJVM](https://en.wikipedia.org/wiki/IJVM) interpreter, written in go.

Implements most of the instructions, excepts INVOKEVIRTUAL, IRETURN, LDC_W and the one for manipulating arrays.

Also added an istruction called DEBUG for printing out status of the stack and the variables.

## Usage
Grab a [release](https://github.com/TheTipo01/goIJVM/releases) from the releases tab for your computer, and drag a file 
containing the instruction for the program you want to run on the executable. 

Or you can always start the program and give the path of the program to run as the first argument.
