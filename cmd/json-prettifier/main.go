package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/conflabermits/json-prettifier/common"
)

type Options struct {
	Url string
}

func parseArgs() (*Options, error) {
	options := &Options{}

	flag.StringVar(&options.Url, "url", "", "url to check")
	//https://www.officedrummerwearswigs.com/api/trpc/songRequest.getLatest
	flag.Usage = func() {
		fmt.Printf("Usage: json-prettifier [options]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	return options, nil
}

func main() {
	options, err := parseArgs()
	if err != nil {
		os.Exit(1)
	}

	if len(options.Url) > 0 {
		response := common.Http_req(options.Url)
		output := common.Parse_json(response)
		fmt.Println(output)
	}

}
