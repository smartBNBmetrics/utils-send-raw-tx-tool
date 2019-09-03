package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stellar/go/clients/horizon"
	"os"
	"context"
	"encoding/hex"
)

type tx struct {
	RawTx string `json:"rawTx"`
}

type sendRawTx func(string) (string, error)

func Send(c *gin.Context) {

	var (
		tx tx
		send sendRawTx
		currency = c.Param("currency")
	)

	err := c.BindJSON(&tx)
	if err != nil{
		c.JSON(404, "bad request")
		return
	}

	switch currency {
	case "eth":
		send = sendETH
	case "etc":
		send = sendETC
	case "xlm":
		send = sendXLM
	}

	hash, err := send(tx.RawTx)
	if err != nil{
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"hash": hash})
}


func sendETH(rawTx string) (string, error) {

	client, err := ethclient.Dial(os.Getenv("ETH"))
	if err != nil {
		return "", err
	}

	rawTxBytes, err := hex.DecodeString(rawTx)

	tx := new(types.Transaction)

	rlp.DecodeBytes(rawTxBytes, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func sendETC(rawTx string) (string, error) {

	client, err := ethclient.Dial(os.Getenv("ETC"))
	if err != nil {
		return "", err
	}

	rawTxBytes, err := hex.DecodeString(rawTx)

	tx := new(types.Transaction)

	rlp.DecodeBytes(rawTxBytes, &tx)

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func sendXLM(rawTx string)(string, error){
	resp, err := horizon.DefaultPublicNetClient.SubmitTransaction(rawTx)
	if err != nil {
		return "", err
	}

	return resp.Hash, nil
}