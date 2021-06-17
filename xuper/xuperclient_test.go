package xuper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuperchain/core/pb"
)

func TestTransfer(t *testing.T) {
	xc, err := New("10.12.199.82:8701")
	if err != nil {
		panic(err)
	}

	memonic := "叫 金 管 逻 亮 栏 贺 扣 导 姚 杂 光"
	// addr := "nGjHu44Dr3ihairvjwkUh3fgTtMxr2LMR"
	me, err := account.RetrieveAccount(memonic, 1)
	tx, err := xc.Transfer(me, "XC1111222233334444@xuper", "10000000", WithFee("10"))

	if err != nil {
		fmt.Printf("%+v \n", tx.Tx)
		fmt.Println(String(tx.Tx))
		fmt.Println(len(tx.Tx.InitiatorSigns))
		fmt.Println(len(tx.Tx.AuthRequire))
		fmt.Println(tx.Tx.AuthRequire)
		fmt.Println(len(tx.Tx.AuthRequireSigns))
		fmt.Println("NNNNNN")
		panic(err)
	}
	fmt.Println(tx)
}

func TestDeployCppWasm(t *testing.T) {

	xc, err := New("10.12.199.82:8701")
	if err != nil {
		panic(err)
	}

	memonic := "叫 金 管 逻 亮 栏 贺 扣 导 姚 杂 光"
	// addr := "nGjHu44Dr3ihairvjwkUh3fgTtMxr2LMR"

	me, err := account.RetrieveAccount(memonic, 1)

	// tx1, err := xc.CreateContractAccount(me, "XC1111222233334444@xuper")
	// if err != nil {
	// 	fmt.Println(err)
	// 	fmt.Printf("%+v \n", tx1.Tx)
	// 	panic(err)
	// }

	me.SetContractAccount("XC1111222233334444@xuper")
	codePath := "~/go/src/github.com/KenianShi/xuper-sdk-go/example/contract_code/counter.wasm"
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		panic(err)
	}
	args := map[string]string{
		"creator": "me",
	}
	tx, err := xc.DeployWasmContract(me, "counter1", code, args)
	if err != nil {
		fmt.Println(err)
		fmt.Printf("%+v \n", tx.Tx)
		fmt.Println(String(tx.Tx))
		fmt.Println(len(tx.Tx.InitiatorSigns))
		fmt.Println(len(tx.Tx.AuthRequire))
		fmt.Println(tx.Tx.AuthRequire)
		fmt.Println(len(tx.Tx.AuthRequireSigns))
		fmt.Println("NNNNNN")
		panic(err)
	}
	fmt.Println(tx)
}
func TestInvokeCpp(t *testing.T) {
	xc, err := New("10.12.199.82:8701")
	if err != nil {
		panic(err)
	}

	memonic := "叫 金 管 逻 亮 栏 贺 扣 导 姚 杂 光"
	// addr := "nGjHu44Dr3ihairvjwkUh3fgTtMxr2LMR"
	me, err := account.RetrieveAccount(memonic, 1)
	args := map[string]string{
		"key": "a",
	}
	tx, err := xc.InvokeWasmContract(me, "counter1", "increase", args)

	if err != nil {
		fmt.Println(err)
		fmt.Printf("%+v \n", tx.Tx)
		fmt.Println(String(tx.Tx))
		fmt.Println(len(tx.Tx.InitiatorSigns))
		fmt.Println(len(tx.Tx.AuthRequire))
		fmt.Println(tx.Tx.AuthRequire)
		fmt.Println(len(tx.Tx.AuthRequireSigns))
		fmt.Println("NNNNNN")
		panic(err)
	}
	fmt.Println(tx)
	fmt.Println(tx.ContractResponse)
}

func String(t *pb.Transaction) string {
	b, err := json.Marshal(*t)
	if err != nil {
		return fmt.Sprintf("%+v", *t)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *t)
	}
	return out.String()
}

func TestNacc(t *testing.T) {
	memonic := "叫 金 管 逻 亮 栏 贺 扣 导 姚 杂 光"
	// addr := "nGjHu44Dr3ihairvjwkUh3fgTtMxr2LMR"
	me, err := account.RetrieveAccount(memonic, 1)
	if err != nil {
		panic(err)
	}
	// me, _ := account.CreateAccount(1, 1)
	fmt.Println(me)

	contractAccount := "XC1111111111111111@1"
	ok, _ := regexp.MatchString(`^XC\d{16}@*`, contractAccount)
	fmt.Println(ok)

}

