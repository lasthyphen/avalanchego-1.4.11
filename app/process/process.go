// (c) 2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package process

import (
	"fmt"
	"path/filepath"

	"github.com/ava-labs/avalanchego/chains"
	"github.com/ava-labs/avalanchego/database/leveldb"
	"github.com/ava-labs/avalanchego/database/manager"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/database/rocksdb"
	"github.com/ava-labs/avalanchego/nat"
	"github.com/ava-labs/avalanchego/node"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/dynamicip"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/version"
)

const (
	Header = `     _____               .__                       .__
    /  _  \___  _______  |  | _____    ____   ____ |  |__   ____    ,_ o
   /  /_\  \  \/ /\__  \ |  | \__  \  /    \_/ ___\|  |  \_/ __ \   / //\,
  /    |    \   /  / __ \|  |__/ __ \|   |  \  \___|   Y  \  ___/    \>> |
  \____|__  /\_/  (____  /____(____  /___|  /\___  >___|  /\___  >    \\
          \/           \/          \/     \/     \/     \/     \/`

	mustUpgradeMsg = `
This version of AvalancheGo requires a database upgrade before running.

To do the database upgrade, restart this node with argument --fetch-only.

This will start the node in fetch only mode. It will bootstrap a new database
version and then stop. By default, this node will attempt to bootstrap from a
node running on the same machine (localhost) with staking port 9651. If no such
node exists, fetch only mode will be unable to complete.

The node in fetch only mode will by default not interfere with the node already
running. When the node in fetch only mode finishes, stop the other node running
on this computer and run without --fetch-only flag to run node normally. Fetch
only mode will not change this node's staking key/certificate.

Note that populating the new database version will approximately double the
amount of disk space required by AvalancheGo. Ensure that this computer has at
least enough disk space available.`

	upgradingMsg = `
Node running in fetch only mode.

Fetch only mode will not change this node's staking key/certificate.

Note that populating the new database version will approximately double the
amount of disk space required by AvalancheGo. Ensure that this computer has at
least enough disk space available.`

	alreadyUpgradedMsg = "fetch only mode done. Restart this node without --fetch-only to run normally"
)

var (
	stakingPortName = fmt.Sprintf("%s-staking", constants.AppName)
	httpPortName    = fmt.Sprintf("%s-http", constants.AppName)
)

// App is a wrapper around a node
type App struct {
	config node.Config
	node   *node.Node
	// log is set in Start()
	log logging.Logger
}

func NewApp(config node.Config) *App {
	return &App{
		config: config,
		node:   &node.Node{},
	}
}

