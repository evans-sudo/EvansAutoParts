package ctrl

import (
	"html/template"
	"net/http"
)


type checkoutController struct {
	template * template.Template
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
	cc.template.Execute(w, nil)
}

func (cc *checkoutController) PostCheckout(w http.ResponseWriter, r *http.Request) {
	cc.template.Execute(w, nil)
}