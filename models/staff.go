package models

type Staff struct {
	StaffID   string `json:"staff_id" firestore:"staff_id"`
	StaffName string `json:"Staff_name" firestore:"Staff_name"`
	Password  string `json:"password" firestore:"password"`
}
