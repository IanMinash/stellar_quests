## Challenge 4: Add multisig to your account and make use of it in a transaction

All Stellar transactions need to have atleaast 1 signature which is usually the public key of the Stellar account. Stellar allows for additional signing keys to be added to an account using the `SetOptions` operation. Extra signing keys should be added with a specific weight which determines the types of transactions that the key can be used to authorize.

An example application of multisig accounts is having a 'company account' where you might want to give some employees the ability to authorize transactions of that account.

_Note:_ You might have to save the generated keys to avoid being locked out of the account depending on how you configure your threshold levels.

### Requirements

Create a `.env` file in the current directory and add the following:

```
SIGN_KEY=<secret key of your stellar account>
```
