package webhandler

import (
	"net/http"
)

type Welfare struct{}

func (this *Welfare) Handler(do string, param map[string]interface{}, w http.ResponseWriter, r *http.Request) (interface{}, int) {

	return "no func name:" + do, 999

}
