package localStore

import "time"

const (
	StatusStored   = "stored"   //Records either pending to commint or commited
	StatusCommited = "commited" //Records that are commited
	StatusPending  = "pending"  //Records pending to commit

	StatusPool = "pool" //Records that are not attached to a date
)

type Record struct {
	Id       string    `json:"id"`
	TaskName string    `json:"taskName"`
	Date     time.Time `json:"date"`
	Hours    float64   `json:"hours"`
	Comments string    `json:"comments"`
}
