package server

import (
	"github.com/gorilla/mux"
	"fmt"
	"time"
	"github.com/gorilla/handlers"
	"github.com/IhorBondartsov/ContenParser/dao/daoInterface"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"encoding/json"
)

type HTTPServer struct {
	Port int32
	Host string
	DB   daoInterface.DAOInterface
}

func NewHTTPServer(port int32, host string, dbClient daoInterface.DAOInterface) *HTTPServer {
	return &HTTPServer{
		Port: port,
		Host: host,
		DB:   dbClient,
	}
}

func (server *HTTPServer) Run() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("HTTPServer Failed")
		}
	}()

	server.DB.Init()

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

// return all content which have in data base
func (server *HTTPServer) getContentHandler(w http.ResponseWriter, r *http.Request) {
	urls, err := server.DB.GetAllURL()

	if err != nil {
		http.Error(w, "Problems with Data Base", 500)
		return
	}

	err = json.NewEncoder(w).Encode(urls)
}
