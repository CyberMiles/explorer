package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/modules/coin"
	"github.com/cosmos/cosmos-sdk/modules/nonce"
	"github.com/cosmos/cosmos-sdk/client/commands"
	"github.com/tendermint/go-wire/data"
	"github.com/ly0129ly/explorer/services/modules/stake"
	"github.com/ly0129ly/explorer/services/modules/db"
	"time"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"context"
	"github.com/tendermint/tmlibs/pubsub/query"
	"encoding/hex"
	"github.com/tendermint/tendermint/types"
	"github.com/spf13/cast"
)

var (
	syncCmd = &cobra.Command{
		Use:  "sync",
		Long: `sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			startSync()
			return nil
		},
	}
)

func prepareSync(){
	url := viper.GetString(MgoUrl)
	db.Mgo.Init(url)
	_,err := db.Mgo.QueryLastedBlock()

	if err != nil {
		//初始化配置表
		tx := db.SyncBlock{
			CurrentPos:1,
			TotalCoinTxs:0,
			TotalStakeTxs:0,
		}
		db.Mgo.Save(tx)
	}
}

func startSync() error {
	prepareSync()
	c := commands.GetNode()
	processSync(c)
	return nil
}

func processSync(c rpcclient.Client){
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	query := query.MustParse("tm.event = 'Tx'")
	txs := make(chan interface{})

	c.Start()
	err := c.Subscribe(ctx, "tx-watch", query, txs)

	if err != nil{
		fmt.Println("got ", err)
	}

	go func() {
		log.Println("listening tx begin")
		for e := range txs {
			block,err := db.Mgo.QueryLastedBlock()
			deliverTxRes := e.(types.TMEventData).Unwrap().(types.EventDataTx)
			height := deliverTxRes.Height

			txb, _ := sdk.LoadTx(deliverTxRes.Tx)
			txtype, tx :=parseTx(txb)
			if (txtype == "coin") {
				coinTx, _ := tx.(db.CoinTx)

				coinTx.TxHash = strings.ToUpper(hex.EncodeToString(deliverTxRes.Tx.Hash()))
				coinTx.Time = queryBlockTime(c,height)
				coinTx.Height = height
				err = db.Mgo.Save(coinTx)
				if(err != nil){
					break
				}
				block.TotalCoinTxs += 1

				log.Printf("watched coin tx,tx_hash=%s",coinTx.TxHash)
			} else if (txtype == "stake") {
				stakeTx, _ := tx.(db.StakeTx)
				stakeTx.TxHash = strings.ToUpper(hex.EncodeToString(deliverTxRes.Tx.Hash()))
				stakeTx.Time = queryBlockTime(c,height)
				stakeTx.Height = height
				err = db.Mgo.Save(stakeTx)
				if(err != nil){
					break
				}
				block.TotalStakeTxs += 1
				log.Printf("watched stake tx,tx_hash=%s",stakeTx.TxHash)
			}
			block.CurrentPos = height
			db.Mgo.UpdateBlock(block)
		}
	}()
}

func parseTx(tx sdk.Tx) (string, interface{}){
	txl, ok := tx.Unwrap().(sdk.TxLayer)
	var txi sdk.Tx
	var coinTx db.CoinTx
	var stakeTx db.StakeTx
	var nonceAddr data.Bytes
	for ok {
		txi = txl.Next()
		switch txi.Unwrap().(type) {
		case coin.SendTx:
			ctx, _ := txi.Unwrap().(coin.SendTx)
			coinTx.From = fmt.Sprintf("%s",ctx.Inputs[0].Address.Address)
			coinTx.To = fmt.Sprintf("%s",ctx.Outputs[0].Address.Address)
			coinTx.Amount = ctx.Inputs[0].Coins
			return "coin", coinTx
		case nonce.Tx:
			ctx, _ := txi.Unwrap().(nonce.Tx)
			nonceAddr = ctx.Signers[0].Address
			break
		case stake.TxUnbond, stake.TxDelegate, stake.TxDeclareCandidacy:
			kind, _ := txi.GetKind()
			stakeTx.From = fmt.Sprintf("%s",nonceAddr)
			stakeTx.Type = strings.Replace(kind, "stake/", "", -1)
			switch kind {
			case "stake/unbond":
				ctx, _ := txi.Unwrap().(stake.TxUnbond)
				stakeTx.Amount.Denom = "fermion"
				stakeTx.Amount.Amount = int64(ctx.Shares)
				break
			case "stake/delegate":
				ctx, _ := txi.Unwrap().(stake.TxDelegate)
				stakeTx.Amount.Denom = ctx.Bond.Denom
				stakeTx.Amount.Amount = ctx.Bond.Amount
				break
			case "stake/declareCandidacy":
				ctx, _ := txi.Unwrap().(stake.TxDeclareCandidacy)
				stakeTx.Amount.Denom = ctx.BondUpdate.Bond.Denom
				stakeTx.Amount.Amount = ctx.BondUpdate.Bond.Amount
				break
			}
			return "stake", stakeTx
		}
		txl, ok = txi.Unwrap().(sdk.TxLayer)
	}
	return "", nil
}

func queryBlockTime(c rpcclient.Client,height int64) time.Time{
	h := cast.ToInt64(height)
	block, err := c.Block(&h)
	if err != nil {
		log.Printf("query block fail ,%d",height)
	}
	return block.BlockMeta.Header.Time
}