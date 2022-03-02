package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/DenrianWeiss/blockTick/config"
	"github.com/DenrianWeiss/blockTick/service"
	"github.com/DenrianWeiss/blockTick/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ScheduledTask struct {
	Rpc         string                       `json:"rpc"`
	Transaction string                       `json:"transaction"`
	Pow         string                       `json:"pow"`
	Type        string                       `json:"type"`
	Secret      string                       `json:"secret"`
	Condition   service.TransactionCondition `json:"condition"`
}

func VerifyPow(difficulty int, address, pow string) bool {
	if difficulty >= 6 {
		difficulty = 6 // Cap difficulty
	}
	powSha := sha256.Sum256([]byte(address + pow))
	for i := 0; i < difficulty; i++ {
		if powSha[i] != 0 {
			return false
		}
	}
	return true
}

func ReceiveTransaction(c *gin.Context) {
	t := ScheduledTask{}
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"reason": "json",
		})
		return
	}
	if config.Config.EnableSecret {
		if t.Secret != config.Config.Secret {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"reason": "secret",
			})
			return
		}
	}
	payload := strings.TrimPrefix(t.Transaction, "0x")
	payloadBytes, err := hex.DecodeString(payload)
	success, addr := service.DecodeAndValidateEthTransactions(payloadBytes)
	addr = utils.AddPrefix(strings.ToLower(addr))

	if !success {
		c.JSON(http.StatusOK, gin.H{
			"status": "error",
			"reason": "tx",
		})
		return
	}
	if config.Config.EnablePow {
		if !VerifyPow(config.Config.PowDifficulty, addr, t.Pow) {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"reason": "pow",
			})
			return
		}
	}
	if config.Config.EnableWhitelist {
		if !config.Config.Whitelist[addr] {
			c.JSON(http.StatusOK, gin.H{
				"status": "error",
				"reason": "whitelist",
			})
			return
		}
	}

}

func MainPageHandler(c *gin.Context) {
	c.String(http.StatusOK, "Visit github.com/DenrianWeiss/blockTick")
}
