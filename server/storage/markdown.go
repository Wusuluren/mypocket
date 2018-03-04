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
	gq       *gquery.GqueryMarkdown
}

func (md *markdown) Load() error {
	file, err := os.Open(md.filepath)
	if err != nil {
		file, err = os.Create(md.filepath)
		if err != nil {
			return err
		}
		fmt.Println("create file:", md.filepath)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	md.gq = gquery.NewMarkdown(string(bytes))
	return nil
}

func (md *markdown) Save() error {
	file, err := os.OpenFile(md.filepath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, node := range md.gq.Gquery(gquery.MdUnorderList) {
		_, err := file.WriteString(fmt.Sprintf("- %s\n", node.Text()))
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
		text := fmt.Sprintf("[%s](%s)", item.Title, item.Url)
		for _, node := range md.gq.Gquery(gquery.MdUnorderList) {
			if text == node.Text() {
				isFind = true
				break
			}
		}
		if !isFind {
			conf := map[string]interface{}{
				"type": gquery.MdUnorderList,
				"text": text,
				"attr": map[string]string{
					"tags": strings.Join(item.Tags, ";"),
				},
			}
			node := gquery.NewMarkdownNode(conf)
			md.gq.TreeRoot().Append(node)
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
		text := fmt.Sprintf("[%s](%s)", item.Title, item.Url)
		list := md.gq.Gquery(gquery.MdUnorderList)
		for _, node := range list {
			if text == node.Text() {
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
