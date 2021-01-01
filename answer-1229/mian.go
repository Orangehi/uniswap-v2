package main

import (
	"bufio"
	"github.com/uniswap-v2/answer-1229/chain"
	"fmt"
	"os"
	"strconv"
)

/*
Question 1:
Provide a coding solution to the following token swap estimations:
1.Estimate the number of DAI's output from any given number of ETH's.
2.Estimate the number of ETH's output from any given number of DAI's.
Requirements:
1.Coding with golang/java/NodeJS.
2.Provide a function to get priceImpact.
3.Provide a function to get minimumAmountOut.
4.Provide a function to get midPrice.
5.Sign and broadcast the tx.
Tips:
1.This glossary can help to understand priceImpact, minimumAmountOut, midPrice concepts. https://uniswap.org/docs/v2/protocol-overview/glossary/
2.UniSwap is an open-source project. https://github.com/Uniswap
3.There could be a cost/gas fee involved if you test your code on Mainnet, so be cautious.
4.Interact directly with the UniSwap web app to get the results are not accepted as a solution. However, feel free to use the underlying UniSwap SDK as the tool to solve the question.
 */

/*
   	1. 可以通过直接调用 uniswap-sdk 库函数来完成上述问题，该库函数由TypeScript语言开发.
	2. 直接调用 uniswap 核心合约函数来解决上述问题.
	3. 直接通过调用uniswap-v2-periphery合约来完成上述问题

	本次回答使用的是上述的第三种解决方案，使用golang作为开发语言，通过与以太坊节点rpc通讯来发送交易，完成Dai币与ETH的互换

 */

const (
	/*
		该地址为uniswap-v2-periphery合约部署在以太坊网络上的地址
		出处：https://uniswap.org/docs/v2/smart-contracts/router02/#address
	 */
	to = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	/*
		该地址为本地通过调用节点RPC函数（personal_newAccount）生成
	 */
	from = "0xa9deb3d33967548093ce3f1661ee92bed3a8bf45"
)


func main() {
	r := bufio.NewReader(os.Stdin)
	//输入eth的数量，打印兑换的dai币数量
	data1, _, _ := r.ReadLine()
	count1, err := strconv.Atoi(string(data1))
	if err != nil {
		return
	}
	// ETH -> Dai
	amounts := EthToDai(4)
	fmt.Println(int64(count1))

	//输入Dai币的数量，打印兑换的eth数量
	data2, _, _ := r.ReadLine()
	count2, err := strconv.Atoi(string(data2))
	if err != nil {
		return
	}
	//Dai -> ETH
	amounts = DaiToEth(int64(count2))
	fmt.Println(amounts)
}

//The difference between the mid-price and the execution price of a trade.
func getPriceImpact(executionPrice uint,midPrice uint)uint{

	if executionPrice > midPrice{
		return executionPrice-midPrice
	}

	if midPrice < executionPrice{
		return midPrice-executionPrice
	}

	return 0
}

//The price between what users can buy and sell tokens at a given moment. In Uniswap this is the ratio of the two ERC20 token reserves.
func getMidPrice(token0,token1 uint)uint{
	return token0/token1
}

/*
minimumAmountOut (since 2.0.4)
minimumAmountOut(slippageTolerance: Percent): TokenAmount
Returns the minimum amount of the output token that should be received from a trade, given the slippage tolerance.
 */
func getMinimumAmountOut(tx string)uint{
	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r,err := chain.GetRewardTransactionReceipt(tx)
	if err != nil{
		return 0
	}

	//
	if len(r.Logs)==0 {
		return 0
	}
	if len(r.Logs[0].Topics) == 0{
		return 0
	}
	s, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	return uint(s)
}




/*
swapExactETHForTokens
function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)
  external
  payable
  returns (uint[] memory amounts);
Copy
Swaps an exact amount of ETH for as many output tokens as possible, along the route determined by the path. The first element of path must be WETH, the last is the output token, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

Name	Type
msg.value (amountIn)	uint	The amount of ETH to send.
amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
to	address	Recipient of the output tokens.
deadline	uint	Unix timestamp after which the transaction will revert.
amounts	uint[] memory	The input token amount and all subsequent output token amounts.

调用智能合约v2中的函数（swapExactETHForTokens）将Eth币转换为Dai,返回可兑换的数量
 */
