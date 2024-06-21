package models

type EmergencyContact struct {
	Name  string `json:"name" firestore:"name"`
	Phone string `json:"phone" firestore:"phone"`
}

type User struct {
	UserID            string             `json:"user_id" firestore:"user_id"`                         //障がい者番号
	UserName          string             `json:"user_name" firestore:"user_name"`                     //名前
	Password          string             `json:"password" firestore:"password"`                       //パスワード
	Age               string             `json:"age" firestore:"age"`                                 //年齢
	Address           string             `json:"address" firestore:"address"`                         //住所
	Photo             string             `json:"photo" firestore:"photo"`                             //本人写真
	BirthDate         string             `json:"birth_date" firestore:"birth_date"`                   //生年月日
	Contact           string             `json:"contact" firestore:"contact"`                         //本人連絡先
	EmergencyContacts []EmergencyContact `json:"emergency_contacts" firestore:"emergency_contacts"`   //緊急連絡先
	ChronicDisease    string             `json:"chronic_disease" firestore:"chronic_disease"`         //病名
	History           string             `json:"history" firestore:"history"`                         //倒れた履歴など
	Date              string             `json:"date" firestore:"date"`                               //日付
	PrimaryCareDoctor string             `json:"primary_care_doctor" firestore:"primary_care_doctor"` //かかりつけ医
	Medication_status string             `json:"medication_statu" firestore:"medication_"`            //服薬中の薬
	Doctor_Comment    string             `json:"doctor_comment" firestore:"doctor_comment"`           //主治医のコメント
}
