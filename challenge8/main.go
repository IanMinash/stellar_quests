package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stellar/go/txnbuild"

	"github.com/joho/godotenv"
	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("An error occured while trying to read from .env")
	}

	// Get private key of receiving account
	signKey, signKeyPresent := os.LookupEnv("SIGN_KEY")
	if !signKeyPresent {
		log.Fatalln("SIGN_KEY is not defined in the environment file. Please define it and try again")
	}

	kp := keypair.MustParse(signKey)

	client := horizonclient.DefaultTestNetClient

	request := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}

	// SRT Asset
	srtAsset := txnbuild.CreditAsset{
		Code:   "SRT",
		Issuer: "GCDNJUBQSX7AJWLJACMJ7I4BC3Z47BQUTMHEICZLE6MU4KQBRYG5JY6B",
	}

	// Create a trustline from the account to the issuer account.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			Operations: []txnbuild.Operation{
				&txnbuild.ChangeTrust{
					Line: srtAsset,
				},
				&txnbuild.PathPaymentStrictSend{
					SendAsset:   txnbuild.NativeAsset{},
					SendAmount:  "150",
					Destination: kp.Address(),
					DestAsset:   srtAsset,
					DestMin:     "15",
				},
			},
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Path Payment Transaction ID:\t%s\n", resp.ID)
}
