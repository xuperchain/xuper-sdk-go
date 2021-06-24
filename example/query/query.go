package main

import (
	"fmt"

	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	// 查询接口有一些返回的是 xuperchain 项目中的 pb 结构，有的是返回 sdk 中的结构。
	// 对于返回 xuperchain 中的 pb 结构，可能会带有 Header，需要先判断 Header 中的 error。
	// 示例如下：
	queryBlockByHeight()

	// 查询余额时，可以指定链名，如果没有创建平行链，默认使用 xuper。
	queryBalance()
}

func queryBlockByHeight() {
	// 示例代码省略了 err 的检查。
	node := "127.0.0.1"
	xclient, _ := xuper.New(node)
	blockResult, _ := xclient.QueryBlockByHeight(8)
	if blockResult.GetHeader().GetError() != 0 {
		// 处理错误。
	} else {
		// 处理区块数据。
		block := blockResult.GetBlock()
		fmt.Println(block.GetBlockid())
		fmt.Println(block.GetHeight())
		fmt.Println(block.GetTxCount())
	}
}

func queryBalance() {
	// 示例代码省略了 err 的检查。
	node := "127.0.0.1"
	xclient, _ := xuper.New(node)

	// 查询账户余额，默认 xuper 链。
	bal, _ := xclient.QueryBalance("nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh")
	fmt.Println(bal)

	// 查询账户余额，在 hello 链。
	bal, _ = xclient.QueryBalance("nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh", xuper.WithQueryBcname("hello"))
	fmt.Println(bal)

	// 查询账户余额详细数据
	balDetails, _ := xclient.QueryBalanceDetail("nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh")
	for _, bd := range balDetails {
		fmt.Println(bd.Balance)
		fmt.Println(bd.IsFrozen)
	}

}
