package ctrl

import (
	"EvansAutoParts/src/pistons/model"
	"EvansAutoParts/src/pistons/vm"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/context"
)

type adminController struct {
	loginTemplate     *template.Template
	menuTemplate      *template.Template
	createEmpTemplate *template.Template
	viewEmpTemplate   *template.Template
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
		adminMenuRoute, _ := adminMenuRegexp.Revert(url.Values{})
		http.Redirect(w, r, adminMenuRoute+"?employeeNumber="+strconv.Itoa(employee.EmployeeNumber), 302)
	}
}

func (ac *adminController) GetMenu(w http.ResponseWriter, r *http.Request) {
	employee := context.Get(r, "employee").(*model.Employee)

	vmodel := vm.Base{Employee: employee}
	ac.menuTemplate.Execute(w, vmodel)
}

func (ac *adminController) GetCreateEmp(w http.ResponseWriter, r *http.Request) {
	employee := context.Get(r, "employee").(*model.Employee)
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
	r.ParseForm()
	
	roleId, _ := strconv.Atoi(r.FormValue("role"))
	role, _ := model.GetRole(roleId)

	newEmployee := new(model.Employee)
	err := decoder.Decode(newEmployee, r.PostForm)
	newEmployee.Role = role

	newEmployee, err = model.CreateEmployee(newEmployee)


	if err != nil {
		employee := context.Get(r, "employee").(*model.Employee)
		roles, _ := model.GetRoles()

		vmodel := vm.AdminCreateEmployee{
			Base:        vm.Base{Employee: employee},
			Roles:       roles,
			NewEmployee: newEmployee,
		}
		ac.createEmpTemplate.Execute(w, vmodel)
	} else {
		adminEmployeeRoute, _ := adminEmployeeRegexp.Revert(url.Values{})
		http.Redirect(w, r, adminEmployeeRoute+"?employeeNumber="+
			strconv.Itoa(newEmployee.EmployeeNumber), 302)
	}
}


func (ac *adminController) GetEmployeeView(w http.ResponseWriter, r *http.Request) {
	employee := context.Get(r, "employee").(*model.Employee)
	vmodel := vm.AdminViewEmployee{Base: vm.Base{Employee: &model.Employee{}},
		ViewedEmployee: employee,
	}

	ac.viewEmpTemplate.Execute(w, vmodel)

}