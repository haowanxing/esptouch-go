package esptouch

import (
	"errors"
	"esptouch/protocol"
	"esptouch/task"
	"esptouch/utils/byteutil"
	"fmt"
	"log"
	"net"
	"time"
)

type EsptouchTask struct {
	parameter             *task.EsptouchParameter
	apSsid                []byte
	apPassword            []byte
	apBssid               []byte
	udpClient             *net.UDPConn
	mEsptouchResultList   []IEsptouchResult
	mBssidTaskSucCountMap map[string]int
	mIsInterrupt          bool
	mIsExecuted           bool
	mIsSuc                bool
}

func NewEsptouchTask(apSsid, apPassword, apBssid []byte) (*EsptouchTask, error) {
	if apSsid == nil || len(apSsid) == 0 {
		return nil, errors.New("SSID can't be empty")
	}
	if apBssid == nil || len(apBssid) != 6 {
		return nil, errors.New("BSSID is empty or length is not 6")
	}
	if apPassword == nil {
		apPassword = []byte("")
	}
	mParameter := task.NewEsptouchParameter()
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   nil,
		Port: mParameter.GetPortListening(),
	})
	if err != nil {
		return nil, err
	}
	return &EsptouchTask{
		parameter:             mParameter,
		apSsid:                apSsid,
		apPassword:            apPassword,
		apBssid:               apBssid,
		udpClient:             conn,
		mEsptouchResultList:   make([]IEsptouchResult, 0),
		mBssidTaskSucCountMap: make(map[string]int),
	}, nil
}

func (p *EsptouchTask) checkTaskValid() {
	if p.mIsExecuted {
		panic("the Esptouch task could be executed only once")
	}
	p.mIsExecuted = true
}
func (p *EsptouchTask) putEsptouchResult(isSuc bool, bssid string, ip net.IP) {
	var count int
	if c, ok := p.mBssidTaskSucCountMap[bssid]; ok {
		count = c
	}
	count++
	p.mBssidTaskSucCountMap[bssid] = count
	if !(count >= p.parameter.GetThresholdSucBroadcastCount()) {
		return
	}
	var isExist = false
	for _, esptouchResultInList := range p.mEsptouchResultList {
		if esptouchResultInList.GetBssid() == bssid {
			isExist = true
			break
		}
	}
	if !isExist {
		esptouchResult := NewEsptouchResult(isSuc, bssid, ip)
		p.mEsptouchResultList = append(p.mEsptouchResultList, esptouchResult)
	}
}

func (p *EsptouchTask) listenAsync(expectDataLen int) {
	var startTime = time.Now()
	var receiveBytes = make([]byte, expectDataLen)
	var receiveOneByte byte = 0
	for {
		if len(p.mEsptouchResultList) >= p.parameter.GetExpectTaskResultCount() || p.mIsInterrupt {
			break
		}
		var expectOneByte = byte(len(p.apSsid) + len(p.apPassword) + 9)
		n, _, err := p.udpClient.ReadFromUDP(receiveBytes)
		if err != nil {
			log.Println("read udp err", err, n)
		}
		if n > 0 {
			receiveOneByte = receiveBytes[0]
		}
		if receiveOneByte == expectOneByte {
			consume := time.Now().Sub(startTime).Milliseconds()
			timeout := p.parameter.GetWaitUdpTotalMillisecond() - consume
			if timeout < 0 {
				break
			} else {
				if len(receiveBytes) > 1 {
					var bssid = byteutil.ParseBssid(receiveBytes, p.parameter.GetEsptouchResultOneLen(), p.parameter.GetEsptouchResultMacLen())
					var inetAddress = byteutil.ParseInetAddr(receiveBytes, p.parameter.GetEsptouchResultOneLen()+p.parameter.GetEsptouchResultMacLen(), p.parameter.GetEsptouchResultIpLen())
					//log.Printf("[Success] Bssid: %s, IP: %s", bssid, inetAddress)
					p.putEsptouchResult(true, bssid, inetAddress)
				}
			}
		}
	}
	p.mIsSuc = len(p.mEsptouchResultList) >= p.parameter.GetExpectTaskResultCount()
	p.interrupt()
}

