package bigip

import (
	bigip "github.com/f5devcentral/go-bigip"
	"testing"
)

func TestConfigClientUnitTC1(t *testing.T) {
	// test string => expected error count
	testConfig := &bigip.Config{}
	testConfig.Address = "192.168.1.1"
	testConfig.Username = "testuser"
	testConfig.Password = "testpasswd"
	testConfig.Port = "443"
	testConfig.CertVerifyDisable = true
	_, _ = Client(testConfig)
}

func TestConfigClientUnitTC2(t *testing.T) {
	// test string => expected error count
	testConfig := &bigip.Config{}
	testConfig.Address = "192.168.1.1"
	testConfig.Username = "testuser"
	testConfig.Password = "testpasswd"
	testConfig.Port = "443"
	testConfig.CertVerifyDisable = false
	_, _ = Client(testConfig)
}

func TestConfigClientUnitTC3(t *testing.T) {
	testConfig := &bigip.Config{}
	testConfig.Address = "192.168.1.1"
	testConfig.Username = "testuser"
	testConfig.Password = "testpasswd"
	testConfig.Port = "443"
	testConfig.CertVerifyDisable = false
	testConfig.TrustedCertificate = folder + "/../examples/servercert.crt"
	_, _ = Client(testConfig)
}

func TestConfigClientUnitTC4(t *testing.T) {
	testConfig := &bigip.Config{}
	testConfig.Address = "192.168.1.1"
	testConfig.Username = "testuser"
	testConfig.Password = "testpasswd"
	testConfig.Port = "443"
	testConfig.LoginReference = "tmos"
	testConfig.CertVerifyDisable = true
	_, _ = Client(testConfig)
}
