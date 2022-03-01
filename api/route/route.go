package route

import (
	"NIXSwag/api/internal/repositories/provider"
	"net/http"
)

func Route(provider provider.Provider, mux *http.ServeMux) {
	mux.HandleFunc("/list", provider.ListOfAllProducts)
}
