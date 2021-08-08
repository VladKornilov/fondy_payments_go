package entities

import "time"

type Purchase struct {
	PurchaseId 		string
	ProductId 		int
	UserId	 		int
	PurchaseTime 	time.Time
	PurchaseStatus 	string
}

const CreatePurchases =
	"CREATE TABLE IF NOT EXISTS purchases (\n" +
	"purchase_id STRING PRIMARY KEY,\n" +
	"product_id INTEGER NOT NULL,\n" +
	"user_id INTEGER NOT NULL,\n" +
	"purchase_time TIMESTAMP NOT NULL \n" +
	"DEFAULT(DATETIME('now','localtime')),\n" +
	"purchase_status VARCHAR(30) NOT NULL DEFAULT('failure'),\n\n" +

	"CONSTRAINT product_id_ref \n" +
	"FOREIGN KEY (product_id) \n" +
	"REFERENCES products(product_id),\n\n" +
	"CONSTRAINT purchase_status_check \n" +
	"CHECK (purchase_status IN ('success', 'failure'))\n" +
	");"