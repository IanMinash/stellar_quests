## Challenge 3: Store some arbitrary data in your Stellar account

Stellar allows you to store arbitrary data in the form of Key / Value pairs. Add a Key of Hello and a Value of World as a data attribute on your account using `ManageData` operation.

_Note:_ To delete an entry from the account issue the `ManageData` operation without providing any value for the entry.

### Requirements

Create a `.env` file in the current directory and add the following:

```
SIGN_KEY=<secret key of your stellar account>
```
