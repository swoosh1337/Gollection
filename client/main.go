package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Item struct {
	Title string
	Body  string
}

type API int

var (
	database = make(map[string]Item)
	mutex    sync.Mutex
)

func (a *API) GetDB(empty string, reply *map[string]Item) error {
	mutex.Lock()
	defer mutex.Unlock()
	*reply = database
	return nil
}

func (a *API) GetByName(title string, reply *Item) error {
	mutex.Lock()
	defer mutex.Unlock()
	item, found := database[title]
	if !found {
		return rpc.ErrShutdown // Using ErrShutdown as an example, define appropriate errors.
	}
	*reply = item
	return nil
}

func (a *API) AddItem(item Item, reply *Item) error {
	mutex.Lock()
	defer mutex.Unlock()
	database[item.Title] = item
	*reply = item
	return nil
}

func (a *API) EditItem(item Item, reply *Item) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, found := database[item.Title]; !found {
		return rpc.ErrShutdown
	}
	database[item.Title] = item
	*reply = item
	return nil
}

func (a *API) DeleteItem(item Item, reply *Item) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, found := database[item.Title]; found {
		delete(database, item.Title)
		*reply = item
		return nil
	}
	return rpc.ErrShutdown
}

func main() {
	api := new(API)
	if err := rpc.Register(api); err != nil {
		log.Fatal("error registering API", err)
	}

	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("Listener error", err)
	}
	defer listener.Close()

	log.Printf("serving rpc on port %d", 4040)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal("error serving: ", err)
	}
}
