package handlers

import (
  "fmt"
  "io"
  "net/http"
  "encoding/hex"

  "github.com/gorilla/mux"
  "github.com/spf13/viper"

  sdk "github.com/cosmos/cosmos-sdk"
  "github.com/cosmos/cosmos-sdk/client/commands"
  "github.com/cosmos/cosmos-sdk/client/commands/query"
  "github.com/cosmos/cosmos-sdk/client/commands/search"
  "github.com/cosmos/cosmos-sdk/modules/coin"
  wire "github.com/tendermint/go-wire"
  ctypes "github.com/tendermint/tendermint/rpc/core/types"
  "github.com/tendermint/tmlibs/common"
)

func doQueryTx(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]

  bkey, err := hex.DecodeString(common.StripHex(txhash))
  if err != nil {
    common.WriteError(w, err)
    return
  }
  prove := !viper.GetBool(commands.FlagTrustNode)
  res, err := getTx(bkey, prove)
  if err != nil {
    common.WriteError(w, err)
    return
  }

  if err := showTx(w, res.Height, res.Proof.Data); err != nil {
    common.WriteError(w, err)
    return
  }
}

func searchTxCoinByHash(w http.ResponseWriter, r *http.Request) {
  args := mux.Vars(r)
  txhash := args["txhash"]

  query := fmt.Sprintf("tx.hash='%s'", txhash)
  prove := !viper.GetBool(commands.FlagTrustNode)
  all, err := search.FindAnyTx(prove, query)
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

func getTx(bkey []byte, prove bool) (*ctypes.ResultTx, error) {
  client := commands.GetNode()
  return client.Tx(bkey, prove)
}

// showTx parses anything that was previously registered as sdk.Tx
func showTx(w io.Writer, height int64, data []byte) error {
  var tx sdk.Tx
  err := wire.ReadBinaryBytes(data, &tx)
  // tx, err := coin.ExtractCoinTx(data)
  if err != nil {
    return err
  }
  return query.FoutputProof(w, tx, height)
}

// mux.Router registrars

func RegisterQueryTx(r *mux.Router) error {
  r.HandleFunc("/tx/{txhash}", doQueryTx).Methods("GET")
  return nil
}

func RegisterQueryTxCoin(r *mux.Router) error {
  r.HandleFunc("/tx/coin/{txhash}", searchTxCoinByHash).Methods("GET")
  return nil
}

// RegisterTx is a convenience function to
// register all the  handlers in this module.
func RegisterTx(r *mux.Router) error {
  funcs := []func(*mux.Router) error{
    RegisterQueryTx,
    RegisterQueryTxCoin,
  }

  for _, fn := range funcs {
    if err := fn(r); err != nil {
      return err
    }
  }
  return nil
}

// End of mux.Router registrars
