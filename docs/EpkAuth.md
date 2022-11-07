TOM-epk {PAYLOAD}

PAYLOAD = base64
  {nonce}:{requestURI}:{IDlib}:{userID}:{fingerprint}:{signature}

nonce       = base64
requestURI  = /path/...
IDlib       = string
userID      = string
fingerprint = string
signature   = base64
