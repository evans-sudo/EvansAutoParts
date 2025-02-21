package ctrl

import (
	"EvansAutoParts/src/pistons/vm"
	"html/template"
	"net/http"
)



type loginController struct {
	loginTemplate *template.Template
}


func (lc *loginController) GetLogin(w http.ResponseWriter, r *http.Request) {
	vmodel := vm.Base{}
	lc.loginTemplate.Execute(w, vmodel)
}