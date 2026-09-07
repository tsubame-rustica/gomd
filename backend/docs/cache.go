package docs

import (
	"os"
	"sync"
	"time"
)

const cacheTTL = time.Hour // コンテンツ更新頻度に合わせて1時間でキャッシュを無効化

// CachedBuilder はドキュメントツリーとファイル本文を TTL 付きでキャッシュするラッパー。
// 初回アクセス時および TTL 超過後にツリーとコンテンツを再ビルドする。
type CachedBuilder struct {
	rootPath string
	mu       sync.RWMutex
	tree     *DocumentNode
	contents map[string]string // filepath -> Markdown本文
	builtAt  time.Time
	buildErr error
}

func NewCachedBuilder(rootPath string) *CachedBuilder {
	return &CachedBuilder{rootPath: rootPath}
}

// build はキャッシュが有効であれば既存の値を返し、期限切れなら再構築する。
func (c *CachedBuilder) build() (*DocumentNode, map[string]string, error) {
	// まず読み取りロックでキャッシュ有効性を確認
	c.mu.RLock()
	if c.tree != nil && time.Since(c.builtAt) < cacheTTL {
		tree, contents := c.tree, c.contents
		c.mu.RUnlock()
		return tree, contents, c.buildErr
	}
	c.mu.RUnlock()

	// 期限切れ or 初回: 書き込みロックで再構築
	c.mu.Lock()
	defer c.mu.Unlock()
	// ダブルチェック: ロック待ち中に他ゴルーチンが更新済みの可能性
	if c.tree != nil && time.Since(c.builtAt) < cacheTTL {
		return c.tree, c.contents, c.buildErr
	}

	tree, err := BuildTree(c.rootPath)
	if err != nil {
		c.buildErr = err
		return nil, nil, err
	}

	// ツリー内の全ファイルを読み込んでコンテンツキャッシュを構築
	contents := make(map[string]string)
	collectContents(tree, contents)

	c.tree = tree
	c.contents = contents
	c.builtAt = time.Now()
	c.buildErr = nil
	return tree, contents, nil
}

// BuildTree はキャッシュされたドキュメントツリーを返す。
func (c *CachedBuilder) BuildTree() (*DocumentNode, error) {
	tree, _, err := c.build()
	return tree, err
}

// FileContents はキャッシュされたファイルパス→本文マップを返す。
// 検索時のディスクI/Oを排除するために使用する。
func (c *CachedBuilder) FileContents() (map[string]string, error) {
	_, contents, err := c.build()
	return contents, err
}

// collectContents はツリーを再帰的に走査してファイル本文をマップに収集する。
func collectContents(node *DocumentNode, m map[string]string) {
	if node == nil {
		return
	}
	if node.IsFile {
		data, err := os.ReadFile(node.Path)
		if err == nil {
			m[node.Path] = string(data)
		}
	}
	for _, child := range node.Children {
		collectContents(child, m)
	}
}
