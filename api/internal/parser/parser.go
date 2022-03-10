package parser

import (
	"NIXSwag/api/internal/database"
	"NIXSwag/api/internal/models"
	pkg2 "NIXSwag/api/pkg"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Parser(conn *sql.DB, err error) error {
	pool := pkg2.NewWorkerPool(4)
	wg := sync.WaitGroup{}
	wg.Add(pool.Count)
	//FILL DATABASE
	for i := 0; i < pool.Count; i++ {
		//log.Println("Starting Routine...")
		go pool.Run(&wg, func(rest models.Restaurant) {
			_, err = conn.Exec(
				"INSERT INTO restaurant VALUE (?, ?, ?, ?, ?, ?)",
				rest.Id, rest.Name,
				rest.Image, rest.Type,
				rest.WorkingHours.Opening,
				rest.WorkingHours.Closing)
			if err != nil {
				if strings.HasPrefix(err.Error(), "Error 1062") {
					return
				} else {
					log.Fatal(err)
				}
			}
			for _, prod := range rest.Menu {
				_ = database.RowId(
					conn,
					"SELECT id FROM product WHERE name = ?",
					"INSERT INTO product VALUE (?, ?, ?, ?, ?)",
					prod.Id, prod.Name, prod.Price, prod.Image, prod.Type)
				if err != nil {
					log.Println(err)
				}

				for _, ing := range prod.Ingredients {
					ingId := database.RowId(
						conn,
						"SELECT id FROM ingredient WHERE name = ?",
						"INSERT INTO ingredient(name) VALUE (?)",
						ing)
					_, err = conn.Exec(
						"INSERT INTO product_ingredient VALUE (?, ?)",
						prod.Id, ingId)
					if err != nil {
						log.Println(err)
					}

				}
			}
		})
	}
	//PARSE JSON
	client := http.DefaultClient
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		"http://foodapi.true-tech.php.nixdev.co/suppliers", nil,
	)
	if err != nil {
		log.Fatalln(err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("%v", err)
	}

	suppliersMap := make(map[string][]models.Restaurant, 0)
	err = json.NewDecoder(res.Body).Decode(&suppliersMap)
	if err != nil {
		log.Fatalln(err)
	}
	err = res.Body.Close()
	if err != nil {
		log.Fatalln(err)
	}
	suppliers := make([]models.Restaurant, 0)
	suppliers = suppliersMap["suppliers"]

	for i := range suppliers {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		req, err = http.NewRequestWithContext(
			ctx, http.MethodGet,
			"http://foodapi.true-tech.php.nixdev.co/suppliers/"+strconv.Itoa(suppliers[i].Id)+"/menu",
			nil,
		)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = client.Do(req)
		if err != nil {
			log.Fatalf("%v", err)
		}

		menuMap := make(map[string][]models.Product, 0)
		err := json.NewDecoder(res.Body).Decode(&menuMap)
		if err != nil {
			return nil
		}
		suppliers[i].Menu = menuMap["menu"]
		pool.Sender <- suppliers[i]
		err = res.Body.Close()
		if err != nil {
			log.Fatalln(err)
		}
		cancel()
	}
	pool.Stop()
	wg.Wait()
	//UPDATE PRICES
	//time.Sleep(time.Minute)
	for i, sup := range suppliers {
		for j, prod := range sup.Menu {
			ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
			req, err = http.NewRequestWithContext(ctx, http.MethodGet,
				"http://foodapi.true-tech.php.nixdev.co/suppliers/"+
					strconv.Itoa(sup.Id)+"/menu/"+strconv.Itoa(prod.Id),
				nil)
			res, err = client.Do(req)
			if err != nil {
				log.Println("Update error: " + err.Error())
			}
			var p models.Product
			err = json.NewDecoder(res.Body).Decode(&p)
			if err != nil {
				log.Println("Update error: " + err.Error())
			}
			if p.Price != prod.Price {
				_, err = conn.Exec(
					"UPDATE product SET price = ? WHERE id = ?",
					p.Price, p.Id)
				if err != nil {
					log.Println(err)
				}
				fmt.Println(p.Name, "price edited", prod.Price, "->", p.Price)
				suppliers[i].Menu[j].Price = p.Price
			} else {
				fmt.Println(p.Name, "not edit with price", p.Price)
			}
			err := res.Body.Close()
			if err != nil {
				log.Fatalln(err)
			}
			cancel()
		}
	}
	return err
}
