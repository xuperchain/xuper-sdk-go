package contract

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/xuperchain/xuperchain/core/pb"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperdata/teesdk"
)

// EncryptArgs call the TEE App to encrypt the value of the args
func (c *WasmContract) EncryptArgs(svn uint32,args map[string]string) (string, error) {
        plainJson, err := json.Marshal(args)
        if err != nil {
                return "", err
        }
        data, err := json.Marshal(teesdk.FuncCaller{
                Method: "store",
                Args: string(plainJson),
                Svn: svn,
                Address: c.Account.Address,
        })
        newCipher, err := c.tfc.Submit("xchaintf", string(data))
        return string(newCipher), err
}

// DecryptArgs call the TEE App to decrypt the value of the args
func (c *WasmContract) DecryptArgs(svn uint32,args map[string]string) (string, error) {
        plainJson, err := json.Marshal(args)
        if err != nil {
                return "", err
        }
        data, err := json.Marshal(teesdk.FuncCaller{
                Method: "debug",
                Args: string(plainJson),
                Svn: svn,
                Address: c.Account.Address,
        })
        newPlain, err := c.tfc.Submit("xchaintf", string(data))
        return string(newPlain), err
}


// QueryWasmContractPlain decrypts QueryWasmContract result to get plaintext value
func (c *WasmContract) DecryptResponse (responseCipher *pb.InvokeRPCResponse) (*pb.InvokeRPCResponse, error) {
        // 取出 key和对应的密文
        respArgs := make(map[string]string)
        for _,res := range responseCipher.GetResponse().GetResponse() {
                respArgs["key"] = string(res)
        }

        // 解密密文得到key和对应明文
        commConfig := config.GetInstance()
        decryptArgs, err := c.DecryptArgs(commConfig.TC.Svn, respArgs)
        if err != nil {
                log.Println("DecryptArgs error,", err)
                return nil, err
        }

        // decryptArgs is key:plain
        args := make(map[string]string)
        err = json.Unmarshal([]byte(decryptArgs), &args)
        if err != nil {
                return nil, err
        }

        // 解密后的明文覆盖原有的密文kv，返回的是明文
        plain := args["key"]
        decodeValueByte, err := base64.StdEncoding.DecodeString(plain)
        if err != nil {
                return nil, err
        }
        resp := responseCipher.GetResponse().GetResponse()
        resp[0] = resp[0][:len(decodeValueByte)]
        copy(resp[0], decodeValueByte)
        responseCipher.GetResponse().Response = resp

        return responseCipher, nil
}

func (c *WasmContract) EncryptWasmArgs(args map[string]string) (map[string]string, error) {
	// preExe
	commConfig := config.GetInstance()
	// TODO fix bug
	if commConfig.TC.Enable {
		encryptedArgs, err := c.EncryptArgs(commConfig.TC.Svn, args)
		if err != nil {
			log.Println("EncryptArgs error,", err)
			return nil, err
		}
		args = map[string]string{}
		err = json.Unmarshal([]byte(encryptedArgs), &args)
		if err != nil {
			return nil, err
		}
		return args, nil
	}
	return args, nil
}
