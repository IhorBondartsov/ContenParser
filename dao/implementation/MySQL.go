package implementation

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"github.com/IhorBondartsov/ContenParser/dao/daoInterface"
	"github.com/Sirupsen/logrus"
)

const (
	tableContent = "content"
)

type SQL struct {
	DB       *sql.DB
	User     string
	Password string
	URI      string
	TypeConn string
	NameDB   string
}

func (base *SQL) Init() (daoInterface.DAOInterface, error) {
	db, err := sql.Open("mysql", base.User+":"+base.Password+"@"+base.TypeConn+"("+base.URI+")/"+base.NameDB)
	if err != nil {
		return nil, err
	}
	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	base.DB = db
	return &SQL{DB: db}, nil
}

func (base *SQL) Close() {
	base.DB.Close()
}

func (base *SQL) SaveConten(url, content string) {

	rows, err := base.DB.Query("INSERT INTO "+tableContent+" VALUES (?,?);", url, content)
	if err != nil {
		logrus.Error("Some problems in SaveContent method ", err)
		return
	}
	defer rows.Close()

}

func (base *SQL) GetAllURL() ([]string, error) {
	var urls []string

	rows, err := base.DB.Query("SELECT url FROM " + tableContent)
	if err != nil {
		logrus.Error("Some problems in SaveContent method ", err)
		return nil, err
	}

	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			logrus.Error("Cant read row")
		}

		urls = append(urls, url)
	}
	defer rows.Close()

	return urls, nil
}
