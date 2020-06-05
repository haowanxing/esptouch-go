package esptouchgenerator

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNewEsptouchGenerator(t *testing.T) {
	eg := NewEsptouchGenerator([]byte("Administrators"), []byte{0xf0, 0xb4, 0x29, 0x5c, 0xea, 0x0b}, []byte("123qweasdzxc"), []byte{192, 168, 123, 196})
	//eg := NewEsptouchGenerator([]byte("360wifi"),[]byte{0xf0,0xb4,0x29,0x5c,0xea,0x0b},[]byte("1234567890"),[]byte{192,168,123,196})
	gc := eg.GetGCBytes2()
	dc := eg.GetDCBytes2()
	t.Log(len(gc))
	for _, v := range gc {
		t.Log("gc-itemLen", len(v))
	}
	t.Log(len(dc))
	tmp := make([]int, 0)
	for _, v := range dc {
		tmp = append(tmp, len(v))
	}
	t.Log(tmp)
}

func TestNewEsptouchGenerator2(t *testing.T) {
	mTimeoutGuideCodeMillisecond := int64(2000)
	mTimeoutDataCodeMillisecond := int64(4000)
	mTimeoutTotalCodeMillisecond := mTimeoutGuideCodeMillisecond + mTimeoutDataCodeMillisecond
	mWaitUdpSendingMillisecond := int64(45000)
	mIntervalGuideCodeMillisecond := int64(8)
	mIntervalDataCodeMillisecond := int64(8)
	mInterrupt := false
	eg := NewEsptouchGenerator([]byte("Administrators"), []byte{0xf0, 0xb4, 0x29, 0x5c, 0xea, 0x0b}, []byte("123qweasdzxc"), []byte{192, 168, 123, 196})
	//eg := NewEsptouchGenerator([]byte("360wifi"),[]byte{0xf0,0xb4,0x29,0x5c,0xea,0x0b},[]byte("1234567890"),[]byte{192,168,123,196})
	gc := eg.GetGCBytes2()
	dc := eg.GetDCBytes2()
	fmt.Println(len(gc), len(dc))

	laddr := &net.UDPAddr{
		IP:   net.IPv4(192, 168, 123, 196),
		Port: 7001,
	}
	raddr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: 7001,
	}
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		t.Error(err)
	}

	r := make(chan bool)

	go func() {
		receiveData(r)
		fmt.Println(1)
	}()
	startTime := time.Now().UnixNano() / 1e6
	currentTime := startTime
	lastTime := currentTime - mTimeoutTotalCodeMillisecond

	index := 0
	for {
		select {
		case b := <-r:
			mInterrupt = b
			t.Log("Finished")
			goto end
		default:
			if !mInterrupt {
				if currentTime-lastTime >= mTimeoutTotalCodeMillisecond {
					//log.Println("send gc code")
					for {
						if !mInterrupt && time.Now().UnixNano()/1e6-currentTime < mTimeoutGuideCodeMillisecond {
							//fmt.Println("send gc code", len(gc))
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
					//fmt.Println("send dc code", len(dc))
					sendData(conn, dc, int64(index), 3, mIntervalDataCodeMillisecond)
					index = (index + 3) % len(dc)
				}
				currentTime = time.Now().UnixNano() / 1e6
				if currentTime-startTime > mWaitUdpSendingMillisecond {
					fmt.Println("Wait udp end.")
					break
				}
			} else {
				break
			}
		}
	}
end:
}

func sendData(conn *net.UDPConn, data [][]byte, offset, count, interval int64) {
	for i := offset; i < offset+count; i++ {
		if len(data[i]) == 0 {
			continue
		}
		_, _ = conn.Write(data[i])
		//log.Println("sData",offset, count, i,len(data[i]))
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func receiveData(r chan<- bool) {
	listen, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 18266,
	})
	if err != nil {
		panic(err)
	}

	//for {
	var data [1024]byte
	n, addr, err := listen.ReadFromUDP(data[:])
	if err != nil {
		fmt.Println("read udp err", err, addr, n)
	}
	//fmt.Printf("read n: %d, addr: %v, data: %s", n, addr, data[:n])
	msg := data[:n]
	fmt.Printf("bssid- %x:%x:%x:%x:%x:%x IP- %d.%d.%d.%d", msg[1], msg[2], msg[3], msg[4], msg[5], msg[6], msg[7], msg[8], msg[9], msg[10])
	//break
	//}
	r <- true
}
