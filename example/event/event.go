package main

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
	"time"
)

//
//	time.Sleep(time.Second * 10)
//	reg.Unregister()
//
//}

func testEvent() error {
	client, err := xuper.New("127.0.0.1:37101")
	if err != nil {
		return err
	}
	watcher := xuper.InitWatcher(client, 10, false)
	filter, err := xuper.NewBlockFilter("xuper")
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
	return nil
}

func main() {
	testEvent()
}
