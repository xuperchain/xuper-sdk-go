package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"
	"google.golang.org/grpc"
)

func QueryBalance(node string, bcname, address string) (*pb.AddressBalanceStatus, error) {

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	addStatus := &pb.AddressBalanceStatus{
		Address: address,
		Tfds:    []*pb.TokenFrozenDetails{{Bcname: bcname}},
	}

	reply, err := client.GetBalanceDetail(ctx, addStatus)
	if err != nil {
		return nil, fmt.Errorf("query balance fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, fmt.Errorf("query balance fail, err: %s", reply.Header.Error.String())
	}
	if reply.Tfds == nil {
		return nil, fmt.Errorf("balance not found")
	}

	reply.Header = nil
	return reply, nil
}
