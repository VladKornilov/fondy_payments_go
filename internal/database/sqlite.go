package database

import (
	"database/sql"
	"errors"
	"github.com/VladKornilov/fondy_payments_go/internal/entities"
	"github.com/VladKornilov/fondy_payments_go/internal/logger"
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

func (d Database) createUsers() error {
	_, err := d.db.Exec(entities.CreateUsers)
	return err
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
	err := d.createUsers()
	if err != nil { return err }
	err = d.createProducts()
	if err != nil { return err }
	err = d.createPurchases()
	if err != nil { return err }
	//err = d.createIndex()
	//if err != nil { return err }
	err = d.InsertProducts()
	return err
}

func (d Database) InsertUser(u entities.User) error {
	_, err := d.db.Exec("INSERT INTO users (uuid) VALUES ($1)", u.Uuid)
	return err
}

func (d Database) InsertProducts() error {
	cnt := 0
	err := d.db.QueryRow("SELECT COUNT(*) FROM products").Scan(&cnt)
	if err != nil { return err }
	if cnt == 0 {
		logger.LogData("Inserting Products list into DB")
		_, err = d.db.Exec(entities.InsertProducts)
	}
	return err
}

func (d Database) InsertPurchase(p entities.Purchase) error {
	var params string
	var values string
	//p.PurchaseTime = time.Now()
	if p.PurchaseStatus != "" {
		params = "(purchase_id, product_id, user_id, purchase_status)"
		values = "($1, $2, $3, $4)"
	} else {
		params = "(purchase_id, product_id, user_id)"
		values = "($1, $2, $3)"
	}

	_, err := d.db.Exec("INSERT INTO purchases " + params + " VALUES " + values,
		p.PurchaseId ,p.ProductId, p.UserId, p.PurchaseStatus)
	return err
}

func (d Database) GetUserByUUID(uuid string) (entities.User, error) {
	row := d.db.QueryRow("SELECT * from users where uuid = $1", uuid)
	u := entities.User{}
	err := row.Scan(&u.UserId, &u.Uuid, &u.Diamonds)
	return u, err
}

func (d Database) UpdateUser(user entities.User) error {
	query := "UPDATE users SET diamonds = $1 WHERE uuid = $2"
	_, err := d.db.Exec(query, user.Diamonds, user.Uuid)
	return err
}

func (d Database) GetProducts() []entities.Product {
	rows, err := d.db.Query("SELECT * from products")
	if err != nil {
		println(err.Error())
		return nil
	}
	defer rows.Close()
	var products []entities.Product

	for rows.Next() {
		p := entities.Product{}
		err = rows.Scan(&p.ProductId, &p.ProductName, &p.Price, &p.Value)
		if err != nil {
			println(err.Error())
			continue
		}
		products = append(products, p)
	}
	return products
}
func (d Database) GetProductById(id string) (entities.Product, error) {
	row := d.db.QueryRow("SELECT * from products where product_id = $1", id)
	p := entities.Product{}
	err := row.Scan(&p.ProductId, &p.ProductName, &p.Price, &p.Value)
	return p, err
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

