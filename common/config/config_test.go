package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	c := GetInstance()
	if c == nil {
		t.Error("GetInstance test failed")
	}

	c.SetGMCrypto()
	if c.Crypto != CRYPTO_GM {
		t.Error("GetInstance test failed")
	}
	c.SetXchainCrypto()
	if c.Crypto != CRYPTO_XCHAIN {
		t.Error("GetInstance test failed")
	}
}

func TestSetConfig(t *testing.T) {
	SetConfig("a", "b", "c", "1", true, true, "1")
	c := GetInstance()
	if c == nil {
		t.Error("GetInstance test failed")
	}

	if c.EndorseServiceHost != "a" {
		t.Error("SetConfig check host failed")
	}

	if c.ComplianceCheck.ComplianceCheckEndorseServiceAddr != "b" {
		t.Error("SetConfig check ComplianceCheckEndorseServiceAddr failed")
	}

	if c.ComplianceCheck.ComplianceCheckEndorseServiceFeeAddr != "c" {
		t.Error("SetConfig check ComplianceCheckEndorseServiceFeeAddr failed")
	}

	if c.ComplianceCheck.ComplianceCheckEndorseServiceFee != 1 {
		t.Error("SetConfig check ComplianceCheckEndorseServiceFee failed")
	}
}
