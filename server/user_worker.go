package server

import (
	"fmt"
	"time"
)

func UserVerificationProcessor() {
	for range time.Tick(time.Second * 30) {
		fmt.Println("running the verification queue")
		if len(UserVerificationQueue) > 0 {
			fmt.Println("about to range through user verification queue")
			for _, u := range UserVerificationQueue {
				if !u.Processed {
					fmt.Printf("about to verify user %d\n", u.Id)
					go verifyUser(u)
				}
			}
		}
	}
}

func verifyUser(x *UserVerification) {
	for _, u := range AllUsers {
		if u.Id == x.Id && u.Id != 1 {
			u.VerificationStatus = true
			x.Processed = true
		}
	}
}

func getUser(userId int) *User {
	for _, u := range AllUsers {
		if u.Id == userId {
			return u
		}
	}
	return nil
}