// Start creates and runs an AvalancheGo node.
// Returns the node's exit code. If [a.Stop()] is called, Start()
// returns 0. This method blocks until the node is done.
func (a *App) Start() int {
	// we want to create the logger after the plugin has started the app
	logFactory := logging.NewFactory(a.config.LoggingConfig)
	defer logFactory.Close()

	var err error
	a.log, err = logFactory.Make("main")
	if err != nil {
		fmt.Printf("starting logger failed with: %s\n", err)
		return 1
	}

	// start the db manager
	var dbManager manager.Manager
	switch a.config.DBName {
	case rocksdb.Name:
		path := filepath.Join(a.config.DBPath, rocksdb.Name)
		dbManager, err = manager.NewRocksDB(path, a.log, version.CurrentDatabase, !a.config.FetchOnly)
	case leveldb.Name:
		dbManager, err = manager.NewLevelDB(a.config.DBPath, a.log, version.CurrentDatabase, !a.config.FetchOnly)
	case memdb.Name:
		dbManager = manager.NewMemDB(version.CurrentDatabase)
	default:
		err = fmt.Errorf(
			"db-type was %q but should have been one of {%s, %s, %s}",
			a.config.DBName,
			leveldb.Name,
			rocksdb.Name,
			memdb.Name,
		)
	}
	if err != nil {
		a.log.Fatal("couldn't create %q db manager at %s: %s", a.config.DBName, a.config.DBPath, err)
		return 1
	}

	// ensure migrations are done
	currentDBBootstrapped, err := dbManager.Current().Database.Has(chains.BootstrappedKey)
	if err != nil {
		a.log.Fatal("couldn't get whether database version %s ever bootstrapped: %s", version.CurrentDatabase, err)
		return 1
	}
	a.log.Info("bootstrapped with current database version: %v", currentDBBootstrapped)
	if a.config.FetchOnly {
		// Flag says to run in fetch only mode
		if currentDBBootstrapped {
			// We have already bootstrapped the current database
			a.log.Info(alreadyUpgradedMsg)
			return constants.ExitCodeDoneMigrating
		}
		a.log.Info(upgradingMsg)
	} else {
		prevDB, exists := dbManager.Previous()
		if !currentDBBootstrapped && exists && prevDB.Version.Compare(version.PrevDatabase) == 0 {
			// If we have the previous database version but not the current one then node
			// must run in fetch only mode (--fetch-only). The default behavior for a node in
			// fetch only mode is to bootstrap from a node on the same machine (127.0.0.1)
			// Tell the user to run in fetch only mode.
			a.log.Fatal(mustUpgradeMsg)
			return 1
		}
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered panic from", r)
		}
	}()

	defer func() {
		if err := dbManager.Close(); err != nil {
			a.log.Warn("failed to close the node's DB: %s", err)
		}
		a.log.StopOnPanic()
		a.log.Stop()
	}()

	// Track if sybil control is enforced
	if !a.config.EnableStaking {
		a.log.Warn("Staking is disabled. Sybil control is not enforced.")
	}

	// Check if transaction signatures should be checked
	if !a.config.EnableCrypto {
		a.log.Warn("transaction signatures are not being checked")
	}
	// TODO: disable crypto verification

	if err := a.config.ConsensusParams.Valid(); err != nil {
		a.log.Fatal("consensus parameters are invalid: %s", err)
		return 1
	}

	// Track if assertions should be executed
	if a.config.LoggingConfig.Assertions {
		a.log.Debug("assertions are enabled. This may slow down execution")
	}

	// SupportsNAT() for NoRouter is false.
	// Which means we tried to perform a NAT activity but we were not successful.
	if a.config.AttemptedNATTraversal && !a.config.Nat.SupportsNAT() {
		a.log.Warn("UPnP or NAT-PMP router attach failed, you may not be listening publicly. " +
			"Please confirm the settings in your router")
	}

	mapper := nat.NewPortMapper(a.log, a.config.Nat)
	defer mapper.UnmapAllPorts()

	// Open staking port we want for NAT Traversal to have the external port
	// (config.StakingIP.Port) to connect to our internal listening port
	// (config.InternalStakingPort) which should be the same in most cases.
	if a.config.StakingIP.IP().Port != 0 {
		mapper.Map(
			"TCP",
			a.config.StakingIP.IP().Port,
			a.config.StakingIP.IP().Port,
			stakingPortName,
			&a.config.StakingIP,
			a.config.DynamicUpdateDuration,
		)
	}

	// Open the HTTP port iff the HTTP server is not listening on localhost
	if a.config.HTTPHost != "127.0.0.1" && a.config.HTTPHost != "localhost" && a.config.HTTPPort != 0 {
		// For NAT Traversal we want to route from the external port
		// (config.ExternalHTTPPort) to our internal port (config.HTTPPort)
		mapper.Map(
			"TCP",
			a.config.HTTPPort,
			a.config.HTTPPort,
			httpPortName,
			nil,
			a.config.DynamicUpdateDuration,
		)
	}

	// Regularly updates our public IP (or does nothing, if configured that way)
	externalIPUpdater := dynamicip.NewDynamicIPManager(
		a.config.DynamicPublicIPResolver,
		a.config.DynamicUpdateDuration,
		a.log,
		&a.config.StakingIP,
	)
	defer externalIPUpdater.Stop()

	if err := a.node.Initialize(&a.config, dbManager, a.log, logFactory); err != nil {
		a.log.Fatal("error initializing node: %s", err)
		return 1
	}

	err = a.node.Dispatch()
	a.log.Debug("node dispatch returned: %s", err)
	return a.node.ExitCode()
}

// Assumes [a.node] is not nil.
// Blocks until [a.node] is done shutting down.
func (a *App) Stop() {
	a.node.Shutdown(0)
	a.node.DoneShuttingDown.Wait()
}
