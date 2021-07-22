package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"strconv"

	"gopkg.in/yaml.v2"
)

// ComplianceCheckConfig endorser config.
type ComplianceCheckConfig struct {
	IsNeedComplianceCheck                bool   `yaml:"isNeedComplianceCheck,omitempty"`
	IsNeedComplianceCheckFee             bool   `yaml:"isNeedComplianceCheckFee,omitempty"`
	ComplianceCheckEndorseServiceFee     int    `yaml:"complianceCheckEndorseServiceFee,omitempty"`
	ComplianceCheckEndorseServiceFeeAddr string `yaml:"complianceCheckEndorseServiceFeeAddr,omitempty"`
	ComplianceCheckEndorseServiceAddr    string `yaml:"complianceCheckEndorseServiceAddr,omitempty"`
}

// CommConfig sdk config.
type CommConfig struct {
	EndorseServiceHost string                `yaml:"endorseServiceHost,omitempty"`
	ComplianceCheck    ComplianceCheckConfig `yaml:"complianceCheck,omitempty"`
	MinNewChainAmount  string                `yaml:"minNewChainAmount,omitempty"`
	Crypto             string                `yaml:"crypto,omitempty"`
	TxVersion          int32                 `yaml:"txVersion,omitempty"`
}

const confPath = "./conf"
const confName = "sdk.yaml"

const CRYPTO_XCHAIN = "xchain"
const CRYPTO_GM = "gm"

var config *CommConfig

// GetInstance get config instance.
func GetInstance() *CommConfig {
	if config == nil {
		var err error
		config, err = GetConfig(filepath.Join(confPath, confName))
		if err != nil {
			log.Printf("no config file in ./conf/sdk.yaml, use default config: %v\n", config)
		}
	}
	return config
}

// SetGMCrypto 使用国密，用这个方法可以不使用配置文件来修改了。
func (c *CommConfig) SetGMCrypto() {
	c.Crypto = CRYPTO_GM
}

// SetXchainCrypto 使用 xchain 加密算法，用这个方法可以不使用配置文件来修改了。
func (c *CommConfig) SetXchainCrypto() {
	c.Crypto = CRYPTO_XCHAIN
}

// GetConfig load config from confFile and new config instance.
func GetConfig(confFile string) (*CommConfig, error) {
	// default config
	commConfig := &CommConfig{
		EndorseServiceHost: "10.144.94.18:8848",
		ComplianceCheck: ComplianceCheckConfig{
			ComplianceCheckEndorseServiceFee:     10,
			ComplianceCheckEndorseServiceFeeAddr: "XBbhR82cB6PvaLJs3D4uB9f12bhmKkHeX",
			ComplianceCheckEndorseServiceAddr:    "TYyA3y8wdFZyzExtcbRNVd7ZZ2XXcfjdw",
		},
		MinNewChainAmount: "100",
		Crypto:            CRYPTO_XCHAIN,
	}

	yamlFile, err := ioutil.ReadFile(confFile)
	if err != nil {
		return commConfig, err
	}

	err = yaml.Unmarshal(yamlFile, commConfig)
	if err != nil {
		return nil, err
	}

	config = commConfig
	return commConfig, nil
}

// SetConfig set config fileds.
func SetConfig(checkHost, checkAddr, checkFeeAddr, checkFee string, isNeedCheck, isNeedCheckFee bool, minNewChainAmount string) {
	commConfig := &CommConfig{
		EndorseServiceHost: "10.144.94.18:8848",
		ComplianceCheck: ComplianceCheckConfig{
			ComplianceCheckEndorseServiceFee:     10,
			ComplianceCheckEndorseServiceFeeAddr: "XBbhR82cB6PvaLJs3D4uB9f12bhmKkHeX",
			ComplianceCheckEndorseServiceAddr:    "TYyA3y8wdFZyzExtcbRNVd7ZZ2XXcfjdw",
		},
		MinNewChainAmount: "100",
		Crypto:            CRYPTO_XCHAIN,
	}
	if checkHost != "" {
		commConfig.EndorseServiceHost = checkHost
	}
	if checkFeeAddr != "" {
		commConfig.ComplianceCheck.ComplianceCheckEndorseServiceFeeAddr = checkFeeAddr
	}
	if checkAddr != "" {
		commConfig.ComplianceCheck.ComplianceCheckEndorseServiceAddr = checkAddr
	}
	if checkFee != "" {
		fee, _ := strconv.Atoi(checkFee)
		commConfig.ComplianceCheck.ComplianceCheckEndorseServiceFee = fee
	}
	if minNewChainAmount != "" {
		commConfig.MinNewChainAmount = minNewChainAmount
	}
	commConfig.ComplianceCheck.IsNeedComplianceCheck = isNeedCheck
	commConfig.ComplianceCheck.IsNeedComplianceCheckFee = isNeedCheckFee

	config = commConfig
}
