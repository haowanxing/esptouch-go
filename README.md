# ESPTouch-Go-SDK
My Learning Project of Smartconfig for ESP8266-WiFi

Looking for my study notes? see: [https://www.imsry.cn/posts/bf40eeb6.html](http://www.imsry.cn/posts/bf40eeb6.html)

Finally. the sdk package born up. （ Code is migrated from https://github.com/EspressifApp/EsptouchForAndroid ）

## SDK-Usage:

```go
package main

import (
	"github.com/haowanxing/esptouch-go"
	"log"
)

func main() {
	task, err := esptouch.NewEsptouchTask([]byte("My-AP"), []byte("400300200"), []byte{0x4c, 0x50, 0x77, 0x73, 0x37, 0xb0})
	if err != nil {
		panic(err)
	}
    // false for multicast, true for broadcast
	task.SetPackageBroadcast(false)
	log.Println("SmartConfig run.")
    // smartconfig device num: 1
	rList := task.ExecuteForResults(1)
	log.Println("Finished", rList)
	return
}
```