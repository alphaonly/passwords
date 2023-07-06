package configuration

import (
	"errors"
	"strings"
)


var (
	ErrAccountOrPasswordIsEmpty = errors.New("account name or password is empty")
	ErrUnknownCommand           = errors.New("unknown command")
)

func CommandParse(command string) (*CommandParseResponseDTO, error) {
	
	commandStrings:=make([]string,10)
	
	commandStrings = strings.Split(command, " ")

	switch commandStrings[0] {
	case "NEW", "UPDATE":
		{
			return &CommandParseResponseDTO{
				Command: commandStrings[0],
				Name:    commandStrings[1],
				Surname: commandStrings[2],
				Phone:   commandStrings[3],
			}, nil

		}

	case "ADD":
		{
			if commandStrings[0] == "" || commandStrings[1] == "" {
				return nil, ErrAccountOrPasswordIsEmpty
			}

			return &CommandParseResponseDTO{
				Command:     commandStrings[0],
				Account:     commandStrings[1],
				Password:    commandStrings[2],
				Description: commandStrings[3],
			}, nil

		}

	case "EDIT":
		{
			if commandStrings[0] == "" || commandStrings[1] == "" {
				return nil, ErrAccountOrPasswordIsEmpty
			}

			return &CommandParseResponseDTO{
				Command:     commandStrings[0],
				Account:     commandStrings[1],
				Password:    commandStrings[2],
				Description: commandStrings[3],
			}, nil

		}
	case "GET":
		{
			if commandStrings[0] == "" || commandStrings[1] == "" {
				return nil, ErrAccountOrPasswordIsEmpty
			}

			return &CommandParseResponseDTO{
				Command:     commandStrings[0],
				Account:     commandStrings[1],
				Password:    commandStrings[2],
				Description: commandStrings[3],
			}, nil

		}
	default:
		{
			if commandStrings[0] != "" {
				return nil, ErrUnknownCommand
			}

			return &CommandParseResponseDTO{}, nil
		}
	}
}
