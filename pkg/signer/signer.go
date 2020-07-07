package signer

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

type Signer interface {
	// 签名
	Sign(src []byte, hash crypto.Hash) ([]byte, error)
	// 验证签名
	Verify(src []byte, sign []byte, hash crypto.Hash) error
}

type rsaClient struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func (client *rsaClient) Sign(src []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, client.PrivateKey, hash, hashed)
}

func (client *rsaClient) Verify(src []byte, sign []byte, hash crypto.Hash) error {
	h := hash.New()
	h.Write(src)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(client.PublicKey, hash, hashed, sign)
}

func New(privateKey, publicKey string) (Signer, error) {
	var priKey *rsa.PrivateKey
	var pubKey *rsa.PublicKey
	var err error
	if privateKey != "" {
		priKey, err = readPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}
	}

	if publicKey != "" {
		pubKey, err = readPublicKey(publicKey)
		if err != nil {
			return nil, err
		}
	}
	return &rsaClient{
		PrivateKey: priKey,
		PublicKey:  pubKey,
	}, nil
}

func VerifyPrivateKey(pri string) error {
	if pri == "" {
		return fmt.Errorf("无效的空字符串")
	}
	_, err := readPrivateKey(pri)
	return err
}

func VerifyPublicKey(pub string) error {
	if pub == "" {
		return fmt.Errorf("无效的空字符串")
	}
	_, err := readPublicKey(pub)
	return err
}

// 读取私钥对象 (pkcs8/pkcs1)
func readPrivateKey(key string) (*rsa.PrivateKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var privateKey *rsa.PrivateKey
	prkI, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		prkI, err = x509.ParsePKCS1PrivateKey(bytes)
		if err != nil {
			return nil, err
		}
	}
	privateKey = prkI.(*rsa.PrivateKey)
	return privateKey, nil
}

// 读取公钥对象
// PKCS8格式单行key处理
func readPublicKey(key string) (*rsa.PublicKey, error) {
	bytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var publicKey *rsa.PublicKey
	pubKI, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, err
	}
	publicKey = pubKI.(*rsa.PublicKey)
	return publicKey, nil
}
