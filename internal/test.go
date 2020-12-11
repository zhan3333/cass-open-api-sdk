package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/method"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err.Error())
	}
	factoryConf := cass.Config{
		URI:             os.Getenv("API_URL"),
		AppId:           os.Getenv("APPID"),
		UserPublicKey:   os.Getenv("PUBLIC_KEY_STR"),
		UserPrivateKey:  os.Getenv("PRIVATE_KEY_STR"),
		SystemPublicKey: os.Getenv("VZHUO_PUBLIC_KEY_STR"),
	}
	fmt.Printf("%v\n", factoryConf)

	f, err := cass.NewClient(factoryConf)
	if err != nil {
		panic(err.Error())
	}
	count := 0
	for true {
		count++
		fmt.Printf("exec %d pay", count)
		request := f.NewRequest(method.M.PayOneBankRemit)
		request.SetBizParams(map[string]interface{}{
			"payChannelK": "1",
			//"payeeChannelType": "2",
			"orderData": [1]interface{}{
				map[string]interface{}{
					"orderSN":          uuid.New().String(),
					"receiptFANO":      "13517210601",
					"payeeAccount":     "詹光",
					"requestPayAmount": "0.01",
					"notifyUrl":        "http://www.baidu.com/a/b?a=b",
					"identityCard":     "420222199212041057",
				},
			},
		})
		response := request.Send().(*cass.Response)
		if response.Error() != nil {
			fmt.Printf("发生错误: %s \n", response.Error().Error())
		} else {
			fmt.Printf("response code: %v \n", response.HTTP.StatusCode)
			fmt.Printf("response: %v \n", response.String())
		}
		time.Sleep(1 * time.Second)
	}
}
