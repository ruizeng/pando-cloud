// rule is used for automated works such as timers, ifttts.
package models

import (
	"time"
)

type Rule struct {
	// inner id
	ID int64
	// which device the rule belongs to
	DeviceID int64
	// rule type, timmer | ifttt
	RuleType string `sql:"type:varchar(20);not null;"`
	// which action triggers the rule
	Trigger string `sql:"type:varchar(200);not null;"`
	// where to send
	Target string `sql:"type:varchar(200);not null;"`
	// what action to take.
	Action string `sql:"type:varchar(200);not null;"`
	// if trigger once
	Once bool `sql:"default(false)"`
	// change history
	CreatedAt time.Time
	UpdatedAt time.Time
}
