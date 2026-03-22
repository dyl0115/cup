package proxy

import (
	"fmt"
	"io"
)

type prefixWriter struct {
	prefix  string
	newline bool
}

func newPrefixWriter(prefix string) io.Writer {
	return &prefixWriter{prefix: prefix, newline: true}
}

func (pw *prefixWriter) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if pw.newline {
			fmt.Print(pw.prefix)
			pw.newline = false
		}
		fmt.Printf("%c", b)
		if b == '\n' {
			pw.newline = true
		}
	}
	return len(p), nil
}
