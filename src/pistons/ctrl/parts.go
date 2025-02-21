package ctrl

import (
	"EvansAutoParts/src/pistons/model"
	"EvansAutoParts/src/pistons/vm"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)



type partController struct { 
    autoMakeTemplate *template.Template
    autoModelTemplate *template.Template
    autoYearTemplate *template.Template
    autoEngineTemplate *template.Template
    searchResultTemplate *template.Template
    searchResultPartialTemplate *template.Template
    partTemplate *template.Template
}

func (pc *partController) GetMake(w http.ResponseWriter, r *http.Request) {
    employee := context.Get(r, "employee").(*model.Employee)
	vmodel := vm.PartMake{Base: vm.Base{Employee: employee}}
	pc.autoMakeTemplate.Execute(w, vmodel)
}

func (pc *partController) AutocompleteMake(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["term"]
	makes, err := model.SearchForMakes(term)
	if err != nil {
		w.WriteHeader(500)
	} else {
		vmodel := make([]vm.Autocomplete, len(makes))
		for idx, make := range makes {
			vmodel[idx] = vm.Autocomplete{
				Label: make.Name,
				Value: make.Name,
				Data:  strconv.Itoa(make.Id),
			}
		}
		w.Header().Add("Content-Type", "application/json")
		resp, _ := json.Marshal(vmodel)
		w.Write(resp)
	}
}

func (pc *partController) GetModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	makeId, _ := strconv.Atoi(vars["makeId"])
	make, err := model.GetMake(makeId)
	if err != nil {
		w.WriteHeader(500)
	} else {
		employee := context.Get(r, "employee").(*model.Employee)
		vmodel := vm.PartModel{Base: vm.Base{Employee: employee}, Make: make}
		pc.autoModelTemplate.Execute(w, vmodel)
	}
}

func (pc *partController) AutocompleteModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	term := vars["term"]
	makeId, _ := strconv.Atoi(vars["makeId"])
	models, err := model.SearchForModels(makeId, term)
	if err != nil {
		w.WriteHeader(500)
	} else {
		vmodel := make([]vm.Autocomplete, len(models))
		for idx, m := range models {
			vmodel[idx] = vm.Autocomplete{
				Label: m.Name,
				Value: m.Name,
				Data:  strconv.Itoa(m.Id),
			}
		}

		w.Header().Add("Content-Type", "application/json")
		resp, _ := json.Marshal(vmodel)
		w.Write(resp)
	}
}


func (pc *partController) GetYear(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	makeId, _ := strconv.Atoi(vars["makeId"])
	modelId, _ := strconv.Atoi(vars["modelId"])
	years, err := model.FindYearsForModel(modelId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
        employee := context.Get(r, "employee").(*model.Employee)
		autoMake, _ := model.GetMake(makeId)
		autoModel, _ := model.GetModel(modelId)

		vmodel := vm.PartYear{Base: vm.Base{Employee: employee},
			Model: autoModel,
			Make:  autoMake,
			Years: years,
		}
		pc.autoYearTemplate.Execute(w, vmodel)
	}
}

func (pc *partController) GetEngine(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	makeId, _ := strconv.Atoi(vars["makeId"])
	modelId, _ := strconv.Atoi(vars["modelId"])
	yearId, _ := strconv.Atoi(vars["yearId"])
	engines, err := model.SearchForEngines(modelId, yearId)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
        employee := context.Get(r, "employee").(*model.Employee)
		autoModel, _ := model.GetModel(modelId)
		autoMake, _ := model.GetMake(makeId)
		autoYear, _ := model.GetYear(yearId)

		vmodel := vm.PartEngine{Base: vm.Base{Employee: employee},
			Make:    autoMake,
			Model:   autoModel,
			Year:    autoYear,
			Engines: engines,
		}
		pc.autoEngineTemplate.Execute(w, vmodel)
	}
}
		
	
func (pc *partController) GetSearch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	makeId, _ := strconv.Atoi(vars["makeId"])
	modelId, _ := strconv.Atoi(vars["modelId"])
	yearId, _ := strconv.Atoi(vars["yearId"])
	engineId, _ := strconv.Atoi(vars["engineId"])
	categories, err := model.GetPartCategories()

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
        employee := context.Get(r, "employee").(*model.Employee)
		autoMake, _ := model.GetMake(makeId)
		autoModel, _ := model.GetModel(modelId)
		autoYear, _ := model.GetYear(yearId)
		autoEngine, _ := model.GetEngine(engineId)

		categoriesJSON, _ := json.Marshal(categories)

		vmodel := vm.SearchResult{Base: vm.Base{Employee: employee},
			Make:           autoMake,
			Model:          autoModel,
			Year:           autoYear,
			Engine:         autoEngine,
			CategoriesJSON: string(categoriesJSON),
		}

		pc.searchResultTemplate.Execute(w, vmodel)
	}
}


func (pc *partController) GetPartSearchPartial(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	modelId, _ := strconv.Atoi(vars["modelId"])
	yearId, _ := strconv.Atoi(vars["yearId"])
	engineId, _ := strconv.Atoi(vars["engineId"])
	typeId, _ := strconv.Atoi(vars["typeId"])
	parts, err := model.SearchForParts(modelId, yearId, engineId, typeId)
	if err != nil {
		log.Print(err)
		w.WriteHeader(500)
	} else {
        employee := context.Get(r, "employee").(*model.Employee)
		vmodel := vm.PartsPartial{Base: vm.Base{Employee: employee}, Parts: parts}
		pc.searchResultPartialTemplate.Execute(w, vmodel)
	}
}

                      
func (pc *partController) GetPart(w http.ResponseWriter, r *http.Request) {
   vars := mux.Vars(r)
    partId, _ := strconv.Atoi(vars["partId"])
    employee := context.Get(r, "employee").(*model.Employee)

	part, err := model.GetPart(partId)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
		vmodel := vm.Part{Base: vm.Base{Employee: employee}, Part: part}

		pc.partTemplate.Execute(w, vmodel)
	}
}

func (pc *partController) PostPart(w http.ResponseWriter, r *http.Request) {
	partId, _ := strconv.Atoi(r.FormValue("PartId"))
	quantity, _ := strconv.Atoi(r.FormValue("Quantity"))
	employee := context.Get(r, "employee").(*model.Employee)

	part, err := model.GetPart(partId)

	session, _ := store.Get(r, "session")
	var cart model.Cart
	if val, _ := session.Values["cart"]; val != nil {
		cart = val.(model.Cart)
	}

	cartItem := model.CartItem{}
	cartItem.Part = part
	cartItem.Quantity = quantity

	cart.Items = append(cart.Items, cartItem)

	session.Values["cart"] = cart
	session.Save(r, w)

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
		vmodel := vm.Part{Base: vm.Base{Employee: employee}, Part: part}

		pc.partTemplate.Execute(w, vmodel)
	}
}
   


