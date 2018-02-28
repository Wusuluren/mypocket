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

func (md *markdown) Load() error {
	file, err := os.Open(md.filepath)
	if err != nil {
		file, err = os.Create(md.filepath)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	md.treeRoot = gquery.ParseMarkdown(string(bytes))
	return nil
}

func (md *markdown) Save() error {
	file, err := os.OpenFile(md.filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, node := range md.treeRoot.Children(gquery.MdUnorderList) {
		_, err := file.WriteString(fmt.Sprintf("- %s\n", node.Text))
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (md *markdown) Add(items ...Item) error {
	needSave := false
	for _, item := range items {
		isFind := false
		for _, node := range md.treeRoot.Children(gquery.MdUnorderList) {
			if fmt.Sprintf("[%s](%s)", item.Title, item.Url) == node.Text {
				isFind = true
				break
			}
		}
		if !isFind {
			md.treeRoot.Append(&gquery.MarkdownNode{
				Type: gquery.MdUnorderList,
				Text: fmt.Sprintf("[%s](%s)", item.Title, item.Url),
				Attribute: map[string]string{
					"tags": strings.Join(item.Tags, ";"),
				},
			})
			needSave = true
		}
	}
	if needSave {
		return md.Save()
	}
	return nil
}

func (md *markdown) Del(items ...Item) error {
	needSave := false
	for _, item := range items {
		list := md.treeRoot.Children(gquery.MdUnorderList)
		for _, node := range list {
			if fmt.Sprintf("[%s](%s)", item.Title, item.Url) == node.Text {
				node.Remove()
				needSave = true
				break
			}
		}
	}
	if needSave {
		return md.Save()
	}
	return nil
}

func newMarkdown(config map[string]string) (Storage, error) {
	strg := &markdown{}
	if filepath, ok := config["filepath"]; ok {
		strg.filepath = filepath
		if err := strg.Load(); err != nil {
			return nil, err
		}
	}
	return strg, nil
}
