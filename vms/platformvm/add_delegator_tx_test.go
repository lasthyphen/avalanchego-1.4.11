// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package platformvm

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto"
)

func TestAddDelegatorTxSyntacticVerify(t *testing.T) {
	vm, _ := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
		vm.ctx.Lock.Unlock()
	}()

	nodeID := keys[0].PublicKey().Address()
	rewardAddress := nodeID

	// Case : tx is nil
	var unsignedTx *UnsignedAddDelegatorTx
	if err := unsignedTx.Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err == nil {
		t.Fatal("should have errored because tx is nil")
	}

	// Case: Wrong network ID
	tx, err := vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(defaultValidateStartTime.Unix()),
		uint64(defaultValidateEndTime.Unix()),
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).NetworkID++
	// This tx was syntactically verified when it was created...pretend it wasn't so we don't use cache
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).syntacticallyVerified = false
	if err := tx.UnsignedTx.(*UnsignedAddDelegatorTx).Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err == nil {
		t.Fatal("should have errored because the wrong network ID was used")
	}

	// Case: Not enough weight
	tx, err = vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(defaultValidateStartTime.Unix()),
		uint64(defaultValidateEndTime.Unix()),
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).Validator.Wght = vm.MinDelegatorStake - 1
	// This tx was syntactically verified when it was created...pretend it wasn't so we don't use cache
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).syntacticallyVerified = false
	if err := tx.UnsignedTx.(*UnsignedAddDelegatorTx).Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err == nil {
		t.Fatal("should have errored because of not enough weight")
	}

	// Case: Validation length is too short
	tx, err = vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(defaultValidateStartTime.Unix()),
		uint64(defaultValidateStartTime.Add(defaultMinStakingDuration).Unix()),
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).Validator.End-- // 1 shorter than minimum stake time
	// This tx was syntactically verified when it was created...pretend it wasn't so we don't use cache
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).syntacticallyVerified = false
	if err = tx.UnsignedTx.(*UnsignedAddDelegatorTx).Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err == nil {
		t.Fatal("should have errored because validation length too short")
	}

	// Case: Validation length is too long
	if tx, err = vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(defaultValidateStartTime.Unix()),
		uint64(defaultValidateStartTime.Add(defaultMaxStakingDuration).Unix()),
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	); err != nil {
		t.Fatal(err)
	}
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).Validator.End++ // 1 longer than maximum stake time
	// This tx was syntactically verified when it was created...pretend it wasn't so we don't use cache
	tx.UnsignedTx.(*UnsignedAddDelegatorTx).syntacticallyVerified = false
	if err := tx.UnsignedTx.(*UnsignedAddDelegatorTx).Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err == nil {
		t.Fatal("should have errored because validation length too long")
	}

	// Case: Valid
	if tx, err = vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(defaultValidateStartTime.Unix()),
		uint64(defaultValidateEndTime.Unix()),
		nodeID,
		rewardAddress,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	); err != nil {
		t.Fatal(err)
	} else if err := tx.UnsignedTx.(*UnsignedAddDelegatorTx).Verify(
		vm.ctx,
		vm.codec,
		vm.MinDelegatorStake,
		defaultMinStakingDuration,
		defaultMaxStakingDuration,
	); err != nil {
		t.Fatal(err)
	}
}

