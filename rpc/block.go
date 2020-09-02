package rpc

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/util"
	"google.golang.org/grpc"
)

//json格式输出
//	iblock := FromInternalBlockPB(block.Block)
//	output, err := json.MarshalIndent(iblock, "", "  ")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(string(output))

func QueryBlockById(node, bcname, blockid string) (*util.InternalBlock, error) {

	rawBlockid, err := hex.DecodeString(blockid)
	if err != nil {
		return nil, fmt.Errorf("blockid invalid, err: %s", err)
	}

	blockIDPB := &pb.BlockID{
		Bcname:      bcname,
		Blockid:     rawBlockid,
		NeedContent: true,
	}

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	client := pb.NewXchainClient(conn)
	reply, err := client.GetBlock(ctx, blockIDPB)
	if err != nil {
		return nil, fmt.Errorf("query block fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, fmt.Errorf("query block fail, err: %s", reply.Header.Error.String())
	}
	if reply.Block == nil {
		return nil, fmt.Errorf("block not found")
	}
	return util.FromInternalBlockPB(reply.Block), nil
}

func QueryBlockByHeight(node, bcname string, height int64) (*util.InternalBlock, error) {

	blockHeightPB := &pb.BlockHeight{
		Bcname: bcname,
		Height: height,
	}

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	client := pb.NewXchainClient(conn)
	reply, err := client.GetBlockByHeight(ctx, blockHeightPB)
	if err != nil {
		return nil, fmt.Errorf("query block fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, fmt.Errorf("query block fail, err: %s", reply.Header.Error.String())
	}
	if reply.Block == nil {
		return nil, fmt.Errorf("block not found")
	}

	return util.FromInternalBlockPB(reply.Block), nil
}
