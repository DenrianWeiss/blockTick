package service

import (
	"bytes"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"log"
	"math/big"
	"strings"
)

const ErrorSelector = "0x08c379a0"

type TransactionCondition struct {
	From        string `json:"from,omitempty"`
	To          string `json:"to"`
	Value       string `json:"value,omitempty"`
	CallData    string `json:"call_data"`
	Expected    string `json:"expected"`
	ExpOperator string `json:"exp_operator,omitempty"`
}

func DecodeAndValidateEthTransactions(i []byte) (bool, string) {
	// rlp Decode Transaction
	inputReader := bytes.NewReader(i)
	tx := new(types.Transaction)
	err := rlp.Decode(inputReader, &tx)
	if err != nil {
		log.Println("DecodeAndValidateEthTransactions: failed to rlp decode:", err.Error())
		return false, ""
	}
	msg, err := tx.AsMessage(types.NewEIP155Signer(tx.ChainId()), big.NewInt(1))
	fromStr := msg.From().String()
	return true, fromStr
}

func IsErrorResponse(r string) bool {
	if !strings.HasPrefix(r, "0x") {
		r = "0x" + r
	}
	return strings.HasPrefix(r, ErrorSelector)
}
