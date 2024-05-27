package models

import "time"

type User struct {
	UserID           string    `json:"user_id" firestore:"user_id"`
	UserName         string    `json:"user_name" firestore:"user_name"`
	MailAddress      string    `json:"mailaddress" firestore:"mailaddress"`
	Password         string    `json:"password" firestore:"password"`
	BirthDate        time.Time `json:"birth_date" firestore:"birth_date"`
	EmergencyContact string    `json:"emergency_contact" firestore:"emergency_contact"`
	WorkContact      string    `json:"work_contact" firestore:"work_contact"`
	BloodType        string    `json:"blood_type" firestore:"blood_type"`
}
