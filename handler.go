package httperr

import "net/http"

type HandlerFunc func(http.ResponseWriter, *http.Request) error

func Handler(hfunc HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := hfunc(w, r); err != nil {
			Abort(w, err)
		}
	}
}
