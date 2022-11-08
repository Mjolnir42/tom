# TOM request authorization scheme

## Objective

The goal of the authorization scheme is to offer a HTTP authentication
scheme that works with public key signatures.

The authentication data is to be sent in the `Authorization` header.

The authentication realm identifier is `TOM-epk`.

The token is a string in base64 encoding. The signature contained in the
token is calculated over multiple fields contained within the token.

The scheme is using Ed25519 keypairs.

## Token components

The cleartext token before base64 encoding consists of multiple fields:

1. nonce
2. timestamp
3. requestURI path
4. key fingerprint
5. user identity library
6. username
7. signature

The token fields are concatenated using a colon character `:` as separator.

### Nonce

6 bytes of per request randomness, encoded as Base64.

### Timestamp

The timestamp of the request, as seconds sind Unix epoch in timezone UTC.
Signatures in scheme TOM-epk are valid for 30 seconds.

### requestURI path

The path of the request that is to be authenticated, without query
parameters. This part should not contain `:` characters, but it must be
handled if it does.

### key fingerprint

The fingerprint of the public key, whose private key was used for signing.
It is the 128 bit Blake2b hash digest of the public key bytes, with the
string `engineroom.machine.tom` used as MAC key.
The fingerprint is the digest in hexadecimal encoding, with lowercase
letters a-f.

### user identity library

The name of the identity library the user account is created in. TOM
supports multiple identity libraries, and accounts are not globally unique.

### username

The username to authenticate the request with.

### signature

The created signature over the token digest hash. The signature is
calculated using the Ed25519 private key of the useraccount over the 128 bit
Blake2b hash of the token fields, without a MAC key:

1. nonce, decoded as binary bytes
2. unix timestamp, 64bit / 8 byte unsigned integer, in big endian byteorder
3. fingerprint, as hexadecimal string bytes
4. requestURI, as string bytes
5. ID library name, as string bytes
6. user name, as string bytes

## Request Body

The request body is signed separately.
