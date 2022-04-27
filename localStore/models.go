package localStore

const (
	StatusStored   = "stored"   //Records either pending to commint or commited
	StatusCommited = "commited" //Records that are commited
	StatusPending  = "pending"  //Records pending to commit

	StatusPool = "pool" //Records that are not attached to a date
)
