package crypto

import (
	"testing"
)

func TestEncryptionAndDecryption(t *testing.T) {
	originalText := "hello world"
	key, _ := GenerateKey() // Assuming you have a function to generate a key

	encryptedText, err := Encrypt(originalText, key)
	if err != nil {
		t.Fatalf("Encrypt function failed: %v", err)
	}
	decryptedText, err := Decrypt(encryptedText, key)
	if err != nil {
		t.Fatalf("Decrypt function failed: %v", err)
	}
	t.Log(decryptedText)
	if decryptedText != originalText {
		t.Errorf("Decrypted text does not match original. got = %s, want = %s", decryptedText, originalText)
	}
}
