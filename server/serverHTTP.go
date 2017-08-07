package server

import (
	"github.com/gorilla/mux"
	"fmt"
	"time"
	"github.com/gorilla/handlers"
	"github.com/IhorBondartsov/ContentParser/dao/daoInterface"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type HTTPServer struct {
	Port     int32
	Host     string
	StopProg chan struct{}
	DB       daoInterface.DAOInterface
}

func NewHTTPServer(port int32, host string, stop chan struct{}) *HTTPServer {
	return &HTTPServer{
		Port:     port,
		Host:     host,
		StopProg: stop,

	}
}

func (server *HTTPServer) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("HTTPServer Failed")
			server.StopProg <- struct{}{}
		}
	}()

	r := mux.NewRouter()
	r.HandleFunc("/getContent", server.getContentHandler).Methods(http.MethodGet)
	//r.HandleFunc("/getArticle", server.getArticleHandler).Methods(http.MethodGet)
	//r.HandleFunc("/getContent", server.getContentHandler).Methods(http.MethodGet)

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../view/")))

	port := fmt.Sprint(server.Port)

	srv := &http.Server{
		Handler:      r,
		Addr:         server.Host + ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	//CORS provides Cross-Origin Resource Sharing middleware
	http.ListenAndServe(server.Host+":"+port, handlers.CORS()(r))

	go log.Fatal(srv.ListenAndServe())
}

