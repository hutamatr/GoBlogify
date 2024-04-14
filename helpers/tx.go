package helpers

import (
	"database/sql"
)

func TxRollbackCommit(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errRollback := tx.Rollback()
		PanicError(errRollback, "failed to rollback transaction")
		panic(err)
	} else {
		errCommit := tx.Commit()
		PanicError(errCommit, "failed to commit transaction")
	}
}
