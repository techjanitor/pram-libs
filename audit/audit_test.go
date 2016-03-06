package audit

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"testing"

	"github.com/eirka/eirka-libs/db"
	e "github.com/eirka/eirka-libs/errors"
)

func TestAudit(t *testing.T) {

	var err error

	mock, err := db.NewTestDb()
	assert.NoError(t, err, "An error was not expected")

	rows := sqlmock.NewRows([]string{"role"}).AddRow(3)

	mock.ExpectPrepare(`INSERT INTO audit \(user_id,ib_id,audit_type,audit_ip,audit_time,audit_action,audit_info\)`).
		ExpectExec().
		WithArgs(1, 1, UserLog, "10.0.0.1", "NOW()", AuditEmailUpdate, "meta info")

	audit := Audit{
		User:   1,
		Ib:     1,
		Type:   UserLog,
		Ip:     "10.0.0.1",
		Action: AuditEmailUpdate,
		Info:   "meta info",
	}

	// submit audit
	err = audit.Submit()
	assert.NoError(t, err, "An error was not expected")

}
