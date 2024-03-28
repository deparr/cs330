package main

import "fmt"

const input = `"hello world"`

func main() {
	var p parser
	p.Init()
	for _, b := range []byte(input) {
		switch p.Write(b) {
		case BadInput:
			fmt.Println("bad input")
		case Success:
			fmt.Println("done")
			return
		}
	}
	fmt.Println("ran out of input")
}

func parseQuoted(read func() byte) bool {
	if read() != '"' {
		return false
	}
	var c byte
	for c != '"' {
		c = read()
		if c == '\\' {
			read()
		}
	}
	return true
}

type Status int

const (
	NeedMoreInput Status = iota
	BadInput
	Success
)

type parser struct {
	resume func(byte) Status
}

func (p *parser) Init() {
	coparse := func(_ byte, yield func(Status) byte) Status {
		read := func() byte { return yield(NeedMoreInput) }
		if !parseQuoted(read) {
			return BadInput
		}
		return Success
	}
	p.resume = coro_New(coparse)
	p.resume(0)
}

func (p *parser) Write(c byte) Status {
	return p.resume(c)
}

func coro_New[In, Out any](f func(In, func(Out) In) Out) (resume func(In) Out) {
	cin := make(chan In)
	cout := make(chan Out)
	resume = func(in In) Out {
		cin <- in
		return <-cout
	}
	yield := func(out Out) In {
		cout <- out
		return <-cin
	}
	go func() {
		cout <- f(<-cin, yield)
	}()
	return resume
}

