package main

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/ly0129ly/explorer/services/modules/db"
)

func TestBeginSync(t *testing.T) {

	block,_ := db.Mgo.QueryLastedBlock()
	c := client.GetNode("http://116.62.62.39:46657")
	Sync(block,c)
}