package scoring

import (
	"encoding/base64"
	"testing"
)

func TestDecryptString(t *testing.T) {
	const (
		ciphertext = "MxCMvqBuYDDc4pzwO8mpZ1GRmTYLSp2DY9TJoMnR4sk=" // fun spoiler
		encodedIv  = "99wtr/3Cfg1L+avJwKQ0bWRzUSFCNtjDZwr3QdOwLHY="
	)

	iv, err := base64.StdEncoding.DecodeString(encodedIv)
	if err != nil {
		t.Fatalf("decode IV: %v", err)
	}

	plaintext, err := DecryptString(ciphertext, iv, encryptionKeyV1)
	if err != nil {
		t.Fatalf("decrypt string: %v", err)
	}
	if plaintext != "100:True:True:False:True" {
		t.Fatalf("plaintext = %q, want %q", plaintext, "100:True:True:False:True")
	}
	t.Logf("decrypted string: %q", plaintext)
}
