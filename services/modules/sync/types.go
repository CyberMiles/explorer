package sync

import (
  "time"

  "github.com/cosmos/cosmos-sdk/modules/coin"
  "github.com/tendermint/go-wire/data"
)

const (
  FlagSyncJson = "sync-json"
  SmallBatchSize = 20
  LargeBatchSize = 5000
  MaxRecentSize = 20
)

type SyncResult struct {
  CurrentPos int64 `json:"current_pos"`
  TotalCoinTxs int64 `json:"total_coin_txs"`
  TotalStakeTxs int64 `json:"total_stake_txs"`
  CoinTxs []CoinTx `json:"coin_txs"`
  StakeTxs []StakeTx `json:"stake_txs"`
}

type CoinTx struct {
  TxHash data.Bytes `json:"tx_hash"`
  Time time.Time `json:"time"`
  Height int64 `json:"height"`
  From data.Bytes `json:"from"`
  To data.Bytes `json:"to"`
}

type StakeTx struct {
  TxHash data.Bytes `json:"tx_hash"` 
  Time time.Time `json:"time"`
  Height int64 `json:"height"`
  From data.Bytes `json:"from"`
  Type string `json:"type"`
  Amount coin.Coin `json:"amount"`
}
