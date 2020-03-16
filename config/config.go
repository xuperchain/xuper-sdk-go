// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package config is related to common conf info
package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type ComplianceCheckConfig struct {
	ComplianceCheckEndorseServiceFee     int    `yaml:"complianceCheckEndorseServiceFee,omitempty"`
	ComplianceCheckEndorseServiceFeeAddr string `yaml:"complianceCheckEndorseServiceFeeAddr,omitempty"`
	ComplianceCheckEndorseServiceAddr    string `yaml:"complianceCheckEndorseServiceAddr,omitempty"`
}

type CommConfig struct {
	EndorseServiceHost string                `yaml:"endorseServiceHost,omitempty"`
	ComplianceCheck    ComplianceCheckConfig `yaml:"complianceCheck,omitempty"`
	Crypto             string                `yaml:"crypto,omitempty"`
}

const confPath = "./conf"
const confName = "sdk.yaml"

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
		EndorseServiceHost: "39.156.69.83:37100",
		ComplianceCheck: ComplianceCheckConfig{
			ComplianceCheckEndorseServiceFee:     100,
			ComplianceCheckEndorseServiceFeeAddr: "aB2hpHnTBDxko3UoP2BpBZRujwhdcAFoT",
			ComplianceCheckEndorseServiceAddr:    "jknGxa6eyum1JrATWvSJKW3thJ9GKHA9n",
		},
		Crypto: "xchain",
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
