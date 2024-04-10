// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	FeeRecoder     = common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff")
	MinerAddr      = common.HexToAddress("0xD2e191f70Bf86fb3f057b93DA1Ef21D03fCf9018")
	OldFounderAddr = common.HexToAddress("0xF4eab9b6fc2AF55c9a4847A49ED7f8c5da6a8FCC")
	NewFounderAddr = common.HexToAddress("0x57d6970aec438C3189F6B8979DbaFb1575461333")
	OldCompAddr    = common.HexToAddress("0x53FF57B949324F976B5a3D933Fa53D13B009c2Db")
	NewCompAddr    = common.HexToAddress("0x0cfaa1e0a1c656414e3eaaf013299bd403d4dc28")
	OldOrderAddr   = common.HexToAddress("0xF2d1531d171655b9342017b62a00f070B8421c67")
	NewOrderAddr   = common.HexToAddress("0x1b50bce4f1795df512a1007c30128a944c550f29")
	OldAddr1       = common.HexToAddress("0x712D0d5BE6c1185F529b4792F9e71055adE9F906")
	NewAddr1       = common.HexToAddress("0x5750d430dabebd1a04db1a3698da7abe7b64e745")
	OldAddr2       = common.HexToAddress("0x8A95EbE76944De758c672f405f870378085B9b82")
	NewAddr2       = common.HexToAddress("0x50c4cd040cb15322fb9195b98e89be1a888420e4")
	OldAddr3       = common.HexToAddress("0x96782C3375307fd7a7A17DA6Bf993D80814eDc19")
	NewAddr3       = common.HexToAddress("0xa47748ed1132563e67abec7c51ecb1e6f4f7b41a")
	OldAddr4       = common.HexToAddress("0x99C4F986a8A01Ee7C240B7FDFec0E0A85637F01B")
	NewAddr4       = common.HexToAddress("0xa736293d70bf0a1718e2bdbe9ab0252c349f8257")
)

// ChainHeaderReader defines a small collection of methods needed to access the local
// blockchain during header verification.
type ChainHeaderReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header
}

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	ChainHeaderReader

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs *[]*types.Transaction,
		uncles []*types.Header, receipts *[]*types.Receipt, systemTxs []*types.Transaction) error

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, []*types.Receipt, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close terminates any background threads maintained by the consensus engine.
	Close() error
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}

// PoSA is a consensus engine based on proof-of-stake-authority.
type PoSA interface {
	Engine

	// PreHandle runs any pre-transaction state modifications (e.g. apply hard fork rules).
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at pre-handling.
	PreHandle(chain ChainHeaderReader, header *types.Header, state *state.StateDB) error

	// IsSysTransaction checks whether a specific transaction is a system transaction.
	IsSysTransaction(sender common.Address, tx *types.Transaction, header *types.Header) (bool, error)

	// CanCreate determines where a given address can create a new contract.
	CanCreate(state StateReader, addr common.Address, height *big.Int) bool

	// ValidateTx do a consensus-related validation on the given transaction at the given header and state.
	ValidateTx(sender common.Address, tx *types.Transaction, header *types.Header, parentState *state.StateDB) error

	// CreateEvmExtraValidator returns a EvmExtraValidator if necessary.
	CreateEvmExtraValidator(header *types.Header, parentState *state.StateDB) types.EvmExtraValidator

	//Methods for debug trace

	// ApplySysTx applies a system-transaction using a given evm,
	// the main purpose of this method is for tracing a system-transaction.
	ApplySysTx(evm *vm.EVM, state *state.StateDB, txIndex int, sender common.Address, tx *types.Transaction) (ret []byte, vmerr error, err error)
}

type StateReader interface {
	GetState(addr common.Address, hash common.Hash) common.Hash
}
