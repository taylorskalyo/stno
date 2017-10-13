package formatter

import (
	"bytes"
	"io"

	toml "github.com/pelletier/go-toml"
)

// TOMLFormatter formats entries as a single toml document.
type TOMLFormatter struct {
	Entries []*toml.Tree
}

// Format combines a slice of TOML Trees and returns the resulting TOML
// document as an io.Reader.
func (tf TOMLFormatter) Format() (io.Reader, error) {
	buf := new(bytes.Buffer)
	for _, e := range tf.Entries {
		_, err := buf.WriteString(e.String())
		if err != nil {
			return buf, err
		}
	}
	return buf, nil
}
