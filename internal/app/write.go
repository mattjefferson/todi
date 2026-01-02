package app

import (
	"fmt"
	"io"
)

func writeLine(w io.Writer, args ...any) {
	if _, err := fmt.Fprintln(w, args...); err != nil {
		return
	}
}
