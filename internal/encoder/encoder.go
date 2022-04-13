package encoder

import (
	"encoding/json"
	"io"

	"cdr.dev/coder-cli/pkg/tablewriter"
	"gopkg.in/yaml.v3"
)

// Encoder can encode generic objects into various output formats.
type Encoder[T any] struct {
	w      io.Writer
	output string
}

// New initializes a new Encoder.
func New[T any](w io.Writer, output string) Encoder[T] {
	return Encoder[T]{w, output}
}

// Encode encodes the generic objects into the output specified during Encoder initialization.
func (e *Encoder[T]) Encode(objects ...T) error {
	var err error
	switch e.output {
	case "json":
		jsonEncoder := json.NewEncoder(e.w)
		jsonEncoder.SetIndent("", "\t")
		jsonEncoder.SetEscapeHTML(true)
		// Don't output as json array if there is only one object
		if len(objects) == 1 {
			err = jsonEncoder.Encode(objects[0])
		} else {
			err = jsonEncoder.Encode(objects)
		}
	case "yaml":
		err = yaml.NewEncoder(e.w).Encode(objects)
	default:
		err = tablewriter.WriteTable(e.w, len(objects),
			func(i int) any {
				return objects[i]
			},
		)
	}
	return err
}
