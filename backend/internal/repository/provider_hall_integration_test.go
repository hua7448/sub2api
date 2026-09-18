//go:build integration

package repository

import "testing"

func TestProviderHallDatabaseContracts(t *testing.T) {
	providerHallDatabaseContracts(t, integrationDB)
	providerHallTargetDatabaseContracts(t, integrationDB)
}
