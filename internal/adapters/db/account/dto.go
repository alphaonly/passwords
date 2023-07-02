package account

import "database/sql"

type DBAccountsDTO struct {
	AccountID sql.NullString
	userID    sql.NullString
	password  sql.NullString
	descr     sql.NullString
	createdAt sql.NullString
}
