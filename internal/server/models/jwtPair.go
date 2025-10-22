package models

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	log "github.com/sirupsen/logrus"
	"gofermart/internal/server/config"
)

type JwtKeysPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

var jwtKeys *JwtKeysPair

func InitJwtPair() {
	decodedBytes, _ := base64.StdEncoding.DecodeString(config.GetConfig().PrivateKey)
	decodedString := string(decodedBytes)
	privateKeyBlock, _ := pem.Decode([]byte(decodedString))
	if privateKeyBlock == nil {
		log.Fatal("Decoding private key failed")
	}
	privateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal("Parsing private key failed")

	}

	encodedPublicKey := config.GetConfig().PublicKey

	decodedBytes, err = base64.StdEncoding.DecodeString(encodedPublicKey)
	if err != nil {
		log.Fatal("Decoding public key failed:", err)
	}

	publicKeyBlock, _ := pem.Decode(decodedBytes)
	if publicKeyBlock == nil {
		log.Fatal("Decoding pem block failed")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		log.Fatal("Parsing public key failed:", err)
	}

	publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Cannot assert type: *ecdsa.PublicKey")
	}

	jwtKeys = &JwtKeysPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

}
func GetJwtPair() *JwtKeysPair {
	return jwtKeys
}
