package app

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// ConvertToMarkdown converts various file formats to Markdown
func ConvertToMarkdown(content []byte, fileType string) (string, error) {
	switch strings.ToLower(fileType) {
	case ".txt":
		return ConvertTxtToMarkdown(content)
	case ".md":
		return string(content), nil
	case ".docx":
		return ConvertDocxToMarkdown(content)
	case ".pdf":
		return ConvertPdfToMarkdown(content)
	default:
		return "", fmt.Errorf("unsupported file type: %s", fileType)
	}
}

// ConvertTxtToMarkdown converts plain text to markdown
func ConvertTxtToMarkdown(content []byte) (string, error) {
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

// pandocConvert is a helper function to call pandoc
func pandocConvert(content []byte, from string, extraArgs ...string) (string, error) {
	args := []string{"--from", from, "--to", "markdown"}
	args = append(args, extraArgs...)
	cmd := exec.Command("pandoc", args...)
	cmd.Stdin = bytes.NewReader(content)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("pandoc conversion from '%s' failed: %w. Stderr: %s", from, err, stderr.String())
	}
	return out.String(), nil
}

// ConvertDocxToMarkdown converts DOCX to markdown using pandoc
func ConvertDocxToMarkdown(content []byte) (string, error) {
	return pandocConvert(content, "docx")
}

// ConvertPdfToMarkdown converts PDF to markdown using pandoc
func ConvertPdfToMarkdown(content []byte) (string, error) {
	// Using --pdf-engine=pdflatex is often more reliable for PDF conversion
	return pandocConvert(content, "pdf", "--pdf-engine=pdflatex")
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
