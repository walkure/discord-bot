package twitter

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var errNotSupported = errors.New("not supported site")

func HandleExpandContent(target string) (*discordgo.MessageSend, error) {
	uri, err := url.Parse(target)

	if err != nil {
		return nil, err
	}

	fmt.Printf("Expanding %q\n", uri)

	if !(uri.Hostname() == "twitter.com" || uri.Hostname() == "x.com") {
		return nil, errNotSupported
	}

	params := strings.Split(uri.Path, "/")

	if len(params) >= 4 && params[3] != "" && params[2] == "status" {
		idStr := sanitizeIdStr(params[3])
		return ExpandContent(idStr)
	}

	return nil, nil
}

func sanitizeIdStr(idStr string) string {
	// remove non-numeric characters and after
	b := make([]byte, len(idStr))
	j := 0
	for i, c := range idStr {
		if c >= '0' && c <= '9' {
			b[i] = byte(c)
			j++
		} else {
			break
		}
	}
	return string(b[:j])
}
