# Go Log Piper

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/golp) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/golp/master/LICENSE) [![Build Status](https://travis-ci.org/rs/golp.svg?branch=master)](https://travis-ci.org/rs/golp) [![Coverage](http://gocover.io/_badge/github.com/rs/golp)](http://gocover.io/github.com/rs/golp)

Go programs sometime generate output you can't easily control like panics and `net/http` recovered panics. By default, those output contains multiple lines with stack traces. This does not play well with most logging systems that will generate one log event per outputed line.

The `golp` is a simple program that reads those kinds of log on its standard input, and merge all lines of a given panic or standard multi-lines Go log message into a single quotted line.

## Usage

Options:

    --json
            Wrap messages to JSON one object per line.
    -json-key string
            The key name to use for the message in JSON mode. (default "message")
    --prefix string
            Go logger prefix set in the application if any.
    --strip
            Strip log line timestamps on output.

Send panics and other program panics to syslog:

    mygoprogram 2>&1 | golp | logger -t mygoprogram -p local7.err

    > Jan  8 16:59:26 host panic: panic: test\n\ngoroutine 1 [running]:\npanic(0x…

Send panics as JSON:

    mygoprogram 2>&1 | golp --json | logger -t mygoprogram -p local7.err

    > Jan  8 16:59:26 host {"message": "panic: panic: test\n\ngoroutine 1 [running]:\npanic(0x…

Add some fields to the JSON output (using [jq](https://stedolan.github.io/jq/)):

    mygoprogram 2>&1 | golp --json | jq -c '. + {"level": "error", "program": "mygoprogram"}'

    > {"level": "error", "program": "mygoprogram", "message": "panic: panic: test\n\ngoroutine 1 [running]:\npanic(0x…

## License

All source code is licensed under the [MIT License](https://raw.github.com/rs/golp/master/LICENSE).