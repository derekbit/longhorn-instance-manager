package proxy

import (
	"os"
	"testing"

	btypes "github.com/longhorn/backupstore/types"
)

func TestGetBackupCredentialRejectsDisallowedEnv(t *testing.T) {
	_, err := getBackupCredential("s3://backupbucket@us-east-1/volume", []string{"LD_PRELOAD=/tmp/evil.so"})
	if err == nil {
		t.Fatal("expected disallowed environment variable to be rejected")
	}
}

func TestGetBackupCredentialOverridesRequestValuesWithoutMutatingProcessEnv(t *testing.T) {
	t.Setenv(btypes.AWSAccessKey, "process-access")
	t.Setenv(btypes.AWSSecretKey, "process-secret")
	t.Setenv(btypes.HTTPProxy, "http://process-proxy")

	credential, err := getBackupCredential("s3://backupbucket@us-east-1/volume", []string{
		btypes.AWSAccessKey + "=request-access",
		btypes.AWSSecretKey + "=request-secret",
		btypes.HTTPProxy + "=http://request-proxy",
	})
	if err != nil {
		t.Fatalf("expected credential override to succeed: %v", err)
	}

	if credential[btypes.AWSAccessKey] != "request-access" {
		t.Fatalf("expected request access key override, got %q", credential[btypes.AWSAccessKey])
	}
	if credential[btypes.AWSSecretKey] != "request-secret" {
		t.Fatalf("expected request secret key override, got %q", credential[btypes.AWSSecretKey])
	}
	if credential[btypes.HTTPProxy] != "http://request-proxy" {
		t.Fatalf("expected request HTTP proxy override, got %q", credential[btypes.HTTPProxy])
	}

	if got := os.Getenv(btypes.AWSAccessKey); got != "process-access" {
		t.Fatalf("expected process environment to remain unchanged, got %q", got)
	}
	if got := os.Getenv(btypes.HTTPProxy); got != "http://process-proxy" {
		t.Fatalf("expected process HTTP proxy to remain unchanged, got %q", got)
	}
}
