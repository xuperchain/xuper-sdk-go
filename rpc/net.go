package rpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"
	"google.golang.org/grpc"
)

func QueryNetUrl(node string) (string, error) {

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return "", fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	client := pb.NewXchainClient(conn)

	reply, err := client.GetNetURL(ctx, &pb.CommonIn{})
	if err != nil {
		return "", fmt.Errorf("query neturl fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return "", fmt.Errorf("query neturl fail, err: %s", reply.Header.Error.String())
	}
	if reply.RawUrl == "" {
		return "", fmt.Errorf("neturl not found")
	}

	i := strings.Index(node, ":")
	return strings.ReplaceAll(reply.RawUrl, "127.0.0.1", node[:i]), nil
}
