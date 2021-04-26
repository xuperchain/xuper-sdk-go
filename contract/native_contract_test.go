package contract

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
)

func TestNativeContract_Deploy(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc) // 对应的地址为：TRxjxi6aoQsbzjX753wL4HLwsrm5oVC6w
	node := "127.0.0.1:37101"
	bcname := "xuper"
	contractAccount := "XC2222222222222222@xuper"
	contractName := "golangcounter8.3"
	runtime := "go"
	fmt.Printf("account address:【%s】\npriv:【%s】\npublic:【%s】\n", acc.Address, acc.PrivateKey, acc.PublicKey)
	sdkClient, err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	path := "./counter"
	nativeContract := InitNativeContract(acc, bcname, contractName, contractAccount, sdkClient)
	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	nativeContract.Fee = "15587500"
	txId, err := nativeContract.DeployNativeContract(args, path, runtime)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(txId)
	//createContractAccount(acc, node, bcname,sdkClient)
}

func TestNativeInvoke(t *testing.T) {
	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	fmt.Println("accounr address:", acc.Address)
	//node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "golangcounter8.3"
	cAccount := "XC2222222222222222@xuper"
	sdkClient, err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	args := map[string]string{
		"key": "test",
	}
	nativeContract := InitNativeContract(acc, bcname, cName, cAccount, sdkClient)
	mName := "Increase"
	r, e := nativeContract.InvokeNativeContract(mName, args)
	if e != nil {
		t.Error(e)
	}
	fmt.Println("invoke sucess:", r)
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
	contractName := "golangcounter8.3"
	//runtime := "go"
	fmt.Printf("account address:【%s】\npriv:【%s】\npublic:【%s】\n", acc.Address, acc.PrivateKey, acc.PublicKey)
	sdkClient, err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	path := "./counter2"
	nativeContract := InitNativeContract(acc, bcname, contractName, contractAccount, sdkClient)
	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	nativeContract.Fee = "15587500"
	txId, err := nativeContract.UpgradeNativeContract(args, path)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(txId)
	//createContractAccount(acc, node, bcname,sdkClient)
}
