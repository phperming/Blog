package route

import (
	"github.com/gorilla/mux"
	"net/http"
)

var route  *mux.Router

func SetRoute(r *mux.Router)  {
	route = r
}

func Name2URL(routeName string,paris ...string) string{
	url, err := route.Get(routeName).URL(paris...)

	if err != nil {
		return ""
	}

	return url.String()
}

func GetRouterVariable(parameter string,r *http.Request)  string {
	vars := mux.Vars(r)
	return vars[parameter]
}
