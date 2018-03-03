package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/cosmos/cosmos-sdk/client/commands/query"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/stack"
	"github.com/tendermint/tmlibs/common"
)

// queryAccount is to query an account by address
func queryAccount(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	address := args["address"]
	actor, err := commands.ParseActor(address)
	if err != nil {
		common.WriteError(w, err)
		return
	}

	actor = coin.ChainAddr(actor)
	key := stack.PrefixedKey(coin.NameCoin, actor.Bytes())
	account := new(coin.Account)
	prove := !viper.GetBool(commands.FlagTrustNode)
	height, err := query.GetParsed(key, account, query.GetHeight(), prove)
	if client.IsNoDataErr(err) {
		err := fmt.Errorf("account bytes are empty for address: %q", address)
		common.WriteError(w, err)
		return
	} else if err != nil {
		common.WriteError(w, err)
		return
	}

	if err := query.FoutputProof(w, account, height); err != nil {
		common.WriteError(w, err)
	}
}

// mux.Router registrars

func RegisterQueryAccount(r *mux.Router) error {
	r.HandleFunc("/account/{address}", queryAccount).Methods("GET")
	return nil
}

// RegisterAccount is a convenience function to
// register all the	handlers in this module.
func RegisterAccount(r *mux.Router) error {
	funcs := []func(*mux.Router) error{
		RegisterQueryAccount,
	}

	for _, fn := range funcs {
		if err := fn(r); err != nil {
			return err
		}
	}
	return nil
}

// End of mux.Router registrars
