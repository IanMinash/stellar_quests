## Challenge 7: Make use of a channel account to make a payment

Payment channels provide a method for submitting transactions to the network at a high rate. A channel is another Stellar account that is used as the source account of the transaction. Channels can allow applications to soak up fees for their users or achieve higher transaction throughput.

### Requirements

Create a `.env` file in the current directory and add the following:

```
SIGN_KEY=<secret key of your stellar account>
RECV_ACC_ID=<public key of the account to receive the payment>
```

Save the generated key if you intend to use the channel account in the future.
