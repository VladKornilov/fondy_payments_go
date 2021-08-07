package entities

type Product struct {
	ProductId	int
	ProductName string
	Price 		int
	Value		int
}

const CreateProducts =
	"CREATE TABLE IF NOT EXISTS products (\n    " +
	"  product_id INTEGER PRIMARY KEY AUTOINCREMENT,\n    " +
	"  product_name VARCHAR(64) UNIQUE NOT NULL,\n    " +
	"  price INTEGER NOT NULL,\n" +
	"  value INTEGER NOT NULL\n" +
	");"

const InsertProducts =
	"INSERT INTO products (product_name, price, value)\n    " +
	"VALUES\n" +
	"  ('Small Crate', 99, 100),\n" +
	"  ('Medium Crate', 399, 500),\n" +
	"  ('Big Crate', 599, 1000);"
