package graph

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"pluralith/pkg/auxiliary"
)

func DecryptPayload(payload string, iv string) (string, error) {
	functionName := "DecryptPayload"

	// Decode initialization vector
	ivDecoded, ivErr := hex.DecodeString(iv)
	if ivErr != nil {
		return "", fmt.Errorf("failed to decode IV -> %v: %w", functionName, ivErr)
	}
	// Decode payload
	payloadDecoded, payloadErr := hex.DecodeString(payload)
	if payloadErr != nil {
		return "", fmt.Errorf("failed to decode payload -> %v: %w", functionName, payloadErr)
	}

	// Decrypt payload
	block, blockErr := aes.NewCipher([]byte(auxiliary.StateInstance.APIKey))
	if blockErr != nil {
		return "", fmt.Errorf("failed to create cipher -> %v: %w", functionName, blockErr)
	}

	mode := cipher.NewCBCDecrypter(block, []byte(ivDecoded))
	mode.CryptBlocks(payloadDecoded, payloadDecoded)

	// Trim decryption appendix and stringify result
	payloadDecoded = payloadDecoded[:len(payloadDecoded)-8]
	payloadString := string(payloadDecoded)

	return payloadString, nil
}
