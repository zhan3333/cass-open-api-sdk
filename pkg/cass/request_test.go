package cass_test

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go-skysharing-openapi/pkg/cass"
	"go-skysharing-openapi/pkg/cass/context"
	"os"
	"strconv"
	"testing"
)

func TestFactory_NewRequest(t *testing.T) {
	assert.Nil(t, godotenv.Load("../../.env"))
	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	var factoryConf = cass.Config{
		URI:             os.Getenv("API_URL"),
		AppId:           os.Getenv("APPID"),
		UserPublicKey:   os.Getenv("PUBLIC_KEY_STR"),
		UserPrivateKey:  os.Getenv("PRIVATE_KEY_STR"),
		SystemPublicKey: os.Getenv("VZHUO_PUBLIC_KEY_STR"),
		Debug:           debug,
	}
	var err error
	cass.F, err = cass.NewClient(factoryConf)
	assert.Nil(t, err)
	assert.NotNil(t, cass.F)
	_ = cass.F.NewRequest(context.GetBalance)
}
