package entity

import (
	"fmt"
	"github.com/AmadlaOrg/hery/entity/query/token"
	"github.com/PaesslerAG/jsonpath"
	"github.com/itchyny/gojq"
	"strings"
	"v.io/v23/glob"
)

// Query
func Query(args []string) {
	for _, arg := range args {
		if arg == "!" {
			continue
		}

		tokens, err := token.Lex(arg)
		if err != nil {
			fmt.Println("Error tokenizing query:", err)
			return
		}

		parsedQueries, order, err := token.Parse(tokens)
		if err != nil {
			fmt.Println("Error parsing tokens:", err)
			return
		}

		for _, key := range order {
			value := parsedQueries[key]

			fmt.Printf("%s: %s\n", key, value)

			if strings.ContainsFunc(key, func(r rune) bool {
				return strings.EqualFold(string(r), "entities")
			}) {
				parse, err := glob.Parse(value)
				if err != nil {
					return
				}
				println(parse)
			} else if strings.EqualFold(key, "jsonpath") {
				//jsonpath.New()
				jsonpath.Language()
			} else if strings.EqualFold(key, "jq") {
				_, err := gojq.Parse("")
				if err != nil {
					return
				}
			}

			// Process each key-value pair as needed
		}
	}
}
