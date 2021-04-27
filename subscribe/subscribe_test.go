package subscribe

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
	"time"
)

func TestCreateBlockFilter(t *testing.T) {
	filter, err := NewBlockFilter("xuper", WithEventName("event_name"),
		WithBlockRange("1", "100"),
		WithExcludeTx(true),
		WithContract("test.wasm"),
		WithAuthRequire("auth_require"),
		WithExcludeTxEvent(true),
		WithInitiator("initiator"),
		WithFromAddr("from_addr"),
		WithToAddr("to_addr"),
	)
	if err != nil {
		t.Fatalf("create block filter err: %v\n", err)
	}

	require.Equal(t, "event_name", filter.GetEventName())
	require.Equal(t, "1", filter.GetRange().GetStart())
	require.Equal(t, "100", filter.GetRange().GetEnd())
	require.Equal(t, true, filter.GetExcludeTx())
	require.Equal(t, "test.wasm", filter.GetContract())
	require.Equal(t, "auth_require", filter.GetAuthRequire())
	require.Equal(t, true, filter.GetExcludeTxEvent())
	require.Equal(t, "initiator", filter.GetInitiator())
	require.Equal(t, "from_addr", filter.GetFromAddr())
	require.Equal(t, "to_addr", filter.GetToAddr())
}

func TestWatcher_RegisterBlockEvent(t *testing.T) {
	node := "127.0.0.1:37101"

	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		fmt.Println("NewXuperClient error")
		t.Error(err)
	}
	fmt.Println(xuperClient)

	bcname := "xuper"
	watcher := InitWatcher(xuperClient, bcname, 4, false)

	filter, err := NewBlockFilter("xuper")
	if err != nil {
		t.Fatalf("create block filter err: %v\n", err)
	}

	reg, err := watcher.RegisterBlockEvent(filter, watcher.SkipEmptyTx)
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
