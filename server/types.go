package server

import "time"

type CreateUserReq struct {
	Name string `json:"name" binding:"required"`
}

type CreateTxReq struct {
	SenderID   int `json:"sender_id" binding:"required"`
	RecieverID int `json:"receiver_id" binding:"required"`
	Amount     int `json:"amount" binding:"required"`
}

type Transaction struct {
	Id         int
	SenderId   int
	ReceiverId int
	Amount     int
	CreatedAt  time.Time
	TxStatus   string
}

type User struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Balance            int    `json:"balance"`
	VerificationStatus bool   `json:"verification_status"`
}

type UserVerification struct {
	Id        int
	Processed bool
}

type TransactionQueue struct {
	TxId      int
	Processed bool
}
