package models

type StaffGroup struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	MemberCount int `json:"member_count"`
}

type StaffShortInfo struct {
    ID       int    `json:"id"`
    FullName string `json:"full_name"`
}

type StaffIds struct {
    StaffIDs []int `json:"staff_ids"`
}
