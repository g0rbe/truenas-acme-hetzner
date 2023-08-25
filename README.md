# TrueNAS-ACME-Hetzner


Based on: https://www.truenas.com/community/threads/acme-dns-authenticator-shell-script.106589/

Relecant source code: https://github.com/truenas/middleware/blob/master/src/middlewared/middlewared/plugins/acme_protocol_/authenticators/shell.py

```
The authenticator script is called two times during the certificate generation:
 
1. The validation record creation which is called in the following way:
   script set domain validation_name validaton_context timeout
2. The validation record deletion which is called in following way:
   script unset domain validation_name validation_context
 
It is up to script implementation to handle both calls and perform the record creation.
```

Example:
```bash
tah set nas.example.com _acme-challenge.nas.example.com validation_token
```

## Install

Download binary:
```bash
wget -O /path/in/pool/tah <link>
```

Make it executable:
```bash
chmod +x /path/in/pool/tah
```

Initialize:
```
/path/in/pool/tah init
```

Set Hetzner DNS API key (**change the API key**) in the current user's $HOME. The file must contains the API key string only!
```bash
echo -n "api-key" > $HOME/.tahtoken
```

Test config (**change the domain**):
```bash
/path/in/pool/tah test nas.example.com
```

## Build

Clone the repo and build:
```bash
go build .
```