package usercommon

import (
	"t1/stocking/common/auditbasic"
)

type UserCommon struct {
	UserId int32
	auditbasic.AuditBasic
}
