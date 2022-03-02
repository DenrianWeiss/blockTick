package service

import (
	"context"
	"github.com/DenrianWeiss/blockTick/config"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"sync"
)

var conn *ethclient.Client
var broadcaster *sync.Cond

func GetConn() *ethclient.Client {
	return conn
}

func GetNewBlockLock() *sync.Cond {
	return broadcaster
}

func NewBlockListener() {
	headers := make(chan *types.Header)
	sub, err := conn.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case err := <-sub.Err():
			log.Println(err)
			sub, err = conn.SubscribeNewHead(context.Background(), headers)
			if err != nil {
				log.Fatal(err)
			}
		case h := <-headers:
			log.Println("New Block: ", h.Number.String())
			broadcaster.Broadcast()
		}
	}
}

func InitAll() {
	log.Println("initializing.")
	c, err := ethclient.Dial(config.Config.DefaultRpc)
	if err != nil {
		panic(err)
	}
	conn = c
	broadcaster = new(sync.Cond)
	go NewBlockListener()
}
