# How To TLS

## Create EC parameters for CA

```
% openssl ecparam -name secp256r1 -out rootEC.pem
```

## Create CA

```
openssl req -x509 -sha256 -days 3650 -nodes -newkey ec:rootEC.pem \
  -subj "/C=DE/CN=tomCA" -keyout rootCA.key -out rootCA.pem \
  -config openssl.cnf
```

## Create EC parameters and key for certificate

```
% openssl ecparam -name  secp256r1 -out localhost.key -genkey
```

## Create CSR configuration

```
cat > csr.conf <<EOF
[ req ]
req_extensions = req_ext
prompt = no
distinguished_name = dn

[ dn ]
C = DE
L = CITY
O = COMPANY
OU = DEPARTMENT
CN = localhost

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
EOF
```

## Create CSR

```
% openssl req -new -key localhost.key -out localhost.csr -config csr.conf
```

## Create certificate configuration

```
% cat > localhost.conf <<EOF
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF
```

## Sign CSR

```
% openssl x509 -req -in localhost.csr -CA rootCA.pem -CAkey rootCA.key \
  -CAcreateserial -out localhost.pem -days 730 -sha256 \
  -extfile localhost.conf
```
