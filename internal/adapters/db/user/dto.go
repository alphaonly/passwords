package user

import "database/sql"

type DBUsersDTO struct {
	userID   sql.NullString
	password sql.NullString
	name     sql.NullString
	surname  sql.NullString
	phone    sql.NullString
	created  sql.NullString
}
