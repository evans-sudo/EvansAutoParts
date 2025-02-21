package vm

import "EvansAutoParts/src/pistons/model"


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

type PartEngine struct {
	Base
	Make *model.Make
	Model *model.Model
	Year *model.Year
	Engines []model.Engine
}


type Part struct {
	Base
	Part *model.Part
}

type PartsPartial struct {
	Base
	Parts []*model.Part
}

type SearchResult struct {
	Base
	Make *model.Make
	Model *model.Model
	Year *model.Year
	Engine *model.Engine
	CategoriesJSON string 
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

type Checkout struct {
	Base
	Cart *model.Cart
}



