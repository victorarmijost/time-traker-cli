package repositories

type DBRecord struct {
	Id     string  `db:"id"`
	Date   string  `db:"date"`
	Status string  `db:"status"`
	Hours  float64 `db:"hours"`
}

type DBOpenRecord struct {
	Id   string `db:"id"`
	Date string `db:"date"`
}
