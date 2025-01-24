package ctrl

import (
	"EvansAutoParts/src/bw/model"
	"EvansAutoParts/src/bw/vm"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)



type partController struct { 
    automakeTemplate *template.Template
    autoModelTemplate *template.Template
    autoYearTemplate *template.Template
    autoEngineTemplate *template.Template
    searchResultTemplate *template.Template
    searchResultPartialTemplate *template.Template
    partTemplate *template.Template
}

func (pc *partController) GetMake(w http.ResponseWriter, r *http.Request) {
    employeeNumber, err := strconv.Atoi(r.FormValue("employeeNumber"))
    if err != nil {
        log.Println(err)
        w.WriteHeader(500)
    } else {
        employee, err := model.GetEmployee(employeeNumber)
        if err != nil {
            log.Print(err)
            http.Redirect(w,r, "/", 307)
        } else {
            vmodel := vm.PartMake{Base: vm.Base{Employee: employee}}
            pc.automakeTemplate.Execute(w, vmodel)
        }
    }
}


func (pc *partController) AutocompleteMake(w http.ResponseWriter, r *http.Request) {
    term := r.URL.Query().Get("term")
    makes, err := model.SearchForMakes(term)
    if err != nil {
        w.WriteHeader(500)
    }

    vmodel := make([]vm.Autocomplete, len(makes))
    for idx, make := range makes {
        vmodel[idx] = vm.Autocomplete{
            Label: make.Name,
            Value: make.Name,
            Data: strconv.Itoa(make.Id),
        }

        w.Header().Add("Content-Type", "application/json")
        resp, _ := json.Marshal(vmodel)
        w.Write(resp)
    }
}



func (pc *partController) PostModel(w http.ResponseWriter, r *http.Request) {
    makeId, err := strconv.Atoi(r.FormValue("make"))
    if err != nil {
        w.WriteHeader(500)
    } else {
        make, err := model.GetMake(makeId)
        if err != nil {
            w.WriteHeader(500)
        } else {
            employeeNumber, err := strconv.Atoi(r.FormValue("employeeNumber"))
            if err != nil {
                log.Println(err)
                w.WriteHeader(500)
            } else {
                employee, _ := model.GetEmployee(employeeNumber)
                vmodel := vm.PartModel{Base: vm.Base{Employee: employee}, Make: make}
                pc.autoModelTemplate.Execute(w, vmodel)
            }
        }
    }
}

func (pc *partController) AutocompleteModel(w http.ResponseWriter, r *http.Request) {
    makeId, err := strconv.Atoi(r.FormValue("make"))
    term := r.URL.Query().Get("term")
    if err != nil {
        w.WriteHeader(500)
        log.Println("Failed to convert make to an integer")
    } else {
        models, err := model.SearchForModels(makeId, term)
        if err != nil {
            w.WriteHeader(500)
        } else {
            vmodel := make([]vm.Autocomplete, len(models))
            for idx, m := range models {
                vmodel[idx] = vm.Autocomplete{
                    Label: m.Name,
                    Value: m.Name,
                    Data: strconv.Itoa(m.Id),
                }
            }

            w.Header().Add("Content=Type", "application/json")
            resp, _ := json.Marshal(vmodel)
            w.Write(resp)
        }
    }
}


func (pc *partController) PostYear(w http.ResponseWriter, r *http.Request) {
    makeId, err := strconv.Atoi(r.FormValue("make")) 
    if err != nil {
        log.Println(err)
        w.WriteHeader(500)
    } else {
        modelId, err := strconv.Atoi(r.FormValue("model"))
        if err != nil {
            log.Println(err)
            w.WriteHeader(500)
        } else {
            years, err := model.FindYearsForModel(modelId)
            if err != nil {
                log.Println(err)
                w.WriteHeader(500)
            } else {
                employeeNumber, err := strconv.Atoi(r.FormValue("employeeNumber"))
                if err != nil {
                    log.Println(err)
                    w.WriteHeader(500)
                } else {
                    employee, _ := model.GetEmployee(employeeNumber)
                    autoMake, _ := model.GetMake(makeId)
                    autoModel, _ := model.GetModel(modelId) 
                    
                    
                    vmodel := vm.PartYear{Base: vm.Base{Employee: employee},
                        Model: autoModel,
                        Make: autoMake,
                         Years: years,
                }
                    pc.autoYearTemplate.Execute(w, vmodel)
                    
                }
            }
        }
    }
}


func (pc *partController) PostEngine(w http.ResponseWriter, r *http.Request) {
    makeId, err := strconv.Atoi(r.FormValue("make"))

    if err != nil {
        log.Println(err)
        w.WriteHeader(500)
    } else {
        modelId, err := strconv.Atoi(r.FormValue("model"))

        if err != nil {
            log.Println(err)
            w.WriteHeader(500)
        } else {
            yearId, err := strconv.Atoi(r.FormValue("year"))


            if err != nil {
                log.Println(err)
                w.WriteHeader(500)
            } else {
                engines, err := model.SearchForEngines(modelId, yearId)
                if err != nil {
                    log.Println(err)
                    w.WriteHeader(500)
                } else {
                    employeeNumber, err := strconv.Atoi(r.FormValue("employeeNumber"))
                    if err != nil {
                        log.Println(err)
                        w.WriteHeader(500)
                    } else {
                        employee, _ := model.GetEmployee(employeeNumber)
                        autoModel, _ := model.GetModel(modelId)
                        autoMake, _ := model.GetMake(makeId)
                        autoYear, _ := model.GetYear(yearId)

                        vmodel := vm.PartEngine{Base: vm.Base{Employee: employee},
                             Make: autoMake,
                             Model: autoModel,
                             Year: autoYear,
                             Engines: engines,                

                    }
                    pc.autoEngineTemplate.Execute(w, vmodel)
                }
            }
        }

        }
    }

}


func (pc *partController) PostSearch(w http.ResponseWriter, r *http.Request) {
    makeId, err := strconv.Atoi(r.FormValue("make"))


    if err != nil {
        log.Println(err)
        w.WriteHeader(500)
    } else {
        modelId, err := strconv.Atoi(r.FormValue("model"))

        if err != nil {
            log.Println(err)
            w.WriteHeader(500)
        } else {
            yearId, err := strconv.Atoi(r.FormValue("year"))

            if err != nil {
                log.Println(err)
                w.WriteHeader(500)
            } else {
                engineId, err := strconv.Atoi(r.FormValue("engine"))

                if err != nil {
                    log.Println(err)
                    w.WriteHeader(500)
                } else {
                    categories, err := model.g
                }
            }
        }
    }
}

