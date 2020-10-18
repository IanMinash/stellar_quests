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

	signKey, signKeyPresent := os.LookupEnv("SIGN_KEY")
	if !signKeyPresent {
		log.Fatalln("SIGN_KEY is not defined in the environment file. Please define it and try again")
	}

	newAccountID, newAccountIDPresent := os.LookupEnv("ACC_ID")
	if !newAccountIDPresent {
		log.Fatalln("ACC_ID is not defined in the environment file. Please define it and try again")
	}

	kp := keypair.MustParse(signKey)

	client := horizonclient.DefaultTestNetClient

	request := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}

	op := txnbuild.CreateAccount{
		Destination: newAccountID,
		Amount:      "1000",
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			IncrementSequenceNum: true,
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	txEnvelope, err := tx.Base64()
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransactionXDR(txEnvelope)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(resp.Account)
}
