package repositories

type DBRecord struct {
	Id    string  `db:"id"`
	Date  string  `db:"date"`
	Hours float64 `db:"hours"`
}

type DBOpenRecord struct {
	Date string `db:"date"`
}

type DBDebt struct {
	StartDate string  `db:"date"`
	Hours     float64 `db:"hours"`
}
