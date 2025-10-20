package util

import (
	"net/url"
	"slices"
	"strings"
)

// ExtractHTTPSURLs extracts URLs contains availableHostSuffixes from text
func ExtractHTTPSURLs(text string, availableHostSuffixes []string) []string {
	var urls []string
	currentText := text
	const prefix = "https://"

	for {
		// 1. 開始点の探索: "https://" の最初の出現位置を探す
		startIndex := strings.Index(currentText, prefix)

		if startIndex < 0 {
			// もう "https://" が見つからなければ終了
			break
		}

		// URLの候補が始まる位置 (prefixを含む)
		potentialURLStart := startIndex

		// URLの終端を探すための文字列を抽出
		searchableText := currentText[potentialURLStart:]

		// 2. 終点の特定: URLの終端を検索する
		endIndex := len(searchableText)

		for i, r := range searchableText {
			// ASCII文字列範囲
			if r > 0x7e || r < 0x21 {
				endIndex = i
				break
			}
		}

		// URLの候補文字列を切り出す
		extractedURLStr := searchableText[:endIndex]

		// 3. URLとして確認
		// net/url.Parse で構造的に有効か確認します。
		if u, err := url.Parse(extractedURLStr); err == nil {
			// https:// から始まっていて、指定したホストsuffixがあれば追加
			if u.Scheme == "https" && slices.ContainsFunc(availableHostSuffixes, func(hostSuffix string) bool {
				return strings.HasSuffix(u.Hostname(), hostSuffix)
			}) {

				urls = append(urls, extractedURLStr)
			}
		}

		// 4. 次の検索へ: 抽出したURLの直後から検索を再開
		// 次の検索開始位置 = 現在のURL開始位置 + 抽出されたURLの長さ (endIndex)
		nextSearchStart := potentialURLStart + endIndex

		if nextSearchStart < len(currentText) {
			currentText = currentText[nextSearchStart:]
		} else {
			// 文字列の終端に達した
			break
		}
	}

	return urls
}
