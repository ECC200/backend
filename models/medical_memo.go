package models

type MedicalMemo struct {
	UserID               string `json:"user_id" firestore:"user_id"`
	ChronicDisease       string `json:"chronic_disease" firestore:"chronic_disease"`
	History              string `json:"history" firestore:"history"`
	HospitalDestination  string `json:"hospital_destination" firestore:"hospital_destination"`
	Department           string `json:"department" firestore:"department"`
	PrimaryCareDoctor    string `json:"primary_care_doctor" firestore:"primary_care_doctor"`
	MedicationManagement string `json:"medication_management" firestore:"medication_management"`
	Allergy              string `json:"allergy" firestore:"allergy"`
	VaccinationStatus    bool   `json:"vaccination_status" firestore:"vaccination_status"`
}
