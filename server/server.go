package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var AllUsers []*User = make([]*User, 0)
var UserVerificationQueue []*UserVerification

var AllTransactions []*Transaction = make([]*Transaction, 0)
var ProcessTransactionQueue []*TransactionQueue

func CreateUser(c *gin.Context) {
	var cReq CreateUserReq
	if err := c.ShouldBindJSON(&cReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &User{
		Id:                 len(AllUsers) + 1,
		Name:               cReq.Name,
		Balance:            1000,
		VerificationStatus: false,
	}

	AllUsers = append(AllUsers, user)

	UserVerificationQueue = append(UserVerificationQueue, &UserVerification{user.Id, false})

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func CreateTransaction(c *gin.Context) {
	var txReq CreateTxReq
	if err := c.ShouldBindJSON(&txReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if txReq.SenderID == txReq.RecieverID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "transfer to self not allowed"})
		return
	}

	sender := getUser(txReq.SenderID)
	if sender == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sender user id "})
		return
	}

	receiver := getUser(txReq.RecieverID)
	if receiver == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver user id "})
		return
	}

	if sender.Balance < txReq.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	deductBalance(sender, txReq.Amount)

	tx := &Transaction{
		Id:         len(AllTransactions) + 1,
		SenderId:   txReq.SenderID,
		ReceiverId: txReq.RecieverID,
		Amount:     txReq.Amount,
		TxStatus:   "Processing",
		CreatedAt:  time.Now(),
	}

	AllTransactions = append(AllTransactions, tx)

	txq := &TransactionQueue{
		TxId:      tx.Id,
		Processed: false,
	}

	ProcessTransactionQueue = append(ProcessTransactionQueue, txq)

	c.JSON(http.StatusOK, gin.H{
		"status": "transaction request created successfully",
	})
}

func GetAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, AllUsers)
}
