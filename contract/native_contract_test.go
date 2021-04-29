package contract

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
	"time"
)

// 部署go 版本的 Native
func TestNativeContract_Deploy_Go(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc) // 对应的地址为：TRxjxi6aoQsbzjX753wL4HLwsrm5oVC6w
	node := "127.0.0.1:37101"
	bcname := "xuper"
	contractAccount := "XC2222222222222222@xuper"
	contractName := "golangcounter5"
	runtime := "go"
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	path := "./test/go/counter1"
	nativeContract := InitNativeContractWithClient(acc, bcname, contractName, contractAccount, xuperClient)
	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	//nativeContract.Fee = "15587500" // 此行注释掉，分别测有fee和nofee的场景。NoFee场景需要超级链的支持
	txId, err := nativeContract.DeployNativeContract(args, path, runtime)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(txId)
	time.Sleep(time.Second * 2)
	txStatus, err := nativeContract.QueryTx(txId)
	if err != nil {
		t.Error(err)
	}
	if txStatus != nil {
		fmt.Printf("txStatus:%d\n", txStatus.Status)
	}
}

func TestNativeContract_Deploy_Java(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc) // 对应的地址为：TRxjxi6aoQsbzjX753wL4HLwsrm5oVC6w
	node := "127.0.0.1:37101"
	bcname := "xuper"
	contractAccount := "XC2222222222222222@xuper"
	contractName := "javacounter1"
	runtime := "java"
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	path := "./test/java/counter"
	nativeContract := InitNativeContractWithClient(acc, bcname, contractName, contractAccount, xuperClient)
	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	nativeContract.Fee = "15587500" // 此行注释掉，分别测有fee和nofee的场景。NoFee场景需要超级链的支持
	txId, err := nativeContract.DeployNativeContract(args, path, runtime)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)
	txStatus, err := nativeContract.QueryTx(txId)
	if err != nil {
		t.Error(err)
	}
	if txStatus != nil {
		fmt.Printf("txStatus:%d\n", txStatus.Status)
	}
}

func TestNativeInvoke(t *testing.T) {
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	fmt.Println("accounr address:", acc.Address)
	node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "golangcounter1"
	cAccount := "XC2222222222222222@xuper"
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	args := map[string]string{
		"key": "test",
	}
	nativeContract := InitNativeContractWithClient(acc, bcname, cName, cAccount, xuperClient)
	mName := "Increase"
	txId, e := nativeContract.InvokeNativeContract(mName, args)
	if e != nil {
		t.Error(e)
	}
	fmt.Printf("txid;%s\n", txId)
	time.Sleep(time.Second * 2)

	txStatus, err := nativeContract.QueryTx(txId)
	if err != nil {
		t.Error(err)
	}
	if txStatus != nil {
		fmt.Printf("txStatus:%d\n", txStatus.Status)
	}
}

func TestNativeContract_Upgrade(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)
	node := "127.0.0.1:37101"
	bcname := "xuper"
	contractAccount := "XC2222222222222222@xuper"
	contractName := "golangcounter1"
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	path := "./test/go/counter2"
	nativeContract := InitNativeContractWithClient(acc, bcname, contractName, contractAccount, xuperClient)
	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	//nativeContract.Fee = "15587500"
	txId, err := nativeContract.UpgradeNativeContract(args, path)
	if err != nil {
		t.Error(err)
	}
	time.Sleep(time.Second * 2)
	txStatus, err := nativeContract.QueryTx(txId)
	if err != nil {
		t.Error(err)
	}
	if txStatus != nil {
		fmt.Printf("txStatus:%d\n", txStatus.Status)
	}
}

func TestNativeContract_QueryNativeContract(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)
	node := "127.0.0.1:37101"
	bcname := "xuper"
	contractAccount := "XC2222222222222222@xuper"
	contractName := "golangcounter1"
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	nativeContract := InitNativeContractWithClient(acc, bcname, contractName, contractAccount, xuperClient)
	args := map[string]string{
		"key": "test",
	}
	resp, err := nativeContract.QueryNativeContract("get", args)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("resp:%+v\n", resp.Response)
}
