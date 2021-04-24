package contract

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"io/ioutil"
	"log"
)

func GenerateDeployInvokeReq(arg map[string]string, codepath string, runtime, contractAccount, contractName, contractType string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		ContractType:contractType,
		Runtime: runtime,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(contractAccount),
		"contract_name": []byte(contractName),
		"contract_code": contractCode,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: "Deploy",
		Args:       args,
	}
}



func generateUpgradeInvokReq(arg map[string]string, codepath string, contractAccount, contractName,contractType string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		ContractType:contractType,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(contractAccount),
		"contract_name": []byte(contractName),
		"contract_code": contractCode,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: "Upgrade",
		Args:       args,
	}
}


func generateInvokeInvokeReq(args map[string]string, methodName, contractName string,moduleName string) *pb.InvokeRequest {
	return &pb.InvokeRequest{
		ModuleName:   moduleName,
		MethodName:   methodName,
		ContractName: contractName,
		Args:         convertToXuperContractArgs(args),
	}
}