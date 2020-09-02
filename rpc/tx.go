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
//	tx := FromPBTx(reply.Tx)
//	output, err := json.MarshalIndent(tx, "", "  ")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(string(output))

func QueryTxById(node, bcname, txid string) (*util.Transaction, error) {

	rawTxid, err := hex.DecodeString(txid)
	if err != nil {
		return nil, fmt.Errorf("txid invalid, err: %s", err)
	}

	txstatus := &pb.TxStatus{
		Bcname: bcname,
		Txid:   rawTxid,
	}

	conn, err := grpc.Dial(node, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		return nil, fmt.Errorf("can not connect to node, err: %s", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	client := pb.NewXchainClient(conn)
	reply, err := client.QueryTx(ctx, txstatus)
	if err != nil {
		return nil, fmt.Errorf("query tx fail, err: %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, fmt.Errorf("query tx fail, err: %s", reply.Header.Error.String())
	}
	if reply.Tx == nil {
		return nil, fmt.Errorf("tx not found")
	}
	return util.FromPBTx(reply.Tx), nil
}
