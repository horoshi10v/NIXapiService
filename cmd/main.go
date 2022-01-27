package main

import (
	"NIXSwag/database"
	"NIXSwag/pkg"
	"context"
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

func main() {
	pool := pkg.NewWorkerPool(4)
	wg := sync.WaitGroup{}
	conn, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	//database.DeleteTables(conn)
	defer conn.Close()

	wg.Add(pool.Count)
	for i := 0; i < pool.Count; i++ {
		fmt.Println("start")
		go pool.Run(&wg, func(rest pkg.Restaurant) {
			_, err = conn.Exec("INSERT INTO restaurant VALUE (?, ?, ?, ?, ?, ?)",
				rest.Id, rest.Name, rest.Image, rest.Type, rest.WorkingHours.Opening,
				rest.WorkingHours.Closing)
			if err != nil {
				if strings.HasPrefix(err.Error(), "Error 1062") {
					return
				} else {
					log.Fatal(err)
				}
			}
			//restID:=database.RowId(conn,"SELECT id FROM restaurant WHERE name = ?",
			//	"INSERT INTO restaurant VALUE (?, ?, ?, ?, ?, ?)",
			//	rest.Id, rest.Name, rest.Image, rest.Type, rest.WorkingHours.Opening,
			//	rest.WorkingHours.Closing)

			for _, prod := range rest.Menu {
				//prodTypeId := database.RowId(conn, "SELECT id FROM product WHERE name = ?",
				//	"INSERT INTO product(type) VALUE (?)", prod.Type)
				_, err = conn.Exec(
					"INSERT INTO product VALUE (?, ?, ?, ?, ?)",
					prod.Id, prod.Name, prod.Price, prod.Image, prod.Type)
				if err != nil {
					if strings.HasPrefix(err.Error(), "Error 1062") {
						continue
					} else {
						log.Fatal(err)
					}
				}

				for _, ing := range prod.Ingredients {
					ingId := database.RowId(conn, "SELECT id FROM ingredient WHERE name = ?",
						"INSERT INTO ingredient(name) VALUE (?)", ing)
					_, err = conn.Exec("INSERT INTO product_ingredient VALUE (?, ?)", prod.Id, ingId)

					if err != nil {
						log.Println(err)
					}
				}
			}
		})
	}

	client := http.DefaultClient

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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

	suppliersMap := make(map[string][]pkg.Restaurant, 0)
	err = json.NewDecoder(res.Body).Decode(&suppliersMap)
	if err != nil {
		fmt.Println(err)
	}

	res.Body.Close()
	suppliers := make([]pkg.Restaurant, 0)
	suppliers = suppliersMap["suppliers"]

	for i, _ := range suppliers {
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)

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
		menuMap := make(map[string][]pkg.Product, 0)
		json.NewDecoder(res.Body).Decode(&menuMap)
		suppliers[i].Menu = menuMap["menu"]
		pool.Sender <- suppliers[i]
		res.Body.Close()
		cancel()
	}
	pool.Stop()
	wg.Wait()

	for {
		time.Sleep(time.Minute)
		for i, sup := range suppliers {
			for j, prod := range sup.Menu {
				ctx, cancel = context.WithTimeout(context.Background(), time.Second)
				req, err = http.NewRequestWithContext(ctx, http.MethodGet,
					"http://foodapi.true-tech.php.nixdev.co/suppliers/"+
						strconv.Itoa(sup.Id)+"/menu/"+strconv.Itoa(prod.Id),
					nil)
				res, err = client.Do(req)
				if err != nil {
					log.Println("update error: " + err.Error())
				}
				var p pkg.Product
				err = json.NewDecoder(res.Body).Decode(&p)
				if err != nil {
					log.Println("update error: " + err.Error())
				}
				if p.Price != prod.Price {
					_, err = conn.Exec("UPDATE product SET price = ? WHERE id = ?", p.Price, p.Id)
					if err != nil {
						log.Println(err)
					}
					fmt.Println(p.Name, " edit from price", prod.Price, " to ", p.Price)
					suppliers[i].Menu[j].Price = p.Price
				} else {
					fmt.Println(p.Name, " not edit with price", p.Price)
				}
				res.Body.Close()
				cancel()
			}
		}
	}
}
