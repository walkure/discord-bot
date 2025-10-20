package twitter

import (
	"errors"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	fetcher "github.com/walkure/slack-unfurler/loader/twitter"
)

var ErrSyndicationAvailable = errors.New("content available in syndication")

func ExpandContent(idStr string) (*discordgo.MessageSend, error) {

	ok, err := checkSyndication(idStr)
	if err != nil {
		return nil, fmt.Errorf("failure to check syndication: %w", err)
	}

	if ok {
		return nil, ErrSyndicationAvailable
	}

	result, err := fetcher.FetchTweetStatus(idStr)
	if err != nil {
		return nil, fmt.Errorf("fetch tweet by API: %w", err)
	}

	legacyTweet := result.Legacy
	noteTweet := result.NoteTweet.NoteTweetResults.Result
	userResult := result.Core.UserResults.Result

	legacyUser := userResult.Legacy
	qtResult := result.QuotedStatusResult.Result

	if result.RestID == "" {
		legacyTweet = result.Tweet.Legacy
		noteTweet = result.Tweet.NoteTweet.NoteTweetResults.Result
		userResult = result.Tweet.Core.UserResults.Result

		legacyUser = userResult.Legacy
		if result.Tweet.QuotedStatusResult.Result != nil {
			qtResult = *result.Tweet.QuotedStatusResult.Result
		}
	}

	// user backfills
	if userResult.Avatar.ImageUrl != "" {
		legacyUser.ProfileImageURLHTTPS = userResult.Avatar.ImageUrl
	}
	if userResult.Core.ScreenName != "" {
		legacyUser.ScreenName = userResult.Core.ScreenName
		legacyUser.Name = userResult.Core.Name
	}
	if legacyUser.IDStr == "" {
		legacyUser.IDStr = userResult.RestID
	}

	var tweetText string
	var entities fetcher.UrlShortenEntities

	if noteTweet.ID != "" {
		tweetText = noteTweet.Text
		// note tweet has no explicit screen_name in conversation
		entities = noteTweet.EntitySet.GetShortenURLs("", "")
	} else {
		tweetText = legacyTweet.FullText
		entities = legacyTweet.Entities.GetShortenURLs(legacyTweet.ConversationIDStr, legacyTweet.InReplyToUserIDStr)
	}

	components := []discordgo.MessageComponent{
		getUserComponent(legacyUser),
		getTweetComponent(tweetText, entities),
	}

	if len(legacyTweet.ExtendedEntities.Media) > 0 {
		mediaItem := make([]discordgo.MediaGalleryItem, 0, len(legacyTweet.ExtendedEntities.Media))
		for _, p := range legacyTweet.ExtendedEntities.Media {
			if it := getMediaURL(p); it != nil {
				mediaItem = append(mediaItem, *it)
			} else {
				fmt.Printf("unknown media")
			}
		}
		components = append(components, discordgo.MediaGallery{
			Items: mediaItem,
		})
	}

	components = append(components, getCreatedAtBlock(legacyTweet.IDStr, time.Time(legacyTweet.CreatedAt)))

	if legacyTweet.QuotedStatusIDStr != "" {
		qtLegacy := qtResult.Legacy
		qtNote := qtResult.NoteTweet.NoteTweetResults.Result
		qtUserResult := qtResult.Core.UserResults.Result
		qtLegacyUser := qtUserResult.Legacy
		if qtResult.RestID == "" {
			qtLegacy = qtResult.Tweet.Legacy
			qtUserResult = qtResult.Tweet.Core.UserResults.Result
			qtLegacyUser = qtUserResult.Legacy
			qtNote = qtResult.Tweet.NoteTweet.NoteTweetResults.Result
		}

		// user backfills
		if qtUserResult.Avatar.ImageUrl != "" {
			qtLegacyUser.ProfileImageURLHTTPS = qtUserResult.Avatar.ImageUrl
		}
		if qtUserResult.Core.ScreenName != "" {
			qtLegacyUser.ScreenName = qtUserResult.Core.ScreenName
			qtLegacyUser.Name = qtUserResult.Core.Name
		}
		if qtLegacyUser.IDStr == "" {
			qtLegacyUser.IDStr = qtUserResult.RestID
		}

		if qtLegacy.IDStr == "" && qtNote.ID == "" {

			components = append(components, &discordgo.Separator{},
				&discordgo.TextDisplay{
					Content: fmt.Sprintf("[%s](%s) (deleted)", legacyTweet.QuotedStatusPermalink.Display,
						legacyTweet.QuotedStatusPermalink.Expanded),
				})

		} else {
			if qtNote.ID != "" {
				tweetText = qtNote.Text
				entities = qtNote.EntitySet.GetShortenURLs("", "")
			} else {
				tweetText = qtLegacy.FullText
				entities = qtLegacy.Entities.GetShortenURLs(qtLegacy.ConversationIDStr, qtLegacy.InReplyToUserIDStr)
			}

			components = append(components,
				&discordgo.Separator{},
				getUserComponent(qtLegacyUser),
				getTweetComponent(tweetText, entities),
			)

			if len(qtLegacy.ExtendedEntities.Media) > 0 {
				mediaItem := make([]discordgo.MediaGalleryItem, 0, len(qtLegacy.ExtendedEntities.Media))
				for _, p := range qtLegacy.ExtendedEntities.Media {
					if it := getMediaURL(p); it != nil {
						mediaItem = append(mediaItem, *it)
					} else {
						fmt.Printf("unknown media")
					}
				}
				components = append(components, discordgo.MediaGallery{
					Items: mediaItem,
				})

			}

			components = append(components, getCreatedAtBlock(qtLegacy.IDStr, time.Time(qtLegacy.CreatedAt)))
		}
	}

	resp := &discordgo.MessageSend{
		Flags: discordgo.MessageFlagsIsComponentsV2,
		Components: []discordgo.MessageComponent{
			discordgo.Container{
				Components: components,
			},
		},
	}

	return resp, nil
}
