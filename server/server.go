package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wusuluren/mypocket/server/storage"
	"io/ioutil"
	"net/http"
)

type BookmarkNode struct {
	User   string
	Passwd string
	Title  string
	Url    string
}

var (
	pUser     = flag.String("user", "", "user name")
	pPasswd   = flag.String("passwd", "", "password")
	pFilepath = flag.String("filepath", "", "filepath")
)

var strg storage.Storage

func processBookmark(pWriter *http.ResponseWriter, pRequest **http.Request, processFunc func(node BookmarkNode)) {
	writer := *pWriter
	request := *pRequest
	defer request.Body.Close()

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
	if *pUser != node.User {
		fmt.Println("bad user")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("bad user"))
		return
	}
	if *pPasswd != node.Passwd {
		fmt.Println("bad passwd")
		writer.WriteHeader(http.StatusForbidden)
		writer.Write([]byte("bad passwd"))
		return
	}
	fmt.Println(node)
	processFunc(node)
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("ok"))
}

func addBookmark(writer http.ResponseWriter, request *http.Request) {
	processBookmark(&writer, &request, func(node BookmarkNode) {
		strg.Add(storage.Item{
			Title: node.Title,
			Url:   node.Url,
		})
	})
}

func removeBookmark(writer http.ResponseWriter, request *http.Request) {
	processBookmark(&writer, &request, func(node BookmarkNode) {
		strg.Del(storage.Item{
			Title: node.Title,
			Url:   node.Url,
		})
	})
}

func main() {
	flag.Parse()
	strg = storage.NewStorage(storage.MarkdownId, map[string]string{
		"filepath": *pFilepath,
	})
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addBookmark)
	mux.HandleFunc("/del", removeBookmark)
	fmt.Println("server is running...")
	fmt.Println(*pUser, *pPasswd, *pFilepath)
	if err := http.ListenAndServe("0.0.0.0:8000", mux); err != nil {
		panic(err)
	}
}
