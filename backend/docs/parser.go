package docs

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseCategoryYAML(dirPath string) (CategoryMeta, error) {
	var c CategoryMeta
	f, err := os.ReadFile(dirPath + "/_category.yml")
	if err != nil {
		return c, err
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func ParseFrontmatter(filePath string) (FrontmatterMeta, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return FrontmatterMeta{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// 1行目が "---" でなければ Frontmatter なし
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return FrontmatterMeta{}, nil
	}

	// 閉じ "---" までの行を収集
	var buf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		buf.WriteString(line)
		buf.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		return FrontmatterMeta{}, err
	}

	var meta FrontmatterMeta
	if err := yaml.Unmarshal(buf.Bytes(), &meta); err != nil {
		return FrontmatterMeta{}, err
	}

	return meta, nil
}
