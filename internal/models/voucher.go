package models

type Voucher struct {
	Id          string      `json:"id"`
	VoucherType string      `json:"type"`
	Properties  interface{} `json:"properties"`
}
