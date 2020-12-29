package chain

import "testing"

func TestGetAddressBalance(t *testing.T) {
	arr, err := EthAccounts()
	if err != nil {
		t.Error(err)
	}

	ba, err := GetAddressBalance(arr[0])
	if err != nil {
		t.Error(err)
	}
	t.Log(ba)
}

func TestEthAccounts(t *testing.T) {
	arr, err := EthAccounts()
	if err != nil {
		t.Error(err)
	}
	t.Log(arr)
}

func TestSendRewardTransaction(t *testing.T) {
	from := "0x8cab3e7cf1f49392460967b52ae0b907b66aeb9c"
	arr, err := EthAccounts()
	if err != nil {
		t.Error(err)
	}

	tx, err := SendRewardTransaction(from, arr[0], 1,nil)
	if err != nil {
		t.Error(err)
	}
	t.Log(tx.String())
}

func TestGetEthPendingTransactions(t *testing.T) {
	length, err := GetEthPendingTransactions()
	if err != nil {
		t.Error(err)
	}
	t.Log(length)
}

func TestGetRewardTransactionReceipt(t *testing.T) {
	//begin := "0xeb89f11867f6c10c93b07070048a62bb57466dc3ef1f2ee010be9c2155c1329d"
	//end := "0x883fd44fea77162cabbe070e66e4e451f2ffa6c1ce8141e79d309a7fb1cbeefe"
	//
	//rec, err := GetRewardTransactionReceipt(begin)
	//if err != nil {
	//	t.Error(err)
	//}
	//if rec != nil {
	//	t.Log(rec.BlockNumber.Int64())
	//}
	//
	//rec, err = GetRewardTransactionReceipt(end)
	//if err != nil {
	//	t.Error(err)
	//}
	//if rec != nil {
	//	t.Log(rec.BlockNumber.Int64())
	//}

	tx :="0xc16c84ff39090a7f854e13c1ccac392d3805bb57c9aed01e1af75966604d3784"

	res,err := GetRewardTransactionReceipt(tx)
	if err != nil{
		t.Error(err)
	}

	t.Log(res)
}
