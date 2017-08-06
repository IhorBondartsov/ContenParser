package main

import (
	log "github.com/Sirupsen/logrus"
	parser2 "github.com/IhorBondartsov/ContentParser/parser"
)



func main() {
	log.Info("Parser has starting")

	parser := parser2.Parser{}
	parser.Connect(basePassword, baseURI,baseUser,baseTypeConn,NameDB)
	parser.StartParse()

}


