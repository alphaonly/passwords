package account

const (
	selectLineAccountsTable      = `SELECT account_id, user_id, password, descr, uploaded_at FROM public.accounts WHERE account_id=$1;`
	selectAllAccountsTableByUser = `SELECT account_id, user_id, password, descr, uploaded_at FROM public.accounts WHERE user_id = $1;`

	createOrUpdateIfExistsAccountsTable = `
	  INSERT INTO public.accounts (account_id, user_id, password,descr,uploaded_at) 
	  VALUES ($1, $2, $3,$4, $5)
	  ON CONFLICT (account_id,user_id) DO UPDATE 
		SET password 	= $3,
		    descr 		= $4,
			uploaded_at = $5; 
		`

	createAccountsTable = `create table public.accounts
	(	account_id varchar(40) not null, 
		user_id varchar(40) not null,
		password varchar(50),		
		descr varchar(200),
		uploaded_at TEXT not null, 
		primary key (account_id,user_id)
	);`

	checkIfAccountsTableExists = `SELECT 'public.accounts'::regclass;`
)

// -d=postgres://postgres:mypassword@localhost:5432/yandex
