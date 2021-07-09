package entities

type Product struct {
	ProductId	int
	ProductName string
	Price 		int
}

const CreateProducts =
	"CREATE TABLE products (\n    " +
	"  product_id INTEGER PRIMARY KEY AUTOINCREMENT,\n    " +
	"  product_name VARCHAR(64) UNIQUE NOT NULL,\n    " +
	"  price INTEGER NOT NULL\n" +
	");"