package dto

type CompanyPatchReq struct {
	ID                string `json:"id,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	AmountOfEmployees int    `json:"amount_of_employees,omitempty"`
	Registered        bool   `json:"registered,omitempty"`
	Type              string `json:"type,omitempty"`
}

type LoginCredentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}
