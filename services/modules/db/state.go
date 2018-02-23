package db

import (
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/spf13/viper"
)
type MgoBackend struct {}


var Mgo = MgoBackend{

}
var url = viper.GetString("mgo-url")
var session, err = mgo.Dial(url)

func (MgoBackend) Init(){
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	index()
}

func index(){
	c := session.DB(DbCosmosTxn).C(TbNmCoinTx)

	index := mgo.Index{
		Key:        []string{"from"}, // 索引字段， 默认升序,若需降序在字段前加-
		Unique:     false,             // 唯一索引 同mysql唯一索引
		DropDups:   false,             // 索引重复替换旧文档,Unique为true时失效
		Background: true,             // 后台创建索引
	}

	c.EnsureIndex(index)

	c = session.DB(DbCosmosTxn).C(TbNmStakeTx)
	c.EnsureIndex(index)
}

func (MgoBackend) Save(tx TxHander) error{
	c := session.DB(DbCosmosTxn).C(tx.TbNm())
	err = c.Insert(tx)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (MgoBackend) QueryCoinTxs() ([]CoinTx)  {
	result := []CoinTx{}
	c := session.DB(DbCosmosTxn).C(TbNmCoinTx)
	err = c.Find(nil).Sort("-time").All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (MgoBackend) QueryStakeTxs() ([]StakeTx)  {
	result := []StakeTx{}
	c := session.DB(DbCosmosTxn).C(TbNmStakeTx)
	err = c.Find(nil).Sort("-time").All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (MgoBackend) QueryCoinTxsByFrom(from string) ([]CoinTx)  {
	result := []CoinTx{}
	c := session.DB(DbCosmosTxn).C(TbNmCoinTx)
	err = c.Find(bson.M{"from": from}).Sort("-time").All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (MgoBackend) QueryStakeTxsByFrom(from string) ([]StakeTx)  {
	result := []StakeTx{}
	c := session.DB(DbCosmosTxn).C(TbNmStakeTx)
	err = c.Find(bson.M{"from": from}).Sort("-time").All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (MgoBackend) QueryPageCoinTxsByFrom(from string,page int)([]CoinTx)  {
	result := []CoinTx{}
	c := session.DB(DbCosmosTxn).C(TbNmCoinTx)
	skip := (page-1) * PageSize
	err = c.Find(bson.M{"from": from}).Sort("-time").Skip(skip).Limit(PageSize).All(&result)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func (MgoBackend) QueryLastedBlock()(SyncBlock,error){
	result := SyncBlock{}
	c := session.DB(DbCosmosTxn).C(TbNmSyncBlock)
	err = c.Find(bson.M{}).One(&result)
	return result,err
}

func (MgoBackend) UpdateBlock(b SyncBlock) error{
	c := session.DB(DbCosmosTxn).C(TbNmSyncBlock)
	return c.Update(nil,b)
}