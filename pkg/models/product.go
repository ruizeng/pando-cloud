// product is a abstract define of same devices made by some vendor
package models

import (
	"time"
)

type Product struct {
	// inner id
	ID int32
	// which vendor
	VendorID int32
	// name
	ProductName string `sql:"type:varchar(200);not null;"`
	// desc
	ProductDescription string `sql:"type:text;not null;"`
	// product key to auth a product
	ProductKey string `sql:"type:varchar(200);not null;unique;key;"`
	// product config string (JSON)
	ProductConfig string `sql:"type:text"; not null;`
	// change history
	CreatedAt time.Time
	UpdatedAt time.Time

	Devices []Device
}
