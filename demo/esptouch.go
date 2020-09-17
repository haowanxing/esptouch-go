package main

import (
	"esptouch"
	"log"
)

func main() {
	task, err := esptouch.NewEsptouchTask([]byte("jiajiajia"), []byte("400302100"), []byte{0x4c, 0x50, 0x77, 0x73, 0x37, 0xb0})
	if err != nil {
		panic(err)
	}
	task.SetPackageBroadcast(false)
	log.Println("SmartConfig run.")
	rList := task.ExecuteForResults(2)
	log.Println("Finished", rList)
	return
}
