package storage_test

import (
	"github.com/wusuluren/mypocket/server/storage"
	"os"
	"testing"
)

func TestMarkdown(t *testing.T) {
	filepath := "test.md"
	config := map[string]string{
		"filepath": filepath,
	}
	file, err := os.Create(filepath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		file.Close()
		os.Remove(filepath)
	}()
	strg := storage.NewStorage(storage.MarkdownId, config)
	testData := []storage.Item{
		{
			Title: "baidu",
			Url:   "http://www.baidu.com",
		},
		{
			Title: "google",
			Url:   "http://www.google.com",
		},
	}
	strg.Add(testData...)
	strg.Del(testData...)
}
