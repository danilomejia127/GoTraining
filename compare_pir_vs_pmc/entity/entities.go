package entity

import "time"

type CustomPMEvent struct {
	ID          int64     `db:"id"`
	Source      string    `db:"source"`
	Request     string    `db:"request"`
	Body        string    `db:"body"`
	Header      string    `db:"header"`
	CodeEvent   string    `db:"code_event"`
	DateCreated time.Time `db:"date_created"`
	SellerID    int64     `db:"seller_id"`
	SiteID      string    `db:"site_id"`
	EntityID    *int64    `db:"entity_id"`
	ScopeName   *string   `db:"scope_name"`
}
