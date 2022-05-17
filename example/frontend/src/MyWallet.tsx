import {
    useConnection,
    useWallet
} from '@solana/wallet-adapter-react';
import {
    WalletDisconnectButton, WalletModalProvider, WalletMultiButton
} from '@solana/wallet-adapter-react-ui';
import bs58 from 'bs58';
import React, { useState } from 'react';
import Swal from 'sweetalert2';
import SolanaLogo from '../public/solana-logo.png';

const MyWallet: React.FC = () => {
    const { connection } = useConnection();
    let walletAddress = "";

    // if you use anchor, use the anchor hook instead
    // const wallet = useAnchorWallet();
    // const walletAddress = wallet?.publicKey.toString();

    const wallet = useWallet();
    if (wallet.connected && wallet.publicKey) {
        walletAddress = wallet.publicKey.toString()
    }

    const { publicKey, signMessage } = useWallet();
    
    // Domain and origin
    const domain = window.location.host;
    const origin = window.location.origin;

    
    let statement = "Sign in with Solana to the app.";

    const [siwsMessage, setSiwsMessage] = useState({});
    const [nonce, setNonce] = useState("");
    const [sign, setSignature] = useState("");

    // Generate a message for signing
    // The nonce is generated on the server side 
    function createSolanaMessage() {

        var raw = JSON.stringify({

                "domain": domain,
                "address": publicKey!.toString(),
                "uri": origin,
                "version": "1",
            "options": {
                "statement": statement,
            }
        });

        var requestOptions = {
            method: 'POST',
            mode: 'cors',
            headers: {
                "Access-Control-Allow-Origin": "*",
                'Content-Type': 'application/json',
            },
            body: raw,
        };

        // @ts-ignore
        fetch("http://localhost:8080/create", requestOptions)
            .then(response => {
                return response.json()
            })
            .then(message => {
                // we need the nonce for verification so getting it in a global variable
                setNonce(message.payload.nonce);
                setSiwsMessage(message);
                console.log(message);
                // @ts-ignore
                requestOptions.body = JSON.stringify(message)
                // @ts-ignore
                fetch("http://localhost:8080/prepareMessage", requestOptions)
                    .then(response => {
                        return response.text()
                    }).then(messageText => {
                        const messageEncoded = new TextEncoder().encode(messageText);
                        signMessage!(messageEncoded).then(resp => setSignature(
                            bs58.encode(resp)));
                    })
                    .catch(error => console.log('error', error));

            })
            .catch(error => console.log('error', error));
    }

    return (
        <>
            {wallet.connected &&
                sign == "" &&
                <span>
                    <p className='center'>Sign Transaction</p>
                    <input className='publicKey' type="text" id="publicKey" value={walletAddress} />
                </span>
            }
            {
                wallet.connected != true &&
                sign=="" &&
                <div>
                    <div className="logo-wrapper">
                        <img className='solana-logo' src={SolanaLogo} />
                    </div>
                    <p className="sign">Sign in With Solana</p>
                </div>
            }
                    
            {wallet.connected &&
                sign == "" &&
                <div>
                    <button className='web3auth' id='w3aBtn' onClick={createSolanaMessage}>Sign-in with Solana</button>
                    <WalletDisconnectButton className='walletButtons' />
                </div>
            }
            {wallet.connected != true &&
                sign == "" &&
                <WalletModalProvider >
                <WalletMultiButton className='walletButtons' />
                </WalletModalProvider>
            }

            {
                sign &&
                <>
                    <p className='center'>Verify Signature</p>
                    <input className='signature' type="text" id="signature" value={sign} onChange={ e=> setSignature(e.target.value)} />
                    <button className='web3auth' id='verify' onClick={e => {
                        const signature = {
                            t: "sip99",
                            s: sign
                        } 
                        const payload = siwsMessage!.payload;
                        var requestOptions = {
                            method: 'POST',
                            mode: 'cors',
                            headers: {
                                "Access-Control-Allow-Origin": "*",
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify({
                                "payload": payload,
                                "signature": signature
                            }),
                        };
                        // @ts-ignore
                        fetch("http://localhost:8080/verify", requestOptions)
                            .then(response => {
                                return response.text()
                            }).then(resp => {
                            if (resp == "true") {
                                new Swal("Success","Signature Verified","success")
                            } else {
                                new Swal("Error",resp,"error")
                            }
                        })
                    }}>Verify</button>
                    <button className='web3auth' id='verify' onClick={e => {
                        setSiwsMessage(null);
                        setNonce("");
                        setSignature("")
                    }}>Back to Wallet</button>
                </>
            }

        </>
    );
};

export default MyWallet;
