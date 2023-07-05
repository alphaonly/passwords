package account

const (
	selectLineAccountsTable      = `SELECT user_id, account_id, login, password, descr, uploaded_at FROM public.accounts WHERE account_id=$1;`
	selectAllAccountsTableByUser = `SELECT user_id, account_id, login, password, descr, uploaded_at FROM public.accounts WHERE user_id = $1;`

	createOrUpdateIfExistsAccountsTable = `
	  INSERT INTO public.accounts (user_id, account_id, login, password,descr,uploaded_at) 
	  VALUES ($1, $2, $3,$4, $5, $6)
	  ON CONFLICT (user_id,account_id) DO UPDATE 
		SET login		= $3
			password 	= $4,
		    descr 		= $5,
			uploaded_at = $6; 
		`

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

// -d=postgres://postgres:mypassword@localhost:5432/yandex
