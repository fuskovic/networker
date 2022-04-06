package table

import (
	"io"

	"cdr.dev/coder-cli/pkg/tablewriter"
)

// Writer describes an interface for writing tables.
type Writer[T any] interface {
	// Write writes a table to the underlying writer.
	// If successfull, it returns the number of rows written and a non-nil error.
	Write([]byte) (int, error)
}

type writer[T any] struct {
	io.Writer
	elements []T
}

// NewWriter initializes a new table writer for writing rows of elements into w.
func NewWriter[T any](w io.Writer, elements []T) Writer[T] {
	tw := new(writer[T])
	tw.Writer = w
	tw.elements = elements
	return tw
}

func (w *writer[T]) Write(b []byte) (int, error) {
	err := tablewriter.WriteTable(w.Writer, len(w.elements),
		func(i int) interface{} {
			return w.elements[i]
		},
	)
	if err != nil {
		return 0, err
	}
	return len(w.elements), nil
}
