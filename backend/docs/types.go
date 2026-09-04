package docs

type CategoryMeta struct {
	Category string `yaml:"label"`
	Order    int    `yaml:"order"`
}

type FrontmatterMeta struct {
	Title string `yaml:"title"`
	Order int    `yaml:"order"`
}

type DocumentNode struct {
	DisplayName string          `json:"displayName"`
	URLPath     string          `json:"urlPath"`
	Path        string          `json:"-"`
	Order       int             `json:"order"`
	IsFile      bool            `json:"isFile"`
	Children    []*DocumentNode `json:"children"`
}

type HierarchyBuilder interface {
	BuildTree(rootPath string) (*DocumentNode, error)
}

type SearchResult struct {
	URLPath     string `json:"urlPath"`
	DisplayName string `json:"displayName"`
	Snippet     string `json:"snippet"`
}

