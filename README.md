# FuckCoolapkTokenV3

Generate random request token for Coolapk Android APP

## Usage

```go
package main

import (
    "fmt"
    "github.com/XiaoMengXinX/FuckCoolapkTokenV3"
)

func main() {
	deviceID := "YOUR_DEVICE_ID"
	appVersion := 2510281

	t, err := token.GetToken(deviceID, appVersion, 1761844580)
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated Token:", t)
}
```

output:

```
Generated Token: v3JDJ5JDEwJE5qa3dNemxrTmpRdlpqRTBaRFU0TmVBV0RxTTl6VTFVckRXNFI1WC91NEFVSUtObXpPczRP
```

## Thanks to

[QQ little ice](https://github.com/qqlittleice233)