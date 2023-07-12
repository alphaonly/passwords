package menu

// menu - package that describes the main function of working with app menu
import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"passwords/internal/client"
	clientService "passwords/internal/client/service"
	configuration "passwords/internal/configuration/client"
	accountDomain "passwords/internal/domain/account"
	userDomain "passwords/internal/domain/user"
	"passwords/internal/pkg/common/logging"
	"strconv"
	"strings"
	"time"
)

// Typical MainMenu interface errors
var (
	ErrAccountNil     = errors.New("account is nil")
	ErrUserNil        = errors.New("user is nil")
	ErrAccountsList   = errors.New("unable to get account lists for menu")
	ErrEditAccount    = errors.New("unable edit account from menu")
	ErrReadNewAccount = errors.New("unable read a new account in menu")
	ErrAddNewAccount  = errors.New("unable add a new account, read from menu")
)

// MainMenu - an interface to gives an API how the client's console menu should be managed
type MainMenu interface {
	ReadNewAccount(ctx context.Context) (*accountDomain.Account, error) // Reads a new account data for user from reader
	EditAccount(ctx context.Context, acc *accountDomain.Account) error  // Edits an account data reading new values from reader
	EditUser(ctx context.Context) error                                 // Edits a user data reading new values from reader
	Start(ctx context.Context) error                                    // Starts the main menu to show
}

// mainMenu  - The MainMenu implementation
type mainMenu struct {
	user          *userDomain.User
	cfg           *configuration.ClientConfiguration
	clientService clientService.Service
	reader        *bufio.Reader
}

// New - a factory for mainMenu implementation receives the user instance which uses the client
func New(
	user *userDomain.User,
	cfg *configuration.ClientConfiguration,
	clientService clientService.Service,
	reader *bufio.Reader,
) MainMenu {

	return &mainMenu{
		user:          user,
		cfg:           cfg,
		clientService: clientService,
		reader:        reader,
	}
}

// ReadNewAccount - reads account's data from reader field by field
func (m mainMenu) ReadNewAccount(ctx context.Context) (acc *accountDomain.Account, err error) {
	acc = &accountDomain.Account{}
	acc.User = m.user.User

	fmt.Println("Enter a new account's data:")
	//read account id
	err = readStringField(m.reader, "account name(without spaces)", &acc.Account)
	if err != nil {
		return nil, err
	}
	acc.Account = strings.TrimSpace(acc.Account)
	//read account login
	err = readStringField(m.reader, "account login(without spaces)", &acc.Login)
	if err != nil {
		return nil, err
	}
	acc.Login = strings.TrimSpace(acc.Login)
	//read account password
	err = readStringField(m.reader, "account password", &acc.Password)
	if err != nil {
		return nil, err
	}
	//read account description
	err = readStringField(m.reader, "account description text", &acc.Description)
	if err != nil {
		return nil, err
	}

	return acc, nil
}

// EditAccount - edits account's data from reader field by field
func (m mainMenu) EditAccount(ctx context.Context, acc *accountDomain.Account) (err error) {
	if acc == nil {
		return ErrAccountNil
	}
	fmt.Printf("You edit the following account %v  data:", acc.Account)
	fmt.Println(acc)

	//read account login
	err = readStringField(m.reader, fmt.Sprintf("login without spaces(current:%v)", acc.Login), &acc.Login)
	if err != nil {
		return err
	}
	acc.Login = strings.TrimSpace(acc.Login)
	//read account password
	err = readStringField(m.reader, fmt.Sprintf("account password(current:%v)", acc.Password), &acc.Password)
	if err != nil {
		return err
	}
	//read account description
	err = readStringField(m.reader, fmt.Sprintf("account description text(current:%v)", acc.Password), &acc.Description)
	if err != nil {
		return err
	}

	return nil
}

