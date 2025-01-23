package ctrl

import (
	"EvansAutoParts/src/bw/model"
	"EvansAutoParts/src/bw/vm"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type adminController struct {
	loginTemplate     *template.Template
	menuTemplate      *template.Template
	createEmpTemplate *template.Template
	viewEmpTemplate   *template.Template
}

func (ac *adminController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ac.GetLogin(w, r)
	case "POST":
		ac.PostLogin(w, r)
	}
}

func (ac *adminController) GetLogin(w http.ResponseWriter, r *http.Request) {
	ac.loginTemplate.Execute(w, nil)
}

func (ac *adminController) PostLogin(w http.ResponseWriter, r *http.Request) {
	employeeNumber, _ := strconv.Atoi(r.FormValue("employeeNumber"))
	password := r.FormValue("password")

	employee, err := model.GetEmployeeWithPassword(employeeNumber, password)

	if err != nil {
		log.Print(err)
		vmodel := vm.Base{Employee: employee}
		ac.loginTemplate.Execute(w, vmodel)
	} else {
		http.Redirect(w, r, "/admin/menu?employeeNumber="+strconv.Itoa(employee.EmployeeNumber), 302)
	}
}

func (ac *adminController) GetMenu(w http.ResponseWriter, r *http.Request) {
	employeeNumber, _ := strconv.Atoi(r.FormValue("employeeNumber"))
	employee, _ := model.GetEmployee(employeeNumber)

	vmodel := vm.Base{Employee: employee}
	ac.menuTemplate.Execute(w, vmodel)
}

func (ac *adminController) HandleCreateEmp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ac.GetCreatEmp(w, r)
	case "POST":
		ac.PostCreateEmp(w, r)
	}
}

func (ac *adminController) GetCreatEmp(w http.ResponseWriter, r *http.Request) {
	employeeNumber, _ := strconv.Atoi(r.URL.Query().Get("employeeNumber"))
	employee, _ := model.GetEmployee(employeeNumber)
	roles, err := model.GetRoles()

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {

		vmodel := vm.AdminCreateEmployee{
			Base:        vm.Base{Employee: employee},
			Roles:       roles,
			NewEmployee: &model.Employee{Role: &model.Role{}},
		}
		ac.createEmpTemplate.Execute(w, vmodel)
	}
}

func (ac *adminController) PostCreateEmp(w http.ResponseWriter, r *http.Request) {
	roleId, _ := strconv.Atoi(r.FormValue("role"))
	role, _ := model.GetRole(roleId)

	payRate, _ := strconv.ParseFloat(r.FormValue("payRate"), 32)
	hireDate, _ := time.Parse("2006-01-02", r.FormValue("hireDate"))
	log.Print(r.FormValue("hireDate"))

	newEmployee := &model.Employee{
		GivenName:  r.FormValue("givenName"),
		Surname:    r.FormValue("surname"),
		Address:    r.FormValue("address"),
		City:       r.FormValue("city"),
		State:      r.FormValue("state"),
		PostalCode: r.FormValue("postalCode"),
		Role:       role,
		HireDate:   hireDate,
		PayRate:    float32(payRate),
	}

	newEmployee, err := model.CreateEmployee(newEmployee)

	if err != nil {
		employeeNumber, _ := strconv.Atoi(r.URL.Query().Get("employeeNumber"))
		employee, _ := model.GetEmployee(employeeNumber)
		roles, _ := model.GetRoles()

		vmodel := vm.AdminCreateEmployee{
			Base:        vm.Base{Employee: employee},
			Roles:       roles,
			NewEmployee: newEmployee,
		}
		ac.createEmpTemplate.Execute(w, vmodel)
	} else {
		http.Redirect(w, r, "/admin/employee?employeeNumber="+
			strconv.Itoa(newEmployee.EmployeeNumber), 302)
	}
}

func (ac *adminController) GetEmployeeView(w http.ResponseWriter, r *http.Request) {
	employeeNumber, err := strconv.Atoi(r.URL.Query().Get("employeeNumber"))

	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
	} else {
		employee, err := model.GetEmployee(employeeNumber)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		} else {
			vmodel := vm.AdminViewEmployee{Base: vm.Base{Employee: &model.Employee{}},
				ViewedEmployee: employee,
			}
			
          ac.viewEmpTemplate.Execute(w, vmodel)
		}
	}

}
