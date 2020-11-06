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

	// Get private key of account
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

	// Generate another signature
	otherKp := keypair.MustRandom()
	fmt.Printf("Other keypair: \n\tAddress:\t%s\n\tSecret Key:\t%s\n", otherKp.Address(), otherKp.Seed())

	// SetOptions operation to add another signer.
	op := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: otherKp.Address(),
			Weight:  txnbuild.Threshold(1),
		},
	}

	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)

	if err != nil {
		log.Printf("%s\n", err)
	}

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full))

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("SetOptions Transaction ID:\t%s\n", resp.ID)

	seqNum, _ := sourceAccount.GetSequenceNumber()

	// Create a new transaction
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.BumpSequence{
					BumpTo: seqNum + 50,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimeout(300),
		},
	)

	// Sign the transaction using the generated signature
	tx, err = tx.Sign(network.TestNetworkPassphrase, otherKp)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err = client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("BumpSequence Transaction ID:\t%s\n", resp.ID)
}
