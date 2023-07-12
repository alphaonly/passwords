package postgres

import (
	"context"
	"log"
)

var Message = []string{
	0:  "postgres client:unable to connect to database",
	1:  "postgres client:%v table has created",
	2:  "postgres client:unable to create %v table",
	3:  "postgres client:createOrUpdateIfExistsUsersTable error",
	4:  "postgres client:QueryRow failed: %v\n",
	5:  "postgres client:RowScan error",
	6:  "postgres client:time cannot be parsed",
	7:  "postgres client:createOrUpdateIfExistsWithdrawalsTable error",
	8:  "postgres client:unable to get postgres conn",
	9:  "postgres client:unable to get postgres conn pull",
	10: "postgres client:unable to create or check tables",
}

const (
	createOrUpdateIfExistsUsersTable = `
	INSERT INTO public.users (user_id, password, name, surname, phone, created) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (user_id) DO UPDATE 
  	SET password 	= $2,
	  	name 		= $3,
		surname 	= $4,
		phone 		= $5,
		created 		= $6; 
  	`
	createUsersTable = `create table public.users
	(	user_id 	varchar(40) not null primary key,
		password  	TEXT not null,
		name 		varchar(40),
		surname 	varchar(40),		
		phone 		varchar(40),
		created 	TEXT not null
	);`

	checkIfUsersTableExists = `SELECT 'public.users'::regclass;`

	createAccountsTable = `create table public.accounts
	(	user_id varchar(40) not null,
		account_id varchar(40) not null, 
		login	varchar(40), 
		password varchar(50),		
		descr varchar(200),
		uploaded_at TEXT not null, 
		primary key (user_id,account_id)
	);`

	checkIfAccountsTableExists = `SELECT 'public.accounts'::regclass;`
)

func CreateTable(ctx context.Context, s postgresClient, checkTableSQL string, createTableSQL string) error {

	resp, err := s.pool.Exec(ctx, checkTableSQL)
	if err != nil {
		log.Println(Message[2] + err.Error())
		//create Table
		resp, err = s.pool.Exec(ctx, createTableSQL)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(Message[1] + resp.String())
	} else {
		log.Println(Message[2] + resp.String())
	}

	return err
}
