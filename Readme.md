
# Sign In With Solana

The primary purpose of this document is to define how Solana accounts authenticate with off-chain services. By signing a standard message format parameterized by scope, session details, and a nonce.

While decentralized identity is not a novel concept, the most common implementations of blockchain-based credentials are either certificate-based or rely on centralized providers. We're proposing an alternative that doesn't require a trusted third party.

### Assumptions

Currently, there is no propsed standard for signing messages on Solana. We are proposing `SIP-99` as a placeholder to conform to CAIP-74.

### Workflow

1. The user connects the wallet to the website.
2. From the frontend pass the domain, address, statement, uri, version, nonce, issuedAt, expirationTime, notBefore, requestId, resources (Array) to the SignInWithSolanaMessage constructor. There is additional regex validation in place as mentioned in the below sections
3. Nonce is needed as a security mechanism from replay attacks and hence it is generated at the server side.
4. The created message needs to be prepared in a wallet friendly format for which <message>.prepareMessage() needs to be called
5. The resultant has to be passed to signMessage method of window.solana.request
6. This function would return the signedMessage


#### Disclaimer :
We haven't undergone a Formal Security Audit yet.
