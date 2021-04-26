package subscribe

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

type Subscribe struct {
	xchain.Xchain
	SkipEmptyTx bool
}

func InitSubscribe(sdkClient *xchain.SDKClient, bcname string, skipEmptyTx bool) *Subscribe {
	return &Subscribe{
		xchain.Xchain{
			ChainName: bcname,
			SDKClient: sdkClient,
		},
		skipEmptyTx,
	}
}

func (s Subscribe) Subscribe(fileter *pb.BlockFilter) (*pb.EventService_SubscribeClient, error) {
	fileter.Bcname = s.ChainName
	buf, err := proto.Marshal(fileter)
	if err != nil {
		return nil, err
	}
	request := &pb.SubscribeRequest{
		Type:   pb.SubscribeType_BLOCK,
		Filter: buf,
	}

	eventClient := *(s.SDKClient.EventClient)
	stream, err := eventClient.Subscribe(context.Background(), request)
	if err != nil {
		return nil, err
	}
	stream.Recv()
	return &stream, nil

}
