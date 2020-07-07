package conf

import (
	"encoding/json"
	"fmt"
	"go-skysharing-openapi/pkg/signer"
	"io/ioutil"
	"os"
)

type Config struct {
	Mode            string `json:"mode"`
	Uri             string `json:"uri"`
	SystemPublicKey string `json:"system_public_key"`
	UserPublicKey   string `json:"user_public_key"`
	UserPrivateKey  string `json:"user_private_key"`
	AppId           string `json:"app_id"`
}

var C Config

func init() {
	err := ReadCache()
	if err != nil {
		fmt.Printf("Read Cache Err: %+v \n", err)
	}
}

func Set(uri, appId, sysPublicKey, userPublicKey, userPrivateKey string) error {
	var err error
	C.Uri = uri
	C.AppId = appId
	if err = signer.VerifyPublicKey(sysPublicKey); err != nil {
		return fmt.Errorf("无效的系统公钥: %s", err.Error())
	}
	if err = signer.VerifyPublicKey(userPublicKey); err != nil {
		return fmt.Errorf("无效的用户公钥: %s", err.Error())
	}
	if err = signer.VerifyPrivateKey(userPrivateKey); err != nil {
		return fmt.Errorf("无效的系统公钥: %s", err.Error())
	}
	C.SystemPublicKey = sysPublicKey
	C.UserPublicKey = userPublicKey
	C.UserPrivateKey = userPrivateKey
	Cache()
	return nil
}

func Cache() {
	f, _ := os.Open("cache.json")
	b, _ := json.Marshal(C)
	_, _ = f.Write(b)
	fmt.Printf("Cache: %+v \n", C)
}

func ReadCache() error {
	var err error
	f, err := os.Open("cache.json")
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	fmt.Printf("Read str: %s \n", string(b))
	err = json.Unmarshal(b, &C)
	fmt.Printf("Read Cache: %+v \n", C)
	if err != nil {
		return err
	}
	return nil
}
