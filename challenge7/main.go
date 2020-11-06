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

	// Get public key of account that we'll send XLM to
	recvAccountID, recvAccountIDPresent := os.LookupEnv("RECV_ACC_ID")
	if !recvAccountIDPresent {
		log.Fatalln("RECV_ACC_ID is not defined in the environment file. Please define it and try again")
	}

	kp := keypair.MustParse(signKey)

	channelKp := keypair.MustRandom()

	client := horizonclient.DefaultTestNetClient

	client.Fund(channelKp.Address())
	fmt.Printf("Channel account %s created successfully. Secret key: %s.\n", channelKp.Address(), channelKp.Seed())

	request := horizonclient.AccountRequest{AccountID: kp.Address()}
	sourceAccount, err := client.AccountDetail(request)
	if err != nil {
		log.Fatalln(err)
	}

	channelRequest := horizonclient.AccountRequest{AccountID: channelKp.Address()}
	channelAccount, err := client.AccountDetail(channelRequest)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a trustline from the account to the issuer account.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &channelAccount,
			IncrementSequenceNum: true,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
			Operations: []txnbuild.Operation{
				&txnbuild.Payment{
					Asset:         txnbuild.NativeAsset{},
					Amount:        "100",
					Destination:   recvAccountID,
					SourceAccount: &sourceAccount,
				},
			},
		},
	)

	tx, err = tx.Sign(network.TestNetworkPassphrase, kp.(*keypair.Full), channelKp)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Channeled Payment Transaction ID:\t%s\n", resp.ID)
}
