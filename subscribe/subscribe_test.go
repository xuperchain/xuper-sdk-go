package subscribe

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"io"
	"testing"
)

func TestSubscribe_Subscribe(t *testing.T) {
	url := "127.0.0.1:37101"
	sdkClient, err := xchain.NewSDKClient(url)
	if err != nil {
		t.Errorf("newSDK error:%s\n", err.Error())
	}
	subscribe := InitSubscribe(sdkClient, "xuper", false)
	filter := &pb.BlockFilter{}
	client, err := subscribe.Subscribe(filter)
	if err != nil {
		t.Errorf("Subscribe error:%s\n", err.Error())
	}
	for {
		event, err := (*client).Recv()
		if err == io.EOF {
			t.Error(err)
		}
		if err != nil {
			t.Error(err)
		}
		var block pb.FilteredBlock
		err = proto.Unmarshal(event.Payload, &block)
		if err != nil {
			t.Error(err)
		}
		if len(block.GetTxs()) == 0 && subscribe.SkipEmptyTx {
			continue
		}
		fmt.Println(&block)
	}
}