func (p *EsptouchTask) getEsptouchResultList() []IEsptouchResult {
	if len(p.mEsptouchResultList) == 0 {
		esptouchResultFail := NewEsptouchResult(false, "", nil)
		p.mEsptouchResultList = append(p.mEsptouchResultList, esptouchResultFail)
	}
	return p.mEsptouchResultList
}

func (p *EsptouchTask) execute(generator *protocol.EsptouchGenerator) bool {
	startTime := time.Now().UnixNano() / 1e6
	currentTime := startTime
	lastTime := currentTime - p.parameter.GetTimeoutTotalCodeMillisecond()

	gc := generator.GetGCBytes2()
	dc := generator.GetDCBytes2()
	index := 0
	for {
		if !p.mIsInterrupt {
			if currentTime-lastTime >= p.parameter.GetTimeoutTotalCodeMillisecond() {
				for {
					if !p.mIsInterrupt && time.Now().UnixNano()/1e6-currentTime < p.parameter.GetTimeoutGuideCodeMillisecond() {
						p.sendData(gc, 0, int64(len(gc)), p.parameter.GetIntervalGuideCodeMillisecond())
					} else {
						index = 0
						break
					}
					if time.Now().UnixNano()/1e6-startTime > p.parameter.GetWaitUdpSendingMillisecond() {
						fmt.Println("Wait udp end.")
						break
					}
				}
				lastTime = currentTime
			} else {
				p.sendData(dc, int64(index), 3, p.parameter.GetIntervalDataCodeMillisecond())
				index = (index + 3) % len(dc)
			}
			currentTime = time.Now().UnixNano() / 1e6
			if currentTime-startTime > p.parameter.GetWaitUdpSendingMillisecond() {
				log.Println("UDP Send Timeout.")
				break
			}
		} else {
			break
		}
	}
	return p.mIsSuc
}

func (p *EsptouchTask) sendData(data [][]byte, offset, count int64, interval int64) {
	for i := offset; i < offset+count; i++ {
		if len(data[i]) == 0 {
			continue
		}
		_, _ = p.udpClient.WriteToUDP(data[i], &net.UDPAddr{
			IP:   net.ParseIP(p.parameter.GetTargetHostname()),
			Port: p.parameter.GetTargetPort(),
		})
		time.Sleep(time.Millisecond * time.Duration(interval))
	}
}

func (p *EsptouchTask) interrupt() {
	if !p.mIsInterrupt {
		p.mIsInterrupt = true
	}
}
func (p *EsptouchTask) Interrupt() {
	p.interrupt()
}

func (p *EsptouchTask) ExecuteForResults(expectTaskResultCount int) []IEsptouchResult {
	p.checkTaskValid()
	p.parameter.SetExpectTaskResultCount(expectTaskResultCount)
	var ipAddress net.IP
	netInterfaces, _ := net.Interfaces()
	for _, v := range netInterfaces {
		if (v.Flags & net.FlagUp) != 0 {
			addrs, _ := v.Addrs()
			for _, address := range addrs {
				if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						ipAddress = ipnet.IP
					}
				}
			}
		}
	}
	generator := protocol.NewEsptouchGenerator(p.apSsid, p.apBssid, p.apPassword, ipAddress)
	go p.listenAsync(p.parameter.GetEsptouchResultTotalLen())
	var isSuc = false
	for i := 0; i < p.parameter.GetTotalRepeatTime(); i++ {
		isSuc = p.execute(generator)
		if isSuc {
			return p.getEsptouchResultList()
		}
	}
	if !p.mIsInterrupt {
		time.Sleep(time.Millisecond * time.Duration(p.parameter.GetWaitUdpReceivingMillisecond()))
	}
	p.interrupt()
	return p.getEsptouchResultList()
}

func (p *EsptouchTask) SetPackageBroadcast(broadcast bool) {
	p.parameter.SetBroadcast(broadcast)
}
