package main

import (
	"fmt"
	// "os"
	"log"
	"net/http"

  // "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	_ "github.com/cosmos/cosmos-sdk/modules/coin"
	_ "github.com/cosmos/cosmos-sdk/modules/nonce"
	"github.com/cosmos/cosmos-sdk/client/commands"

	"github.com/cybermiles/explorer/services/handlers"
)

var (
	restServerCmd = &cobra.Command{
		Use:   "rest-server",
		Long:  `presents  a nice (not raw hex) interface to the gaia blockchain structure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdRestServer(cmd, args)
		},
	}

	flagPort = "port"
)

func prepareRestServerCommands() {
	commands.AddBasicFlags(restServerCmd)
	restServerCmd.PersistentFlags().IntP(flagPort, "p", 8998, "port to run the server on")
}

func AddV1Routes(r *mux.Router) {
  AddRoutes(r)
}

func AddRoutes(r *mux.Router) {
	routeRegistrars := []func(*mux.Router) error{
		handlers.RegisterStatus,
		handlers.RegisterBlock,
		handlers.RegisterAccount,
		handlers.RegisterTx,
	}

	for _, routeRegistrar := range routeRegistrars {
		if err := routeRegistrar(r); err != nil {
			log.Fatal(err)
		}
	}
}

func cmdRestServer(cmd *cobra.Command, args []string) error {
	router := mux.NewRouter()
  // latest
  AddRoutes(router)
  // v1
  AddV1Routes(router.PathPrefix("/v1").Subrouter())

	addr := fmt.Sprintf(":%d", viper.GetInt(flagPort))

	log.Printf("Serving on %q", addr)
	return http.ListenAndServe(addr, router)
	// return http.ListenAndServe(addr,
 //        handlers.LoggingHandler(os.Stdout, handlers.CORS(
 //            handlers.AllowedMethods([]string{"GET"}),
 //            handlers.AllowedOrigins([]string{"*"}),
 //            handlers.AllowedHeaders([]string{"X-Requested-With"}))(s)))
}
