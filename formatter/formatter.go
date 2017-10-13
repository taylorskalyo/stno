package formatter

import "io"

// Formatter prepares a list of entries for output.
type Formatter interface {
	Format() (io.Reader, error)
}
