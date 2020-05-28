package leaky

import (
	"net/http"
)

func Ping() bool {
	res, err := http.Get("https://jameshunt.us/")
	if err != nil {
		return false
	}

	return res.StatusCode == 200
}
