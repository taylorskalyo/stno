package journal

import "github.com/taylorskalyo/hj/datastore"

// Journal represents a journal with entries.
type Journal struct {
	Entries   []Entry
	DataStore *datastore.DataStore
}

// Entry represents a journal entry.
type Entry struct {
	uuid string
}
