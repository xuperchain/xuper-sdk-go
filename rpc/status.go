package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/util"
	"google.golang.org/grpc"
)

func QueryStatus(node string, bcname ...string) (*util.SystemStatus, error) {

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	reply := &pb.SystemsStatusReply{
		SystemsStatus: &pb.SystemsStatus{
			BcsStatus: make([]*pb.BCStatus, 0),
		},
	}

	//if not input bcname then query all chain
	switch bcname == nil {
	//query all chain
	case true:
		reply, err = client.GetSystemStatus(ctx, &pb.CommonIn{})

	//query one chain
	case false:
		var bcStatus *pb.BCStatus
		bcStatus, err = client.GetBlockChainStatus(ctx, &pb.BCStatus{Bcname: bcname[0]})
		if bcStatus != nil {
			reply.Header = bcStatus.Header
			reply.SystemsStatus.BcsStatus = append(reply.SystemsStatus.BcsStatus, bcStatus)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("query node status fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		msg := reply.Header.Error.String()
		if msg == "CONNECT_REFUSE" {
			msg = "the chain not exist"
		}
		return nil, fmt.Errorf("query node status fail, err: %s", msg)
	}
	//if reply.SystemsStatus == nil {
	//	return nil, fmt.Errorf("status not found")
	//}
	return util.FromSystemStatusPB(reply.SystemsStatus), nil
}
