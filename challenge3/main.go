package main

import (
	"encoding/json"
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

	// Get private key of account to send XLM from
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

	dataJSON, err := json.MarshalIndent(sourceAccount.Data, "", "\t")
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		fmt.Printf("Initial data entries for %s\n", sourceAccount.AccountID)
		fmt.Printf("%s\n", dataJSON)
	}

	op := txnbuild.ManageData{
		Name:  "Hello",
		Value: []byte("World"),
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

	fmt.Printf("Transaction ID:\t%s\n", resp.ID)
	sourceAccount, err = client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}
	dataJSON, err = json.MarshalIndent(sourceAccount.Data, "", "\t")
	if err != nil {
		log.Printf("%s\n", err)
	} else {
		fmt.Printf("Current data entries for %s\n", sourceAccount.AccountID)
		fmt.Printf("%s\n", dataJSON)
	}
}
