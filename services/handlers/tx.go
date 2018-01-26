package handlers

import (
  "fmt"
  "bytes"
  "strings"
  "net/http"
  "io/ioutil"
  "encoding/hex"
  "encoding/base64"
  "encoding/json"

  "github.com/gorilla/mux"
  "github.com/spf13/viper"

  sdk "github.com/cosmos/cosmos-sdk"
  "github.com/cosmos/cosmos-sdk/client/commands"
  "github.com/cosmos/cosmos-sdk/client/commands/search"
  "github.com/cosmos/cosmos-sdk/modules/coin"
  "github.com/cosmos/cosmos-sdk/modules/fee"

  wire "github.com/tendermint/go-wire"
  "github.com/tendermint/tmlibs/common"
  ctypes "github.com/tendermint/tendermint/rpc/core/types"

  "github.com/cybermiles/explorer/services/modules/stake"
  "github.com/cybermiles/explorer/services/modules/sync"
)

type resp struct {
  Height int64       `json:"height"`
  Tx   interface{} `json:"tx"`
  TxHash string   `json:"txhash"`
}

// queryRawTx is to query a raw transaction by txhash
func queryRawTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]
  raw := true

  err := getTx(w, txhash, raw)
  if err != nil {
    common.WriteError(w, err)
  }
}

// queryTx is to query "inner" transaction by txhash
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
  wrap, err := formatTx(res.Height, res.Proof.Data, raw, txhash)
  if err != nil {
    return err
  }

  // display
  return printResult(w, wrap)
}

func formatTx(height int64, data []byte, raw bool, txhash string) (interface{}, error) {
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
        case fee.Fee, coin.SendTx, stake.TxDelegate, stake.TxDeclareCandidacy, stake.TxUnbond:
          tx = txi
          break loop
      }
      txl, ok = txi.Unwrap().(sdk.TxLayer)
    }
  }
  wrap := &resp{height, tx, strings.ToUpper(txhash)}
  return wrap, nil
}

// searchTxByBlock is to search for inner transaction by block height
func searchTxByBlock(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  height := args["height"]
  query := fmt.Sprintf("height=%s", height)

  wrap, err := searchTx(w, query)
  if err != nil {
    common.WriteError(w, err)
  }
  // display
  printResult(w, wrap)
}

// searchCoinTxByAccount is to search for
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

  wrap, err := searchTx(w, findSender, findReceiver)
  if err != nil {
    common.WriteError(w, err)
  }
  // display
  printResult(w, wrap)
}

func searchTx(w http.ResponseWriter, queries ...string) ([]interface{}, error) {
  prove := !viper.GetBool(commands.FlagTrustNode)

  all, err := search.FindAnyTx(prove, queries ...)
  if err != nil {
    return nil, err
  }

  // format
  return formatSearch(all)
}

func formatSearch(res []*ctypes.ResultTx) ([]interface{}, error) {
  out := make([]interface{}, 0, len(res))
  for _, r := range res {
    wrap, err := formatTx(r.Height, r.Tx, false, hex.EncodeToString(r.Tx.Hash()))
    if err != nil {
      return nil, err
    }
    out = append(out, wrap)
  }
  return out, nil
}

// decodeRaw is to decode tx string
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


// queryRecentCoinTx is to get recent coin transactions
func queryRecentCoinTx(w http.ResponseWriter, r *http.Request) {
  file := viper.GetString(sync.FlagSyncJson)
  raw, err := ioutil.ReadFile(file)
  if err != nil {
    common.WriteError(w, err)
    return
  }
  var syncResult sync.SyncResult
  json.Unmarshal(raw, &syncResult)

  // display
  printResult(w, syncResult.CoinTxs)
}

// queryRecentStakeTx is to get recent stake transactions
func queryRecentStakeTx(w http.ResponseWriter, r *http.Request) {
  file := viper.GetString(sync.FlagSyncJson)
  raw, err := ioutil.ReadFile(file)
  if err != nil {
    common.WriteError(w, err)
    return
  }
  var syncResult sync.SyncResult
  json.Unmarshal(raw, &syncResult)

  // display
  printResult(w, syncResult.StakeTxs)
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

func RegisterQueryRecentCoinTx(r *mux.Router) error {
  r.HandleFunc("/txs/recentcoin", queryRecentCoinTx).Methods("GET")
  return nil
}

func RegisterQueryRecentStakeTx(r *mux.Router) error {
  r.HandleFunc("/txs/recentstake", queryRecentStakeTx).Methods("GET")
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
    RegisterQueryRecentCoinTx,
    RegisterQueryRecentStakeTx,
  }

  for _, fn := range funcs {
    if err := fn(r); err != nil {
      return err
    }
  }
  return nil
}

// End of mux.Router registrars
