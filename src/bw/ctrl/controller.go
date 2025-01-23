package ctrl

import (
	"html/template"
	"net/http"
)

var (
	login    *loginController    = new(loginController)
	parts    *partController     = new(partController)
	checkout *checkoutController = new(checkoutController)
	admin    *adminController    = new(adminController)
)

func Setup(tc *template.Template) {
	SetTemplateCache(tc)
	createResourceServer()

	http.HandleFunc("/", login.GetLogin)
	http.HandleFunc("/parts/makes", parts.GetMake)
	http.HandleFunc("/parts/models", parts.PostModel)
	http.HandleFunc("/parts/years", parts.PostYear)
	http.HandleFunc("/parts/engines", parts.PostEngine)
	http.HandleFunc("/parts/searchresults", parts.PostSearch)
	http.HandleFunc("/parts", parts.GetPartSearchPartial)
	http.HandleFunc("/parts/detail", parts.GetPart)
	http.HandleFunc("/checkout", checkout.HandleCheckout)
	http.HandleFunc("/admin", admin.HandleLogin)
	http.HandleFunc("/admin/menu", admin.GetMenu)
	http.HandleFunc("/admin/employees/new", admin.HandleCreateEmp)
	http.HandleFunc("/admin/employee", admin.GetEmployeeView)

	http.HandleFunc("/api/makes", parts.AutocompleteMake)
	http.HandleFunc("/api/models", parts.AutocompleteModel)
}


func createResourceServer(tc *template.Template) {

}


func SetTemplateCache(tc *template.Template) {
	login.loginTemplate = tc.Lookup("login.html")

	parts.autoMakeTemplate = tc.Lookup("make.html")
	parts.autoModelTemplate = tc.Lookup("model.html")
	parts.autoYearTemplate = tc.Lookup("year.html")
	parts.autoEngineTemplate = tc.Lookup("engine.html")
	parts.searchResultTemplate = tc.Lookup("search_results.html")
	parts.partTemplate = tc.Lookup("part.html")
	parts.searchResultPartialTemplate = tc.Lookup("_result.html")

	checkout.template = tc.Lookup("checkout.html")

	admin.loginTemplate = tc.Lookup("admin_login.html")
	admin.menuTemplate = tc.Lookup("admin_menu.html")
	admin.createEmpTemplate = tc.Lookup("admin_create_emp.html")
	admin.viewEmpTemplate = tc.Lookup("admin_employee.html")
}

