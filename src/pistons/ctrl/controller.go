package ctrl

import (
	"EvansAutoParts/src/pistons/model"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"reflect"
	"encoding/gob"


	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/reverse"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)


var decoder *schema.Decoder
var cookieCodec *securecookie.SecureCookie

var store *sessions.filesystemStore = sessions.NewFilesystemStore("", []byte("12345678901234567890123456789012"))


const hashKey = "12345678901234567890123456789012"
const blockKey = "12345678901234567890123456789012"

var (
	login    *loginController    = new(loginController)
	parts    *partController     = new(partController)
	checkout *checkoutController = new(checkoutController)
	admin    *adminController    = new(adminController)
	warehouse *warehouseController = new(warehouseController)
)

var (
	TemplateFuncs template.FuncMap

	partMakesRegexp         *reverse.Regexp
	partModelsRegexp        *reverse.Regexp
	partYearsRegexp         *reverse.Regexp
	partEnginesRegexp       *reverse.Regexp
	partSearchRegexp        *reverse.Regexp
	partSearchPartialRegexp *reverse.Regexp
	partDetailRegexp        *reverse.Regexp
	checkoutRegexp          *reverse.Regexp

	adminRegexp            *reverse.Regexp
	adminMenuRegexp        *reverse.Regexp
	adminEmployeeNewRegexp *reverse.Regexp
	adminEmployeeRegexp    *reverse.Regexp

	apiMakesRegexp  *reverse.Regexp
	apiModelsRegexp *reverse.Regexp
)

func init() {
	decoder = schema.NewDecoder()

	decoder.RegisterConverter(time.Now(), func(value string) reflect.Value {
		result := reflect.Value{}
		if t, err := time.Parse("2006-01-02", value); err == nil {
			result = reflect.ValueOf(t)
		}

		return result
	})

	gob.Register(model.Cart{})

	cookieCodec = securecookie.New([]byte(hashKey), []byte(blockKey))

	partMakesRegexp, _ = reverse.CompileRegexp(`/parts/makes`)
	partModelsRegexp, _ = reverse.CompileRegexp(`/parts/makes/(?P<makeId>.+)/models`)
	partYearsRegexp, _ = reverse.CompileRegexp(`/parts/makes/(?P<makeId>.+)/models/(?P<modelId>.+)/years`)
	partEnginesRegexp, _ = reverse.CompileRegexp(`/parts/makes/(?P<makeId>.+)/models/(?P<modelId>.+)/years/(?P<yearId>.*)/engines`)
	partSearchRegexp, _ = reverse.CompileRegexp(`/parts/search`)
	partSearchPartialRegexp, _ = reverse.CompileRegexp(`/parts`)
	partDetailRegexp, _ = reverse.CompileRegexp(`/parts/detail`)

	checkoutRegexp, _ = reverse.CompileRegexp(`/checkout`)

	adminRegexp, _ = reverse.CompileRegexp(`/admin/login`)
	adminMenuRegexp, _ = reverse.CompileRegexp(`/admin/menu`)
	adminEmployeeNewRegexp, _ = reverse.CompileRegexp(`/admin/employees/new`)
	adminEmployeeRegexp, _ = reverse.CompileRegexp(`/admin/employees`)

	apiMakesRegexp, _ = reverse.CompileRegexp(`/api/makes`)
	apiModelsRegexp, _ = reverse.CompileRegexp(`/api/models`)

	TemplateFuncs = template.FuncMap{
		"partMakes":        reverter(partMakesRegexp),
		"partModels":       reverter(partModelsRegexp),
		"partYears":        reverter(partYearsRegexp),
		"partEngines":      reverter(partEnginesRegexp),
		"partSearch":       reverter(partSearchRegexp),
		"partSearchDetail": reverter(partSearchPartialRegexp),
		"partDetail":       reverter(partDetailRegexp),

		"checkout": reverter(checkoutRegexp),

		"admin":            reverter(adminRegexp),
		"adminMenu":        reverter(adminMenuRegexp),
		"adminEmployeeNew": reverter(adminEmployeeNewRegexp),
		"adminEmploye":     reverter(adminEmployeeRegexp),

		"apiMakes":  reverter(apiMakesRegexp),
		"apiModels": reverter(apiModelsRegexp),
	}
}


func reverter(regexp *reverse.Regexp) func(...interface{}) (string, error) {
	return func(params ...interface{}) (string, error) {
		values := url.Values{}
		for i := 0; i < len(params); i += 2 {
			var val string
			t := params[i+1]
			switch t.(type) {
			case string:
				val = params[i+1].(string)
			case int:
				val = strconv.Itoa(params[i+1].(int))
			}

			values[params[i].(string)] = []string{val}
		}
		return regexp.Revert(values)
	}
}

func Setup(tc *template.Template) {
	SetTemplateCache(tc)
	createResourceServer()

	r := mux.NewRouter()
	partMakesRoute, _ := partMakesRegexp.Revert(url.Values{})
	partModelsRoute, _ := partModelsRegexp.Revert(url.Values{"makeId": {"{makeId:[0-9]+}"}})
	partYearsRoute, _ := partYearsRegexp.Revert(url.Values{
		"makeId":  {"{makeId:[0-9]+}"},
		"modelId": {"{modelId:[0-9]+}"},
	})
	partEnginesRoute, _ := partEnginesRegexp.Revert(url.Values{
		"makeId":  {"{makeId:[0-9]+}"},
		"modelId": {"{modelId:[0-9]+}"},
		"yearId":  {"{yearId:[0-9]+}"},
	})
	partSearchRoute, _ := partSearchRegexp.Revert(url.Values{})
	partSearchPartialRoute, _ := partSearchPartialRegexp.Revert(url.Values{})
	partDetailRoute, _ := partDetailRegexp.Revert(url.Values{})

	checkoutRoute, _ := checkoutRegexp.Revert(url.Values{})

	adminRoute, _ := adminRegexp.Revert(url.Values{})
	adminMenuRoute, _ := adminMenuRegexp.Revert(url.Values{})
	adminEmployeeNewRoute, _ := adminEmployeeNewRegexp.Revert(url.Values{})
	adminEmployeeRoute, _ := adminEmployeeRegexp.Revert(url.Values{})

	apiMakesRoute, _ := apiMakesRegexp.Revert(url.Values{})
	apiModelsRoute, _ := apiModelsRegexp.Revert(url.Values{})



	r.HandleFunc("/", login.GetLogin)
	r.HandleFunc(partMakesRoute, parts.GetMake)
	r.HandleFunc(partModelsRoute, parts.GetModel)
	r.HandleFunc(partYearsRoute, parts.GetYear)
	r.HandleFunc(partEnginesRoute, parts.GetEngine)
	r.HandleFunc(partSearchRoute, parts.GetSearch).
		Queries("makeId", "{makeId:[0-9]+}").
		Queries("modelId", "{modelId:[0-9]+}").
		Queries("yearId", "{yearId:[0-9]+}").
		Queries("engineId", "{engineId:[0-9]+}")
	r.HandleFunc(partSearchPartialRoute, parts.GetPartSearchPartial).
		Queries("typeId", "{typeId:[0-9]+}").
		Queries("modelId", "{modelId:[0-9]+}").
		Queries("yearId", "{yearId:[0-9]+}").
		Queries("engineId", "{engineId:[0-9]+}")
	r.HandleFunc(partDetailRoute, parts.GetPart).
		Queries("partId", "{partId:[0-9]+}")
	r.HandleFunc(checkoutRoute, checkout.HandleCheckout)
	r.HandleFunc(adminRoute, admin.GetLogin).Methods("GET")
	r.HandleFunc(adminRoute, admin.PostLogin).Methods("POST")
	r.HandleFunc(adminMenuRoute, admin.GetMenu)
	r.HandleFunc(adminEmployeeNewRoute, admin.GetCreateEmp).
		Methods("GET")
	r.HandleFunc(adminEmployeeNewRoute, admin.PostCreateEmp).
		Methods("POST")
	r.HandleFunc(adminEmployeeRoute, admin.GetEmployeeView)

	r.HandleFunc(apiMakesRoute, parts.AutocompleteMake).
		Queries("term", "{term}")
	r.HandleFunc(apiModelsRoute, parts.AutocompleteModel).
		Queries("term", "{term}").
		Queries("makeId", "{makeId:[0-9]+}")

		r.HandleFunc("/api/inventory/{partNumber:.+}", warehouse.CheckInventory)

		http.Handle("/admin/login", r)

	
		http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			if cookie, err := req.Cookie("employee"); err != nil {
				if empNoRaw := req.URL.Query().Get("employeeNumber"); empNoRaw != "" {
					empNo, _ := strconv.Atoi(empNoRaw)
					employee, err := model.GetEmployee(empNo)
					if err != nil {
						if req.URL.Path != "/" {
							http.Redirect(w, req, "/", 307)
						}
					} else {
						context.Set(req, "employee", employee)
						defer context.Clear(req)
	
						if encoded, err := cookieCodec.Encode("employee", employee); err == nil {
							cookie = &http.Cookie{
								Name:  "employee",
								Value: encoded,
								Path:  "/",
							}
	
							http.SetCookie(w, cookie)
						}
	
					}
				} else {
					if req.URL.Path != "/" {
						http.Redirect(w, req, "/", 307)
					}
				}
			} else {
				var employee *model.Employee
				if err := cookieCodec.Decode("employee", cookie.Value, &employee); err == nil {
					context.Set(req, "employee", employee)
				}
			}
	
			r.ServeHTTP(w, req)
		})
	}

func createResourceServer() {
	http.Handle("/res/", http.StripPrefix("/res", http.FileServer(http.Dir("res"))))
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

