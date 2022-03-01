package provider

import (
	"NIXSwag/api/internal/repositories"
	"encoding/json"
	"net/http"
)

type Provider struct {
	ProductRepository *repositories.ProductDBRepository
}

func (rp Provider) ListOfAllProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listAllProducts := rp.ProductRepository.GetAllProducts()
		data, _ := json.Marshal(listAllProducts)
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Write(data)

	default:
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
	}
}
