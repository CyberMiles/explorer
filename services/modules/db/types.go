package db

import (
	"time"
	"github.com/cosmos/cosmos-sdk/modules/coin"
)

const (
	DbCosmosTxn = "cosmos_txn"
	TbNmCoinTx = "coin_tx"
	TbNmStakeTx = "stake_tx"
	TbNmSyncBlock = "sync_block"
	PageSize = 10
)

type TxHander interface {
	TbNm() string
}

type CoinTx struct {
	TxHash string `json:"tx_hash"`
	Time time.Time `json:"time"`
	Height int64 `json:"height"`
	From string `json:"from"`
	To string `json:"to"`
	Amount coin.Coins `json:"coins"`
}

type StakeTx struct {
	TxHash string `json:"tx_hash"`
	Time time.Time `json:"time"`
	Height int64 `json:"height"`
	From string `json:"from"`
	Type string `json:"type"`
	Amount coin.Coin `json:"amount"`
}

type SyncBlock struct {
	CurrentPos int64 `json:"current_pos"`
	TotalCoinTxs int64 `json:"total_coin_txs"`
	TotalStakeTxs int64 `json:"total_stake_txs"`
	State int64 //0：同步完成，1：同步中
}

func(c CoinTx) TbNm() string{
	return TbNmCoinTx
}

func(c StakeTx) TbNm() string{
	return TbNmStakeTx
}

func(c SyncBlock) TbNm() string{
	return TbNmSyncBlock
}