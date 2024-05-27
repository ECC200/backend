package models

type MedicalMemo struct {
	UserID               string `firestore:"user_id"`
	ChronicDisease       string `firestore:"chronic_disease"`
	History              string `firestore:"history"`
	HospitalDestination  string `firestore:"hospital_destination"`
	Department           string `firestore:"department"`
	PrimaryCareDoctor    string `firestore:"primary_care_doctor"`
	MedicationManagement string `firestore:"medication_management"`
	Allergy              string `firestore:"allergy"`
	VaccinationStatus    bool   `firestore:"vaccination_status"`
}
