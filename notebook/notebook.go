package notebook

import (
	"github.com/taylorskalyo/stno/datastore"
)

// Notebook wraps the underlying data store of a stno notebook and provides
// methods for manipulating the contents.
type Notebook struct {
	datastore.DataStore
}
