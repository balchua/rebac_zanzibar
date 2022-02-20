package deal

import "time"

type DealInfo struct {
	DealId        string
	ClientName    string
	ClientId      string
	ContractId    string
	Date          time.Time
	Supplementary DealSupplementaryInfo
	Servicing     DealServicingInfo
	Origin        string
}

type DealSupplementaryInfo struct {
	SupplInfo1 string
	SupplInfo2 string
	SupplInfo3 string
}

type DealServicingInfo struct {
	Servicing1 string
	Servicing2 string
	Servicing3 string
}
