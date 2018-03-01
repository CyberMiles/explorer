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
	"encoding/hex"
	"github.com/tendermint/tendermint/types"
	"github.com/spf13/cast"
)

var (
	syncCmd = &cobra.Command{
		Use:  "sync",
		Long: `sync`,
		RunE: func(cmd *cobra.Command, args []string) error {
			startWatch()
			return nil
		},
	}
)

func prepareSync(){
	url := viper.GetString(MgoUrl)
	db.Mgo.Init(url)
	block,err := db.Mgo.QueryLastedBlock()

	if err != nil {
		//初始化配置表
		tx := db.SyncBlock{
			CurrentPos:1,
			TotalCoinTxs:0,
			TotalStakeTxs:0,
		}
		db.Mgo.Save(tx)
		return
	}

	//开始漏单查询
	go sync(block)
}

func startWatch() error {
	c := commands.GetNode()
	prepareSync()


	processSync(c)
	return nil
}

func processSync(c rpcclient.Client){
	log.Printf("watched Transactions start")

	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	//query := query.MustParse("tm.event = 'Tx'")
	txs := make(chan interface{})

	c.Start()
	err := c.Subscribe(ctx, "tx-watch", types.EventQueryTx, txs)

	if err != nil{
		fmt.Println("got ", err)
	}

	go func() {
		log.Println("listening tx begin")
		for e := range txs {
			block,_ := db.Mgo.QueryLastedBlock()
			deliverTxRes := e.(types.TMEventData).Unwrap().(types.EventDataTx)
			height := deliverTxRes.Height

			txb, _ := sdk.LoadTx(deliverTxRes.Tx)
			txtype, tx :=parseTx(txb)
			if (txtype == "coin") {
				coinTx, _ := tx.(db.CoinTx)

				coinTx.TxHash = strings.ToUpper(hex.EncodeToString(deliverTxRes.Tx.Hash()))
				coinTx.Time = queryBlockTime(c,height)
				coinTx.Height = height
				if db.Mgo.Save(coinTx) != nil {
					break
				}
				block.TotalCoinTxs += 1

				log.Printf("watched coinTx,tx_hash=%s",coinTx.TxHash)
			} else if (txtype == "stake") {
				stakeTx, _ := tx.(db.StakeTx)
				stakeTx.TxHash = strings.ToUpper(hex.EncodeToString(deliverTxRes.Tx.Hash()))
				stakeTx.Time = queryBlockTime(c,height)
				stakeTx.Height = height
				if db.Mgo.Save(stakeTx) != nil {
					break
				}
				block.TotalStakeTxs += 1
				log.Printf("watched stakeTx,tx_hash=%s",stakeTx.TxHash)
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

func sync(curBlock db.SyncBlock) {
	log.Printf("sync Transactions start")
	c := commands.GetNode()

	current := curBlock.CurrentPos
	latest := int64(0)

	log.Printf("last block heigth：%d",current)

	for ok := true; ok; ok = current < latest {
		blocks, err := c.BlockchainInfo(current, current + 20)
		if err != nil {
			log.Fatal(err)
		}
		for _, block := range blocks.BlockMetas {
			if block.Header.NumTxs > 0 {
				txhash := block.Header.DataHash
				prove := !viper.GetBool(commands.FlagTrustNode)
				res, _ := c.Tx(txhash, prove)
				txs, _ := sdk.LoadTx(res.Proof.Data)
				txtype, tx := parseTx(txs)
				if  txtype == "coin" {
					coinTx, _ := tx.(db.CoinTx)
					coinTx.TxHash = strings.ToUpper(fmt.Sprintf("%s",txhash))
					coinTx.Time = block.Header.Time
					coinTx.Height = block.Header.Height
					if db.Mgo.Save(coinTx) == nil {
						curBlock.TotalCoinTxs += 1
						log.Printf("sync coinTx,tx_hash=%s",coinTx.TxHash)
					}

				} else if txtype == "stake" {
					stakeTx, _ := tx.(db.StakeTx)
					stakeTx.TxHash = strings.ToUpper(fmt.Sprintf("%s",txhash))
					stakeTx.Time = block.Header.Time
					stakeTx.Height = block.Header.Height
					if db.Mgo.Save(stakeTx) == nil {
						curBlock.TotalStakeTxs += 1
						log.Printf("sync stakeTx,tx_hash=%s",stakeTx.TxHash)
					}
				}
			}
		}
		current = blocks.BlockMetas[0].Header.Height + 1
		latest = blocks.LastHeight
	}

	curBlock.CurrentPos = current
	db.Mgo.UpdateBlock(curBlock)
	log.Printf("sync Transactions end,current block height:%d",current)
}