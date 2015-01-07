package main

import (
	"testing"
)

func TestAPICongigParameters(t *testing.T) {
	analytics_client := AnalyticsClient{}
	analytics_client.Init()

	if &analytics_client.apiVersion == nil {
		t.Errorf("version was expected when initializing analytics api client")
	}

	if &analytics_client.apiKey == nil {
		t.Errorf("key was expected when initializing analytics api client")
	}

	if &analytics_client.apiSecret == nil {
		t.Errorf("secret was expected when initializing analytics api client")
	}

}

func TestCallWithCorrectCreds(t *testing.T) {
	apiSecret := []byte("123456")
	calculatedMac := mac([]byte("abcdefg"), apiSecret)
	expectedMac := "04b09561db1d5aa5e8a6ed6228d9f8a4f933ad1e"
	if calculatedMac != expectedMac {
		t.Errorf("unexpected result from mac function. expected %s got %s", expectedMac, calculatedMac)
	}
}
