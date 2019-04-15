package token

// TODO: get params from param keeper
// unit okb
const (
	FeeIssue    = "20000"
	FeeMint     = "2000"
	FeeBurn     = "10"
	FeeFreeze   = "0.1"
	FeeUnfreeze = "0.1"
	FeeTransfer = "0.0125"
)

type FeeDetail struct {
	Address   string `gorm:"index;type:varchar(80)" json:"address"`
	Fee       string `json:"fee"`
	FeeType   string `json:"feeType"` // transfer, deal, etc. see common/const.go
	Timestamp int64  `gorm:"index;type:int64" json:"timestamp"`
}
