package account

import "database/sql"

type DBAccountsDTO struct {
	UserID    sql.NullString
	AccountID sql.NullString
	Login     sql.NullString
	Password  sql.NullString
	Descr     sql.NullString
	CreatedAt sql.NullString
}
