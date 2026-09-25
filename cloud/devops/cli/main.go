package main

import (
	"atonixcorp/cloud/devops/cli/router"
	"os"
)

func main() {
	router.Dispatch(os.Args)
}
