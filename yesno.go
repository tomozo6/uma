package uma

import (
	"log"

	"github.com/manifoldco/promptui"
)

func YesNo(label string) bool {
	prompt := promptui.Select{
		Label: label + "[Yes/No]",
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}
