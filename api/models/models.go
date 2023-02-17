package models

type Company struct {
	ID                string `json:"id,omitempty" db:"id"  valid:"uuidv4,required"`
	Name              string `json:"name" db:"name" valid:"stringlength(1|15),required"`
	Description       string `json:"description,omitempty" db:"description" valid:"maxstringlength(3000)"`
	AmountOfEmployees int    `json:"amount_of_employees" db:"amount_of_employees" valid:"required"`
	Registered        bool   `json:"registered" db:"registered" valid:"required"`
	Type              string `json:"type" db:"type" valid:"in(Corporations|NonProfit|Cooperative|Sole Proprietorship),required"`
}
