package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wusuluren/mypocket/server/storage"
	"io/ioutil"
	"net/http"
	"os"
)

type BookmarkNode struct {
	User   string
	Passwd string
	Title  string
	Url    string
}

type Config struct {
	User     string
	Passwd   string
	Filepath string
	Port     int
}

var strg storage.Storage
var config Config

func processBookmark(writer http.ResponseWriter, request *http.Request, processFunc func(node BookmarkNode) error) {
	defer request.Body.Close()
	if request.Method != "POST" {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	node := BookmarkNode{}
	err = json.Unmarshal(bytes, &node)
	if err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if config.User != node.User {
		fmt.Println("bad user")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("bad user"))
		return
	}
	if config.Passwd != node.Passwd {
		fmt.Println("bad passwd")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("bad passwd"))
		return
	}
	fmt.Println(node)
	if err = processFunc(node); err != nil {
		fmt.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
}

func addBookmark(writer http.ResponseWriter, request *http.Request) {
	processBookmark(writer, request, func(node BookmarkNode) error {
		return strg.Add(storage.Item{
			Title: node.Title,
			Url:   node.Url,
		})
	})
}

func removeBookmark(writer http.ResponseWriter, request *http.Request) {
	processBookmark(writer, request, func(node BookmarkNode) error {
		return strg.Del(storage.Item{
			Title: node.Title,
			Url:   node.Url,
		})
	})
}

func usage() {
	fmt.Println(`Usage:	server -config=[config file]
config file like this:
{
  "user": "root",
  "passwd": "toor",
  "filepath": "test.md",
  "port": 8000
}`)
	os.Exit(0)
}

func main() {
	var err error
	pConfig := flag.String("config", "config.json", "config file")
	flag.Parse()
	bytes, err := ioutil.ReadFile(*pConfig)
	if err != nil {
		usage()
	}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}
	if config.Port == 0 {
		config.Port = 8000
	}
	fmt.Println(config.User, config.Passwd, config.Filepath)
	strg, err = storage.NewStorage(storage.MarkdownId, map[string]string{
		"filepath": config.Filepath,
	})
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addBookmark)
	mux.HandleFunc("/del", removeBookmark)
	fmt.Println("server is running...")
	err = http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", config.Port), mux)
	if err != nil {
		panic(err)
	}
}
