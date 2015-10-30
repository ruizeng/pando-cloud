// vendor is those who make products
package models

import (
	"time"
)

type Vendor struct {
	// inner id
	ID int32
	// vendor name
	VendorName string `sql:"type:varchar(200);not null;"`
	// vendor key
	VendorKey string `sql:"type:varchar(200);not null;key;"`
	// vendor description
	VendorDescription string `sql:"type:text;not null;"`
	// change history
	CreatedAt time.Time
	UpdatedAt time.Time

	Products []Product
}
