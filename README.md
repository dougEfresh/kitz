# Logzio for go-kit logger
Send go-kit logs to Logzio

## Getting Started

### Get Logzio token
1. Go to Logzio website
2. Sign in with your Logzio account
3. Click the top menu gear icon (Account)
4. The Logzio token is given in the account page

### Initialize Logger
```go
package main

import (
        "github.com/dougEfresh/kitz"
        "github.com/go-kit/kit/log"
)

const LOGZIO_TOKEN = "123456789"

func main() {
        logger, err := kitz.WithDefaults(LOGZIO_TOKEN)
        if err != nil {
                panic(err)
        }
        // message is required 
        logger.Log("message", "hello!")
}
```

**NOTE**: Set `LOGZIO_TOKEN` to the Logzio token as mentioned in `Get Logzio token`.