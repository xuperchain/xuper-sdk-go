package xuper

import (
	"testing"
)

func TestEventOpts(t *testing.T) {
	xclient := &XClient{}
	watcher, err := xclient.newWatcher(
		WithBlockEventBcname("xuper"),
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
	if watcher.opt.blockFilter.AuthRequire != "a" ||
		watcher.opt.blockFilter.Contract != "counter" {
		t.Error("Event opts assert failed")
	}
}
