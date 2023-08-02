package main

import (
	"log"
	"os"
	"speedypleath/deemix/deezer"
)

func main() {
	log.SetOutput(os.Stdout)
	session := deezer.Session{}
	deezer.Ping(&session)
	deezer.UserData(&session)
	trackTokens := deezer.GetListData([]string{"355777961"}, &session)
	r := deezer.GetStreamUrl(trackTokens, session)
	log.Println(r)
}
