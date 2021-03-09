// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package config 配置信息
package config

import (
	"io/ioutil"
	"log"

	"strconv"

	"gopkg.in/yaml.v2"
)

// ComplianceCheckConfig 背书检查相关配置。
type ComplianceCheckConfig struct {
	IsNeedComplianceCheck                bool   `yaml:"isNeedComplianceCheck,omitempty"`
	IsNeedComplianceCheckFee             bool   `yaml:"isNeedComplianceCheckFee,omitempty"`
	ComplianceCheckEndorseServiceFee     int    `yaml:"complianceCheckEndorseServiceFee,omitempty"`
	ComplianceCheckEndorseServiceFeeAddr string `yaml:"complianceCheckEndorseServiceFeeAddr,omitempty"`
	ComplianceCheckEndorseServiceAddr    string `yaml:"complianceCheckEndorseServiceAddr,omitempty"`
}

// CommConfig SDK 配置
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

// GetInstance 获取配置实例。
func GetInstance() *CommConfig {
	if config == nil {
		config = GetConfig(confPath, confName)
	}
	return config
}

// GetConfig 根据配置文件加载配置信息，如果没有配置文件默认不需要背书检查。
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

// SetConfig 设置配置信息。
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
