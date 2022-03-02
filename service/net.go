package service

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"strings"
)

type RawTransactions struct {
	JSONRPCVersion string   `json:"jsonrpc"`
	Method         string   `json:"method"`
	Id             int      `json:"id"`
	Params         []string `json:"params"`
}

type CallParams struct {
	From  string `json:"from,omitempty"`
	To    string `json:"to"`
	Value string `json:"value"`
	Data  string `json:"data"`
}

type Call struct {
	JSONRPCVersion string       `json:"jsonrpc"`
	Method         string       `json:"method"`
	Id             int          `json:"id"`
	Params         []CallParams `json:"params"`
}

type TxResponse struct {
	JSONRPCVersion string `json:"jsonrpc"`
	Id             int    `json:"id"`
	Result         string `json:"result"`
}

type CallResponse struct {
	ID      int    `json:"id"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    string `json:"data"`
	} `json:"error"`
	Result string `json:"result"`
}

func NewRawTransactions(tx string) *RawTransactions {
	return &RawTransactions{
		JSONRPCVersion: "2.0",
		Method:         "eth_sendRawTransaction",
		Id:             114514,
		Params:         []string{tx},
	}
}

func NewCall(condition TransactionCondition) *Call {
	return &Call{
		JSONRPCVersion: "2.0",
		Method:         "eth_call",
		Id:             0,
		Params: []CallParams{
			{
				From:  condition.From,
				To:    condition.To,
				Value: condition.Value,
				Data:  condition.CallData,
			},
		},
	}
}

func SendTxn(rpc, tx string) string {
	if !strings.HasPrefix(tx, "0x") {
		tx = "0x" + tx
	}
	txn, err := json.Marshal(NewRawTransactions(tx))
	if err != nil {
		log.Println("Failed to Marshal transaction: err", err.Error())
		return ""
	}
	msg, err := http.Post(rpc,
		"application/json",
		bytes.NewBuffer(txn),
	)
	resp := TxResponse{}
	rData, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Println("Failed to parse:", err.Error())
		return ""
	}
	err = json.Unmarshal(rData, &resp)
	if err != nil {
		log.Println("Response parse failed.", err.Error())
		return ""
	}
	return resp.Result
}

func SendCall(rpc string, tx TransactionCondition) (r *CallResponse) {
	txn, err := json.Marshal(NewCall(tx))
	if err != nil {
		log.Println("Failed to Marshal transaction: err", err.Error())
		return
	}
	msg, err := http.Post(rpc,
		"application/json",
		bytes.NewBuffer(txn),
	)
	resp := CallResponse{}
	rData, err := ioutil.ReadAll(msg.Body)
	if err != nil {
		log.Println("Failed to parse:", err.Error())
		return
	}
	err = json.Unmarshal(rData, &resp)
	if err != nil {
		log.Println("Response parse failed.", err.Error())
		return
	}
	return &resp
}

func TestCall(rpc string, tx TransactionCondition) bool {
	// Send Call.
	r := SendCall(rpc, tx)
	if r == nil {
		return false
	}
	if r.Error.Code != 0 {
		return false
	}
	if IsErrorResponse(r.Result) {
		return false
	}
	// Convert Response to big int.
	resultInt := big.NewInt(0)
	resultInt.SetString(r.Result, 0)
	requestInt := big.NewInt(0)
	requestInt.SetString(tx.Expected, 0)
	if tx.ExpOperator == ">" {
		return resultInt.Cmp(requestInt) > 0
	}
	if tx.ExpOperator == ">=" {
		return resultInt.Cmp(requestInt) >= 0
	}
	if tx.ExpOperator == "<" {
		return resultInt.Cmp(requestInt) < 0
	}
	if tx.ExpOperator == "<=" {
		return resultInt.Cmp(requestInt) <= 0
	}
	if tx.ExpOperator == "!" {
		return true
	}
	return resultInt.Cmp(requestInt) == 0
}
