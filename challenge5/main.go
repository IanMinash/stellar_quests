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

	// Create keypair for issuing account
	issuerKp := keypair.MustRandom()
	fmt.Printf("Issuer keypair: \n\tAddress:\t%s\n\tSecret Key:\t%s\nPlease store these keys.\n", issuerKp.Address(), issuerKp.Seed())

	// Create the issuing account by funding it with the testnet's friendbot
	client.Fund(issuerKp.Address())

	// Our custom token
	thaToken := txnbuild.CreditAsset{
		Code:   "THA",
		Issuer: issuerKp.Address(),
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
					Line: thaToken,
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

	fmt.Printf("Transaction ID:\t%s\n", resp.ID)

	// Get Issuer account
	issuerRequest := horizonclient.AccountRequest{AccountID: issuerKp.Address()}
	issuerAccount, err := client.AccountDetail(issuerRequest)
	if err != nil {
		log.Fatalln(err)
	}

	// Send the new tokens from the issuer to the account
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &issuerAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Destination: kp.Address(),
					Amount:      "1",
					Asset:       thaToken,
				},
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimeout(300),
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, issuerKp)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err = client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Transaction ID:\t%s\n", resp.ID)
}
