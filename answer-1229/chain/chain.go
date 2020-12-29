package chain

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
)

//SendTxArgs represents the arguments to sumbit a new transaction into the transaction pool.
type SendTxArgs struct {
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`
}

//组装交易并向以太坊节点发送交易
func SendRewardTransaction(from string, addr string, reward int64, input *hexutil.Bytes) (common.Hash, error) {
	//这里应该输入以太坊主网节点
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return common.Hash{}, err
	}
	defer cli.Close()

	to := common.HexToAddress(addr)
	tx := SendTxArgs{
		From: common.HexToAddress(from),
		To:   &to,
		Data: input,
	}
	rewardWei := new(big.Int).Mul(big.NewInt(1e+15), new(big.Int).SetInt64(reward))
	tx.Value = (*hexutil.Big)(rewardWei)
	var txHash common.Hash
	con := context.Background()
	err = cli.CallContext(con, &txHash, "eth_sendTransaction", tx)
	if err != nil {
		log.Error("err", err)
		return common.Hash{}, err
	}
	return txHash, nil
}

//查询账户余额
func GetAddressBalance(address string) (string, error) {
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return "", err
	}
	defer cli.Close()

	to := common.HexToAddress(address)

	con := context.Background()
	var res hexutil.Big
	err = cli.CallContext(con, &res, "eth_getBalance", to, "latest")
	if err != nil {
		log.Error("err", err)
		return "", err
	}
	balance := res.ToInt().String()
	if len(balance) > 18 {
		return balance[:len(balance)-18] + "." + balance[len(balance)-18:], nil
	}
	return "0" + "." + balance, nil
}

func GetRewardTransactionReceipt(hash string) (*types.Receipt, error) {
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	defer cli.Close()

	to := common.HexToHash(hash)
	con := context.Background()
	var res *types.Receipt

	err = cli.CallContext(con, &res, "eth_getTransactionReceipt", to)
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	return res, nil
}

func EthAccounts() ([]string, error) {
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	defer cli.Close()

	con := context.Background()
	var address []string

	err = cli.CallContext(con, &address, "eth_accounts")
	if err != nil {
		log.Error("err", err)
		return nil, err
	}
	return address, nil
}

func UnlockAccount(address string) error {
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return err
	}
	defer cli.Close()

	con := context.Background()
	var status bool
	err = cli.CallContext(con, &status, "personal_unlockAccount", address, "", 0)
	if err != nil {
		log.Error("err", err)
		return err
	}
	return nil
}

//查询账户余额
func GetEthPendingTransactions() (int, error) {
	cli, err := rpc.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Error("err", err)
		return 0, err
	}
	defer cli.Close()

	con := context.Background()
	var res []interface{}
	err = cli.CallContext(con, &res, "eth_pendingTransactions")
	if err != nil {
		log.Error("err", err)
		return 0, err
	}
	log.Info("res", res)
	return len(res), nil
}
