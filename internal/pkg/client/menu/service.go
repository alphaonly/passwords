package menu


type Service interface{
	AuthorizeUser(User string,Password string) error
	AddNewUser()error	
	AddNewAccount
}





