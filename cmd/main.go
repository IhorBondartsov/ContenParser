package main

import (
	log "github.com/Sirupsen/logrus"
	parser2 "github.com/IhorBondartsov/ContenParser/parser"
	"github.com/IhorBondartsov/ContenParser/server"
	"github.com/IhorBondartsov/ContenParser/dao/implementation"
)

func main() {
	log.Info("Parser has starting")

	dbClient := implementation.SQL{Password: basePassword, URI: baseURI,
				User:baseUser, TypeConn: baseTypeConn, NameDB: NameDB}

	parser := parser2.Parser{}
	parser.Connect(&dbClient)
	parser.StartParse()

	// Start HTTP Server
	httpServer := server.NewHTTPServer(httpPort, httpHost, &dbClient)
	httpServer.Run()
}
