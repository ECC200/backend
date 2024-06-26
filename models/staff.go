package models

type Staff struct {
	StaffID       string `json:"staff_id" firestore:"staff_id"`
	StaffName     string `json:"Staff_name" firestore:"Staff_name"`
	Password      string `json:"password" firestore:"password"`
	Department    string `json:"department" firestore:"department"`         //部門
	Position      string `json:"position" firestore:"position"`             //役職
	Date          string `json:"date" firestore:"date"`                     //入社日
	Boss          string `json:"boss" firestore:"boss"`                     //上司
	DoctorMessage string `json:"Doctor_message" firestore:"Doctor_message"` //主治医メッセージ
	WorkStatus    string `json:"workstatus" firestore:"workstatus"`         //仕事中か否か
	//管理レベルは一旦保留
}
