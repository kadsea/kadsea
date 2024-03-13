package systemcontract

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/congress/vmcaller"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math"
	"math/big"
)

const (
	punishV2Code = "0x608060405234801561001057600080fd5b50600436106100f55760003560e01c8063ab3b5d4b11610097578063d6eb3de111610066578063d6eb3de114610241578063e0d8ea531461024a578063ea7221a114610252578063f62af26c1461026557600080fd5b8063ab3b5d4b146101da578063be645692146101ec578063bfba503d146101fd578063c967f90f1461022657600080fd5b806332f3c17f116100d357806332f3c17f146101525780633a061bd3146101a95780635e7ba496146101b25780638129fc1c146101d257600080fd5b8063158ef93e146100fa5780631b5e358c1461011c5780631e8877271461013d575b600080fd5b6000546101079060ff1681565b60405190151581526020015b60405180910390f35b61012561f20081565b6040516001600160a01b039091168152602001610113565b61015061014b366004610bc1565b610278565b005b610165610160366004610bc1565b6102db565b6040516101139190600060a0820190508251825260208301516020830152604083015160408301526060830151606083015260808301511515608083015292915050565b61012561f10081565b6101c56101c0366004610bea565b61036b565b6040516101139190610c0c565b610150610490565b6006545b604051908152602001610113565b6101de69021e19e0c9bab240000081565b6101de61020b366004610bc1565b6001600160a01b031660009081526005602052604090205490565b61022e601581565b60405161ffff9091168152602001610113565b61012561f30081565b6002546101de565b610150610260366004610bc1565b6104f0565b610125610273366004610c6e565b610975565b3361f100146102ce5760405162461bcd60e51b815260206004820152601860248201527f56616c696461746f727320636f6e7472616374206f6e6c79000000000000000060448201526064015b60405180910390fd5b6102d78161099f565b5050565b61030f6040518060a00160405280600081526020016000815260200160008152602001600081526020016000151581525090565b506001600160a01b0316600090815260016020818152604092839020835160a0810185528154815292810154918301919091526002810154928201929092526003820154606082015260049091015460ff161515608082015290565b60606103778383610c9d565b67ffffffffffffffff81111561038f5761038f610cb4565b6040519080825280602002602001820160405280156103ed57816020015b6103da604051806060016040528060006001600160a01b0316815260200160008152602001600081525090565b8152602001906001900390816103ad5790505b509050825b82811015610489576006818154811061040d5761040d610cca565b600091825260209182902060408051606081018252600390930290910180546001600160a01b03168352600181015493830193909352600290920154918101919091528261045b8684610c9d565b8151811061046b5761046b610cca565b6020026020010181905250808061048190610ce0565b9150506103f2565b5092915050565b60005460ff16156104d95760405162461bcd60e51b8152602060048201526013602482015272105b1c9958591e481a5b9a5d1a585b1a5e9959606a1b60448201526064016102c5565b600080546001600160a81b03191662f10001179055565b33411461052c5760405162461bcd60e51b815260206004820152600a6024820152694d696e6572206f6e6c7960b01b60448201526064016102c5565b60005460ff1661056d5760405162461bcd60e51b815260206004820152600c60248201526b139bdd081a5b9a5d081e595d60a21b60448201526064016102c5565b4360009081526003602052604090205460ff16156105c05760405162461bcd60e51b815260206004820152601060248201526f105b1c9958591e481c1d5b9a5cda195960821b60448201526064016102c5565b436000908152600360209081526040808320805460ff191660019081179091556001600160a01b038516845290915290206004015460ff1661066c57600280546001600160a01b038316600081815260016020819052604082208082018590558185019095557f405787fa12a823e0f2b7631cc41b3ba8828b3321ca811111fa75cd3aa3bb5ace90930180546001600160a01b03191683179055526004909101805460ff191690911790555b6001600160a01b038116600090815260016020526040812080549161069083610ce0565b90915550506001600160a01b03811660009081526001602081905260409091205414156106d7576001600160a01b0381166000908152600160205260409020426002909101555b6001600160a01b03811660009081526001602052604081206003018054916106fe83610ce0565b90915550506001600160a01b0381166000908152600160205260409020600201541580159061075657506001600160a01b03811660009081526001602052604090206002015462093a80906107539042610c9d565b10155b15610972576000546040516380d2dde160e01b81526001600160a01b03838116600483015269010f0cf064dd592000006024830181905292610100900416906380d2dde190604401600060405180830381600087803b1580156107b857600080fd5b505af11580156107cc573d6000803e3d6000fd5b5050604080516060810182526001600160a01b0386811680835260208084018881524385870190815260068054600181018255600091825287517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d3f600390920291820180546001600160a01b031916919098161790965591517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d40860155517ff652222313e28459528d920b65115c16c04f3efc82aaedc97be59f3f377c0d41909401939093559082526005905291909120549092506108ac915083610b5b565b6001600160a01b03848116600081815260056020526040808220949094555492516340a141ff60e01b8152600481019190915261010090920416906340a141ff90602401600060405180830381600087803b15801561090a57600080fd5b505af115801561091e573d6000803e3d6000fd5b5050505061092b8361099f565b50604080518381524260208201526001600160a01b038516917f275fc6ed9b56be0a0510b84cac142c60905435f8dac675a63ea794679daa29fd910160405180910390a250505b50565b6002818154811061098557600080fd5b6000918252602090912001546001600160a01b0316905081565b6001600160a01b038116600090815260016020526040812054156109d7576001600160a01b0382166000908152600160205260408120555b6001600160a01b03821660009081526001602052604090206004015460ff168015610a03575060025415155b15610b5357600254610a1790600190610c9d565b6001600160a01b0383166000908152600160208190526040909120015414610ae8576002805460009190610a4d90600190610c9d565b81548110610a5d57610a5d610cca565b60009182526020808320909101546001600160a01b0386811684526001928390526040909320909101546002805493909216935083928110610aa157610aa1610cca565b600091825260208083209190910180546001600160a01b0319166001600160a01b039485161790558583168252600190819052604080832082015494909316825291902001555b6002805480610af957610af9610cfb565b60008281526020808220830160001990810180546001600160a01b03191690559092019092556001600160a01b03841682526001908190526040822090810182905560048101805460ff1916905560028101829055600301555b506001919050565b600080610b688385610d11565b905083811015610bba5760405162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f77000000000060448201526064016102c5565b9392505050565b600060208284031215610bd357600080fd5b81356001600160a01b0381168114610bba57600080fd5b60008060408385031215610bfd57600080fd5b50508035926020909101359150565b602080825282518282018190526000919060409081850190868401855b82811015610c6157815180516001600160a01b0316855286810151878601528501518585015260609093019290850190600101610c29565b5091979650505050505050565b600060208284031215610c8057600080fd5b5035919050565b634e487b7160e01b600052601160045260246000fd5b600082821015610caf57610caf610c87565b500390565b634e487b7160e01b600052604160045260246000fd5b634e487b7160e01b600052603260045260246000fd5b6000600019821415610cf457610cf4610c87565b5060010190565b634e487b7160e01b600052603160045260246000fd5b60008219821115610d2457610d24610c87565b50019056fea264697066735822122001676204fe99b7748e8346d4a8e53e5de46d995ba86297fd04b2a1bb8149dfa764736f6c63430008080033"
)

type hardForkPunishV2 struct {
}

func (s *hardForkPunishV2) GetName() string {
	return PunishV2ContractName
}

func (s *hardForkPunishV2) Update(config *params.ChainConfig, height *big.Int, state *state.StateDB) (err error) {
	contractCode := common.FromHex(punishV2Code)

	//write code to sys contract
	state.SetCode(PunishV2ContractAddr, contractCode)
	log.Debug("Write code to system contract account", "addr", PunishV2ContractAddr.String(), "code", punishV2Code)

	return
}

func (s *hardForkPunishV2) Execute(state *state.StateDB, header *types.Header, chainContext core.ChainContext, config *params.ChainConfig) (err error) {
	// initialize v1 contract
	method := "initialize"
	data, err := GetInteractiveABI()[s.GetName()].Pack(method)
	if err != nil {
		log.Error("Can't pack data for initialize", "error", err)
		return err
	}

	msg := vmcaller.NewLegacyMessage(header.Coinbase, &PunishV2ContractAddr, 0, new(big.Int), math.MaxUint64, new(big.Int), data, false)
	_, err = vmcaller.ExecuteMsg(msg, state, header, chainContext, config)

	return
}
