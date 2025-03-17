package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/troykessler/hyperlane-cosmos/util"
	warpTypes "github.com/troykessler/hyperlane-cosmos/x/warp/types"
)

func Decode(messageStr string) error {
	messageBytes, err := util.DecodeEthHex(messageStr)
	if err != nil {
		return err
	}

	message, err := util.ParseHyperlaneMessage(messageBytes)
	if err != nil {
		return err
	}

	fmt.Println("### Message ###")
	fmt.Printf("Version: \t%d\n", message.Version)
	fmt.Printf("Nonce:  \t%d\n", message.Nonce)
	fmt.Printf("Origin: \t%d\n", message.Origin)
	fmt.Printf("Sender: \t%s\n", message.Sender.String())
	fmt.Printf("Destination: \t%d\n", message.Destination)
	fmt.Printf("Recipient: \t%s\n", message.Recipient.String())
	fmt.Printf("Body: \t\t0x%s\n", hex.EncodeToString(message.Body))

	return nil
}

func GenerateWarpTransfer(senderContract string, recipientContract string, recipientUser string, amount uint64) error {
	var bz []byte
	if strings.HasPrefix(recipientUser, "kyve") {
		dbz, err := sdk.GetFromBech32(recipientUser, "kyve")
		if err != nil {
			panic(err)
		}
		bz = dbz
	} else if strings.HasPrefix(recipientUser, "0x") && len(recipientUser) == 66 {
		dbz, err := util.DecodeHexAddress(recipientUser)
		if err != nil {
			panic(err)
		}
		bz = dbz.Bytes()
	}

	payload, err := warpTypes.NewWarpPayload(bz, *big.NewInt(int64(amount)))
	if err != nil {
		panic(err)
	}

	recipient, err := util.DecodeHexAddress(recipientContract)
	if err != nil {
		panic(err)
	}

	sender, err := util.DecodeHexAddress(senderContract)
	if err != nil {
		panic(err)
	}

	msg := util.HyperlaneMessage{
		Version:     1,
		Nonce:       3,
		Origin:      1,
		Sender:      sender,
		Destination: 0,
		Recipient:   recipient,
		Body:        payload.Bytes(),
	}

	fmt.Println(msg.String())

	return nil
}
