package signer_test

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/signer"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")
	m.Run()
}

func TestReadPrivateKey(t *testing.T) {
	var key string
	key = os.Getenv("PRIVATE_KEY_STR")
	bytes, err := base64.StdEncoding.DecodeString(key)
	assert.Nil(t, err)
	assert.NotEmpty(t, bytes)
	var privateKey *rsa.PrivateKey
	prkI, err := x509.ParsePKCS8PrivateKey(bytes)
	assert.Nil(t, err)
	assert.NotEmpty(t, prkI)
	privateKey = prkI.(*rsa.PrivateKey)
	assert.NotEmpty(t, privateKey)
}

func TestReadPublicKey(t *testing.T) {
	key := os.Getenv("PUBLIC_KEY_STR")
	bytes, err := base64.StdEncoding.DecodeString(key)
	assert.Nil(t, err)
	var publicKey *rsa.PublicKey
	pubKI, err := x509.ParsePKIXPublicKey(bytes)
	publicKey = pubKI.(*rsa.PublicKey)
	assert.Nil(t, err)
	assert.NotNil(t, publicKey)
}

func TestNew(t *testing.T) {
	_, err := signer.New(os.Getenv("PRIVATE_KEY_STR"), os.Getenv("PUBLIC_KEY_STR"))
	assert.Nil(t, err)
}

func TestRsaClient_Sign(t *testing.T) {
	client, _ := signer.New(os.Getenv("PRIVATE_KEY_STR"), os.Getenv("PUBLIC_KEY_STR"))
	signBytes, err := client.Sign([]byte("123456"), crypto.SHA256)
	assert.Nil(t, err)
	assert.NotNil(t, signBytes)
	err = client.Verify([]byte("123456"), signBytes, crypto.SHA256)
	assert.Nil(t, err)
	err = client.Verify([]byte("654321"), signBytes, crypto.SHA256)
	assert.NotNil(t, err)
}

func TestRsaClient_Sign2(t *testing.T) {
	client, err := signer.New(os.Getenv("VZHUO_PRIVATE_KEY_STR"), os.Getenv("VZHUO_PUBLIC_KEY_STR"))
	assert.Nil(t, err)
	signBytes, err := client.Sign([]byte("123456"), crypto.SHA256)
	assert.Nil(t, err)
	assert.NotNil(t, signBytes)
	err = client.Verify([]byte("123456"), signBytes, crypto.SHA256)
	assert.Nil(t, err)
	err = client.Verify([]byte("654321"), signBytes, crypto.SHA256)
	assert.NotNil(t, err)
}

func TestJson(t *testing.T) {
	var data = map[string]interface{}{
		"name": "詹光",
	}
	_, err := json.Marshal(data)
	assert.Nil(t, err)
}

func TestVerifyPrivateKey(t *testing.T) {
	err := signer.VerifyPrivateKey(os.Getenv("PRIVATE_KEY_STR"))
	assert.Nil(t, err)
	err = signer.VerifyPrivateKey("123456")
	assert.NotNil(t, err)
}

func TestVerifyPublicKey(t *testing.T) {
	err := signer.VerifyPublicKey(os.Getenv("PUBLIC_KEY_STR"))
	assert.Nil(t, err)
	err = signer.VerifyPublicKey("123456")
	assert.NotNil(t, err)
}
