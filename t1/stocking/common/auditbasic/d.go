package auditbasic

import (
	"t1/stocking/common/createbasic"
)

type AuditBasic struct {
	createbasic.CreateBasic
	AuditTime string
	AuditUser string
}
