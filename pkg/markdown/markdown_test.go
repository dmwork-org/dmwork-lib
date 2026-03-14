package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func TestToHtml(t *testing.T) {

	htm := ToHtml("a\n```go\n /** test **/ func Test(v []byte) (error){ fmt.Println(\"zdsdsdsd\")}\n```\nb `测试`")

	fmt.Println("htm--->", htm)

}

func TestToHtml_XSSPrevention(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldNotContain []string
		shouldContain    []string
	}{
		{
			name:  "inline code with script tag",
			input: "test `<script>alert('xss')</script>` end",
			shouldNotContain: []string{"<script>", "</script>"},
			shouldContain:    []string{"&lt;script&gt;", "&lt;/script&gt;"},
		},
		{
			name:  "inline code with img onerror",
			input: "test `<img src=x onerror=\"alert(1)\">` end",
			shouldNotContain: []string{"<img "},
			shouldContain:    []string{"&lt;img"},
		},
		{
			name:  "inline code with angle brackets",
			input: "test `<div onclick=\"evil()\">` end",
			shouldNotContain: []string{"<div "},
			shouldContain:    []string{"&lt;div"},
		},
		{
			name:  "normal inline code unchanged",
			input: "test `fmt.Println(\"hello\")` end",
			shouldContain: []string{"fmt.Println"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToHtml(tt.input)

			for _, s := range tt.shouldNotContain {
				if strings.Contains(result, s) {
					t.Errorf("output should not contain %q, got: %s", s, result)
				}
			}

			for _, s := range tt.shouldContain {
				if !strings.Contains(result, s) {
					t.Errorf("output should contain %q, got: %s", s, result)
				}
			}
		})
	}
}
