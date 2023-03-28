package server

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex

func TransactionProcessor() {
	for range time.Tick(time.Second * 30) {
		fmt.Println("running the tx processor queue")
		if len(ProcessTransactionQueue) > 0 {
			fmt.Println("about to range through process tx queue")
			for _, txq := range ProcessTransactionQueue {
				if !txq.Processed {
					fmt.Printf("about to process tx %d\n", txq.TxId)
					go ProcessTransaction(txq)
				}
			}
		}
	}
}

func ProcessTransaction(txq *TransactionQueue) {
	tx, err := FetchTransaction(txq.TxId)
	if err != nil {
		fmt.Printf("error occured in fetching tx: %d, due to: %v\n", txq.TxId, err)
		return
	}
	uSender := getUser(tx.SenderId)
	if uSender != nil && !uSender.VerificationStatus {
		UserVerificationQueue = append(UserVerificationQueue, &UserVerification{tx.SenderId, false})
	}

	uReceiver := getUser(tx.ReceiverId)
	if uReceiver != nil && !uReceiver.VerificationStatus {
		UserVerificationQueue = append(UserVerificationQueue, &UserVerification{tx.ReceiverId, false})
	}

	if uSender != nil && uReceiver != nil && uSender.VerificationStatus {
		if uSender.Balance > 0 && uSender.Balance < tx.Amount {
			fmt.Printf("balance not sufficient: userId: %d, balance: %d, txAmount: %d\n", uSender.Id, uSender.Balance, tx.Amount)
			mutex.Unlock()
			return
		}

		err := PerformTransfer(uReceiver.Id, tx.Amount)
		if err != nil {
			fmt.Printf("error occured in performing credit to userId:%d, due to: %v\n", uReceiver.Id, err)
			//credit sender
			creditUser(uSender.Id, tx.Amount)
			return
		}

		txq.Processed = true
		tx.TxStatus = "Processed"
	}
}

func FetchTransaction(txId int) (*Transaction, error) {
	for _, tx := range AllTransactions {
		if tx.Id == txId {
			return tx, nil
		}
	}
	return nil, errors.New("tx not found")
}

func PerformTransfer(receiverId int, amount int) error {
	mutex.Lock()
	receiver := getUser(receiverId)
	fmt.Printf("crediting amount:%d to userId:%d\n", amount, receiver.Id)
	receiver.Balance = receiver.Balance + amount
	mutex.Unlock()

	//return error if any
	return nil
}
