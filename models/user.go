package models

import "time"

type EmergencyContact struct {
	Name  string `json:"name" firestore:"name"`
	Phone string `json:"phone" firestore:"phone"`
}

type User struct {
	UserID               string             `json:"user_id" firestore:"user_id"`
	UserName             string             `json:"user_name" firestore:"user_name"`
	MailAddress          string             `json:"mailaddress" firestore:"mailaddress"`
	Password             string             `json:"password" firestore:"password"`
	BirthDate            time.Time          `json:"birth_date" firestore:"birth_date"`
	WorkContact          string             `json:"work_contact" firestore:"work_contact"`
	BloodType            string             `json:"blood_type" firestore:"blood_type"`
	EmergencyContacts    []EmergencyContact `json:"emergency_contacts" firestore:"emergency_contacts"`
	ChronicDisease       string             `json:"chronic_disease" firestore:"chronic_disease"`
	History              string             `json:"history" firestore:"history"`
	HospitalDestination  string             `json:"hospital_destination" firestore:"hospital_destination"`
	Department           string             `json:"department" firestore:"department"`
	PrimaryCareDoctor    string             `json:"primary_care_doctor" firestore:"primary_care_doctor"`
	MedicationManagement string             `json:"medication_management" firestore:"medication_management"`
	Allergy              string             `json:"allergy" firestore:"allergy"`
	VaccinationStatus    bool               `json:"vaccination_status" firestore:"vaccination_status"`
}
