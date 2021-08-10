package dto

import "time"

type CardDTO struct {
	Id      int64     `json:"id"`
	Number  string    `json:"number"`
	Balance int64     `json:"balance"`
	Issuer  string    `json:"issuer"`
	Holder  string    `json:"holder"`
	OwnerId int64     `json:"owner_id"`
	Status  string    `json:"status"`
	Created time.Time `json:"created"`
}

type TransactionDTO struct {
	Id             int64     `json:"id"`
	CardId         int64     `json:"card_id"`
	Amount         int64     `json:"amount"`
	Created        time.Time `json:"created"`
	Status         string    `json:"status"`
	MccId          int64     `json:"mcc_id"`
	Description    string    `json:"description"`
	SupplierIconId int64     `json:"supplier_icon_id"`
}

type ErrDTO struct {
	Err string `json:"error"`
}
