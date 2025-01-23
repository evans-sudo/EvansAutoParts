package vm

import "EvansAutoParts/src/bw/model"


type Base struct {
	Employee *model.Employee
}


type PartMake struct {
	Base
}

type PartModel struct {
	Base
	Make *model.Make
}

type PartYear struct {
	Base
	Make  *model.Make
	Model *model.Model
	Years []model.Year
}


type Autocomplete struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Data string `json:"data"`
}



type AdminCreateEmployee struct {
	Base
	Roles       []*model.Role
	NewEmployee *model.Employee
}

type AdminViewEmployee struct {
	Base
	ViewedEmployee *model.Employee
}