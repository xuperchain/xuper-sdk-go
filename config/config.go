// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package config is related to common conf info
package config

import (
	"io/ioutil"
	"log"

	"strconv"

	"gopkg.in/yaml.v2"
)

type ComplianceCheckConfig struct {
	IsNeedComplianceCheck                bool   `yaml:"isNeedComplianceCheck,omitempty"`
	IsNeedComplianceCheckFee             bool   `yaml:"isNeedComplianceCheckFee,omitempty"`
	ComplianceCheckEndorseServiceFee     int    `yaml:"complianceCheckEndorseServiceFee,omitempty"`
	ComplianceCheckEndorseServiceFeeAddr string `yaml:"complianceCheckEndorseServiceFeeAddr,omitempty"`
	ComplianceCheckEndorseServiceAddr    string `yaml:"complianceCheckEndorseServiceAddr,omitempty"`
}

type CommConfig struct {
	EndorseServiceHost string                `yaml:"endorseServiceHost,omitempty"`
	ComplianceCheck    ComplianceCheckConfig `yaml:"complianceCheck,omitempty"`
	MinNewChainAmount  string                `yaml:"minNewChainAmount,omitempty"`
	Crypto             string                `yaml:"crypto,omitempty"`
}

const confPath = "./conf"
const confName = "sdk.yaml"

const CRYPTO_XCHAIN = "xchain"
const CRYPTO_GM = "gm"

var config *CommConfig

func GetInstance() *CommConfig {
	if config == nil {
		config = GetConfig(confPath, confName)
	}
	return config
}

func GetConfig(configPath string, confName string) *CommConfig {
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

	filename := configPath + "/" + confName
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Config yamlFile get error #%v", err)
	}

	err = yaml.Unmarshal(yamlFile, commConfig)
	if err != nil {
		log.Fatalf("Config Unmarshal error: %v", err)
	}

	log.Printf("GetConfig: %v\n", commConfig)
	return commConfig
}

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
