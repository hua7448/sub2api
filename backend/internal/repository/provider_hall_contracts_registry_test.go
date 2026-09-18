//go:build integration || providerhall_localdb

package repository

import (
	"database/sql"
	"sort"
	"testing"
)

// providerHallLocalDBContract is one batch's database contract suite. Each
// batch registers its own from an init() in its *_db_test.go file so the
// shared cluster harness never needs editing.
type providerHallLocalDBContract struct {
	name string
	fn   func(t *testing.T, db *sql.DB)
}

var providerHallLocalDBContracts []providerHallLocalDBContract

func registerProviderHallLocalDBContract(name string, fn func(t *testing.T, db *sql.DB)) {
	providerHallLocalDBContracts = append(providerHallLocalDBContracts, providerHallLocalDBContract{name: name, fn: fn})
	sort.SliceStable(providerHallLocalDBContracts, func(i, j int) bool {
		return providerHallLocalDBContracts[i].name < providerHallLocalDBContracts[j].name
	})
}
