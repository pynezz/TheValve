package tui

import (
	"bufio"
	"fmt"
	"unicode/utf8"

	vault "github.com/pynezz/thevalve/vault"
)

// ╭ ─ ╮
//   │
// ╰  ╯

// The menu struct
type menu struct {
	title   string
	options []string
}

func (m *menu) printMenu() {

	// Calculate the maximum width among options
	maxWidth := 0
	for _, option := range m.options {
		if len(option) > maxWidth {
			maxWidth = utf8.RuneCountInString(option)
		}
	}

	// Calculate padding for the title to center it, considering the side borders
	titleLength := utf8.RuneCountInString(m.title)
	totalWidth := maxWidth + 4 // Adding padding for the sides
	if titleLength > maxWidth {
		totalWidth = titleLength + 4
	}
	titlePadding := (totalWidth-titleLength)/2 - 2 // Adjust for the "─ " prefix and " ─╮" suffix

	// Print the top border
	fmt.Print("╭─")
	for i := 0; i < titlePadding; i++ {
		fmt.Print("─")
	}
	fmt.Print(m.title)
	for i := 0; i < titlePadding; i++ {
		fmt.Print("─")
	}
	if (titleLength+titlePadding*2)%2 != 0 { // Adjust for odd lengths
		fmt.Print("─")
	}
	fmt.Println("─╮")

	// Print each option with padding
	for _, option := range m.options {
		fmt.Printf("│ %s", option)
		for i := utf8.RuneCountInString(option); i < maxWidth; i++ {
			fmt.Print(" ")
		}
		fmt.Println(" │")
	}

	// Print the bottom border
	fmt.Print("╰")
	for i := 0; i < totalWidth-2; i++ {
		fmt.Print("─")
	}
	fmt.Println("╯")
}

// The TUI is the text user interface
type TUI struct {
}

// NewTUI creates a new TUI
func NewTUI() TUI {
	return TUI{}
}

// Start starts the TUI
func (t *TUI) Start() {
	fmt.Println("Starting TUI")
}

func (t *TUI) Greet() int {

	var banner_ []string = []string{
		"______________________________________________________________________",
		"████████╗██╗  ██╗███████╗    ██╗   ██╗ █████╗ ██╗    ██╗   ██╗███████╗",
		"╚══██╔══╝██║  ██║██╔════╝    ██║   ██║██╔══██╗██║    ██║   ██║██╔════╝",
		"   ██║   ███████║█████╗      ██║   ██║███████║██║    ██║   ██║█████╗  ",
		"   ██║   ██╔══██║██╔══╝      ╚██╗ ██╔╝██╔══██║██║    ╚██╗ ██╔╝██╔══╝  ",
		"   ██║   ██║  ██║███████╗     ╚████╔╝ ██║  ██║███████╗╚████╔╝ ███████╗",
		"   ╚═╝   ╚═╝  ╚═╝╚══════╝      ╚═══╝  ╚═╝  ╚═╝╚══════╝ ╚═══╝  ╚══════╝",
		"______________________________________________________________________",
	}

	for _, line := range banner_ {
		fmt.Printf("%s%s%s\n", "\033[36m ", line, " \033[0m")
	}

	m := menu{
		title: " Welcome to The Valve ",
		options: []string{
			"1. | Authenticate                 ",
			"2. | List all owners (for debug)  ",
			"3. | New vault                    ",
			" ──────────────────────────────── ",
			"0. | Exit                         ",
		},
	}
	m.printMenu()
	fmt.Print("Enter your choice: ")
	var choice int
	fmt.Scan(&choice)
	return choice
}

func (t *TUI) DisplayVault() {
	m := menu{
		title:   "Vault",
		options: []string{},
	}

	m.printMenu()
}

func (t *TUI) NewSection(v *vault.Vault) {
	if v.GetVaultCount() > 5 {
		fmt.Println("Please upgrade to premium version to add more sections")
		return
	}
	reader := bufio.Reader{}

	fmt.Print("Enter the name of the section: ")
	sectionName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the section name")
	}

	fmt.Print("Enter the name of the owner: ")
	ownerName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading the owner name")
	}

	fmt.Print("Enter the secret key: ")
	secretKey, err := reader.ReadString('\n') // Password
	if err != nil {
		fmt.Println("Error reading the secret key")
	}

	v.NewSection(sectionName, ownerName, secretKey)
}
