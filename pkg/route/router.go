package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

var Router  *mux.Router

func Initialize()  {
	Router = mux.NewRouter()
}

func Name2URL(routeName string,paris ...string) string{
	url, err := Router.Get(routeName).URL(paris...)

	if err != nil {
		return ""
	}

	return url.String()
}

func GetRouterVariable(parameter string,r *http.Request)  string {
	vars := mux.Vars(r)
	return vars[parameter]
}
