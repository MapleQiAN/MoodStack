package app

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ConvertToMarkdown converts various file formats to Markdown
func ConvertToMarkdown(content []byte, fileType string) (string, error) {
	switch strings.ToLower(fileType) {
	case ".txt":
		return convertTxtToMarkdown(content)
	case ".md":
		return string(content), nil
	case ".docx":
		return convertDocxToMarkdown(content)
	case ".pdf":
		return convertPdfToMarkdown(content)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// convertTxtToMarkdown converts plain text to markdown
func convertTxtToMarkdown(content []byte) (string, error) {
	text := string(content)

	// Simple conversion: wrap content in markdown code block if it looks like code
	// Otherwise, just return as is with some basic formatting
	lines := strings.Split(text, "\n")
	var markdownLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			markdownLines = append(markdownLines, "")
			continue
		}

		// If line looks like a title (short and no period at end)
		if len(line) < 100 && !strings.HasSuffix(line, ".") && !strings.HasSuffix(line, "!") && !strings.HasSuffix(line, "?") {
			if len(markdownLines) == 0 || markdownLines[len(markdownLines)-1] == "" {
				markdownLines = append(markdownLines, "## "+line)
				continue
			}
		}

		markdownLines = append(markdownLines, line)
	}

	return strings.Join(markdownLines, "\n"), nil
}

// convertDocxToMarkdown converts DOCX to markdown (simplified version)
func convertDocxToMarkdown(content []byte) (string, error) {
	// For MVP, we'll provide a placeholder
	// In a full implementation, you would use a library like "github.com/nguyenthenguyen/docx"
	return "# 从DOCX文件导入的内容\n\n*注意：完整的DOCX解析功能将在后续版本中实现*\n\n```\n" +
		"原始文件内容需要专门的DOCX解析库来处理\n" +
		"```", nil
}

// convertPdfToMarkdown converts PDF to markdown (simplified version)
func convertPdfToMarkdown(content []byte) (string, error) {
	// For MVP, we'll provide a placeholder
	// In a full implementation, you would use a library like "github.com/ledongthuc/pdf"
	return "# 从PDF文件导入的内容\n\n*注意：完整的PDF解析功能将在后续版本中实现*\n\n```\n" +
		"原始文件内容需要专门的PDF解析库来处理\n" +
		"```", nil
}

// ExtractTitle extracts a title from content
func ExtractTitle(content string, filename string) string {
	lines := strings.Split(content, "\n")

	// Look for the first non-empty line as title
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			// Remove markdown headers if present
			line = strings.TrimPrefix(line, "# ")
			line = strings.TrimPrefix(line, "## ")
			line = strings.TrimPrefix(line, "### ")

			// Limit title length
			if len(line) > 50 {
				line = line[:47] + "..."
			}

			return line
		}
	}

	// If no content found, use filename
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}
