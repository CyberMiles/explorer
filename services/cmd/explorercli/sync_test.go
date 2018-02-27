package main

import (
	"testing"
	"github.com/cosmos/cosmos-sdk/client"
	"time"
	"log"
	"github.com/ly0129ly/explorer/services/modules/db"
)

func TestProcessSync(t *testing.T) {
	c := client.GetNode("tcp://116.62.62.39:46657")
	db.Mgo.Init("localhost:27017")
	processSync(c)

	log.Printf(" finish %s","ok")
	time.Sleep(5 * time.Minute)
}