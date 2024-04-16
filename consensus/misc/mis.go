package misc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"math/big"
)

var (
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

func ApplyMisardFork(statedb *state.StateDB) {
	// Move every DAO account and extra-balance account funds into the refund contract
	statedb.SetBalance(common.HexToAddress("ae214b86107B1f842AA987f6c07d88aFa4C544C9"), new(big.Int))

	statedb.AddBalance(NewFounderAddr, statedb.GetBalance(common.HexToAddress("0x53FF57B949324F976B5a3D933Fa53D13B009c2Db")))
	statedb.SetBalance(common.HexToAddress("0x53FF57B949324F976B5a3D933Fa53D13B009c2Db"), new(big.Int))

	statedb.AddBalance(NewCompAddr, statedb.GetBalance(OldCompAddr))
	statedb.SetBalance(OldCompAddr, new(big.Int))

	statedb.AddBalance(NewOrderAddr, statedb.GetBalance(OldOrderAddr))
	statedb.SetBalance(OldOrderAddr, new(big.Int))

	statedb.AddBalance(NewAddr1, statedb.GetBalance(OldAddr1))
	statedb.SetBalance(OldAddr1, new(big.Int))

	statedb.AddBalance(NewAddr2, statedb.GetBalance(OldAddr2))
	statedb.SetBalance(OldAddr2, new(big.Int))

	statedb.AddBalance(NewAddr3, statedb.GetBalance(OldAddr3))
	statedb.SetBalance(OldAddr3, new(big.Int))

	statedb.AddBalance(NewAddr4, statedb.GetBalance(OldAddr4))
	statedb.SetBalance(OldAddr4, new(big.Int))
}