func newClient() *XClient {
	client, err := New("127.0.0.1:37101")
	if err != nil {
		fmt.Printf("newClient err:%s\n", err.Error())
		panic(err)
	}
	return client
}

func TestXClient_RegisterBlockEvent(t *testing.T) {
	client := newClient()
	watcher := InitWatcher(client, 10, false)

	filter, err := NewBlockFilter("xuper")
	if err != nil {
		t.Fatalf("create block filter err: %v\n", err)
	}

	reg, err := watcher.RegisterBlockEvent(filter, watcher.SkipEmptyTx)
	if err != nil {
		t.Error("RegisterBlockEvent")
		t.Error(err)
	}

	go func() {
		for {
			b := <-reg.FilteredBlockChan
			fmt.Printf("%+v\n", b)
		}
	}()

	time.Sleep(time.Second * 10)
	reg.Unregister()
}

func TestXClient_QueryTxByID(t *testing.T) {
	txID := "b1ae1868c4e46651657b5aa9be20ab284d36161c9cc311787a9e81e391dc2bed"
	client := newClient()

	tx, err := client.QueryTxByID(txID)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("tx:%v\n", tx)
}

func TestXClient_QueryBlockByHeight(t *testing.T) {
	client := newClient()
	height := int64(213822)

	block, err := client.QueryBlockByHeight(height)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("block:%v\n", block)
}

func TestXClient_QueryBlockByID(t *testing.T) {
	client := newClient()
	blockID := "b4213754f1f645a9bbcf5cfe13de65d019dd843b76929782d78045fff3cceace"

	block, err := client.QueryBlockByID(blockID)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("block:%+v\n", block)
}

func TestXClient_QueryAccountAcl(t *testing.T) { //todo
	client := newClient()
	account := "XC1111111111111111@xuper"

	acl, err := client.QueryAccountACL(account)
	if err != nil {
		t.Fatalf("err:%v\n", err.Error())
	}

	fmt.Printf("acl:%v\n", acl)
}

func TestXClient_QueryMethodAcl(t *testing.T) { // todo
	name := "test0611"
	method := "testUint256Event"
	client := newClient()

	acl, err := client.QueryMethodACL(name, method)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", acl)
}

func TestXClient_QueryAccountContracts(t *testing.T) {
	account := "XC1111111111111111@xuper"
	client := newClient()

	contracts, err := client.QueryAccountContracts(account)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	for _, contract := range contracts {
		fmt.Printf("%v\n", contract)
	}
}

func TestXClient_QueryAddressContracts(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	contracts, err := client.QueryAddressContracts(address)
	if err != nil {
		t.Fatalf("err:%s\n", err)
	}

	for _, contract := range contracts {
		fmt.Printf("%v\n", contract)
	}
}

func TestXClient_QueryBalance(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	addrStatus, err := client.queryBalance(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", addrStatus)
}

func TestXClient_QueryBalanceDetail(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	addrBalanceStatus, err := client.queryBalanceDetail(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", addrBalanceStatus)
}

func TestXClient_QuerySystemStatus(t *testing.T) {
	client := newClient()

	reply, err := client.QuerySystemStatus()
	if err != nil {
		t.Fatalf("err:%s\n", err)
	}

	fmt.Printf("reply:%v\n", reply)
}

func TestXClient_QueryBlockChains(t *testing.T) {
	client := newClient()

	bcStatus, err := client.QueryBlockChains()
	if err != nil {
		t.Fatalf("err:%v", err)
	}

	fmt.Printf("bcStatus:%v\n", bcStatus)
}

func TestXClient_QueryBlockChainStatus(t *testing.T) {
	bcName := "xuper"
	client := newClient()

	bcStatus, err := client.QueryBlockChainStatus(bcName)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("bcStatus:%v\n", bcStatus)
}

func TestXClient_QueryNetURL(t *testing.T) {
	client := newClient()

	rawURL, err := client.QueryNetURL()
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("rawURL:%v\n", rawURL)
}

func TestXClient_QueryAccountByAK(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	resp, err := client.QueryAccountByAK(address)
	if err != nil {
		t.Fatalf("err:%v\n", err)
	}

	fmt.Printf("resp:%v\n", resp)
}
