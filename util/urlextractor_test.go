package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walkure/discord-unfurler/util"
)

func TestExtractHTTPSURLs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		text                  string
		availableHostSuffixes []string
		want                  []string
	}{
		{
			name:                  "Extracts All URLs",
			text:                  "日本語のサイトはhttps://www.example.com/jp/index.htmlです。次はhttps://www.google.com/search?q=test,の検索結果。非URLです。もう一つはhttps://www-jp.example.com/。https://www.example.org/終わり。",
			availableHostSuffixes: []string{"example.com", "google.com", "example.org"},
			want:                  []string{"https://www.example.com/jp/index.html", "https://www-jp.example.com/", "https://www.google.com/search?q=test,", "https://www.example.org/"},
		},
		{
			name:                  "some non-https URL",
			text:                  "SSHはssh://www.example.com/jp/index.htmlです。次はhttps://www.google.com/search?q=test,の検索結果。非URLです。もう一つはhttps://www-jp.example.com/。https://www.example.org/終わり。",
			availableHostSuffixes: []string{"example.com", "google.com", "example.org"},
			want:                  []string{"https://www-jp.example.com/", "https://www.google.com/search?q=test,", "https://www.example.org/"},
		},
		{
			name:                  "ignore some hosts",
			text:                  "日本語のサイトはhttps://www.example.com/jp/index.htmlです。次はhttps://www.google.com/search?q=test,の検索結果。非URLです。もう一つはhttps://www-jp.example.com/。https://www.example.org/終わり。",
			availableHostSuffixes: []string{"example.com"},
			want:                  []string{"https://www.example.com/jp/index.html", "https://www-jp.example.com/"},
		},
		{
			name:                  "Contains URLs only",
			text:                  "https://www.example.com/jp/index.html\nhttps://www.google.com/search?q=test, \n https://www-jp.example.com/",
			availableHostSuffixes: []string{"example.com", "google.com"},
			want:                  []string{"https://www.example.com/jp/index.html", "https://www-jp.example.com/", "https://www.google.com/search?q=test,"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := util.ExtractHTTPSURLs(tt.text, tt.availableHostSuffixes)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
