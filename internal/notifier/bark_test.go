package notifier

import (
	"testing"
)

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "English text shorter than max",
			input:  "Hello World",
			maxLen: 20,
			want:   "Hello World",
		},
		{
			name:   "English text longer than max",
			input:  "This is a very long sentence that needs to be truncated",
			maxLen: 10,
			want:   "This is a ...",
		},
		{
			name:   "Chinese text shorter than max",
			input:  "你好世界",
			maxLen: 10,
			want:   "你好世界",
		},
		{
			name:   "Chinese text longer than max",
			input:  "这是一个很长的中文句子需要被截断处理",
			maxLen: 10,
			want:   "这是一个很长的中文句...",
		},
		{
			name:   "Mixed English and Chinese",
			input:  "Hello 世界 this is a test 测试",
			maxLen: 15,
			want:   "Hello 世界 this i...",
		},
		{
			name:   "Emoji support",
			input:  "Hello 👋 World 🌍",
			maxLen: 10,
			want:   "Hello 👋 Wo...",
		},
		{
			name:   "Empty string",
			input:  "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "String with spaces",
			input:  "   trimmed   ",
			maxLen: 20,
			want:   "trimmed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxBytes int
		wantLen  bool // whether result should be valid UTF-8
	}{
		{
			name:     "Chinese text with byte limit",
			input:    "你好世界测试",
			maxBytes: 10,
			wantLen:  true,
		},
		{
			name:     "Mixed content with byte limit",
			input:    "Hello世界",
			maxBytes: 8,
			wantLen:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncateBytes(tt.input, tt.maxBytes)
			if len(got) > tt.maxBytes+3 { // +3 for "..."
				t.Errorf("truncateBytes() result too long: got %d bytes, max %d", len(got), tt.maxBytes)
			}
			// Check if result is valid UTF-8
			for _, r := range got {
				if r == '\uFFFD' {
					// Find if this replacement character was in the original
					foundInOriginal := false
					for _, origR := range tt.input {
						if origR == '\uFFFD' {
							foundInOriginal = true
							break
						}
					}
					if !foundInOriginal {
						t.Errorf("truncateBytes() produced invalid UTF-8")
						break
					}
				}
			}
		})
	}
}
