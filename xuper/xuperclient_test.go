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

func TestNewXClient(t *testing.T) {
	x, err := New("127.0.0.1:9999")
	if err != nil {
		t.Error("New xuperclient asser failed:", err)
	}
	x.Close()

	_, err = New("127.0.0.1:9999", WithConfigFile("./conf/sdk.yaml"))
	if err == nil {
		t.Error("New xuperclient asser failed:", err)
	}

	x, err = New("127.0.0.1:9999", WithGrpcGZIP())
	if err != nil {
		t.Error("New xuperclient asser failed:", err)
	}
	x.Close()
	_, err = New("127.0.0.1:9999", WithGrpcGZIP(), WithGrpcTLS("aaa", "aaa", "aaa", "aaa"))
	if err == nil {
		t.Error("New xuperclient asser failed:", err)
	}
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

	_, err = client.RegisterBlockEvent(filter, false)
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

	bal, err := client.QueryBalance(address)
	if err != nil {
		t.Fatalf("err:%s\n", err.Error())
	}

	fmt.Printf("%v\n", bal)
}

func TestXClient_QueryBalanceDetail(t *testing.T) {
	address := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	client := newClient()

	addrBalanceStatus, err := client.QueryBalanceDetail(address)
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

func TestDeployContract(t *testing.T) {

	type Case struct {
		from                    *account.Account
		name                    string
		runtime                 string
		abi                     []byte
		code                    []byte
		args                    map[string]string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: GoRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "部署 go native 合约。",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: JavaRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts:                    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			desc:                    "部署 java native 合约，account 支付 fee，fee=10。",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: EvmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "部署 evm 合约，需要背书，不需要背书手续费。",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: CRuntime,
			opts:    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "部署 wasm 合约，背书手续费100",
		},
		{
			from:    aks[1],
			name:    "hello",
			code:    []byte("code"),
			runtime: CRuntime,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "部署 wasm 合约，背书手续费100，account 支付手续费，fee=10。",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock 不检查余额。
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			switch c.runtime {
			case JavaRuntime:
				_, err := xc.DeployNativeJavaContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.UpgradeNativeContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case GoRuntime:
				_, err := xc.DeployNativeGoContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

			case CRuntime:
				_, err := xc.DeployWasmContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.UpgradeWasmContract(c.from, c.name, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case EvmContractModule:
				_, err := xc.DeployEVMContract(c.from, c.name, c.abi, c.code, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			default:
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestInvokeContract(t *testing.T) {
	type Case struct {
		from                    *account.Account
		name                    string
		module                  string
		method                  string
		args                    map[string]string
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: NativeContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "调用 native 合约。",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts:                    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			desc:                    "调用 wasm 合约，account 支付 fee，fee=10。",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "调用 wasm 合约，需要背书，不需要背书手续费。",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: WasmContractModule,
			opts: []RequestOption{WithContractInvokeAmount("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "调用 wasm 合约，背书手续费100",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			module: EvmContractModule,
			opts:   []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "调用 wasm 合约，背书手续费100，account 支付手续费，fee=10。",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock 不检查余额。
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			switch c.module {
			case NativeContractModule:
				_, err := xc.InvokeNativeContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}

				_, err = xc.QueryNativeContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case WasmContractModule:
				_, err := xc.InvokeWasmContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
				_, err = xc.QueryWasmContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			case EvmContractModule:
				_, err := xc.InvokeEVMContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
				_, err = xc.QueryEVMContract(c.from, c.name, c.method, c.args, c.opts...)
				if err != nil {
					t.Error(err)
				}
			default:
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestACLSet(t *testing.T) {
	type Case struct {
		from                    *account.Account
		name                    string
		method                  string
		acl                     *ACL
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
	}

	cases := []Case{
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "设置合约方法 ACL。",
		},
		{
			from: aks[1],
			acl:  getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "设置合约账户 ACL。",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc:          true,
			onlyFeeFromContractAccc: true,
			opts: []RequestOption{WithFeeFromAccount(), WithFee("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			desc: "设置方法 ACL，account 支付 fee，fee=10。",
		},
		{
			from: aks[1],
			acl:  getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: true,
			desc:           "设置合约账户 ACL，需要背书，不需要背书手续费。",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			opts:   []RequestOption{WithContractInvokeAmount("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "调用 wasm 合约，背书手续费100",
		},
		{
			from:   aks[1],
			name:   "hello",
			method: "a",
			acl:    getDefaultACL("a"),
			opts:   []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasContractAcc: true,
			desc:           "调用 wasm 合约，背书手续费100，account 支付手续费，fee=10。",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock 不检查余额。
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			if c.name != "" {
				_, err := xc.SetMethodACL(c.from, c.name, c.method, c.acl, c.opts...)
				if err != nil {
					t.Error(err)
				}
			} else {
				_, err := xc.SetAccountACL(c.from, c.acl, c.opts...)
				if err != nil {
					t.Error(err)
				}
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestCreateAccount(t *testing.T) {
	type Case struct {
		from                    *account.Account
		account                 string
		acl                     *ACL
		fee                     *big.Int
		opts                    []RequestOption
		cfg                     *config.CommConfig
		hasContractAcc          bool
		onlyFeeFromContractAccc bool
		desc                    string
		hasError                bool
	}

	cases := []Case{
		{
			from:    aks[1],
			account: "",
			acl:     getDefaultACL("a"),
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasError: true,
			desc:     "创建合约账户，账户参数为空。",
		},
		{
			from:     aks[1],
			account:  "XC1234567812345678@xuper",
			acl:      getDefaultACL("a"),
			hasError: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			hasContractAcc: true,
			desc:           "创建合约账户，from 设置了合约账户。",
		},
		{
			from:     aks[1],
			account:  "hello",
			acl:      getDefaultACL("a"),
			hasError: true,
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{},
			},
			onlyFeeFromContractAccc: true,
			opts: []RequestOption{WithFeeFromAccount(), WithFee("10"), WithBcname("xuper"),
				WithDesc("haha"), WithOtherAuthRequires([]string{"a", "b"})},
			desc: "创建合约账户，账户不符合规则。",
		},
		{
			from:    aks[1],
			acl:     getDefaultACL("a"),
			account: "XC1234567812345678@xuper",
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			hasContractAcc: false,
			desc:           "创建合约账户 ok。",
		},
		{
			from:    aks[1],
			account: "XC1234567812345678@xuper",
			acl:     getDefaultACL("a"),
			opts:    []RequestOption{WithContractInvokeAmount("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			desc: "创建合约账户，背书手续费100",
		},
		{
			from:    aks[1],
			account: "XC1234567812345678@xuper",
			acl:     getDefaultACL("a"),
			opts:    []RequestOption{WithFeeFromAccount(), WithFee("10")},
			cfg: &config.CommConfig{
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck:                true,
					IsNeedComplianceCheckFee:             true,
					ComplianceCheckEndorseServiceFee:     100,
					ComplianceCheckEndorseServiceFeeAddr: "aaa",
					ComplianceCheckEndorseServiceAddr:    "bbb",
				},
			},
			hasError: true,
			desc:     "创建合约账户，背书手续费100，account 支付手续费，fee=10。",
		},
	}

	for _, c := range cases {
		if testNode != "" {
			panic("Please implement me")
		} else {
			// mock 不检查余额。
			xc := newClient()
			xc.cfg = c.cfg
			if c.hasContractAcc {
				c.from.SetContractAccount(contractAcc[1])
			}

			tx, err := xc.CreateContractAccount(c.from, c.account, c.opts...)
			if c.hasError {
				if err == nil {
					t.Error("Create contract assert err filed", c.desc)
				}
			} else {
				if err != nil {
					t.Error("Create contract assert err filed", err, c.desc)
				}

				if len(tx.Tx.GetTxid()) == 0 {
					t.Error("Create contract assert tx filed", c.desc)
				}
			}
			c.from.RemoveContractAccount()
			xc.cfg = nil
		}
	}
}

func TestRequest(t *testing.T) {
	r := new(Request)
	r.SetArgs(map[string][]byte{"a": []byte("a")})
	r.SetContractName("counter")
	acc, _ := account.CreateAccount(1, 1)
	r.SetInitiatorAccount(acc)
	r.SetModule("xx")
	r.SetTransferAmount("10")
	r.SetTransferTo("bob")

	if r.contractName != "counter" {
		t.Error("Request set assert failed")
	}

	if r.module != "xx" {
		t.Error("Request set assert failed")
	}
}

// func TestMy(t *testing.T) {
// 	contractAcc := "XC9999999988888888@xuper"
// 	me, _ := account.RetrieveAccount(mnemonic, 1)
// 	me.SetContractAccount(contractAcc)
// 	// you := "nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh"
// 	// yourMnemonic := "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
// 	// you, _ := account.RetrieveAccount(yourMnemonic, 1)
// 	node := "10.12.199.82:8701"
// 	xclient, err := New(node /**, WithGrpcGZIP()**/)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// tx, err := xclient.CreateContractAccount(me, contractAcc)
// 	// if err != nil {
// 	// 	fmt.Println(String(tx.Tx))
// 	// 	panic(err)
// 	// }
// 	// fmt.Println(String(tx.Tx))

// 	// 转账
// 	fmt.Println("转账前：")
// 	fmt.Println(xclient.queryBalance(me.GetContractAccount()))
// 	fmt.Println(xclient.queryBalance(me.Address))
// 	// fmt.Println(xclient.queryBalance(you.Address))

// 	// tx, err := xclient.Transfer(me, you.Address, "100" /**WithFee("100"),**/, WithFeeFromAccount())
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	// fmt.Println("fee:", tx.Fee)

// 	// fmt.Println("转账后：")
// 	// fmt.Println(xclient.queryBalance(me.GetContractAccount()))
// 	// fmt.Println(xclient.queryBalance(me.Address))
// 	// fmt.Println(xclient.queryBalance(you.Address))

// 	// 调用合约

// 	tx, err := xclient.InvokeWasmContract(me, "counter1", "increase", map[string]string{"key": "a"}, WithFee("100"), WithFeeFromAccount(), WithContractInvokeAmount("100"))
// 	// if err != nil {
// 	// 	fmt.Println(String(tx.Tx))
// 	// 	panic(err)
// 	// }
// 	// // fmt.Println(String(tx.Tx))
// 	// fmt.Println("contract resp:", tx.ContractResponse)
// 	// fmt.Println("fee:", tx.Fee)
// 	// fmt.Println("gas used:", tx.GasUsed)

// 	// fmt.Println("调用后：")
// 	// fmt.Println(xclient.queryBalance(me.GetContractAccount()))
// 	// fmt.Println(xclient.queryBalance(me.Address))

// 	// abi := []byte(`[{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":true,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"checkHash","outputs":[{"internalType":"uint256","name":"code","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"getEvidence","outputs":[{"internalType":"uint256","name":"code","type":"uint256"},{"internalType":"uint256","name":"createTime","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getUsers","outputs":[{"internalType":"address[]","name":"users","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"save","outputs":[{"internalType":"uint256","name":"code","type":"uint256"},{"internalType":"uint256","name":"createTime","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`)
// 	// bin := []byte(`60806040526000805560018055600280556003805534801561002057600080fd5b506109a1806100306000396000f3fe608060405234801561001057600080fd5b506004361061004b5760003560e01c8062ce8e3e1461005057806338e48f06146100af578063b16c6ee714610185578063e670f7cd1461025b575b600080fd5b61005861032a565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561009b578082015181840152602081019050610080565b505050509050019250505060405180910390f35b610168600480360360208110156100c557600080fd5b81019080803590602001906401000000008111156100e257600080fd5b8201836020820111156100f457600080fd5b8035906020019184600183028401116401000000008311171561011657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506103b8565b604051808381526020018281526020019250505060405180910390f35b61023e6004803603602081101561019b57600080fd5b81019080803590602001906401000000008111156101b857600080fd5b8201836020820111156101ca57600080fd5b803590602001918460018302840111640100000000831117156101ec57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610699565b604051808381526020018281526020019250505060405180910390f35b6103146004803603602081101561027157600080fd5b810190808035906020019064010000000081111561028e57600080fd5b8201836020820111156102a057600080fd5b803590602001918460018302840111640100000000831117156102c257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610778565b6040518082815260200191505060405180910390f35b606060058054806020026020016040519081016040528092919081815260200182805480156103ae57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610364575b5050505050905090565b6000806000600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050600081856040518082805190602001908083835b602083106104355780518252602082019150602081019050602083039250610412565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090506000816001015414156104de5760053390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b848160000190805190602001906104f6929190610840565b50428160010181905550338160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020866040518082805190602001908083835b602083106105b75780518252602082019150602081019050602083039250610594565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020600082018160000190805460018160011615610100020316600290046106119291906108c0565b50600182015481600101556002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506000548160010154935093505050915091565b6000806000600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020846040518082805190602001908083835b6020831061071157805182526020820191506020810190506020830392506106ee565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090506000816001015414156107655760035460008090509250925050610773565b600054816001015492509250505b915091565b600080600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020836040518082805190602001908083835b602083106107ee57805182526020820191506020810190506020830392506107cb565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060010154141561083557600354905061083b565b60015490505b919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061088157805160ff19168380011785556108af565b828001600101855582156108af579182015b828111156108ae578251825591602001919060010190610893565b5b5090506108bc9190610947565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108f95780548555610936565b8280016001018555821561093657600052602060002091601f016020900482015b8281111561093557825482559160010191906001019061091a565b5b5090506109439190610947565b5090565b61096991905b8082111561096557600081600090555060010161094d565b5090565b9056fea265627a7a72315820df60b56cc52d0f2db57231c429413ce141b5f81aa7a355a6de0fa13f9b7a958964736f6c63430005110032`)

// 	// tx, err := xclient.DeployEVMContract(me, "evidence", abi, bin, nil)
// 	// tx, err := xclient.InvokeEVMContract(me, "evidence", "save", map[string]string{"fileHashHex": "112233"}, WithFee("100"), WithFeeFromAccount())
// 	if err != nil {
// 		// fmt.Println(String(tx.Tx))
// 		panic(err)
// 	}
// 	// fmt.Println(String(tx.Tx))
// 	fmt.Println("contract resp:", tx.ContractResponse)
// 	fmt.Println("fee:", tx.Fee)
// 	fmt.Println("gas used:", tx.GasUsed)

// 	fmt.Println("调用后：")
// 	fmt.Println(xclient.queryBalance(me.GetContractAccount()))
// 	fmt.Println(xclient.queryBalance(me.Address))
// }
