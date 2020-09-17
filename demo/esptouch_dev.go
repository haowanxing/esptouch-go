package main

import (
	"esptouch/protocol"
	"fmt"
	"log"
	"net"
	"time"
)

var (
	datagramCount = 0
	broadcast     bool
)

func main1() {
	broadcast = false
	//broadcast = true
	mIntervalGuideCodeMillisecond := int64(8)
	mIntervalDataCodeMillisecond := int64(8)
	mTimeoutGuideCodeMillisecond := int64(2000)
	mTimeoutDataCodeMillisecond := int64(4000)
	mTimeoutTotalCodeMillisecond := mTimeoutGuideCodeMillisecond + mTimeoutDataCodeMillisecond
	mWaitUdpSendingMillisecond := int64(45000)
	mInterrupt := false
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("Administrators"), []byte{0xf0, 0xb4, 0x29, 0x5c, 0xea, 0x0b}, []byte("123qweasdzxc"), []byte{192, 168, 123, 196})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("WiWide"), []byte{0x00,0x1f,0x7a,0x7b,0xed,0x70}, []byte("wiwide123456"), []byte{10, 11, 98, 45})
	eg := protocol.NewEsptouchGenerator([]byte("jiajiajia"), []byte{0x4c, 0x50, 0x77, 0x73, 0x37, 0xb0}, []byte("400302100"), []byte{172, 16, 104, 145})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("jiajiajia"), []byte{0x4c,0x50,0x77,0x73,0x37,0xb0}, []byte("400302100"), []byte{192 ,168, 3, 90})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("jiajiajia"), []byte{0x4c,0x50,0x77,0x73,0x37,0xb0}, []byte("wiwide"), []byte{172, 16, 104, 145})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("wihidden"), []byte{0x00, 0x1f, 0x7a, 0x7b, 0xed, 0x70}, []byte("wiwide123456"), []byte{10, 11, 98, 45})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("Wiwide-Office"),[]byte{0x00,0x1f,0x7a,0x71,0x93,0xb0},[]byte("4006500311"),[]byte{10,11,146,45})
	//eg := esptouchgenerator.NewEsptouchGenerator([]byte("wihidden2"), []byte{0x00, 0x1f, 0x7a, 0x59, 0x4b, 0x08}, []byte("12345678"), []byte{172, 16, 235, 219})
	gc := eg.GetGCBytes2()
	dc := eg.GetDCBytes2()

	laddr := &net.UDPAddr{
		IP:   nil,
		Port: 18266,
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		panic(err)
	}

	r := make(chan bool)

	go func() {
		receiveData(r, conn)
		fmt.Println(1)
	}()
	startTime := time.Now().UnixNano() / 1e6
	currentTime := startTime
	lastTime := currentTime - mTimeoutTotalCodeMillisecond
	log.Println("SmartConfig start.")
	index := 0
	for {
		select {
		case b := <-r:
			mInterrupt = b
			log.Println("Finished")
			goto end
		default:
			if !mInterrupt {
				if currentTime-lastTime >= mTimeoutTotalCodeMillisecond {
					for {
						if !mInterrupt && time.Now().UnixNano()/1e6-currentTime < mTimeoutGuideCodeMillisecond {
							sendData(conn, gc, 0, int64(len(gc)), mIntervalGuideCodeMillisecond)
						} else {
							index = 0
							break
						}
						if time.Now().UnixNano()/1e6-startTime > mWaitUdpSendingMillisecond {
							fmt.Println("Wait udp end.")
							break
						}
					}
					lastTime = currentTime
				} else {
					sendData(conn, dc, int64(index), 3, mIntervalDataCodeMillisecond)
					index = (index + 3) % len(dc)
				}
				currentTime = time.Now().UnixNano() / 1e6
				if currentTime-startTime > mWaitUdpSendingMillisecond {
					log.Println("UDP Send Timeout.")
					goto end
				}
			} else {
				goto end
			}
		}
	}
end:
	//_ = conn.Close()
}

func sendData(conn *net.UDPConn, data [][]byte, offset, count, interval int64) {
	target := getTargetIP()
	for i := offset; i < offset+count; i++ {
		if len(data[i]) == 0 {
			continue
		}
		_, _ = conn.WriteToUDP(data[i], &net.UDPAddr{
			IP:   target,
			Port: 7001,
		})
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func getNextDatagramCount() byte {
	data := byte(1 + datagramCount%100)
	datagramCount += 1
	return data
}

func getTargetIP() net.IP {
	if broadcast {
		return net.IPv4(255, 255, 255, 255)
	}
	addrByte := getNextDatagramCount()
	return net.IPv4(234, addrByte, addrByte, addrByte)
}

func receiveData(r chan<- bool, conn *net.UDPConn) {
	var data [128]byte
	n, _, err := conn.ReadFromUDP(data[:])
	if err != nil {
		log.Println("read udp err", err, n)
	}
	msg := data[:n]
	log.Printf("bssid- %x:%x:%x:%x:%x:%x IP- %d.%d.%d.%d \n", msg[1], msg[2], msg[3], msg[4], msg[5], msg[6], msg[7], msg[8], msg[9], msg[10])
	r <- true
}
