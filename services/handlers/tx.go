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
  "github.com/cosmos/cosmos-sdk/modules/fee"
  "github.com/cybermiles/explorer/services/modules/stake"

  wire "github.com/tendermint/go-wire"
  "github.com/tendermint/tmlibs/common"
  ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type resp struct {
  Height int64       `json:"height"`
  Tx   interface{} `json:"tx"`
}

// queryRawTx is the HTTP handlerfunc to query a raw tx by txhash
func queryRawTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]
  raw := true

  err := getTx(w, txhash, raw)
  if err != nil {
    common.WriteError(w, err)
  }
}

// queryTx is the HTTP handlerfunc to query an "useful" tx by txhash
func queryTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]
  raw := false

  err := getTx(w, txhash, raw)
  if err != nil {
    common.WriteError(w, err)
  }
}

func getTx(w http.ResponseWriter, txhash string, raw bool) error {
  prove := !viper.GetBool(commands.FlagTrustNode)
  bkey, err := hex.DecodeString(common.StripHex(txhash))
  if err != nil {
    return err
  }

  client := commands.GetNode()
  res, err := client.Tx(bkey, prove)
  if err != nil {
    return err
  }

  // format
  wrap, err := formatTx(res.Height, res.Proof.Data, raw)
  if err != nil {
    return err
  }

  // display
  return printResult(w, wrap)
}

func formatTx(height int64, data []byte, raw bool) (interface{}, error) {
  tx, err := sdk.LoadTx(data)
  if err != nil {
    return tx, err
  }
  if (!raw) {
    txl, ok := tx.Unwrap().(sdk.TxLayer)
    var txi sdk.Tx
    loop: for ok {
      txi = txl.Next()
      switch txi.Unwrap().(type) {
        case fee.Fee, coin.SendTx, stake.TxDeclareCandidacy:
          tx = txi
          break loop
      }
      txl, ok = txi.Unwrap().(sdk.TxLayer)
    }
  }
  wrap := &resp{height, tx}
  return wrap, nil
}

// searchTxByBlock is the HTTP handlerfunc to search for
// "useful" transaction by block height
func searchTxByBlock(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["height"]
  query := fmt.Sprintf("height=%s", txhash)

  err := searchTx(w, query)
  if err != nil {
    common.WriteError(w, err)
  }
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

  err = searchTx(w, findSender, findReceiver)
  if err != nil {
    common.WriteError(w, err)
  }
}

func searchTx(w http.ResponseWriter, queries ...string) error {
  prove := !viper.GetBool(commands.FlagTrustNode)

  all, err := search.FindAnyTx(prove, queries ...)
  if err != nil {
    return err
  }

  // format
  wrap, err := formatSearch(all)
  if err != nil {
    return err
  }

  // display
  return printResult(w, wrap)
}

func formatSearch(res []*ctypes.ResultTx) ([]interface{}, error) {
  out := make([]interface{}, 0, len(res))
  for _, r := range res {
    wrap, err := formatTx(r.Height, r.Tx, false)
    if err != nil {
      return nil, err
    }
    out = append(out, wrap)
  }
  return out, nil
}

// decodeRaw is the HTTP handlerfunc to decode tx string
func decodeRaw(w http.ResponseWriter, r *http.Request) {
  buf := new(bytes.Buffer)
  buf.ReadFrom(r.Body)
  body := buf.String()

  err := decode(w, body)
  if err != nil {
    common.WriteError(w, err)
  }
}

func decode(w http.ResponseWriter, body string) error {
  data, err := base64.StdEncoding.DecodeString(body)
  if err != nil {
    return err
  }

  var tx sdk.Tx
  err = wire.ReadBinaryBytes([]byte(data), &tx)
  if err != nil {
    return err
  }

  // display
  return printResult(w, tx)
}

// mux.Router registrars

func RegisterQueryTx(r *mux.Router) error {
  r.HandleFunc("/tx/{txhash}", queryTx).Methods("GET")
  return nil
}

func RegisterQueryRawTx(r *mux.Router) error {
  r.HandleFunc("/tx/{txhash}/raw", queryRawTx).Methods("GET")
  return nil
}

func registerSearchTxByBlock(r *mux.Router) error {
  r.HandleFunc("/block/{height}/tx", searchTxByBlock).Methods("GET")
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
    RegisterQueryRawTx,
    registerSearchTxByBlock,
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
