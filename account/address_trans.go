package account

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/btcsuite/btcutil/base58"
	"github.com/hyperledger/burrow/binary"
	"github.com/hyperledger/burrow/crypto"
)

const (
	evmAddressFiller = "-"

	contractNamePrefixs    = "1111"
	contractAccountPrefixs = "1112"

	accountPrefix       = "XC"
	accountBcnameSep    = "@"
	accountSize         = 16
	contractNameMaxSize = 16
	contractNameMinSize = 4

	// XchainAddrType xchain AK 地址类型
	XchainAddrType = "xchain"

	// ContractNameType 合约名字地址类型
	ContractNameType = "contract-name"

	// ContractAccountType 合约账户地址类型
	ContractAccountType = "contract-account"
)

var (
	contractNameRegex = regexp.MustCompile("^[a-zA-Z_]{1}[0-9a-zA-Z_.]+[0-9a-zA-Z_]$")
)

// XchainToEVMAddress xchain address transfer to evm address: xchainAddr can be xchain contract account, AK address, xchain contract name.
//
// Return: evm address, address type, error.
func XchainToEVMAddress(xchainAddr string) (string, string, error) {
	var addr crypto.Address
	var addrType string
	var err error
	if determineContractAccount(xchainAddr) {
		addr, err = contractAccountToEVMAddress(xchainAddr)
		addrType = ContractAccountType
	} else if determineContractName(xchainAddr) == nil {
		addr, err = contractNameToEVMAddress(xchainAddr)
		addrType = ContractNameType
	} else {
		addr, err = xchainAKToEVMAddress(xchainAddr)
		addrType = XchainAddrType
	}

	if err != nil {
		return "", "", err
	}

	return addr.String(), addrType, nil
}

// EVMToXchainAddress evm address transfer to xchain address: evmAddr can be evm contract account, AK address, xchain contract name.
//
// Return: xchain address, address type, error.
func EVMToXchainAddress(evmAddr string) (string, string, error) {
	eAddr, err := crypto.AddressFromHexString(evmAddr)
	if err != nil {
		return "", "", err
	}

	evmAddrWithPrefix := eAddr.Bytes()
	evmAddrStrWithPrefix := string(evmAddrWithPrefix)

	var addr, addrType string
	if evmAddrStrWithPrefix[0:4] == contractAccountPrefixs {
		addr, err = evmAddressToContractAccount(eAddr)
		addrType = ContractAccountType
	} else if evmAddrStrWithPrefix[0:4] == contractNamePrefixs {
		addr, err = evmAddressToContractName(eAddr)
		addrType = ContractNameType
	} else {
		addr, err = evmAddressToXchain(eAddr)
		addrType = XchainAddrType
	}

	if err != nil {
		return "", "", err
	}

	return addr, addrType, nil
}

// transfer xchain address to evm address
func xchainAKToEVMAddress(addr string) (crypto.Address, error) {
	rawAddr := base58.Decode(addr)
	if len(rawAddr) < 21 {
		return crypto.ZeroAddress, errors.New("bad address")
	}
	ripemd160Hash := rawAddr[1:21]

	return crypto.AddressFromBytes(ripemd160Hash)
}

// transfer evm address to xchain address
func evmAddressToXchain(evmAddress crypto.Address) (string, error) {
	addrType := 1
	nVersion := uint8(addrType)
	bufVersion := []byte{byte(nVersion)}

	outputRipemd160 := evmAddress.Bytes()

	strSlice := make([]byte, len(bufVersion)+len(outputRipemd160))
	copy(strSlice, bufVersion)
	copy(strSlice[len(bufVersion):], outputRipemd160)

	checkCode := DoubleSha256(strSlice)
	simpleCheckCode := checkCode[:4]
	slice := make([]byte, len(strSlice)+len(simpleCheckCode))
	copy(slice, strSlice)
	copy(slice[len(strSlice):], simpleCheckCode)

	return base58.Encode(slice), nil
}

// transfer contract name to evm address
func contractNameToEVMAddress(contractName string) (crypto.Address, error) {
	contractNameLength := len(contractName)
	var prefixStr string
	for i := 0; i < binary.Word160Length-contractNameLength-4; i++ {
		prefixStr += evmAddressFiller
	}
	contractName = prefixStr + contractName
	contractName = contractNamePrefixs + contractName

	return crypto.AddressFromBytes([]byte(contractName))
}

// transfer evm address to contract name
func evmAddressToContractName(evmAddr crypto.Address) (string, error) {
	contractNameWithPrefix := evmAddr.Bytes()
	contractNameStrWithPrefix := string(contractNameWithPrefix)
	prefixIndex := strings.LastIndex(contractNameStrWithPrefix, evmAddressFiller)

	return contractNameStrWithPrefix[prefixIndex+1:], nil
}

// transfer contract account to evm address
func contractAccountToEVMAddress(contractAccount string) (crypto.Address, error) {
	contractAccountValid := contractAccount[2:18]
	contractAccountValid = contractAccountPrefixs + contractAccountValid

	return crypto.AddressFromBytes([]byte(contractAccountValid))
}

// transfer evm address to contract account
func evmAddressToContractAccount(evmAddr crypto.Address) (string, error) {
	contractNameWithPrefix := evmAddr.Bytes()
	contractNameStrWithPrefix := string(contractNameWithPrefix)

	return accountPrefix + contractNameStrWithPrefix[4:] + "@xuper", nil
}

// determine whether it is a contract account
func determineContractAccount(account string) bool {
	if isAccount(account) != 1 {
		return false
	}

	return strings.Index(account, "@xuper") != -1
}

func isAccount(name string) int {
	if name == "" {
		return -1
	}
	if !strings.HasPrefix(name, accountPrefix) {
		return 0
	}
	prefix := strings.Split(name, "@")[0]
	prefix = prefix[len(accountPrefix):]
	if err := validRawAccount(prefix); err != nil {
		return 0
	}

	return 1
}

// ValidRawAccount validate account number
func validRawAccount(accountName string) error {
	// param absence check
	if accountName == "" {
		return fmt.Errorf("invoke NewAccount failed, account name is empty")
	}

	// account naming rule check
	if len(accountName) != accountSize {
		return fmt.Errorf("invoke NewAccount failed, account name length expect %d, actual: %d", accountSize, len(accountName))
	}

	for i := 0; i < accountSize; i++ {
		if accountName[i] >= '0' && accountName[i] <= '9' {
			continue
		} else {
			return fmt.Errorf("invoke NewAccount failed, account name expect continuous %d number", accountSize)
		}
	}

	return nil
}

// determine whether it is a contract name
func determineContractName(contractName string) error {
	return validContractName(contractName)
}

func validContractName(contractName string) error {
	// param absence check
	// contract naming rule check
	contractSize := len(contractName)
	contractMaxSize := contractNameMaxSize
	contractMinSize := contractNameMinSize

	if contractSize > contractMaxSize || contractSize < contractMinSize {
		return fmt.Errorf("contract name length expect [%d~%d], actual: %d", contractMinSize, contractMaxSize, contractSize)
	}

	if !contractNameRegex.MatchString(contractName) {
		return fmt.Errorf("contract name does not fit the rule of contract name")
	}

	return nil
}

// UsingSha256 get the hash result of data using SHA256
func UsingSha256(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	out := h.Sum(nil)

	return out
}

// DoubleSha256 执行2次SHA256，这是为了防止SHA256算法被攻破。
func DoubleSha256(data []byte) []byte {
	return UsingSha256(UsingSha256(data))
}
