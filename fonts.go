package writeascli

import (
	"fmt"
)

// postFont represents a valid post appearance value in the API.
type postFont string

// Valid appearance types for posts.
const (
	PostFontNormal postFont = "norm"
	PostFontSans            = "sans"
	PostFontMono            = "mono"
	PostFontWrap            = "wrap"
	PostFontCode            = "code"
)

var postFontMap = map[string]postFont{
	"norm":      PostFontNormal,
	"normal":    PostFontNormal,
	"serif":     PostFontNormal,
	"sans":      PostFontSans,
	"sansserif": PostFontSans,
	"mono":      PostFontMono,
	"monospace": PostFontMono,
	"wrap":      PostFontWrap,
	"code":      PostFontCode,
}

func getFont(code bool, font string) string {
	if code {
		if font != "" && font != DefaultFont {
			fmt.Printf("A non-default font '%s' and --code flag given. 'code' type takes precedence.\n", font)
		}
		return "code"
	}

	// Validate font value
	if f, ok := postFontMap[font]; ok {
		return string(f)
	}

	fmt.Printf("Font '%s' invalid. Using default '%s'\n", font, DefaultFont)
	return string(DefaultFont)
}
