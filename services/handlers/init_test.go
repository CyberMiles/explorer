package handlers

import (
	"os"
	"testing"

	"github.com/adams-sarah/prettytest"
	"github.com/adams-sarah/test2doc/test"
	"github.com/adams-sarah/test2doc/vars"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	_ "github.com/cosmos/cosmos-sdk/modules/coin"
	_ "github.com/cosmos/cosmos-sdk/modules/fee"
	_ "github.com/cosmos/cosmos-sdk/modules/nonce"

	_ "github.com/CyberMiles/explorer/services/modules/stake"
)

var router *mux.Router
var server *test.Server

type mainSuite struct {
	prettytest.Suite
}

const (
	Address = "7334A4B2668DE1CEF0DD7DBA695C29449EC3A0D0"
	Height  = 10546
	TxHash  = "640BEECFF7D4035EA7C9ABCC7F83C0DCA933024C"
	RawTx   = "FgMBBmdhaWEtMgAAAAAAAAAAaQAAAAIBAQABBHNpZ3MBFHM0pLJmjeHO8N19umlcKUSew6DQIAEBAAEEc2lncwEUczSksmaN4c7w3X26aVwpRJ7DoNABAQEHZmVybWlvbgAAAAAAAAABAQEAAQRzaWdzARRqmuoDMVmHmdXwCeybfWNbuPNO/wEBAQdmZXJtaW9uAAAAAAAAAAEB5pH7Kschidpq7JobqzmfyiiVaYBT7JBXleIDrF6kbIwjoJzFv6zhuumL1G2yBNFN1ecjuKixcrNYeNTH3epODQFDGBaojUlBJwZnojLRIJGmwWB92kp+Vm9HD+2dU9xwGw=="
)

func TestRunner(t *testing.T) {
	var err error

	homeDir := os.ExpandEnv("$HOME/.explorer-cli")
	viper.Set("home", homeDir)
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(homeDir)  // search root directory
	//viper.Set(sync.FlagSyncJson, "../sync.json")
	err = viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}

	router = mux.NewRouter()
	RegisterStatus(router)
	RegisterBlock(router)
	RegisterAccount(router)
	RegisterTx(router)

	test.RegisterURLVarExtractor(vars.MakeGorillaMuxExtractor(router))

	server, err = test.NewServer(router)
	if err != nil {
		panic(err.Error())
	}
	defer server.Finish()

	prettytest.RunWithFormatter(
		t,
		new(prettytest.TDDFormatter),
		new(mainSuite),
	)
}
