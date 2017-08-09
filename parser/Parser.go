package parser

import (
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"bytes"
	"strings"

	log "github.com/Sirupsen/logrus"

	. "github.com/IhorBondartsov/ContenParser/models"
	"github.com/IhorBondartsov/ContenParser/dao/daoInterface"
)

type Parser struct {
	MapURL        MapURL
	RemovableTags []string
	MainUrls      []string
	DB            daoInterface.DAOInterface
}

func (parser *Parser) Init() {
	parser.MapURL.URLs = make(map[string]bool)
	parser.RemovableTags = []string{"head", "iframe", "footer", "header", "form",
									"script", "button", "style", "img", "nav", "panel"}

}

func (parser *Parser) Connect(dbClient daoInterface.DAOInterface) *Parser {
	parser.Init()
	parser.DB = dbClient
	return parser
}

func (parser *Parser) StartParse() {
	urls, err := parser.ReadConfigFile("config")
	if err != nil {
		log.Infof("StartParse. Cant read config. %v", err)
	}

	for _, val := range urls {
		if strings.Contains(val, "http") {
			parser.ParsePage(val)
		}
	}
}

func (parser *Parser) ReadConfigFile(fileName string) ([]string, error) {
	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	file := bytes.NewBuffer(b).String()
	arrayUrl := strings.Split(file, "\n")

	for key, val := range arrayUrl {
		val = strings.Replace(val, "\r", "", 1)
		arrayUrl[key]= val
	}

	log.Info("Config file was read")
	return arrayUrl, nil
}

func (parser Parser) ParsePage(url string) {

	dataBase, err := parser.DB.Init()
	if err != nil {
		log.Error(err)
	}
	defer dataBase.Close()

	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Errorf("ParsePage1. URL:%v ERROR: %v ", url, err)
	}

	urls := parser.ParseUrls(doc)
	log.Info(urls)

	content := parser.SelectCMS(doc)

	log.Info("ParsePage. Content was Parsed.  URL:", url, ".")

	parser.MapURL.AddURL(url)

	dataBase.SaveConten(url, content)

	for _, val := range urls {
		if !parser.MapURL.CheckURL(val) {
			go parser.ParsePage(val)
		}
	}
}

func (parser Parser) ParseUrls(doc *goquery.Document) []string {
	mainURl := doc.Url.String()
	var urls []string = []string{mainURl}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {

		val, _ := s.Attr("href")
		if strings.Contains(val, mainURl) && !strings.Contains(val, "#") {
			urls = append(urls, val)
		}
	})
	return urls
}

func (parser Parser) RemoveUnnecessaryTag(doc *goquery.Document) {
	for _, val := range parser.RemovableTags {
		doc.Find(val).Remove()
	}
}

func (parser Parser) ParseContentDefault(doc *goquery.Document) string {
	parser.RemoveUnnecessaryTag(doc)
	val, _ := doc.Html()
	return val
}

func (parser Parser) SelectCMS(doc *goquery.Document) string {
	// logic which is define CMS which is using for this site
	result := "default"

	switch result {
	case "default":
		content := parser.ParseContentDefault(doc)
		return content
	}
	return ""
}


