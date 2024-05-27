package models

import "time"

type User struct {
	UserID           string    `firestore:"user_id"`
	UserName         string    `firestore:"user_name"`
	MailAddress      string    `firestore:"mailaddress"`
	Password         string    `firestore:"password"`
	BirthDate        time.Time `firestore:"birth_date"`
	EmergencyContact string    `firestore:"emergency_contact"`
	WorkContact      string    `firestore:"work_contact"`
	BloodType        string    `firestore:"blood_type"`
}
