package main

import (
	"log"
	"os"
	"speedypleath/deemix/deezer"
)

func main() {
	log.SetOutput(os.Stdout)
	deezer.GetListData([]string{"355777961"})
}
