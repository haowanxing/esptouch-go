package task

import "fmt"

type EsptouchParameter struct {
	datagramCount                 int
	mIntervalGuideCodeMillisecond int64
	mIntervalDataCodeMillisecond  int64
	mTimeoutGuideCodeMillisecond  int64
	mTimeoutDataCodeMillisecond   int64
	mTotalRepeatTime              int
	mEsptouchResultOneLen         int
	mEsptouchResultMacLen         int
	mEsptouchResultIpLen          int
	mEsptouchResultTotalLen       int
	mPortListening                int
	mTargetPort                   int
	mWaitUdpReceivingMillisecond  int64
	mWaitUdpSendingMillisecond    int64
	mThresholdSucBroadcastCount   int
	mExpectTaskResultCount        int
	mBroadcast                    bool
}

func NewEsptouchParameter() *EsptouchParameter {
	return &EsptouchParameter{
		datagramCount:                 0,
		mIntervalGuideCodeMillisecond: 8,
		mIntervalDataCodeMillisecond:  8,
		mTimeoutGuideCodeMillisecond:  2000,
		mTimeoutDataCodeMillisecond:   4000,
		mTotalRepeatTime:              1,
		mEsptouchResultOneLen:         1,
		mEsptouchResultMacLen:         6,
		mEsptouchResultIpLen:          4,
		mEsptouchResultTotalLen:       1 + 6 + 4,
		mPortListening:                18266,
		mTargetPort:                   7001,
		mWaitUdpReceivingMillisecond:  15000,
		mWaitUdpSendingMillisecond:    45000,
		mThresholdSucBroadcastCount:   1,
		mExpectTaskResultCount:        1,
		mBroadcast:                    true,
	}
}

func (p *EsptouchParameter) GetIntervalGuideCodeMillisecond() int64 {
	return p.mIntervalGuideCodeMillisecond
}

func (p *EsptouchParameter) GetIntervalDataCodeMillisecond() int64 {
	return p.mIntervalDataCodeMillisecond
}

func (p *EsptouchParameter) GetTimeoutGuideCodeMillisecond() int64 {
	return p.mTimeoutGuideCodeMillisecond
}

func (p *EsptouchParameter) GetTimeoutDataCodeMillisecond() int64 {
	return p.mTimeoutDataCodeMillisecond
}
func (p *EsptouchParameter) GetTimeoutTotalCodeMillisecond() int64 {
	return p.mTimeoutGuideCodeMillisecond + p.mTimeoutDataCodeMillisecond
}

func (p *EsptouchParameter) GetTotalRepeatTime() int {
	return p.mTotalRepeatTime
}

func (p *EsptouchParameter) GetEsptouchResultOneLen() int {
	return p.mEsptouchResultOneLen
}

func (p *EsptouchParameter) GetEsptouchResultMacLen() int {
	return p.mEsptouchResultMacLen
}

func (p *EsptouchParameter) GetEsptouchResultIpLen() int {
	return p.mEsptouchResultIpLen
}

func (p *EsptouchParameter) GetEsptouchResultTotalLen() int {
	return p.mEsptouchResultTotalLen
}

func (p *EsptouchParameter) GetPortListening() int {
	return p.mPortListening
}

func (p *EsptouchParameter) GetTargetPort() int {
	return p.mTargetPort
}

func (p *EsptouchParameter) GetWaitUdpReceivingMillisecond() int64 {
	return p.mWaitUdpReceivingMillisecond
}

func (p *EsptouchParameter) GetWaitUdpSendingMillisecond() int64 {
	return p.mWaitUdpSendingMillisecond
}

func (p *EsptouchParameter) GetWaitUdpTotalMillisecond() int64 {
	return p.mWaitUdpReceivingMillisecond + p.mWaitUdpSendingMillisecond
}

func (p *EsptouchParameter) GetThresholdSucBroadcastCount() int {
	return p.mThresholdSucBroadcastCount
}
func (p *EsptouchParameter) GetExpectTaskResultCount() int {
	return p.mExpectTaskResultCount
}

func (p *EsptouchParameter) SetExpectTaskResultCount(mExpectTaskResultCount int) {
	p.mExpectTaskResultCount = mExpectTaskResultCount
}

func (p *EsptouchParameter) GetBroadcast() bool {
	return p.mBroadcast
}
func (p *EsptouchParameter) SetBroadcast(mBroadcast bool) {
	p.mBroadcast = mBroadcast
}

func (p *EsptouchParameter) nextDatagramCount() int {
	p.datagramCount++
	return 1 + (p.datagramCount-1)%100
}
func (p *EsptouchParameter) GetTargetHostname() string {
	if p.mBroadcast {
		return "255.255.255.255"
	} else {
		count := p.nextDatagramCount()
		return fmt.Sprintf("234.%d.%d.%d", count, count, count)
	}
}
