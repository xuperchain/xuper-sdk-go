package xuper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuperchain/core/pb"
)

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

func newClient() *XClient {
	if testNode == "" {
		return &XClient{
			xc:  &MockXClient{},
			ec:  &MockEClient{},
			esc: &MockESClient{},
		}
	}
	return xclient
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

	if testNode != "" {
		go func() {
			for {
				b := <-reg.FilteredBlockChan
				fmt.Printf("%+v\n", b)
			}
		}()
		time.Sleep(time.Second * 1)
		go func() {
			reg.Unregister()
		}()
	}
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

	bal, err := client.queryBalance(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", bal)
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

var (
	testNode           string = "" // 如果搭建了 xchain 节点，URL 写在这里就会发送交易到这个节点。//10.12.199.82:8701
	wasmCppCounterFile        = "" // 还需要编译counter wasm合约然后拷贝过来。
)

var (
	// 地址：QZcRXTPwtYo2UtZXNRppEC3Mto2AvmgKZ，有节点的话别忘了给这个账户转钱。
	mnemonic    = "刮 腰 秀 算 图 称 洛 校 新 韩 哈 门"
	aks         = []*account.Account{}
	contractAcc = []string{
		"XC0000000000000000@xuper",
		"XC1111111111111111@xuper",
		"XC2222222222222222@xuper",
		"XC3333333333333333@xuper",
		"XC4444444444444444@xuper",
		"XC5555555555555555@xuper",
		"XC6666666666666666@xuper",
		"XC7777777777777777@xuper",
		"XC8888888888888888@xuper",
		"XC9999999999999999@xuper",
	}

	xclient *XClient = nil
)

// 创建10个AK，每个AK创建一个合约账户从0-9.
// 如果启动了节点，需要先手动给 QZcRXTPwtYo2UtZXNRppEC3Mto2AvmgKZ 这个地址转账大约100000000000。
// 助记词：刮 腰 秀 算 图 称 洛 校 新 韩 哈 门
func init() {
	for i := 0; i < len(contractAcc); i++ {
		acc, _ := account.CreateAccount(1, 1)
		aks = append(aks, acc)
		if testNode != "" {
			var err error
			richAcc, _ := account.RetrieveAccount(mnemonic, 1)
			if xclient == nil {
				xclient, err = New(testNode)
				if err != nil {
					panic(err)
				}
			}

			// 转账
			_, err = xclient.Transfer(richAcc, acc.Address, "1000000000")
			if err != nil {
				panic(err)
			}
			_, err = xclient.CreateContractAccount(acc, contractAcc[i])
			if err != nil {
				panic(err)
			}
			// 转账给合约账户
			_, err = xclient.Transfer(richAcc, contractAcc[i], "1000000000")
			if err != nil {
				panic(err)
			}
		}
	}
	if testNode != "" {
		fmt.Println("Test real node client start...")
	} else {
		fmt.Println("Test mock client start...")
	}
}

func TestTransfers(t *testing.T) {
	type Case struct {
		from                    *account.Account
		to                      string
		amount                  string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[0],
			to:     aks[1].Address,
			amount: "100",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			desc: "正常转账100，没有 fee。",
		},
		{
			from:   aks[0],
			to:     aks[1].Address,
			amount: "100",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:  big.NewInt(100),
			desc: "正常转账 fee 100",
			opts: []RequestOption{WithFee("100")},
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:  big.NewInt(100),
			desc: "有合约账户，转账 fee 100",
			opts: []RequestOption{WithFee("100")},
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			fee:                     big.NewInt(100),
			desc:                    "有合约账户，转账 fee 100",
			opts:                    []RequestOption{WithFee("100"), WithFeeFromAccount()},
			onlyFeeFromContractAccc: true,
		},
		{
			from:           aks[0],
			to:             aks[1].Address,
			amount:         "100",
			hasContractAcc: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			fee:  big.NewInt(100),
			desc: "正常转账 fee 100",
			opts: []RequestOption{WithFee("100")},
		},
	}

	for _, c := range cases {
		if testNode != "" {
			// 转账前先查询两个账户余额
			bal, err := xclient.QueryBalance(c.from.Address)
			if err != nil {
				t.Error(err)
			}
			fmt.Println(bal)

			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[0])
			}
			xclient.Transfer(c.from, c.to, c.amount, c.opts...)
			c.from.RemoveContractAccount()
			time.Sleep(time.Millisecond * 500)
			// 转账后再查询两个账户余额
		} else {
			// mock 不检查余额。
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[0])
			}
			_, err := xc.Transfer(c.from, c.to, c.amount, c.opts...)
			if err != nil {
				t.Error(err)
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func string2Bigint(v string) *big.Int {
	b, ok := big.NewInt(0).SetString(v, 10)
	if !ok {
		panic("string 2 bigint failed:" + v)
	}
	return b
}

func TestDeployNativeGoContract(t *testing.T) {

}
