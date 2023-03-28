package server

func deductBalance(user *User, amount int) {
	user.Balance = user.Balance - amount
}

func creditUser(userId, amount int) {
	mutex.Lock()
	user := getUser(userId)
	user.Balance = user.Balance + amount
	mutex.Unlock()
}
