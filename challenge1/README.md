## Challenge 1: Create and fund a Stellar account

Create an account on the testnet and fund it with 1000XLM

### Requirements

For this implementation, you'll need to have an initial account with some Stellar Lumens. You can create one from the [Stellar Laboratory](https://laboratory.stellar.org/#account-creator?network=test) and use the Friendbot to fund it with the initial balance. This account will be used to issue the `CreateAccount` operation.

Create a `.env` file in the current directory and add the following:

```
ACC_ID=<public key of the account to be created>
SIGN_KEY=<secret key of the account created on the Stellar Laboratory>
```
