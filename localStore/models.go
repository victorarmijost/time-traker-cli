package localStore

import "varmijo/time-tracker/bairestt"

const (
	StatusStored   = "stored"   //Records either pending to commint or commited
	StatusCommited = "commited" //Records that are commited
	StatusPending  = "pending"  //Records pending to commit

	StatusPool = "pool" //Records that are not attached to a date
)

type Record struct {
	Id string `json:"id"`
	bairestt.TimeRecord
}
