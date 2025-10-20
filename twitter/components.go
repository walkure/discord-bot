package twitter

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	fetcher "github.com/walkure/slack-unfurler/loader/twitter"
)

func getTweetComponent(tweet string, shortenEntities fetcher.UrlShortenEntities) discordgo.MessageComponent {
	txt := shortenEntities.ExtractShortenURLs(tweet, fetcher.ExtractUrlMarkDown)
	return &discordgo.TextDisplay{
		Content: txt,
	}
}

func getUserComponent(user fetcher.UserEntity) discordgo.MessageComponent {
	return &discordgo.Section{
		Components: []discordgo.MessageComponent{
			discordgo.TextDisplay{
				Content: fmt.Sprintf("**%s** [@%s](https://twitter.com/i/user/%s)",
					user.Name,
					user.ScreenName, user.IDStr,
				),
			},
		},
		Accessory: discordgo.Thumbnail{
			Media: discordgo.UnfurledMediaItem{
				URL: user.ProfileImageURLHTTPS,
			},
			Description: &([]string{user.ScreenName})[0],
		},
	}
}

func getCreatedAtBlock(tweetId string, createdAt time.Time) discordgo.MessageComponent {
	return &discordgo.TextDisplay{
		Content: fmt.Sprintf("[Tw](https://twitter.com/i/status/%s) %s", tweetId, time.Time(createdAt).Local().Format(time.UnixDate)),
	}
}

func getMediaURL(media fetcher.MediaEntity) *discordgo.MediaGalleryItem {
	if media.Type == "photo" {
		return &discordgo.MediaGalleryItem{
			Media: discordgo.UnfurledMediaItem{
				URL: media.MediaURLHTTPS,
			},
			Description: &([]string{media.DisplayURL})[0],
		}
	}

	if media.Type == "video" || media.Type == "animated_gif" {
		videoURL := ""
		bitrate := 0
		for _, v := range media.VideoInfo.Variants {
			if v.ContentType != "video/mp4" {
				continue
			}

			// use best bitrate
			if bitrate <= v.Bitrate {
				bitrate = v.Bitrate
				videoURL = v.URL
				//fmt.Printf("use:%d %s\n", bitrate, videoURL)
			}
		}

		return &discordgo.MediaGalleryItem{
			Media: discordgo.UnfurledMediaItem{
				URL: videoURL,
			},
			Description: &([]string{media.DisplayURL})[0],
		}

	}

	return nil
}
