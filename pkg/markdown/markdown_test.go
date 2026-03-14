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
		name          string
		input         string
		shouldContain string
		shouldNotContain string
	}{
		{
			name:             "inline code with script tag",
			input:            "`<script>alert(1)</script>`",
			shouldContain:    "&lt;script&gt;alert(1)&lt;/script&gt;",
			shouldNotContain: "<script>alert(1)</script>",
		},
		{
			name:             "inline code with img onerror",
			input:            "`<img src=x onerror=alert(1)>`",
			shouldContain:    "&lt;img src=x onerror=alert(1)&gt;",
			shouldNotContain: "<img src=x onerror=alert(1)>",
		},
		{
			name:             "inline code with html entities",
			input:            "`<div>&amp;&lt;&gt;</div>`",
			shouldContain:    "&lt;div&gt;&amp;amp;&amp;lt;&amp;gt;&lt;/div&gt;",
			shouldNotContain: "<div>&amp;&lt;&gt;</div>",
		},
		{
			name:          "normal inline code",
			input:         "`hello world`",
			shouldContain: "hello world",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToHtml(tt.input)
			if tt.shouldContain != "" && !strings.Contains(result, tt.shouldContain) {
				t.Errorf("expected output to contain %q, got %q", tt.shouldContain, result)
			}
			if tt.shouldNotContain != "" && strings.Contains(result, tt.shouldNotContain) {
				t.Errorf("output should not contain unescaped %q, got %q", tt.shouldNotContain, result)
			}
		})
	}
}