func TestAddDelegatorTxSemanticVerify(t *testing.T) {
	nodeID := keys[0].PublicKey().Address()
	rewardAddress := nodeID

	factory := crypto.FactorySECP256K1R{}
	keyIntf, err := factory.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	newValidatorKey := keyIntf.(*crypto.PrivateKeySECP256K1R)
	newValidatorID := newValidatorKey.PublicKey().Address()
	newValidatorStartTime := uint64(defaultValidateStartTime.Add(5 * time.Second).Unix())
	newValidatorEndTime := uint64(defaultValidateEndTime.Add(-5 * time.Second).Unix())

	// [addValidator] adds a new validator to the primary network's pending validator set
	addValidator := func(vm *VM) {
		tx, err := vm.newAddValidatorTx(
			vm.MinValidatorStake,                    // stake amount
			newValidatorStartTime,                   // start time
			newValidatorEndTime,                     // end time
			newValidatorID,                          // node ID
			rewardAddress,                           // Reward Address
			PercentDenominator,                      // subnet
			[]*crypto.PrivateKeySECP256K1R{keys[0]}, // key
			ids.ShortEmpty,                          // change addr
		)
		if err != nil {
			t.Fatal(err)
		}

		vm.internalState.AddCurrentStaker(tx, 0)
		vm.internalState.AddTx(tx, Committed)
		if err := vm.internalState.Commit(); err != nil {
			t.Fatal(err)
		}
		if err := vm.internalState.(*internalStateImpl).loadCurrentValidators(); err != nil {
			t.Fatal(err)
		}
	}

	freshVM, _ := defaultVM()
	currentTimestamp := freshVM.internalState.GetTimestamp()

	type test struct {
		stakeAmount   uint64
		startTime     uint64
		endTime       uint64
		nodeID        ids.ShortID
		rewardAddress ids.ShortID
		feeKeys       []*crypto.PrivateKeySECP256K1R
		setup         func(vm *VM)
		shouldErr     bool
		description   string
	}

	tests := []test{
		{
			freshVM.MinDelegatorStake,
			uint64(defaultValidateStartTime.Unix()),
			uint64(defaultValidateEndTime.Unix()) + 1,
			nodeID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			nil,
			true,
			"validator stops validating primary network earlier than subnet",
		},
		{
			freshVM.MinDelegatorStake,
			uint64(currentTimestamp.Add(maxFutureStartTime + time.Second).Unix()),
			uint64(currentTimestamp.Add(maxFutureStartTime * 2).Unix()),
			nodeID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			nil,
			true,
			fmt.Sprintf("validator should not be added more than (%s) in the future", maxFutureStartTime),
		},
		{
			freshVM.MinDelegatorStake,
			uint64(defaultValidateStartTime.Unix()),
			uint64(defaultValidateEndTime.Unix()) + 1,
			nodeID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			nil,
			true,
			"end time is after the primary network end time",
		},
		{
			freshVM.MinDelegatorStake,
			uint64(defaultValidateStartTime.Add(5 * time.Second).Unix()),
			uint64(defaultValidateEndTime.Add(-5 * time.Second).Unix()),
			newValidatorID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			nil,
			true,
			"validator not in the current or pending validator sets of the subnet",
		},
		{
			freshVM.MinDelegatorStake,
			newValidatorStartTime - 1, // start validating subnet before primary network
			newValidatorEndTime,
			newValidatorID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			addValidator,
			true,
			"validator starts validating subnet before primary network",
		},
		{
			freshVM.MinDelegatorStake,
			newValidatorStartTime,
			newValidatorEndTime + 1, // stop validating subnet after stopping validating primary network
			newValidatorID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			addValidator,
			true,
			"validator stops validating primary network before subnet",
		},
		{
			freshVM.MinDelegatorStake,
			newValidatorStartTime, // same start time as for primary network
			newValidatorEndTime,   // same end time as for primary network
			newValidatorID,
			rewardAddress,
			[]*crypto.PrivateKeySECP256K1R{keys[0]},
			addValidator,
			false,
			"valid",
		},
		{
			freshVM.MinDelegatorStake, // weight
			uint64(currentTimestamp.Unix()),
			uint64(defaultValidateEndTime.Unix()),
			nodeID,                                  // node ID
			rewardAddress,                           // Reward Address
			[]*crypto.PrivateKeySECP256K1R{keys[0]}, // tx fee payer
			nil,
			true,
			"starts validating at current timestamp",
		},
		{
			freshVM.MinDelegatorStake,               // weight
			uint64(defaultValidateStartTime.Unix()), // start time
			uint64(defaultValidateEndTime.Unix()),   // end time
			nodeID,                                  // node ID
			rewardAddress,                           // Reward Address
			[]*crypto.PrivateKeySECP256K1R{keys[1]}, // tx fee payer
			func(vm *VM) { // Remove all UTXOs owned by keys[1]
				utxoIDs, err := vm.internalState.UTXOIDs(keys[1].PublicKey().Address().Bytes(), ids.Empty, math.MaxInt32)
				if err != nil {
					t.Fatal(err)
				}
				for _, utxoID := range utxoIDs {
					vm.internalState.DeleteUTXO(utxoID)
				}
				if err := vm.internalState.Commit(); err != nil {
					t.Fatal(err)
				}
			},
			true,
			"tx fee paying key has no funds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			vm, _ := defaultVM()
			vm.ctx.Lock.Lock()
			defer func() {
				if err := vm.Shutdown(); err != nil {
					t.Fatal(err)
				}
				vm.ctx.Lock.Unlock()
			}()

			tx, err := vm.newAddDelegatorTx(
				tt.stakeAmount,
				tt.startTime,
				tt.endTime,
				tt.nodeID,
				tt.rewardAddress,
				tt.feeKeys,
				ids.ShortEmpty, // change addr
			)
			if err != nil {
				t.Fatalf("couldn't build tx: %s", err)
			}
			if tt.setup != nil {
				tt.setup(vm)
			}
			if _, _, _, _, err := tx.UnsignedTx.(UnsignedProposalTx).SemanticVerify(vm, vm.internalState, tx); err != nil && !tt.shouldErr {
				t.Fatalf("shouldn't have errored but got %s", err)
			} else if err == nil && tt.shouldErr {
				t.Fatalf("expected test to error but got none")
			}
		})
	}
}

