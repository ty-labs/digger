package utils

import (
	"fmt"
	"strings"

	"github.com/diggerhq/digger/libs/orchestrator"
	"github.com/diggerhq/digger/libs/orchestrator/scheduler"
)

func GetTerraformOutputAsCollapsibleComment(summary string, open bool) func(string) string {
	var openTag string
	if open {
		openTag = "open=\"true\""
	} else {
		openTag = ""
	}

	return func(comment string) string {
		return fmt.Sprintf(`<details %v><summary>`+summary+`</summary>

`+"```terraform"+`
`+comment+`
`+"```"+`
</details>`, openTag)
	}
}

func GetTerraformOutputAsComment(summary string) func(string) string {
	return func(comment string) string {
		return summary + "\n```terraform\n" + comment + "\n```"
	}
}

func AsCollapsibleComment(summary string, open bool) func(string) string {
	return func(comment string) string {
		return fmt.Sprintf(`<details><summary>` + summary + `</summary>
  ` + comment + `
</details>`)
	}
}

func AsComment(summary string) func(string) string {
	return func(comment string) string {
		return summary + "\n" + comment
	}
}

func CreateTableComment[K scheduler.SerializedJob | orchestrator.Job](headers []string, rows []K, rowTransformer func(int, K) []string) string {
	table := ""

	table += "| " + strings.Join(headers, " |") + " |\n"
	table += "| " + strings.Repeat("---|", len(headers)) + "\n"

	// turn each row col values into MD table row
	for i, row := range rows {
		cols := rowTransformer(i, row)
		table += "| " + strings.Join(cols, " |") + " |\n"
	}

	return table
}
