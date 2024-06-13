package models

type Hosoital struct {
	HospitalID   string `json:"hospital_id" firestore:"hospital_id"`
	HospitalName string `json:"hospital_name" firestore:"hosiital_name"`
	Password     string `json:"password" firestore:"password"`
}
