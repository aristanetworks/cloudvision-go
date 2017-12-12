* Steps:

** Customer side:

*** Have/Create private key (the customer will do this, not us):

```sh
openssl genrsa -out private.pem 2048
```

*** Generate CSR from that private key (the customer will do this, not us):

```sh
openssl req -new -key private.pem -out private.csr
```

*** Send this CSR to us

** Aeris side

*** Sign the CSR


```sh
./aeris-cli cert enroll \
    --storage hbasenano \
    --storageoption zkquorum=zookeeper-0.zookeeper.default.svc.internal.ovh-bhs.prod.arista.io:2181 \
    --storageoption table=aeris \
    --csr customers/arista-alpha/private.csr \
    --org arista-alpha \
    --signer-crt ./aeris/ca.pem \
    --signer-key ./aeris/ca-key.pem
```

This will generate a file `enroll.crt` which will be sent to the customer.
The customer will use this file to bootstrap TerminAttr, by adding the flag `--autocert=true --ingestauth=certs,enroll.crt,<private.pem>` to the TerminAttr command.

Those certs should be in `/persist/secure` (ToBeValidated) on the switch
