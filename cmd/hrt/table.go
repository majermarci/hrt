package main

import (
	"fmt"
	"strings"
)

func printTable(responses []Response) {
	// Determine the width of each column
	testWidth := len("Test")
	statusWidth := len("Status Code")
	bodyWidth := len("Response Body")

	for _, response := range responses {
		if len(response.Test) > testWidth {
			testWidth = len(response.Test)
		}
		if len(response.StatusCode) > statusWidth {
			statusWidth = len(response.StatusCode)
		}
		lines := strings.Split(response.ResponseBody, "\n")
		for _, line := range lines {
			if len(line) > bodyWidth {
				bodyWidth = len(line)
			}
		}
	}

	// Print the table header
	printLine(testWidth, statusWidth, bodyWidth, "╭", "┬", "╮")
	fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, "Test", statusWidth, "Status Code", bodyWidth, "Response Body")
	printLine(testWidth, statusWidth, bodyWidth, "├", "┼", "┤")

	// For each response, print a table row
	for i, response := range responses {
		lines := strings.Split(response.ResponseBody, "\n")
		for j, line := range lines {
			if j == 0 {
				fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, response.Test, statusWidth, response.StatusCode, bodyWidth, line)
			} else {
				fmt.Printf("│ %-*s │ %-*s │ %-*s │\n", testWidth, "", statusWidth, "", bodyWidth, line)
			}
		}
		// Only print the line separator if it's not the last response
		if i < len(responses)-1 {
			printLine(testWidth, statusWidth, bodyWidth, "├", "┼", "┤")
		}
	}

	// Print the table footer
	printLine(testWidth, statusWidth, bodyWidth, "╰", "┴", "╯")
}

func printLine(testWidth, statusWidth, bodyWidth int, start, middle, end string) {
	fmt.Print(start)
	fmt.Print(strings.Repeat("─", testWidth+2))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", statusWidth+2))
	fmt.Print(middle)
	fmt.Print(strings.Repeat("─", bodyWidth+2))
	fmt.Println(end)
}
