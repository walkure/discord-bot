package twitter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func checkSyndication(idStr string) (bool, error) {
	target := fmt.Sprintf("https://cdn.syndication.twimg.com/tweet-result?id=%s&token=x", idStr)

	resp, err := http.Get(target)
	if err != nil {
		return false, fmt.Errorf("http get: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusNotFound:
		return false, errors.New("http:not found")
	case http.StatusTooManyRequests:
		return false, errors.New("http:too many requests")
	}
	defer resp.Body.Close()

	tweet := struct {
		Typename string `json:"__typename"`
	}{}

	if err := json.NewDecoder(resp.Body).Decode(&tweet); err != nil {
		return false, fmt.Errorf("json decode: %w", err)
	}

	switch tweet.Typename {
	case "TweetTombstone":
		return false, nil
	case "Tweet":
		return true, nil
	default:
		return false, fmt.Errorf("unknown type: %q", tweet.Typename)
	}
}
