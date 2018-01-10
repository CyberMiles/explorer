package handlers

import (
  "fmt"
  "bytes"
  "net/http"
  "encoding/hex"
  "encoding/base64"

  "github.com/gorilla/mux"
  "github.com/spf13/viper"

  sdk "github.com/cosmos/cosmos-sdk"
  "github.com/cosmos/cosmos-sdk/client/commands"
  "github.com/cosmos/cosmos-sdk/client/commands/search"
  "github.com/cosmos/cosmos-sdk/modules/coin"
  wire "github.com/tendermint/go-wire"
  "github.com/tendermint/tmlibs/common"
)

// queryTx is the HTTP handlerfunc to query a raw tx by txhash
func queryTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]

  bkey, err := hex.DecodeString(common.StripHex(txhash))
  if err != nil {
    common.WriteError(w, err)
    return
  }
  prove := !viper.GetBool(commands.FlagTrustNode)

  getTx(w, prove, bkey)
}

// searchCoinTx is the HTTP handlerfunc to search for
// a SendTx transaction by txhash
func searchCoinTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]

  query := fmt.Sprintf("tx.hash='%s'", txhash)
  prove := !viper.GetBool(commands.FlagTrustNode)

  searchTx(w, prove, query)
}

// searchCoinTxByBlock is the HTTP handlerfunc to search for
// all SendTx transactions by block height
func searchCoinTxByBlock(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["height"]

  query := fmt.Sprintf("height=%s", txhash)
  prove := !viper.GetBool(commands.FlagTrustNode)

  searchTx(w, prove, query)
}

// searchCoinTxByAccount is the HTTP handlerfunc to search for
// all SendTx transactions with this account as sender
// or receiver
func searchCoinTxByAccount(w http.ResponseWriter, r *http.Request) {
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

  searchTx(w, prove, findSender, findReceiver)
}

func decodeRaw(w http.ResponseWriter, r *http.Request) {
  buf := new(bytes.Buffer)
  buf.ReadFrom(r.Body)
  body := buf.String()

  decodeTx(w, body)
}

type proof struct {
  Height int64       `json:"height"`
  Data   interface{} `json:"data"`
}
// getTx parses anything that was previously registered as sdk.Tx
func getTx(w http.ResponseWriter, prove bool, bkey []byte) {
  client := commands.GetNode()
  res, err := client.Tx(bkey, prove)
  if err != nil {
    common.WriteError(w, err)
    return
  }
  var tx sdk.Tx
  err = wire.ReadBinaryBytes(res.Proof.Data, &tx)
  // tx, err := coin.ExtractCoinTx(data)
  if err != nil {
    common.WriteError(w, err)
    return
  }
  wrap := &proof{res.Height, tx}
  // display
  if err := printResult(w, wrap); err != nil {
    common.WriteError(w, err)
  }
}

func searchTx(w http.ResponseWriter, prove bool, queries ...string) {
  all, err := search.FindAnyTx(prove, queries ...)
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

func decodeTx(w http.ResponseWriter, body string) {
  data, err := base64.StdEncoding.DecodeString(body)
  if err != nil {
    common.WriteError(w, err)
    return
  }

  var tx sdk.Tx
  err = wire.ReadBinaryBytes([]byte(data), &tx)
  if err != nil {
    common.WriteError(w, err)
    return
  }
  // display
  if err := printResult(w, tx); err != nil {
    common.WriteError(w, err)
  }
}

// mux.Router registrars

func RegisterQueryTx(r *mux.Router) error {
  r.HandleFunc("/tx/{txhash}", queryTx).Methods("GET")
  return nil
}

func RegisterSearchCoinTx(r *mux.Router) error {
  r.HandleFunc("/tx/coin/{txhash}", searchCoinTx).Methods("GET")
  return nil
}

func registerSearchCoinTxByBlock(r *mux.Router) error {
  r.HandleFunc("/block/{height}/tx/coin", searchCoinTxByBlock).Methods("GET")
  return nil
}

func RegisterSearchCoinTxByAccount(r *mux.Router) error {
  r.HandleFunc("/account/{address}/tx/coin", searchCoinTxByAccount).Methods("GET")
  return nil
}

func RegisterDecodeRaw(r *mux.Router) error {
  r.HandleFunc("/tx/decode", decodeRaw).Methods("POST")
  return nil
}

// RegisterTx is a convenience function to
// register all the  handlers in this module.
func RegisterTx(r *mux.Router) error {
  funcs := []func(*mux.Router) error{
    RegisterQueryTx,
    RegisterSearchCoinTx,
    registerSearchCoinTxByBlock,
    RegisterSearchCoinTxByAccount,
    RegisterDecodeRaw,
  }

  for _, fn := range funcs {
    if err := fn(r); err != nil {
      return err
    }
  }
  return nil
}

// End of mux.Router registrars
