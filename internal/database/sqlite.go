package database

import (
	"database/sql"
	"errors"
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

type Database struct {
	db *sql.DB
}
func (d Database) Close() error {
	return d.db.Close()
}
func (d Database) createProducts() error {
	_, err := d.db.Exec(entities.CreateProducts)
	return err
}
func (d Database) createPurchases() error {
	_, err := d.db.Exec(entities.CreatePurchases)
	return err
}
func (d Database) createIndex() error {
	_, err := d.db.Exec("CREATE INDEX purchases_products_idx ON purchases (product_id);")
	return err
}
func (d Database) Create() error {
	err := d.createProducts()
	if err != nil { return err }
	err = d.createPurchases()
	if err != nil { return err }
	err = d.createIndex()
	return err
}
func (d Database) InsertProduct(p entities.Product) error {
	_, err := d.db.Exec("INSERT INTO products VALUES ($1, $2, $3)", p.ProductId, p.ProductName, p.Price)
	return err
}
func (d Database) InsertPurchase(p entities.Purchase) error {
	var params string
	var values string
	//p.PurchaseTime = time.Now()
	if p.PurchaseStatus != "" {
		params = "(product_id, purchase_status)"
		values = "($1, $2)"
	} else {
		params = "(product_id)"
		values = "($1)"
	}

	_, err := d.db.Exec("INSERT INTO purchases " + params + " VALUES " + values,
		p.ProductId, p.PurchaseStatus, p.PurchaseTime)
	return err
}

func (d Database) GetPurchases() []entities.Purchase {
	rows, err := d.db.Query("SELECT * from purchases")
	if err != nil {
		println(err.Error())
		return nil
	}
	defer rows.Close()
	var purchases []entities.Purchase

	for rows.Next() {
		p := entities.Purchase{}
		err = rows.Scan(&p.PurchaseId, &p.ProductId, &p.PurchaseTime, &p.PurchaseStatus)
		if err != nil {
			println(err.Error())
			continue
		}
		purchases = append(purchases, p)
	}
	return purchases
}



func OpenDatabase() (Database, error){
	url, exists := os.LookupEnv("DATABASE_URL")
	if !exists {
		return Database{}, errors.New("missing DATABASE_URL variable")
	}

	path := strings.TrimPrefix(url, "sqlite://")

	// if not exist, create and initialize a new db
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return Database{}, err
		}
		db, err := sql.Open("sqlite3", path)
		base := Database{db}
		err = base.Create()
		return base, err
	}

	db, err := sql.Open("sqlite3", path)
	return Database{db}, err
}