func EthToDai(amounts int64)uint{

	//TODO input字段数据需要通过该智能合约的abi生成，我是用智能合约编辑器（http://remix.hubwiz.com/#optimize=false&version=soljson-v0.6.6+commit.6c089d02.js）来生成的数据
	tx ,err := chain.SendRewardTransaction(from,to,amounts,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r,err := chain.GetRewardTransactionReceipt(tx.String())
	if err != nil{
		return 0
	}

	//
	if len(r.Logs)==0 {
		return 0
	}
	if len(r.Logs[0].Topics) == 0{
		return 0
	}
	 s, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	 if err != nil{
		 return 0
	 }

	return uint(s)
}

/*
swapExactTokensForETH
function swapExactTokensForETH(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline)
  external
  returns (uint[] memory amounts);
Copy
Swaps an exact amount of tokens for as much ETH as possible, along the route determined by the path. The first element of path is the input token, the last must be WETH, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

If the to address is a smart contract, it must have the ability to receive ETH.
Name	Type
amountIn	uint	The amount of input tokens to send.
amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
to	address	Recipient of the ETH.
deadline	uint	Unix timestamp after which the transaction will revert.
amounts	uint[] memory	The input token amount and all subsequent output token amounts.

调用智能合约v2中的函数（swapExactTokensForETH）将Dai币转换为ETH,返回可兑换的数量
 */
func DaiToEth(amounts int64)uint{
	//TODO input字段测试数据需要通过该智能合约的abi生成，我是用智能合约编辑器（http://remix.hubwiz.com/#optimize=false&version=soljson-v0.6.6+commit.6c089d02.js）来生成的数据
	tx ,err := chain.SendRewardTransaction(from,to,amounts,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r,err := chain.GetRewardTransactionReceipt(tx.String())
	if err != nil{
		return 0
	}

	//
	if len(r.Logs)==0 {
		return 0
	}
	if len(r.Logs[0].Topics) == 0{
		return 0
	}
	s, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	return uint(s)
}


/*
Question 2:
Create your own DEX forked from UniSwap DEX and then provide your own liquidity.
Requirements:
1.Deploy your own DEX on Ropsten testNet.
2.Provide your own liquidity for DAI and ETH pair.
3.Swap the tokens in your own liquidity pool, sign and broadcast it. (Same as question 1)
Tips:
1.You do not have to do too much coding.
2.You need to understand the contract source code.
3.EIP-1014 will be helpful.
https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1014.md
 */

/*
	本题可以通过上面的两个函数DaiToEth与EthToDai来发送交易
 */




/*
Question 3 (Optional):
Find the best route path for swapping USDT to CRO (UniSwap)
1.Estimate the number of CRO’s output from any given number of USDT’s.
2.Estimate the number of USDT’s output from any given number of CRO’s.
Requirements:
1.Avoid the loss of precision as much as possible.
2.Find the best route path.
3.Sign and broadcast tx.
Tips:
1.This project will be helpful: https://github.com/Uniswap/uniswap-interface
2.This class `sac/hooks/Trades.ts` is the key point.

  	USDT属于erc20 token，为稳定币与美元挂钩
	CRO同样属于ERC20 token，由Crypto.com发行的erc20 token
	两者之间转换没有直接关系，需要先转换为WETH,然后在转换为对于ERC20 token
	1.	USDT -> CRO
	首先 通过调用合约方法（swapExactTokensForETH）将USDT转换为WETH,
	然后将转换过来的WETH通过调用合约方法（swapExactETHForTokens）转换为CRO

	2. 	CRO -> USDT
	首先 通过调用合约方法（swapExactTokensForETH）将CRO转换为WETH,
	然后将转换过来的WETH通过调用合约方法（swapExactETHForTokens）转换为USDT
 */

//CRO -> USDT
func CroToUsdt(amounts int64)uint{
	/*
		先将CRO转换为WETH
		swapExactTokensForETH
			function swapExactTokensForETH(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline)
	  			external
	 		 returns (uint[] memory amounts);
			Swaps an exact amount of tokens for as much ETH as possible, along the route determined by the path. The first element of path is the input token, the last must be WETH, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

			If the to address is a smart contract, it must have the ability to receive ETH.
			Name	Type
			amountIn	uint	The amount of input tokens to send.
			amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
			path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
			to	address	Recipient of the ETH.
			deadline	uint	Unix timestamp after which the transaction will revert.
			amounts	uint[] memory	The input token amount and all subsequent output token amounts.
	 */
	tx1 ,err := chain.SendRewardTransaction(from,to,amounts,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r1,err := chain.GetRewardTransactionReceipt(tx1.String())
	if err != nil{
		return 0
	}

	//
	if len(r1.Logs)==0 {
		return 0
	}
	if len(r1.Logs[0].Topics) == 0{
		return 0
	}
	//这里拿到了WETH的数量
	s1, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	/*
		然后将WETH转换为USDT
				swapExactETHForTokens
			function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)
			  external
			  payable
			  returns (uint[] memory amounts);
			Swaps an exact amount of ETH for as many output tokens as possible, along the route determined by the path. The first element of path must be WETH, the last is the output token, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

			Name	Type
			msg.value (amountIn)	uint	The amount of ETH to send.
			amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
			path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
			to	address	Recipient of the output tokens.
			deadline	uint	Unix timestamp after which the transaction will revert.
			amounts	uint[] memory	The input token amount and all subsequent output token amounts.
	 */
	tx2 ,err := chain.SendRewardTransaction(from,to, s1,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r2,err := chain.GetRewardTransactionReceipt(tx2.String())
	if err != nil{
		return 0
	}

	//
	if len(r2.Logs)==0 {
		return 0
	}
	if len(r2.Logs[0].Topics) == 0{
		return 0
	}
	//这里拿到了USDT的数量
	s2, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	return uint(s2)
}

//USDT -> CRO
func UsdtToCro(amounts int64)uint{
	/*
			先将USDT转换为WETH
			swapExactTokensForETH
				function swapExactTokensForETH(uint amountIn, uint amountOutMin, address[] calldata path, address to, uint deadline)
		  			external
		 		 returns (uint[] memory amounts);
				Swaps an exact amount of tokens for as much ETH as possible, along the route determined by the path. The first element of path is the input token, the last must be WETH, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

				If the to address is a smart contract, it must have the ability to receive ETH.
				Name	Type
				amountIn	uint	The amount of input tokens to send.
				amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
				path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
				to	address	Recipient of the ETH.
				deadline	uint	Unix timestamp after which the transaction will revert.
				amounts	uint[] memory	The input token amount and all subsequent output token amounts.
	*/
	tx1 ,err := chain.SendRewardTransaction(from,to,amounts,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r1,err := chain.GetRewardTransactionReceipt(tx1.String())
	if err != nil{
		return 0
	}

	//
	if len(r1.Logs)==0 {
		return 0
	}
	if len(r1.Logs[0].Topics) == 0{
		return 0
	}
	//这里拿到了WETH的数量
	s1, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	/*
		然后将WETH转换为CRO
				swapExactETHForTokens
			function swapExactETHForTokens(uint amountOutMin, address[] calldata path, address to, uint deadline)
			  external
			  payable
			  returns (uint[] memory amounts);
			Swaps an exact amount of ETH for as many output tokens as possible, along the route determined by the path. The first element of path must be WETH, the last is the output token, and any intermediate elements represent intermediate pairs to trade through (if, for example, a direct pair does not exist).

			Name	Type
			msg.value (amountIn)	uint	The amount of ETH to send.
			amountOutMin	uint	The minimum amount of output tokens that must be received for the transaction not to revert.
			path	address[] calldata	An array of token addresses. path.length must be >= 2. Pools for each consecutive pair of addresses must exist and have liquidity.
			to	address	Recipient of the output tokens.
			deadline	uint	Unix timestamp after which the transaction will revert.
			amounts	uint[] memory	The input token amount and all subsequent output token amounts.
	*/
	tx2 ,err := chain.SendRewardTransaction(from,to, s1,nil)
	if err != nil{
		return 0
	}

	//golang没有原生的web3库函数,智能通过修改智能合约结构(增加emit事件将结果返回)，然后通过解析交易回执的logs中的Topics拿到amounts字段
	r2,err := chain.GetRewardTransactionReceipt(tx2.String())
	if err != nil{
		return 0
	}

	//
	if len(r2.Logs)==0 {
		return 0
	}
	if len(r2.Logs[0].Topics) == 0{
		return 0
	}
	//这里拿到了CRO的数量
	s2, err := strconv.ParseInt(r.Logs[0].Topics[0].String(), 16, 32)
	if err != nil{
		return 0
	}

	return uint(s2)
}
