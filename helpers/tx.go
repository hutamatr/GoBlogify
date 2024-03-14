package helpers

import (
	"database/sql"
)

func TxRollbackCommit(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errRollback := tx.Rollback()
		PanicError(errRollback)
		panic(err)
	} else {
		errCommit := tx.Commit()
		PanicError(errCommit)
	}
}
