package models

import "sync"

type MapURL struct{
	sync.Mutex
	URLs map[string]bool
}

func (mapURL *MapURL) AddURL(url string){
	mapURL.Lock()
	mapURL.URLs[url] = true
	mapURL.Unlock()
}

func (mapURL *MapURL) CheckURL(url string) bool{
	state, _ := mapURL.URLs[url]
	return state
}
