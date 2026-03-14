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

func TestToHtml_XSS_InlineCode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		mustNot  []string // strings that must NOT appear in output (unescaped)
		mustHave []string // strings that must appear in output (escaped)
	}{
		{
			name:    "script tag in inline code",
			input:   "test `<script>alert('xss')</script>` end",
			mustNot: []string{"<script>", "</script>"},
			mustHave: []string{"&lt;script&gt;", "&lt;/script&gt;"},
		},
		{
			name:    "img onerror XSS",
			input:   "test `<img src=x onerror=\"alert(document.cookie)\">` end",
			mustNot: []string{"<img src=x"},
			mustHave: []string{"&lt;img"},
		},
		{
			name:    "svg onload XSS",
			input:   "test `<svg onload=\"alert(1)\">` end",
			mustNot: []string{"<svg onload"},
			mustHave: []string{"&lt;svg"},
		},
		{
			name:    "HTML entities ampersand",
			input:   "test `a < b && c > d` end",
			mustNot: []string{},
			mustHave: []string{"&lt;", "&gt;", "&amp;&amp;"},
		},
		{
			name:    "iframe injection",
			input:   "test `<iframe src=\"javascript:alert(1)\"></iframe>` end",
			mustNot: []string{"<iframe", "</iframe>"},
			mustHave: []string{"&lt;iframe"},
		},
		{
			name:    "event handler injection",
			input:   "test `<div onclick=\"alert(1)\">click</div>` end",
			mustNot: []string{"<div onclick"},
			mustHave: []string{"&lt;div"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToHtml(tt.input)

			for _, bad := range tt.mustNot {
				if strings.Contains(result, bad) {
					t.Errorf("XSS vulnerability: output contains unescaped %q\nInput: %s\nOutput: %s", bad, tt.input, result)
				}
			}

			for _, good := range tt.mustHave {
				if !strings.Contains(result, good) {
					t.Errorf("Missing escaped content %q\nInput: %s\nOutput: %s", good, tt.input, result)
				}
			}
		})
	}
}
