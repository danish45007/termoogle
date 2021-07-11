/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/danish45007/my-google.com/domains"
	"github.com/danish45007/my-google.com/languages"
	"github.com/danish45007/my-google.com/searchengine"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// termoogleCmd represents the termoogle command
var termoogleCmd = &cobra.Command{
	Use:   "termoogle",
	Short: "google.com at your terminal",
	Long:  `termoogle lets you query at your terminal with various options so that you don't have to leave your terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		createSearchQuery()
	},
}

func Execute() {
	if err := termoogleCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type promptContent struct {
	errorMsg string
	label    string
}

func promptGetInput(pc promptContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func promptGetInputWithoutValidation(pc promptContent) (string, bool) {

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	if len(result) == 0 {
		return "", true
	}

	return result, false
}

func promptGetSelect(pc promptContent, items []string) string {
	index := -1
	var result string
	var err error

	for index < 0 {
		prompt := promptui.SelectWithAdd{
			Label:    pc.label,
			Items:    items,
			AddLabel: "Other",
		}

		index, result, err = prompt.Run()

		if index == -1 {
			items = append(items, result)
		}
	}

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func createSearchQuery() {
	searchTerm := promptContent{
		"Please provide a search term.",
		"What word would you like to search about ?",
	}
	page := promptContent{
		"Please provide a page number.",
		"What page number would you like to see?",
	}
	resultCount := promptContent{
		"Please provide result count.",
		"How many results would you like to see ?",
	}
	proxy := promptContent{
		"Please provide a proxy",
		"Would you like pass proxy IP defult it will use your IP address (press enter with empty input) ?",
	}
	domain := promptContent{
		"Please provide a domain.",
		fmt.Sprintf("What domain you want to search %s?", searchTerm),
	}

	language := promptContent{
		"Please provide a language.",
		fmt.Sprintf("What language you want to search %s in?", searchTerm),
	}

	enteredSearchTerm := promptGetInput(searchTerm)
	enteredPage := promptGetInput(page)
	enteredResultCount := promptGetInput(resultCount)
	enteredProxy, empty := promptGetInputWithoutValidation(proxy)
	enteredDomain := promptGetSelect(domain, domains.GetListOfDomains())
	enteredLanguage := promptGetSelect(language, languages.GetListOfLanguages())
	var Proxy interface{}
	if empty {
		Proxy = nil
	} else {
		Proxy = enteredProxy
	}
	languageCode, _ := languages.GetGoogleLanguageCode(enteredLanguage)
	intPage, _ := strconv.Atoi(enteredPage)
	intResultCount, _ := strconv.Atoi(enteredResultCount)
	searchengine.SearchEngine(enteredSearchTerm, enteredDomain, languageCode, intPage, intResultCount, Proxy, 5)

}
