package models

type EmployeeStatus struct {
	ID       int    `json:"id"`
    FullName string `json:"full_name"`
	Status string `json:"status"`
}

type Status struct {
    Status string `json:"status"`
}
