package docs

import "sync"

// CachedBuilder は BuildTree の結果を sync.Once でキャッシュするラッパー。
// サーバー起動時に1回だけツリーを構築し、以降は同じ結果を返す。
type CachedBuilder struct {
	rootPath string
	once     sync.Once
	cached   *DocumentNode
	err      error
}

func NewCachedBuilder(rootPath string) *CachedBuilder {
	return &CachedBuilder{rootPath: rootPath}
}

// BuildTree は初回呼び出し時にツリーを構築してキャッシュし、2回目以降はキャッシュを返す。
func (c *CachedBuilder) BuildTree() (*DocumentNode, error) {
	c.once.Do(func() {
		c.cached, c.err = BuildTree(c.rootPath)
	})
	return c.cached, c.err
}
