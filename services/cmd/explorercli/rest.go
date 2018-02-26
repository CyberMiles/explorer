package main

import (
	"fmt"
	// "os"
	"log"
	"net/http"

	"github.com/gorilla/mux"
  // "github.com/gorilla/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	_ "github.com/cosmos/cosmos-sdk/modules/coin"
	_ "github.com/cosmos/cosmos-sdk/modules/nonce"
	_ "github.com/cosmos/cosmos-sdk/modules/fee"
	"github.com/cosmos/cosmos-sdk/client/commands"

	_ "github.com/ly0129ly/explorer/services/modules/stake"
	services "github.com/ly0129ly/explorer/services/handlers"
	"github.com/gorilla/handlers"
	"os"
)

const (
	FlagPort = "port"
	MgoUrl = "mgo-url"
	Cron = "cron"
)

var (
	restServerCmd = &cobra.Command{
		Use:   "rest-server",
		Long:  `presents  a nice (not raw hex) interface to the gaia blockchain structure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdRestServer(cmd, args)
		},
	}
)

func prepareRestServerCommands() {
	commands.AddBasicFlags(restServerCmd)
	restServerCmd.PersistentFlags().String(MgoUrl, "localhost:27017", "url of MongoDB")
	restServerCmd.PersistentFlags().String(Cron, "@every 60s", "a new crontab schedule for sync block")
	restServerCmd.PersistentFlags().IntP(FlagPort, "p", 8998, "port to run the server on")
}

func AddV1Routes(r *mux.Router) {
  AddRoutes(r)
}

func AddRoutes(r *mux.Router) {
	routeRegistrars := []func(*mux.Router) error{
		services.RegisterStatus,
		services.RegisterBlock,
		services.RegisterAccount,
		services.RegisterTx,
	}

	for _, routeRegistrar := range routeRegistrars {
		if err := routeRegistrar(r); err != nil {
			log.Fatal(err)
		}
	}
}

func cmdRestServer(cmd *cobra.Command, args []string) error {
	startSync()
	router := mux.NewRouter()
  // latest
  AddRoutes(router)


  // v1
  AddV1Routes(router.PathPrefix("/v1").Subrouter())

	addr := fmt.Sprintf(":%d", viper.GetInt(FlagPort))

	log.Printf("Serving on %q", addr)

	// loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	//return http.ListenAndServe(addr, router)
	 return http.ListenAndServe(addr,
         handlers.LoggingHandler(os.Stdout, handlers.CORS(
             handlers.AllowedMethods([]string{"GET"}),
             handlers.AllowedOrigins([]string{"*"}),
             handlers.AllowedHeaders([]string{"X-Requested-With"}))(router)))
}
