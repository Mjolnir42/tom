TOM-epk {PAYLOAD}

PAYLOAD = base64
  {nonce}:{time}:{requestURI}:{fingerprint}:{IDlib}:{userID}:{signature}

nonce       = base64
time        = int64
requestURI  = /path/...
fingerprint = string
IDlib       = string
userID      = string
signature   = base64
