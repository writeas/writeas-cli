package posts

import (
	"strings"
)

func ExtractTitle(content string) (title string, body string) {
	if hashIndex := strings.Index(content, "# "); hashIndex == 0 {
		eol := strings.IndexRune(content, '\n')
		// First line should start with # and end with \n
		if eol != -1 {
			body = strings.TrimLeft(content[eol:], " \t\n\r")
			title = content[len("# "):eol]
			return
		}
	}
	body = content
	return
}
