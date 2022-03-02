package service

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestDecodeTx(t *testing.T) {
	txn := "f86c0a85046c7cfe0083016dea94d1310c1e038bc12865d3d3997275b3e" +
		"4737c6302880b503be34d9fe80080269fc7eaaa9c21f59adf8ad43ed66cf5ef9" +
		"ee1c317bd4d32cd65401e7aaca47cfaa0387d79c65b90be6260d09dcfb780f29" +
		"dd8133b9b1ceb20b83b7e442b4bfc30cb"
	txnBytes, _ := hex.DecodeString(txn)
	_, sender := DecodeAndValidateEthTransactions(txnBytes)
	if strings.ToLower(sender) != "0x67835910d32600471f388a137bbff3eb07993c04" {
		panic("Decode function failed.")
	}
}
