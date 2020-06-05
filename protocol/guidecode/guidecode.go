package guidecode

const (
	GUIDE_CODE_LEN = 4
)

type GuideCode struct {
}

func NewGuideCode() *GuideCode {
	return &GuideCode{}
}

func (g *GuideCode) GetU8s() []uint16 {
	u8s := make([]uint16, GUIDE_CODE_LEN)
	u8s[0] = 515
	u8s[1] = 514
	u8s[2] = 513
	u8s[3] = 512
	return u8s
}
