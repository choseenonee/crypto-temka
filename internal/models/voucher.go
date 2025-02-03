package models

const (
	PerOnceVoucher = "once"
)

type Voucher struct {
	Id          string      `json:"id"`
	VoucherType string      `json:"type"`
	Properties  interface{} `json:"properties"`
}
