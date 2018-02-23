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
	"github.com/ly0129ly/explorer/services/modules/sync"
	"github.com/robfig/cron"
	"github.com/ly0129ly/explorer/services/modules/db"
	"time"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
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
	db.Mgo.Init()
	_,err := db.Mgo.QueryLastedBlock()

	if err != nil {
		//初始化配置表
		tx := db.SyncBlock{
			CurrentPos:1,
			TotalCoinTxs:0,
			TotalStakeTxs:0,
			State:0,
		}
		db.Mgo.Save(tx)
	}
}

func startSync() error {
	mgo := fmt.Sprintf("%s", viper.GetString(MgoUrl))

	log.Printf("Mgo on %s",mgo)

	prepareSync()

	log.Printf("sync Transactions start")
	c := commands.GetNode()

	spec := "@every 60s"
	cron := cron.New()
	cron.AddFunc(spec, func() {
		block,err := db.Mgo.QueryLastedBlock()
		//只有同步完成，才可以执行下一次同步任务
		if block.State == 0 && err == nil{
			Sync(block,c)
		}
	})
	cron.Start()

	return nil
}

func Sync(b db.SyncBlock,c rpcclient.Client){
	//更新状态为同步中
	b.State = 1
	db.Mgo.UpdateBlock(b)

	current := b.CurrentPos
	latest := int64(0)

	log.Printf("sync pos: %d",current)

	for ok := true; ok; ok = current < latest {
		blocks, err := c.BlockchainInfo(current, current+sync.SmallBatchSize)
		if err != nil {
			log.Fatal(err)
		}
		for _, block := range blocks.BlockMetas {
			if (block.Header.NumTxs > 0) {
				txhash := block.Header.DataHash
				txtype, tx := parseTx(txhash,c)
				if (txtype == "coin") {
					coinTx, _ := tx.(db.CoinTx)
					coinTx.TxHash = fmt.Sprintf("%s",txhash)
					coinTx.Time = block.Header.Time
					coinTx.Height = block.Header.Height
					//coinTxs = append([]db.CoinTx{coinTx}, coinTxs...)
					err = db.Mgo.Save(coinTx)
					if(err != nil){
						break
					}
					b.TotalCoinTxs = b.TotalCoinTxs + 1
				} else if (txtype == "stake") {
					stakeTx, _ := tx.(db.StakeTx)
					stakeTx.TxHash = fmt.Sprintf("%s",txhash)
					stakeTx.Time = block.Header.Time
					stakeTx.Height = block.Header.Height
					//stakeTxs = append([]db.StakeTx{stakeTx}, stakeTxs...)
					err = db.Mgo.Save(stakeTx)
					if(err != nil){
						break
					}
					b.TotalStakeTxs = b.TotalStakeTxs + 1
				}
			}
		}
		current = blocks.BlockMetas[0].Header.Height + 1
		latest = blocks.LastHeight
	}
	b.CurrentPos = current
	b.State = 0
	//更新状态为已完成
	db.Mgo.UpdateBlock(b)

	if (current >= latest) {
		time.Sleep(time.Second * 60)
	}
}



//func cmdSync(cmd *cobra.Command, args []string) error {
//	// load current syncing progress from file
//	file := viper.GetString(sync.FlagSyncJson)
//	raw, err := ioutil.ReadFile(file)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var syncResult sync.SyncResult
//	json.Unmarshal(raw, &syncResult)
//
//	for {
//		syncResult = batch(syncResult)
//	}
//
//	return nil
//}



//func batch(syncResult sync.SyncResult) sync.SyncResult {
//	current := syncResult.CurrentPos
//	max := current + sync.LargeBatchSize
//	latest := int64(0)
//	c := commands.GetNode()
//	for ok := true; ok; ok = (current < latest && current < max) {
//		blocks, err := c.BlockchainInfo(current, current+sync.SmallBatchSize)
//		if err != nil {
//			log.Fatal(err)
//		}
//		for _, block := range blocks.BlockMetas {
//			if (block.Header.NumTxs > 0) {
//				txhash := block.Header.DataHash
//				txtype, tx := parseTx(txhash)
//				if (txtype == "coin") {
//					coinTx, _ := tx.(sync.CoinTx)
//					coinTx.TxHash = txhash
//					coinTx.Time = block.Header.Time
//					coinTx.Height = block.Header.Height
//					// prepend
//					syncResult.CoinTxs = append([]sync.CoinTx{coinTx}, syncResult.CoinTxs...)
//					if (len(syncResult.CoinTxs) > sync.MaxRecentSize) {
//						// remove last one
//						syncResult.CoinTxs = syncResult.CoinTxs[:len(syncResult.CoinTxs)-1]
//					}
//					// increase count
//					syncResult.TotalCoinTxs = syncResult.TotalCoinTxs + 1
//				} else if (txtype == "stake") {
//					stakeTx, _ := tx.(sync.StakeTx)
//					stakeTx.TxHash = txhash
//					stakeTx.Time = block.Header.Time
//					stakeTx.Height = block.Header.Height
//					syncResult.StakeTxs = append([]sync.StakeTx{stakeTx}, syncResult.StakeTxs...)
//					// if (len(syncResult.StakeTxs) > sync.MaxRecentSize) {
//					//   syncResult.StakeTxs = syncResult.StakeTxs[:len(syncResult.StakeTxs)-1]
//					// }
//					syncResult.TotalStakeTxs = syncResult.TotalStakeTxs + 1
//				}
//			}
//		}
//		current = blocks.BlockMetas[0].Header.Height + 1
//		latest = blocks.LastHeight
//	}
//
//	// save batch process result into file
//	syncResult.CurrentPos = current
//	json, err := data.ToJSON(syncResult)
//	if err != nil {
//		log.Fatal(err)
//	}
//	file := viper.GetString(sync.FlagSyncJson)
//	ioutil.WriteFile(file, json, 0644)
//	fmt.Printf("%d scanned\n", current)
//
//	// stop if it's latest block
//	if (current >= latest) {
//		time.Sleep(time.Second * 60)
//	}
//
//	return syncResult
//}

func parseTx(bkey []byte,client rpcclient.Client) (string, interface{}) {
	// load tx by hash
	prove := !viper.GetBool(commands.FlagTrustNode)
	//client := commands.GetNode()
	res, err := client.Tx(bkey, prove)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := sdk.LoadTx(res.Proof.Data)
	if err != nil {
		log.Fatal(err)
	}

	// parse
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
