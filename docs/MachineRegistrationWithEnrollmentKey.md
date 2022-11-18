# Requirements

1. IdentityLibrary has been created
2. IdentityLibrary has attribute isMachineLibrary = true
3. IdentityLibrary has attribute isSelfEnrollmentEnabled = false
4. IdentityLibrary has attribute enrolment-key set to a base64 encoded
   Ed25519 public key

# Authenticated Server Enrollment Process

1. generate Ed25519 public/private key pair for the server
2. generate fingerprint of the public key:
   - hash algorithm Blake2b
   - output size 16 bytes (128bit)
   - key {$libraryname}.machine.tom
   - input raw keybytes
   use hexadecimal output of the hash as machine fingerprint uid:
   - e37d7ba79800b25cd2f063a64bb126ae
3. generate CSR
```
  user-name:           calculated machine fingerprint uid
  identity-library:    name of the machine library
  fqdn:                FQDN of the machine
  public-key:          Base64 encoding of the machine Ed25519 public key
  enrolment-key:       Base64 encoding of the Ed25519 enrolment public key
  valid-from:          RFC3339 formatted datetime string since when this CSR is valid
  valid-until:         RFC3339 formatted datetime string until when this CSR is valid
  signature.hash:      Base64 encoded Blake2b/512 hash of the CSR fields
  signature.signature: Base64 encoded signature of the hash value, signed with the Ed25519 private enrolment key

  {
    "user-name":
    "identity-library":
    "fqdn":
    "public-key":
    "enrolment-key":
    "valid-from":
    "valid-until":
    "signature": {
      "hash":
      "signature":
    }
  }
```
4. generate request JSON body
```
  {
    "user": {
      "library-name": "${library}",
      "user-name":    "e37d7ba79800b25cd2f063a64bb126ae",
      "first-name":   "lxjpernfuss10",
      "last-name":    "lxjpernfuss10.united.domain",
      "credential": {
        "category":   "public-key",
        "value":      "NP6v2h3uLischfgnOb3IonF64CpN1LuG6nbk5Eai/MM="
      }
    },
    "authorization": {
      "timestamp":    "2022-10-21T14:01:05+02:00",
      "userID":       "e37d7ba79800b25cd2f063a64bb126ae",
      "csr": {
        "user-name":        "e37d7ba79800b25cd2f063a64bb126ae",
        "identity-library": "${library}",
        "fqdn":             "lxjpernfuss10.united.domain",
        "public-key":       "NP6v2h3uLischfgnOb3IonF64CpN1LuG6nbk5Eai/MM=",
        "enrolment-key":    "wGMl6OwdOd4+i+1BsvRfnXwpK205TgP8wCxRXHyUWh4=",
        "valid-from":       "2022-10-21T14:01:00+02:00",
        "valid-until":      "2022-10-21T15:01:00+02:00",
        "signature": {
          "hash":           "jPtwFqhCMW0H70+7U0q1jXxqJTWznAkkfxeP4jhIpFL8ByutTXJW3GOrLXbtBZaStk3JBUMZdmNtML2Yo2xLZA==",
          "signature":      "Cbxj3Sl6QvUizQA7dITdE6iph0jibvGXXmOedT1neOp80XZ/8NRXpd85FZPXsxSA73KlvGUj0rtEGGJzjZ0NAg=="
        }
      }
      "signature": {
        "hash":             "5M0GdIpS/RsRiY5e6IztDaXhSTmbjrc7olUJufGntTyIEjH2B5Cvm1ScC9kv0miGBxn1xVe9Qh/UiTxUDS8QNg==",
        "signature":        "xsOTCcMz8vKtjPQI4DrhkdjO10QFmk3VD2vICTAgwZqQHZVlLnTzzQkniuQ6OPvuUPjguU36EL5SJ2oA3of4CA=="
      }
    }
  }
```
5. send the request body as `PUT` request to `/machine/${uid}`
   - /machine/f748a40a96e9c1500aa4f251c5d27b89