// EditUser - edits  user's data reading new values from reader and passes it to user variable
func (m mainMenu) EditUser(ctx context.Context) (err error) {
	if m.user == nil {
		return ErrUserNil
	}
	fmt.Println("You edit the following user's  data:")
	fmt.Println(m.user)

	fmt.Println("Enter the user's data(put sign * to remain field with the current value):")
	//edit user password
	err = readStringField(m.reader, fmt.Sprintf("password (current:%v)", m.user.Password), &m.user.Password)
	if err != nil {
		return err
	}
	//edit username
	err = readStringField(m.reader, fmt.Sprintf("name (current:%v)", m.user.Name), &m.user.Name)
	if err != nil {
		return err
	}
	//edit user surname
	err = readStringField(m.reader, fmt.Sprintf("surname (current:%v)", m.user.Surname), &m.user.Surname)
	if err != nil {
		return err
	}
	//edit user surname
	err = readStringField(m.reader, fmt.Sprintf("phone number (current:%v)", m.user.Phone), &m.user.Phone)
	if err != nil {
		return err
	}
	m.user.Phone = strings.TrimSpace(m.user.Phone)

	//renew password in conf
	m.cfg.Password = m.user.Password
	return nil
}

func (m mainMenu) Start(ctx context.Context) error {
	if m.reader == nil {
		m.reader = bufio.NewReader(os.Stdin)

	}

	menuOptions := []string{
		"1: Add new account",
		"2: Edit account data",
		fmt.Sprintf("3: Edit user %v data", m.user.User),
		"4: Exit",
	}
	for {
		client.ClearScreen()
		fmt.Printf("User %v menu:\n\n", m.user.User)

		for _, option := range menuOptions {
			fmt.Println(option)
		}
		fmt.Print("Enter the number: ")
		input, err := m.reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)

		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch choice {
		case 1: //Add new account
			{
				//read new account's data
				acc, err := m.ReadNewAccount(ctx)
				if err != nil {
					return fmt.Errorf(ErrReadNewAccount.Error()+"%w", err)
				}
				err = m.clientService.AddNewAccount(ctx, *acc)
				if err != nil {
					return fmt.Errorf(ErrAddNewAccount.Error()+" %w", err)
				}
				fmt.Printf("Account %v added for user %v succesfully", acc.Account, acc.User)
				time.Sleep(3 * time.Second)
			}

		case 2: //Edit account data
			{
				//Get all accounts of user
				accounts, err := m.clientService.GetAllAccounts(ctx, m.user.User)
				if err != nil {
					er := fmt.Errorf("show all accounts error: %w", ErrAccountsList)
					return fmt.Errorf(er.Error()+" %w", err)
				}
				//Show all accounts of the user
				fmt.Println(accounts)
				//Pick the account by entering the account ID
				fmt.Print("Enter the account ID to correct:")
				pickedAccountID, err := m.reader.ReadString('\n')
				if err != nil {
					fmt.Println("Invalid input. Please enter an account ID.")
				}
				pickedAccountID = strings.Replace(pickedAccountID, "\r\n", "", 1)
				account, ok := accounts[pickedAccountID]
				if !ok {
					fmt.Println("Invalid input. Please, enter an account ID from the given list.")
				}
				//Edit picked account's data
				err = m.EditAccount(ctx, &account)
				if err != nil {
					return fmt.Errorf(ErrEditAccount.Error()+" %w", err)
				}
				//Save edited data
				err = m.clientService.AddNewAccount(ctx, account)
				if err != nil {
					return fmt.Errorf(ErrEditAccount.Error()+" %w", err)
				}
				fmt.Printf("account %v edited successfully", account.Account)
				time.Sleep(3 * time.Second)
			}
		case 3: //Edit user data
			{
				err = m.EditUser(ctx)
				if err != nil {
					return err
				}
				//save edited data
				err = m.clientService.AddNewUser(ctx, *m.user)
				if err != nil {
					logging.LogFatal(err)
				}
				fmt.Printf("Edited user's new data saved:")
				fmt.Println(m.user)
				time.Sleep(3 * time.Second)

			}
		case 4: //Exit menu program
			{
				fmt.Println("Exiting menu...")
				// Exit the menu
				return nil
			}
		default:
			{
				fmt.Println("Invalid choice. Please select a number from the menu.")
			}
		}

	}

	return nil
}

// readStringField - read a string field value from reader and pass to variable
func readStringField(reader *bufio.Reader, name string, value *string) (err error) {
	//fmt.Print("Name(without spaces):")
	fmt.Print(name + ":")

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading value %v : %v", value, err)
		return err
	}
	input = strings.Replace(input, "\r\n", "", 1)
	if input != "*" {
		*value = input
	}
	fmt.Println()
	return err
}
