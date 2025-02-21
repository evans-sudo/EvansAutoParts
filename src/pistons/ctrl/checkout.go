package ctrl

import (
	"EvansAutoParts/src/pistons/model"
	"EvansAutoParts/src/pistons/vm"
	"html/template"
	"net/http"

	"github.com/gorilla/context"
)


type checkoutController struct {
	template *template.Template
}


func (cc *checkoutController) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cc.GetCheckout(w,r)
	case "POST":
		cc.PostCheckout(w,r)
	}
}



func (cc *checkoutController) GetCheckout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	employee, _ := context.Get(r, "employee").(*model.Employee)

	var cart model.Cart
	if val, _ := session.Values["cart"]; val != nil {
		cart = val.(model.Cart)
	}


	viewModel := vm.Checkout{}
	viewModel.Employee = employee
	viewModel.Cart = &cart
	
	cc.template.Execute(w, viewModel)
}

func (cc *checkoutController) PostCheckout(w http.ResponseWriter, r *http.Request) {
	cc.template.Execute(w, nil)
}