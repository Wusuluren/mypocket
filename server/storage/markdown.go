package storage

import (
	"fmt"
	"github.com/wusuluren/gquery"
	"io/ioutil"
	"os"
	"strings"
)

type markdown struct {
	filepath string
	treeRoot *gquery.MarkdownNode
}

func (md *markdown) Load() {

}

func saveNodetreeToFile(file *os.File, node *gquery.MarkdownNode) {
	str := node.RawText + "\n"
	for _, v := range node.Value {
		str += (v + "\n")
	}
	_, err := file.WriteString(str)
	if err != nil {
		fmt.Println(err)
	}
	for _, child := range node.Children(gquery.MdNone) {
		saveNodetreeToFile(file, child)
	}
}

func (md *markdown) Save() {
	file, err := os.OpenFile(md.filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	saveNodetreeToFile(file, md.treeRoot)
}

func (md *markdown) Add(items ...Item) {
	for _, item := range items {
		isFind := false
		for _, node := range md.treeRoot.Children(gquery.MdUnorderList) {
			if item.Title == node.Text {
				//update
				if len(node.Value) > 0 {
					node.Value[0] = item.Url
				}
				isFind = true
				break
			}
		}
		if !isFind {
			//add
			md.treeRoot.Append(&gquery.MarkdownNode{
				Type:    gquery.MdUnorderList,
				Text:    item.Title,
				RawText: "- " + item.Title,
				Value:   []string{item.Url},
				Attribute: map[string]string{
					"tags": strings.Join(item.Tags, ";"),
				},
			})
		}
	}
	md.Save()
}

func (md *markdown) Del(items ...Item) {
	for _, item := range items {
		list := md.treeRoot.Children(gquery.MdUnorderList)
		for _, node := range list {
			if item.Title == node.Text {
				node.Remove()
				break
			}
		}
	}
	md.Save()
}

func newMarkdown(config map[string]string) Storage {
	strg := &markdown{}
	if value, ok := config["filepath"]; ok {
		filepath := value
		file, err := os.Open(filepath)
		if err != nil {
			return strg
		}
		defer file.Close()
		bytes, err := ioutil.ReadAll(file)
		if err != nil {
			return strg
		}
		strg.filepath = filepath
		strg.treeRoot = gquery.ParseMarkdown(string(bytes))
	}
	return strg
}
