package menu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)



func mainMenu() {

	reader := bufio.NewReader(os.Stdin)

	menuOptions := []string{
		"1: Add new account",
		"2: Edit user data",
		"3: Edit account data",
		"4: Exit",
	}

	for {
		fmt.Printf("user %v menu:","testUser")
		for _, option := range menuOptions {
			fmt.Println(option)
		}

		fmt.Print("Enter the number: ")
		input, err := reader.ReadString('n')

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
		case 1:
			fmt.Println("You selected option 1.")
		// Do something for option 1
		case 2:
			fmt.Println("You selected option 2.")
		// Do something for option 2
		case 3:
			fmt.Println("Exiting.")
			// Exit the program
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please select a number from the menu.")
		}
	}

}
