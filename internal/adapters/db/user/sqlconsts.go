package user

const (
	selectLineUsersTable = `SELECT user_id, password, name, surname,phone, created  FROM public.users WHERE user_id=$1;`

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
		created 	TEXT not null, 
		primary key (user_id)
	);`

	checkIfUsersTableExists = `SELECT 'public.users'::regclass;`
)

// -d=postgres://postgres:mypassword@localhost:5432/yandex
