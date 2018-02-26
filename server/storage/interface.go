package storage

type Item struct {
	Title string
	Url   string
	Tags  []string
}

type Storage interface {
	Add(items ...Item)
	Del(items ...Item)
}

const (
	MarkdownId = iota
)

func NewStorage(id int, config map[string]string) Storage {
	var strg Storage
	switch id {
	case MarkdownId:
		strg = newMarkdown(config)
	}
	return strg
}
