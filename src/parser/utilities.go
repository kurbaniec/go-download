package parser

import (
	"net/url"
	"strings"
)

// Returns the deciphered URL.
func buildStreamUrl(cipher string, operations *CipherOperations) string {
	cipherMap := getCipherMap(cipher)
	signature := cipherMap["s"]
	decipheredSignature := operations.decipher(signature)
	url := cipherMap["url"] + "&sig=" + decipheredSignature
	return url
}

// Parses the cipher tag.
func getCipherMap(cipher string) map[string]string {
	cipherMap := make(map[string]string)
	params := strings.Split(cipher, "&")
	for _, param := range params {
		cipherDecoded, err := url.QueryUnescape(param)
		errorHandler(err)
		equalsPos := strings.Index(cipherDecoded, "=")
		key := cipherDecoded[:equalsPos]
		value := cipherDecoded[equalsPos+1:]
		cipherMap[key] = value
	}
	return cipherMap
}
