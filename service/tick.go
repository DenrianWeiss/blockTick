package service

import (
	"sync"
	"time"
)

func TimedTransition(rpc, data string, trigger int64) {
	go func() {
		time.Sleep(time.Until(time.Unix(trigger - 1, 800000000)))
		SendTxn(rpc, data)
	}()
}

func ConditionalTransition(rpc, data string, trigger TransactionCondition, signal *sync.Cond) {
	go func() {
		for !TestCall(rpc, trigger) {
			signal.Wait()
		}
		SendTxn(rpc, data)
	}()
}
