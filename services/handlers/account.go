package handlers

import (
  "fmt"
  "net/http"

  "github.com/gorilla/mux"
  "github.com/spf13/viper"

  "github.com/cosmos/cosmos-sdk/client"
  "github.com/cosmos/cosmos-sdk/client/commands"
  "github.com/cosmos/cosmos-sdk/client/commands/query"
  "github.com/cosmos/cosmos-sdk/client/commands/search"
  "github.com/cosmos/cosmos-sdk/modules/coin"
  "github.com/cosmos/cosmos-sdk/stack"
  "github.com/tendermint/tmlibs/common"
)

// doQueryAccount is the HTTP handlerfunc to query an account
// It expects a query string with
func doQueryAccount(w http.ResponseWriter, r *http.Request) {
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

// doSearchTxCoin is the HTTP handlerfunc to search for
// all SendTx transactions with this account as sender
// or receiver
func doSearchTxCoin(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  account := args["address"]
  actor, err := commands.ParseActor(account)
  if err != nil {
    common.WriteError(w, err)
    return
  }

  findSender := fmt.Sprintf("coin.sender='%s'", actor)
  findReceiver := fmt.Sprintf("coin.receiver='%s'", actor)
  prove := !viper.GetBool(commands.FlagTrustNode)
  all, err := search.FindAnyTx(prove, findSender, findReceiver)
  if err != nil {
    common.WriteError(w, err)
    return
  }

  // format....
  output, err := search.FormatSearch(all, coin.ExtractCoinTx)
  if err != nil {
    common.WriteError(w, err)
    return
  }

  // display
  if err := printResult(w, output); err != nil {
    common.WriteError(w, err)
  }
}

// mux.Router registrars

func RegisterQueryAccount(r *mux.Router) error {
  r.HandleFunc("/account/{address}", doQueryAccount).Methods("GET")
  return nil
}

func RegisterSearchTxCoin(r *mux.Router) error {
  r.HandleFunc("/account/{address}/tx/coin", doSearchTxCoin).Methods("GET")
  return nil
}

// RegisterAccount is a convenience function to
// register all the  handlers in this module.
func RegisterAccount(r *mux.Router) error {
  funcs := []func(*mux.Router) error{
    RegisterQueryAccount,
    RegisterSearchTxCoin,
  }

  for _, fn := range funcs {
    if err := fn(r); err != nil {
      return err
    }
  }
  return nil
}

// End of mux.Router registrars
