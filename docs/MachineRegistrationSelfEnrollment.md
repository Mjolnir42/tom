# Requirements

1. IdentityLibrary has been created
2. IdentityLibrary has attribute isMachineLibrary = true
3. IdentityLibrary has attribute isSelfEnrollmentEnabled = true

# Self-Enrollment Process

1. generate Ed25519 public/private key pair
2. generate hash of the public key:
   - hash algorithm Blake2b
   - output size 16 bytes (128bit)
   - key engineroom.machine.tom
   - input raw keybytes
3. use hexadecimal output of the hash as machine fingerprint uid:
   - f748a40a96e9c1500aa4f251c5d27b89
   this leads to the following tomID
   => f748a40a96e9c1500aa4f251c5d27b89.engineroom.machine.tom
4. generate JSON body
```
  library-name:                         'engineroom' identity library
  user-name:                            machine fingerprint uid
  first-name:                           hostname (short)
  last-name:                            fqdn
  credential.category:                  magic value `public-key` for machines
  credential.value:                     base64 encoding of the generated public key
  authorization.timestamp:              request timestamp
  authorization.userID:                 userID that the request is authorized with
  authorization.fingerprint:            fingerprint of the pubKey used for authorization
  authorization.nonce:                  base64 encoding of 6 bytes of per-request randomness
  authorization.signature.datahash:     hash of the request serialization (not transmitted)
  authorization.signature.signature:    privatekey signature of the untransmitted datahash

  {
    "user": {
      "library-name": "engineroom",
      "user-name":    "f748a40a96e9c1500aa4f251c5d27b89",
      "first-name":   "lxjpernfuss",
      "last-name":    "lxjpernfuss.united.domain",
      "credential": {
        "category":   "public-key",
        "value":      "tGSjX+stPPxfRfb7CZ9LPLY5JV4z4NPyF3/NNLEP2ns="
      }
    },
    "authorization": {
      "timestamp":    "2022-10-21T14:01:05+02:00",
      "userID":       "f748a40a96e9c1500aa4f251c5d27b89",
      "fingerprint":  "f748a40a96e9c1500aa4f251c5d27b89",
      "nonce":        "fNGq3Ifu",
      "signature": {
        "signature":    "uEnpAqJU3PuTYEbhH6ie8+AYQ4NGoyA9IJAT+d2nymbO/SDsVtjovCP9zUCKyZ9uwqADExBln3tFdLqh4IhqBw=="
      }
    }
  }
```
5. send the request body as `PUT` request to `/machine/${tomID}`
   - /machine/f748a40a96e9c1500aa4f251c5d27b89.engineroom.machine.tom
