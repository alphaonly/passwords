package common

import "fmt"

func PrintBuildTags(version string, date string) {
	fmt.Printf("Build version: %v\n", version)
	fmt.Printf("Build date: %v\n", date)
}
