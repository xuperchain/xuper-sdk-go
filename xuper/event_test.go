package xuper

import (
	"fmt"
	"testing"
	"time"
)

func TestInitWatcher(t *testing.T) {

	client := newClient()

	watcher := InitWatcher(client, 10, false)

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

	time.Sleep(time.Second * 1)
	if testNode != "" {
		reg.Unregister()
	}
}

func TestEventOpts(t *testing.T) {
	filter, err := NewBlockFilter("xuper",
		WithAuthRequire("a"),
		WithBlockRange("1", "10"),
		WithContract("counter"),
		WithEventName("event"),
		WithExcludeTx(true),
		WithExcludeTxEvent(true),
		WithFromAddr("a"),
		WithInitiator("a"),
		WithToAddr("b"))
	if err != nil {
		t.Error(err)
	}
	if filter.AuthRequire != "a" ||
		filter.Contract != "counter" {
		t.Error("Event opts assert failed")
	}
}
