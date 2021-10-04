// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
package keystore

import (
	"github.com/ava-labs/avalanchego/chains/atomic"
	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/manager"
	"github.com/ava-labs/avalanchego/database/prefixdb"
	"github.com/ava-labs/avalanchego/version"
)

// Perform a database migration if required
func (ks *keystore) migrate(dbManager manager.Manager) error {
	prevDB, prevDBExists := dbManager.Previous()
	if !prevDBExists {
		// There is nothing to migrate
		return nil
	}
	prevDBVersion := prevDB.Version
	currentDB := dbManager.Current()
	currentDBVersion := currentDB.Version
	// Right now the only valid migration is from database version 1.0.0 to 1.4.5
	if prevDBVersion.Compare(version.DatabaseVersion1_0_0) == 0 &&
		currentDBVersion.Compare(version.DatabaseVersion1_4_5) == 0 {
		return ks.migrate1_0_0To1_4_5(prevDB, currentDB)
	}
	return nil
}

func (ks *keystore) migrate1_0_0To1_4_5(prevDB, currentDB *manager.VersionedDatabase) error {
	migrated, err := currentDB.Database.Has(migratedKey)
	if err != nil {
		return err
	} else if migrated {
		// Skip migration if previously done.
		return nil
	}
	ks.log.Info("migrating keystore from database version %s to %s", prevDB.Version, currentDB.Version)

	previousUserDB := prefixdb.New(usersPrefix, prevDB.Database)
	previousBCDB := prefixdb.New(bcsPrefix, prevDB.Database)
	userIterator := previousUserDB.NewIterator()
	defer userIterator.Release()
	for userIterator.Next() {
		username := userIterator.Key()
		exists, err := ks.userDB.Has(username)
		if err != nil {
			return err
		} else if exists { // already have this username in the keystore
			continue
		}

		userBatch := ks.userDB.NewBatch()
		if err := userBatch.Put(username, userIterator.Value()); err != nil {
			return err
		}

		currentUserBCDB := prefixdb.New(username, ks.bcDB)
		previousUserBCDB := prefixdb.New(username, previousBCDB)
		bcsBatch := currentUserBCDB.NewBatch()
		if err := migrateUserBCDB(previousUserBCDB, bcsBatch, userBatch); err != nil {
			return err
		}
	}

	if err := userIterator.Error(); err != nil {
		return err
	}
	if err := currentDB.Database.Put(migratedKey, []byte(prevDB.Version.String())); err != nil {
		return err
	}
	ks.log.Info("finished migrating keystore from database version %s to %s", prevDB.Version, currentDB.Version)
	return nil
}

func migrateUserBCDB(previousUserBCDB database.Database, bcsBatch database.Batch, userBatch database.Batch) error {
	iterator := previousUserBCDB.NewIterator()
	defer iterator.Release()

	for iterator.Next() {
		if err := bcsBatch.Put(iterator.Key(), iterator.Value()); err != nil {
			return err
		}
	}

	if err := iterator.Error(); err != nil {
		return err
	}
	return atomic.WriteAll(userBatch, bcsBatch)
}
