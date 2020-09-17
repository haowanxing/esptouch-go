package esptouch

import (
	"fmt"
	"net"
)

type IEsptouchResult interface {
	IsSuc() bool
	GetBssid() string
	GetIP() net.IP
	String() string
}

type EsptouchResult struct {
	mIsSuc       bool
	mBssid       string
	mInetAddress net.IP
}

func NewEsptouchResult(isSuc bool, bssid string, inetAddress net.IP) *EsptouchResult {
	return &EsptouchResult{
		mIsSuc:       isSuc,
		mBssid:       bssid,
		mInetAddress: inetAddress,
	}
}
func (r EsptouchResult) IsSuc() bool {
	return r.mIsSuc
}
func (r EsptouchResult) GetBssid() string {
	return r.mBssid
}
func (r EsptouchResult) GetIP() net.IP {
	return r.mInetAddress
}
func (r EsptouchResult) String() string {
	return fmt.Sprintf("bssid=%s, address=%s, suc=%v", r.mBssid, r.mInetAddress, r.mIsSuc)
}
