## Challenge 5: Create a custom asset and send it to your account

The ability to issue assets is a core feature of Stellar. Any asset can be tokenized, and, once tokenized, transferred or traded over the Stellar network quickly and cheaply.

_Note_: In production environments, its better practice to have a separate distribution account as a proxy instead of issuing tokens directly from the issuer recipients. See [Stellar Docs](https://developers.stellar.org/docs/issuing-assets/how-to-issue-an-asset/#why-have-separate-accounts-for-issuing-and-distribution)

### Requirements

Create a `.env` file in the current directory and add the following:

```
SIGN_KEY=<secret key of your stellar account>
```

Save the generated key if you intend to use the issuing account in the future.
