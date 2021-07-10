package entities

type User struct {
	UserId	 int
	Uuid	 string
	Diamonds int
}

const CreateUsers =
	"CREATE TABLE users (\n    " +
		"  user_id INTEGER PRIMARY KEY AUTOINCREMENT,\n    " +
		"  uuid VARCHAR(64) UNIQUE,\n    " +
		"  diamonds INTEGER NOT NULL DEFAULT(0)\n" +
		");"
