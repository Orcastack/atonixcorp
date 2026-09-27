package main

import (
	"atonixcorp/cloud/devtools/cli/router"
	"os"
)

func main() {
	router.Dispatch(os.Args)
}
