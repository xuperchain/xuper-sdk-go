package main

import (
	"fmt"
	"time"

	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	testEvent()
}

func testEvent() error {
	// 创建节点客户端
	client, err := xuper.New("127.0.0.1:37101")
	if err != nil {
		return err
	}
	watcher := xuper.InitWatcher(client, 10, false)
	filter, err := xuper.NewBlockFilter("xuper") // 此处可以添加其他顾虑条件：xuper.WithContract() 等。
	if err != nil {
		return err
	}

	reg, err := watcher.RegisterBlockEvent(filter, watcher.SkipEmptyTx)
	if err != nil {
		return err
	}

	go func() {
		for {
			b := <-reg.FilteredBlockChan
			fmt.Printf("%+v\n", b)
		}
	}()

	time.Sleep(time.Second * 10)

	reg.Unregister()
	client.Close()
	return nil
}