func TestAddDelegatorTxOverDelegatedRegression(t *testing.T) {
	vm, _ := defaultVM()
	vm.ctx.Lock.Lock()
	defer func() {
		if err := vm.Shutdown(); err != nil {
			t.Fatal(err)
		}
		vm.ctx.Lock.Unlock()
	}()

	validatorStartTime := defaultGenesisTime.Add(syncBound).Add(1 * time.Second)
	validatorEndTime := validatorStartTime.Add(360 * 24 * time.Hour)
	key, err := vm.factory.NewPrivateKey()
	if err != nil {
		t.Fatal(err)
	}
	id := key.PublicKey().Address()

	// create valid tx
	addValidatorTx, err := vm.newAddValidatorTx(
		vm.MinValidatorStake,
		uint64(validatorStartTime.Unix()),
		uint64(validatorEndTime.Unix()),
		id,
		id,
		PercentDenominator,
		[]*crypto.PrivateKeySECP256K1R{keys[0]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}

	// trigger block creation
	if err := vm.mempool.IssueTx(addValidatorTx); err != nil {
		t.Fatal(err)
	}
	addValidatorBlockIntf, err := vm.BuildBlock()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the proposed block
	if err := addValidatorBlockIntf.Verify(); err != nil {
		t.Fatal(err)
	}

	// Assert preferences are correct
	addValidatorBlock := addValidatorBlockIntf.(*ProposalBlock)
	options, err := addValidatorBlock.Options()
	if err != nil {
		t.Fatal(err)
	}

	// verify the commit block
	commit := options[0].(*CommitBlock)
	if err := commit.Verify(); err != nil {
		t.Fatal(err)
	}

	// Accept the proposal block and the commit block
	if err := addValidatorBlock.Accept(); err != nil {
		t.Fatal(err)
	}
	if err := commit.Accept(); err != nil {
		t.Fatal(err)
	}

	vm.clock.Set(validatorStartTime)

	firstAdvanceTimeBlockIntf, err := vm.BuildBlock()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the proposed block
	if err := firstAdvanceTimeBlockIntf.Verify(); err != nil {
		t.Fatal(err)
	}

	// Assert preferences are correct
	firstAdvanceTimeBlock := firstAdvanceTimeBlockIntf.(*ProposalBlock)
	options, err = firstAdvanceTimeBlock.Options()
	if err != nil {
		t.Fatal(err)
	}

	// verify the commit block
	commit = options[0].(*CommitBlock)
	if err := commit.Verify(); err != nil {
		t.Fatal(err)
	}

	// Accept the proposal block and the commit block
	if err := firstAdvanceTimeBlock.Accept(); err != nil {
		t.Fatal(err)
	}
	if err := commit.Accept(); err != nil {
		t.Fatal(err)
	}

	firstDelegatorStartTime := validatorStartTime.Add(syncBound).Add(1 * time.Second)
	firstDelegatorEndTime := firstDelegatorStartTime.Add(vm.MinStakeDuration)

	// create valid tx
	addFirstDelegatorTx, err := vm.newAddDelegatorTx(
		4*vm.MinValidatorStake, // maximum amount of stake this delegator can provide
		uint64(firstDelegatorStartTime.Unix()),
		uint64(firstDelegatorEndTime.Unix()),
		id,
		keys[0].PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{keys[0], keys[1]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}

	// trigger block creation
	if err := vm.mempool.IssueTx(addFirstDelegatorTx); err != nil {
		t.Fatal(err)
	}
	addFirstDelegatorBlockIntf, err := vm.BuildBlock()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the proposed block
	if err := addFirstDelegatorBlockIntf.Verify(); err != nil {
		t.Fatal(err)
	}

	// Assert preferences are correct
	addFirstDelegatorBlock := addFirstDelegatorBlockIntf.(*ProposalBlock)
	options, err = addFirstDelegatorBlock.Options()
	if err != nil {
		t.Fatal(err)
	}

	// verify the commit block
	commit = options[0].(*CommitBlock)
	if err := commit.Verify(); err != nil {
		t.Fatal(err)
	}

	// Accept the proposal block and the commit block
	if err := addFirstDelegatorBlock.Accept(); err != nil {
		t.Fatal(err)
	}
	if err := commit.Accept(); err != nil {
		t.Fatal(err)
	}

	vm.clock.Set(firstDelegatorStartTime)

	secondAdvanceTimeBlockIntf, err := vm.BuildBlock()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the proposed block
	if err := secondAdvanceTimeBlockIntf.Verify(); err != nil {
		t.Fatal(err)
	}

	// Assert preferences are correct
	secondAdvanceTimeBlock := secondAdvanceTimeBlockIntf.(*ProposalBlock)
	options, err = secondAdvanceTimeBlock.Options()
	if err != nil {
		t.Fatal(err)
	}

	// verify the commit block
	commit = options[0].(*CommitBlock)
	if err := commit.Verify(); err != nil {
		t.Fatal(err)
	}

	// Accept the proposal block and the commit block
	if err := secondAdvanceTimeBlock.Accept(); err != nil {
		t.Fatal(err)
	}
	if err := commit.Accept(); err != nil {
		t.Fatal(err)
	}

	secondDelegatorStartTime := firstDelegatorEndTime.Add(2 * time.Second)
	secondDelegatorEndTime := secondDelegatorStartTime.Add(vm.MinStakeDuration)

	vm.clock.Set(secondDelegatorStartTime.Add(-10 * syncBound))

	// create valid tx
	addSecondDelegatorTx, err := vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(secondDelegatorStartTime.Unix()),
		uint64(secondDelegatorEndTime.Unix()),
		id,
		keys[0].PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{keys[0], keys[1], keys[3]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}

	// trigger block creation
	if err := vm.mempool.IssueTx(addSecondDelegatorTx); err != nil {
		t.Fatal(err)
	}
	addSecondDelegatorBlockIntf, err := vm.BuildBlock()
	if err != nil {
		t.Fatal(err)
	}

	// Verify the proposed block
	if err := addSecondDelegatorBlockIntf.Verify(); err != nil {
		t.Fatal(err)
	}

	// Assert preferences are correct
	addSecondDelegatorBlock := addSecondDelegatorBlockIntf.(*ProposalBlock)
	options, err = addSecondDelegatorBlock.Options()
	if err != nil {
		t.Fatal(err)
	}

	// verify the commit block
	commit = options[0].(*CommitBlock)
	if err := commit.Verify(); err != nil {
		t.Fatal(err)
	}

	// Accept the proposal block and the commit block
	if err := addSecondDelegatorBlock.Accept(); err != nil {
		t.Fatal(err)
	}
	if err := commit.Accept(); err != nil {
		t.Fatal(err)
	}

	thirdDelegatorStartTime := firstDelegatorEndTime.Add(-time.Second)
	thirdDelegatorEndTime := thirdDelegatorStartTime.Add(vm.MinStakeDuration)

	// create valid tx
	addThirdDelegatorTx, err := vm.newAddDelegatorTx(
		vm.MinDelegatorStake,
		uint64(thirdDelegatorStartTime.Unix()),
		uint64(thirdDelegatorEndTime.Unix()),
		id,
		keys[0].PublicKey().Address(),
		[]*crypto.PrivateKeySECP256K1R{keys[0], keys[1], keys[4]},
		ids.ShortEmpty, // change addr
	)
	if err != nil {
		t.Fatal(err)
	}

	// trigger block creation
	if err := vm.mempool.IssueTx(addThirdDelegatorTx); err != nil {
		t.Fatal(err)
	}

	// Verify the proposed tx is invalid
	if _, err := vm.BuildBlock(); err == nil {
		t.Fatalf("should have marked the delegator as being over delegated")
	}
}
