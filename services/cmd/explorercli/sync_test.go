package main

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/client"
	"fmt"

	"github.com/cosmos/cosmos-sdk"

	"github.com/ly0129ly/explorer/services/modules/db"
	"strings"
	"time"
	"context"
	"github.com/tendermint/tmlibs/pubsub/query"
	//"github.com/tendermint/types"
	"log"
	"github.com/tendermint/tendermint/types"
	"encoding/hex"
)

func TestBeginSync(t *testing.T) {
	c := client.GetNode("tcp://116.62.62.39:46657")
	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	query := query.MustParse("tm.event = 'Tx'")
	txs := make(chan interface{})
	//
	c.Start()
	err := c.Subscribe(ctx, "tx-watch", query, txs)
	//
	if err != nil{
		fmt.Println("got ", err)
	}
	go func(){
		for e := range txs {
			deliverTxRes := e.(types.TMEventData).Unwrap().(types.EventDataTx)
			txb, _ := sdk.LoadTx(deliverTxRes.Tx)
			_, tx :=parseTx(txb)
			log.Printf("from=%s",tx.(db.CoinTx).From)
			log.Printf("tx_hash=%s",strings.ToUpper(hex.EncodeToString(deliverTxRes.Tx.Hash())))
		}
	}()

	log.Printf(" finish %s","ok")
	time.Sleep(5 * time.Minute)
}

func TestPaseTx(t *testing.T)  {
	fmt.Printf("hellow")
}