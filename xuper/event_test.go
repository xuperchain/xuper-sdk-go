package xuper

import (
	"fmt"
	"testing"
	"time"
)

func TestInitWatcher(t *testing.T) {
	client,err := New("127.0.0.1:37101")
	if err != nil {
		fmt.Printf("err:%s\n",err.Error())

		t.Error(err)
	}
	watcher := InitWatcher(client,"xuper",10,false)

	filter, err := NewBlockFilter("xuper")
	if err != nil {
		t.Fatalf("create block filter err: %v\n", err)
	}

	reg,err := watcher.RegisterBlockEvent(filter,watcher.SkipEmptyTx)
	if err != nil {
		t.Error("RegisterBlockEvent")
		t.Error(err)
	}
	go func() {
		for {
			b := <-reg.FilteredBlockChan
			fmt.Printf("%+v\n", b)
		}
	}()

	time.Sleep(time.Second * 10)
	reg.Unregister()

}
