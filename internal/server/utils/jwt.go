package utils

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type JwtKeysPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

func InitJwtPair(privateKeyString, publicKeyString string) JwtKeysPair {
	decodedBytes, _ := base64.StdEncoding.DecodeString(privateKeyString)
	decodedString := string(decodedBytes)
	privateKeyBlock, _ := pem.Decode([]byte(decodedString))
	if privateKeyBlock == nil {
		log.Fatal("Decoding private key failed")
	}
	privateKey, err := x509.ParseECPrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		log.Fatal("Parsing private key failed")

	}
	decodedBytes, err = base64.StdEncoding.DecodeString(publicKeyString)
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

	return JwtKeysPair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

}

func (jwtKp JwtKeysPair) IsValidToken(signedToken, tokenType string) bool {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return false, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKp.PublicKey, nil
	})

	if err != nil {
		log.Error("JWT parse", "err", err.Error())
		return false
	}
	claims := token.Claims.(jwt.MapClaims)

	return claims["type"] == tokenType
}

func (jwtKp JwtKeysPair) GenerateJWT(userID int64) (TokenPair, error) {
	var tokenPair TokenPair
	accessToken := jwt.New(jwt.SigningMethodES256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["iat"] = time.Now().Unix()
	claims["user_id"] = userID
	claims["type"] = "access"
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(30)).Unix()
	signedToken, err := accessToken.SignedString(jwtKp.PrivateKey)
	if err != nil {
		return tokenPair, err
	}

	refreshToken := jwt.New(jwt.SigningMethodES256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["user_id"] = userID
	refreshClaims["iat"] = time.Now().Unix()
	refreshClaims["exp"] = time.Now().Add(time.Minute * time.Duration(120)).Unix()
	refreshClaims["type"] = "refresh"
	signedRefreshToken, err := refreshToken.SignedString(jwtKp.PrivateKey)
	if err != nil {
		return tokenPair, err
	}
	tokenPair.AccessToken = signedToken
	tokenPair.RefreshToken = signedRefreshToken
	return tokenPair, nil

}

func (jwtKp JwtKeysPair) GetUserID(tokenString string) (int64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKp.PublicKey, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userIDStr := fmt.Sprint(claims["user_id"])
	return strconv.ParseInt(userIDStr, 10, 64)
}
