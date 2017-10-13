package query

// Query stores query expressions and a notebook data store to run them
// against.
type Query struct {
	WhereExpr string
	OrderExpr string
	DS        datastore.Datastore
}

// Result holds an entry's contents as a toml Tree and the entry's UID.
type Result struct {
	UID  string
	Tree *toml.Tree
}

// Run performs a query and returns the results.
func (q Query) Run() ([]Result, error) {
	uids, err := ds.ListEntries()
	if err != nil {
		return out, err
	}
	abort := make(chan struct{})
	defer close(abort)

	entriesc := openEntries(abort, uids)
}

func (q Query) openEntries(abort <-chan struct{}, uids []string) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for _, uid := range uids {
			rc, err := a.DS.NewEntryReadCloser(uid)
			if err != nil {
				continue
			}
			tree, err := toml.LoadReader(rc)
			if err != nil {
				continue
			}
			select {
			case out <- Result{UID: uid, Tree: tree}:
			case <-abort:
				return
			}
		}
	}()
	return out
}

// Where filters entries based on an expression.
func Where(abort <-chan struct{}, in <-chan Result) <-chan Result {
	out := make(chan Result)
	go func() {
		close(out)
		for _, e := range in {
			if true {
				out <- e
			}
		}
	}()
	return out
}
