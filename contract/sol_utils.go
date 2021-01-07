package contract

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/hyperledger/burrow/execution/evm/abi"
)

func convertToXuper3Args(args map[string]interface{}) (map[string][]byte, error) {
	argmap := make(map[string][]byte)
	for k, v := range args {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("bad key %s, expect string value, got %v", k, v)
		}
		argmap[k] = []byte(s)
	}

	return argmap, nil
}

func convertToEvmArgsWithAbiFile(abiFile string, method string, args map[string]interface{}) (map[string][]byte, []byte, error) {
	buf, err := ioutil.ReadFile(abiFile)
	if err != nil {
		return nil, nil, err
	}

	return convertToEvmArgsWithAbiData(buf, method, args)
}

func convertToEvmArgsWithAbiData(abiData []byte, method string, args map[string]interface{}) (map[string][]byte, []byte, error) {
	enc, err := New(abiData)
	if err != nil {
		return nil, nil, err
	}
	input, err := enc.Encode(method, args)
	if err != nil {
		return nil, nil, err
	}
	ret := map[string][]byte{
		"input": input,
	}

	return ret, abiData, nil
}

// SolContractCallResult contract call return result.
type SolContractCallResult struct {
	Index string
	Value string
}

// ParseRespWithAbiForEVM parse contract preExe response.
func ParseRespWithAbiForEVM(abiData, methodName string, resp []byte) ([]SolContractCallResult, error) {
	Variables, err := abi.DecodeFunctionReturn(abiData, methodName, resp)
	if err != nil {
		return nil, err
	}

	result := make([]SolContractCallResult, 0, len(Variables))
	for i := range Variables {
		sr := SolContractCallResult{}
		sr.Index = Variables[i].Name
		if len(Variables[i].Value) == 32 {
			sr.Value = hex.EncodeToString([]byte(Variables[i].Value))
		} else {
			sr.Value = Variables[i].Value
		}
		result = append(result, sr)
	}

	return result, nil
}
