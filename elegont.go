package main

import (
	"errors"
	"fmt"
	"regexp"
	//"encoding/json"
	"os"

	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "Elegont"
	app.Usage = "2009 is the new 1984"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "lang, l",
			Value: "english",
			Usage: "language for the greeting",
			//Destination: &language, // language = c.String("lang")
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration file from `FILE`",
		},
	}

	app.Action = func(c *cli.Context) error {
		fmt.Println("boom! I say! CANONICAL")
		for _, arg := range c.Args() {
			fmt.Println(arg)
		}

		return nil
	}

	app.Run(os.Args)

}

func Dissect(ego *string, syntax Syntax) (string, error) {
	var (
		procesedLines int    = 0
		output        string = ""
	)

	for !empty(*ego) {
		whiteSpaces := nextWhiteSpaces(ego)
		output += whiteSpaces
		procesedLines += countLines(&whiteSpaces)
		code, err := nextComponent(ego, syntax)

		if err, ok := err.(*SyntaxError); ok {
			err.line = procesedLines + 1
			err.file = "fileName.ego"
			err.code = (*ego)[:regexp.MustCompile(`^[^\n]+`).FindAllStringIndex(*ego, 1)[0][1]]
			return "", err
		}

		output += fixIdent(code, whiteSpaces)
		procesedLines += countLines(&code)
	}

	return output, nil
}

/**
 * Cuts first (\s or \n)'s consecutives chars
 * @param str *string
 * @return cutted string 	removed string from str.
 */
func nextWhiteSpaces(str *string) (cutted string) {
	whiteSpaces := regexp.MustCompile(`^(\s|\n)*`)
	pos := whiteSpaces.FindStringIndex(*str)

	if pos != nil {
		cutted = (*str)[0:pos[1]]
	} else {
		cutted = ""
	}

	*str = (*str)[pos[1]:len(*str)]
	return cutted
}

/**
 * Cuts ego's first `Component`, according to `syntax`
 * Idea: Iterates on syntax[i][j] and stops once regex match's index is 0.
 * @param  str *string 	ego code
 * @return code str 		go code
 */
func nextComponent(ego *string, syntax Syntax) (code string, err error) {

	if *ego == "" {
		return "", errors.New("Empty string")
	}

	for _, variants := range syntax {
		for _, variant := range variants {
			if pos := variant.Test(ego); pos != nil {
				isFirst := pos[0] == 0
				if isFirst {
					return variant.Get(ego, pos[1], syntax)
				}
			}
		}
	}

	return "", &SyntaxError{}
}

func fixIdent(code string, whiteSpaces string) string {
	var (
		empty           = len(code) == 0
		isLastCharCurly = !empty && code[len(code)-1] == '}'
	)

	if isLastCharCurly {
		return code[:len(code)-1] + getLastLine(whiteSpaces) + "}\n"

	} else {
		return code
	}
}
