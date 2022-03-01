package repositories

import (
	"NIXSwag/api/internal/models"
	"database/sql"
	logger "github.com/horoshi10v/loggerNIX/v4"
	"log"
)

type ProductDBRepository struct {
	conn   *sql.DB
	TX     *sql.Tx
	logger *logger.Logger
}

func (p ProductDBRepository) GetAllProducts() []models.Product {
	var product models.Product
	var listProd []models.Product

	rows, err := p.conn.Query("SELECT id, name, price, image, type FROM product")
	if err != nil {
		p.logger.Error("Can't query products")
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Image, &product.Type)
		if err != nil {
			log.Println(err)
			return listProd

		}
		listProd = append(listProd, product)
	}
	return listProd
}
func NewProductRepo(conn *sql.DB, TX *sql.Tx, logger *logger.Logger) *ProductDBRepository {
	return &ProductDBRepository{conn: conn, TX: TX, logger: logger}
}
