package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/nicolasverle/casestudy/pkg/extractor"
	"github.com/spf13/cobra"
)

var (
	output string

	// add a root command to launch the extraction of links
	// more infos about the library : https://github.com/spf13/cobra/blob/main/site/content/user_guide.md
	linkExtractor = &cobra.Command{
		Use:   "linkextractor",
		Short: "Command that will parse a set URLs and extract all the links from their contents",
		Args:  validateArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			parser := extractor.New(args)

			links, err := parser.ExtractAllLinks()
			if err != nil {
				exit(err)
			}

			if output == "json" {
				out, err := parser.ToJSON(links)
				if err != nil {
					exit(err)
				}

				fmt.Println(string(out))
			} else {
				fmt.Println(links)
			}

			return nil
		},
	}
)

// validateArgs will ensure that all the positional arguments are correct URLs...
func validateArgs(cmd *cobra.Command, args []string) error {
	errs := []error{}

	for _, a := range args {
		_, err := url.ParseRequestURI(a)
		if err != nil {
			errs = append(errs, fmt.Errorf("unable to parse Url %s, %s", a, err.Error()))
		}
	}

	if len(errs) > 0 {
		errStr := "Validation errors : \n"

		for _, e := range errs {
			errStr = fmt.Sprintf("%s\n* %s", errStr, e.Error())
		}

		return errors.New(errStr)
	}

	return nil
}

func main() {
	// execute the root command
	if err := linkExtractor.Execute(); err != nil {
		exit(err)
	}
}

func init() {
	// add the output option to change the format of the result
	linkExtractor.PersistentFlags().StringVarP(&output, "output", "o", "stdout", "format of the output")
}

func exit(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
