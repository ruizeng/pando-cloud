// application is app who will use the cloud api
package models

import (
	"time"
)

type Application struct {
	// inner id
	ID int32
	// App-Key for api
	AppKey string `sql:"type:varchar(200);not null;"`
	// App-Token for web hook
	AppToken string `sql:"type:varchar(200);not null;"`
	// Report Url for web hook
	ReportUrl string `sql:"type:varchar(200);not null;"`
	// name
	AppName string `sql:"type:varchar(200);"`
	// desc
	AppDescription string `sql:"type:text;"`
	// app domain which allows wildcard string like "/*", "/vendors/12/*", "/products/10"
	AppDomain string `sql:"type:varchar(200);not null;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
