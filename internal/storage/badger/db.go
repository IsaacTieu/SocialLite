package badgerdb

import "github.com/dgraph-io/badger/v4"

func Open(path string) (*badger.DB, error) {
	opts := badger.DefaultOptions(path).WithLogger(nil)
	return badger.Open(opts)
}