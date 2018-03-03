package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/commands"
	_ "github.com/cosmos/cosmos-sdk/modules/auth"
	_ "github.com/cosmos/cosmos-sdk/modules/base"
	_ "github.com/cosmos/cosmos-sdk/modules/coin"
	_ "github.com/cosmos/cosmos-sdk/modules/fee"
	_ "github.com/cosmos/cosmos-sdk/modules/nonce"

	services "github.com/CyberMiles/explorer/services/handlers"
	_ "github.com/CyberMiles/explorer/services/modules/stake"
)

const (
	FlagPort = "port"
	MgoUrl   = "mgo-url"
)

var (
	restServerCmd = &cobra.Command{
		Use:  "rest-server",
		Long: `presents  a nice (not raw hex) interface to the gaia blockchain structure.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmdRestServer(cmd, args)
		},
	}
)

func prepareRestServerCommands() {
	commands.AddBasicFlags(restServerCmd)
	restServerCmd.PersistentFlags().String(MgoUrl, "localhost:27017", "url of MongoDB")
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
	startWatch()
	router := mux.NewRouter()
	// latest
	AddRoutes(router)
	// v1
	AddV1Routes(router.PathPrefix("/v1").Subrouter())

	addr := fmt.Sprintf(":%d", viper.GetInt(FlagPort))

	log.Printf("Serving on %q", addr)

	// loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return http.ListenAndServe(addr, router)
	//return http.ListenAndServe(addr,
	//	handlers.LoggingHandler(os.Stdout, handlers.CORS(
	//		handlers.AllowedMethods([]string{"GET"}),
	//		handlers.AllowedOrigins([]string{"*"}),
	//		handlers.AllowedHeaders([]string{"X-Requested-With"}))(s)))
}
