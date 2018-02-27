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

func NewStorage(id int, config map[string]string) (Storage, error) {
	var strg Storage
	var err error
	switch id {
	case MarkdownId:
		strg, err = newMarkdown(config)
	}
	return strg, err
}
