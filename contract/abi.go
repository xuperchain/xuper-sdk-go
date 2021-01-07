package contract

import (
	"fmt"

	"github.com/hyperledger/burrow/execution/evm/abi"
)

// ABI abi.Spec wrapper.
type ABI struct {
	spec *abi.Spec
}

// LoadFile load abi file.
func LoadFile(fpath string) (*ABI, error) {
	spec, err := abi.LoadPath(fpath)
	if err != nil {
		return nil, err
	}
	return newABI(spec), nil
}

// New new ABI.
func New(buf []byte) (*ABI, error) {
	spec, err := abi.ReadSpec(buf)
	if err != nil {
		return nil, err
	}
	return newABI(spec), nil
}

func newABI(spec *abi.Spec) *ABI {
	return &ABI{
		spec: spec,
	}
}

// Encode encode method name and args to []byte.
func (a *ABI) Encode(methodName string, args map[string]interface{}) ([]byte, error) {
	if methodName == "" {
		if a.spec.Constructor != nil {
			return a.encodeMethod(a.spec.Constructor, args)
		}
		return nil, nil
	}
	method, ok := a.spec.Functions[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s not found", methodName)
	}

	return a.encodeMethod(method, args)
}

func (a *ABI) encodeMethod(method *abi.FunctionSpec, args map[string]interface{}) ([]byte, error) {
	var inputs []interface{}
	for _, input := range method.Inputs {
		v, ok := args[input.Name]
		if !ok {
			return nil, fmt.Errorf("arg name %s not found", input.Name)
		}
		inputs = append(inputs, v)
	}
	out, _, err := a.spec.Pack(method.Name, inputs...)

	return out, err
}
