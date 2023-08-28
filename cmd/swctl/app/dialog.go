package app

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
)

type PromptContent struct {
	Label    string
	MsgError string
	// attribute items is needed only for select and multi-select
	Items []string
}

func (pc *PromptContent) PromptGetInput() []string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.MsgError)
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green}}",
		Invalid: "{{ . | red}}",
		Success: "{{ . | bold}}",
	}

	prompt := promptui.Prompt{
		Label:     pc.Label,
		Validate:  validate,
		Templates: templates,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	fmt.Printf("Input: %s\n", result)
	return []string{result}
}

func (pc *PromptContent) PromptGetSelect() []string {
	index := -1
	var items []string
	var result string
	var err error

	for index < 0 {
		prompt := promptui.Select{
			Label: pc.Label,
			Items: pc.Items,
		}
		index, result, err = prompt.Run()
		if index == -1 {
			items = append(items, result)
		}

	}
	if err != nil {
		fmt.Println("Prompt failed %v\n", err)
	}
	fmt.Printf("Input: %s\n", result)
	return []string{result}
}

type localItem struct {
	Id         string
	IsSelected bool
}

func (pc *PromptContent) PromptGetMultiSelect() []string {
	fmt.Print(pc.Label)
	const DoneId = "done"
	var usedStrings = pc.Items

	var items []*localItem
	for _, name := range usedStrings {
		items = append(items, &localItem{Id: name})
	}
	items = append(items, &localItem{Id: DoneId})
	selectedItems := localPromptGetMultiSelect(0, items, pc.Label)
	var stringerizedSelectedItems []string
	for _, buf := range selectedItems {
		stringerizedSelectedItems = append(stringerizedSelectedItems, buf.Id)
	}
	return stringerizedSelectedItems
}

func localPromptGetMultiSelect(selectedPos int, allItems []*localItem, label string) []*localItem {
	const doneId = "done"
	templates := &promptui.SelectTemplates{
		Label: `{{ if .IsSelected }}
		         [#]   
				{{ else }}
				 [#]
		        {{ end }} {{ .Id }}`,
		Active:   "{{if .IsSelected}} [#] {{else}} [ ] {{end}}{{ .Id | red }}",
		Inactive: "{{if .IsSelected}} [#] {{else}} [ ] {{end}}{{ .Id | cyan }}",
	}

	prompt := promptui.Select{
		Label:        label,
		Items:        allItems,
		Templates:    templates,
		Size:         5,
		CursorPos:    selectedPos,
		HideSelected: true,
	}

	selectionIdx, buf, err := prompt.Run()
	_ = selectionIdx
	_ = buf
	if err != nil {
		return nil
	}
	chosenItem := allItems[selectionIdx]

	if chosenItem.Id != doneId {
		// If the user selected something other than "Done",
		// toggle selection on this item and run the function again.
		chosenItem.IsSelected = !chosenItem.IsSelected
		return localPromptGetMultiSelect(selectionIdx, allItems, label)
	}

	var selectedItems []*localItem
	for _, item := range allItems {
		if item.IsSelected {
			selectedItems = append(selectedItems, item)
		}
	}
	return selectedItems

}
