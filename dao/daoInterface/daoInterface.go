package daoInterface

type DAOInterface interface {
	Init() (DAOInterface, error)
	Close()

	SaveConten(url, content string)
}
