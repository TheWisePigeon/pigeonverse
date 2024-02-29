package helpers

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

type Frontmatter struct {
	Title    string `yaml:"title"`
	Slug     string `yaml:"slug"`
	PostedAt string `yaml:"posted_at"`
	TLDR     string `yaml:"tldr"`
}

func ExtractFrontmatter(filePath string) (*Frontmatter, string, error) {
	frontmatter := new(Frontmatter)
	content, err := os.ReadFile(filePath)
	if err != nil {
		return frontmatter, "", fmt.Errorf("Error while reading file: %w", err)
	}
	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return frontmatter, "", fmt.Errorf("Error while extracting frontmatter")
	}
	frontmatterContent := parts[1]
	err = yaml.Unmarshal([]byte(frontmatterContent), frontmatter)
	if err != nil {
		return frontmatter, "", fmt.Errorf("Error while unmarshalling frontmatter: %w", err)
	}
	return frontmatter, parts[2], nil
}
