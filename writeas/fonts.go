package main

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
