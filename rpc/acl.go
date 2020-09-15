package rpc

import (
	"context"
	"fmt"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"
	"google.golang.org/grpc"
)

func QueryAccountAcl(node, bcname, contractAccount string) (*pb.Acl, error) {
	aclStatus := &pb.AclStatus{
		Bcname:      bcname,
		AccountName: contractAccount,
	}
	return queryAcl(node, aclStatus)
}

func QueryMethodAcl(node, bcname, contractName, methodName string) (*pb.Acl, error) {
	aclStatus := &pb.AclStatus{
		Bcname:       bcname,
		ContractName: contractName,
		MethodName:   methodName,
	}
	return queryAcl(node, aclStatus)
}

func queryAcl(node string, aclStatus *pb.AclStatus) (*pb.Acl, error) {
	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	c := pb.NewXchainClient(conn)
	reply, err := c.QueryACL(ctx, aclStatus)
	if err != nil {
		return nil, fmt.Errorf("query acl fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, fmt.Errorf("query acl fail, err: %s", reply.Header.Error.String())
	}
	if reply.Acl == nil {
		return nil, fmt.Errorf("acl not found")
	}
	return reply.Acl, nil
}
