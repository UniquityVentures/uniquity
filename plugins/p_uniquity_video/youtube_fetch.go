package p_uniquity_video

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// YouTubeSnippetMeta holds display strings from YouTube Data API v3 (videos.list).
type YouTubeSnippetMeta struct {
	Title          string
	PublishedAt    string
	UploadStatus   string
	ViewCount      string
	LikeCount      string
	CommentCount   string
}

// FetchYouTubeSnippetMeta loads snippet, status, and statistics for videoID using an API key.
// videoID must be a normalized 11-character id (see [YouTubeWatchURL] / clean path).
func FetchYouTubeSnippetMeta(ctx context.Context, apiKey, videoID string) (*YouTubeSnippetMeta, error) {
	apiKey = strings.TrimSpace(apiKey)
	videoID = strings.TrimSpace(videoID)
	if apiKey == "" {
		return nil, fmt.Errorf("youtube api key not configured")
	}
	if !ytVideoIDRe.MatchString(videoID) {
		return nil, fmt.Errorf("invalid youtube video id")
	}

	svc, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("youtube service: %w", err)
	}

	resp, err := svc.Videos.List([]string{"snippet", "status", "statistics"}).Id(videoID).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("video not found or not visible with this API key")
	}
	v := resp.Items[0]
	out := &YouTubeSnippetMeta{}
	if v.Snippet != nil {
		out.Title = strings.TrimSpace(v.Snippet.Title)
		raw := strings.TrimSpace(v.Snippet.PublishedAt)
		if raw != "" {
			if t, e := time.Parse(time.RFC3339, raw); e == nil {
				out.PublishedAt = t.Format(time.RFC3339)
			} else {
				out.PublishedAt = raw
			}
		}
	}
	if v.Status != nil {
		out.UploadStatus = strings.TrimSpace(v.Status.UploadStatus)
	}
	if v.Statistics != nil {
		st := v.Statistics
		out.ViewCount = strconv.FormatUint(st.ViewCount, 10)
		out.LikeCount = strconv.FormatUint(st.LikeCount, 10)
		out.CommentCount = strconv.FormatUint(st.CommentCount, 10)
	}
	return out, nil
}

// FormatPublishedAtForTZ parses an API RFC3339 instant and formats it for loc (falls back to UTC).
func FormatPublishedAtForTZ(publishedRFC3339 string, loc *time.Location) string {
	publishedRFC3339 = strings.TrimSpace(publishedRFC3339)
	if publishedRFC3339 == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, publishedRFC3339)
	if err != nil {
		return publishedRFC3339
	}
	if loc == nil {
		loc = time.UTC
	}
	return t.In(loc).Format("2006-01-02 15:04 MST")
}

// FormatUploadStatusLabel normalizes YouTube uploadStatus for display.
func FormatUploadStatusLabel(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	return strings.ReplaceAll(s, "_", " ")
}

// Context keys for published video detail (set by youtubePublishedMetaLayer).
// Keys must not contain '.' — [getters.Key] treats dots as path separators.
const (
	ctxYouTubeSnippetTitle       = "youtubeSnippetTitle"
	ctxYouTubePublishedAtDisplay = "youtubePublishedAtDisplay"
	ctxYouTubeUploadStatus       = "youtubeUploadStatus"
	ctxYouTubeViewCount          = "youtubeViewCount"
	ctxYouTubeLikeCount          = "youtubeLikeCount"
	ctxYouTubeCommentCount       = "youtubeCommentCount"
)

func attachYouTubeMetaToContext(ctx context.Context, pv PublishedVideo) context.Context {
	clearStats := func(c context.Context) context.Context {
		c = context.WithValue(c, ctxYouTubeSnippetTitle, "")
		c = context.WithValue(c, ctxYouTubePublishedAtDisplay, "")
		c = context.WithValue(c, ctxYouTubeUploadStatus, "")
		c = context.WithValue(c, ctxYouTubeViewCount, "")
		c = context.WithValue(c, ctxYouTubeLikeCount, "")
		c = context.WithValue(c, ctxYouTubeCommentCount, "")
		return c
	}

	requireYouTubeAPIKey(strings.TrimSpace(VideoPluginConfig.YouTubeAPIKey))
	key := strings.TrimSpace(VideoPluginConfig.YouTubeAPIKey)
	id := strings.TrimSpace(pv.YouTubeVideoID)
	if !ytVideoIDRe.MatchString(id) {
		return clearStats(ctx)
	}

	callCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	meta, err := FetchYouTubeSnippetMeta(callCtx, key, id)
	if err != nil {
		slog.Warn("p_uniquity_video: youtube metadata", "videoID", id, "error", err)
		return clearStats(ctx)
	}

	loc, _ := ctx.Value("$tz").(*time.Location)
	publishedDisplay := FormatPublishedAtForTZ(meta.PublishedAt, loc)

	ctx = context.WithValue(ctx, ctxYouTubeSnippetTitle, meta.Title)
	ctx = context.WithValue(ctx, ctxYouTubePublishedAtDisplay, publishedDisplay)
	ctx = context.WithValue(ctx, ctxYouTubeUploadStatus, FormatUploadStatusLabel(meta.UploadStatus))
	ctx = context.WithValue(ctx, ctxYouTubeViewCount, meta.ViewCount)
	ctx = context.WithValue(ctx, ctxYouTubeLikeCount, meta.LikeCount)
	ctx = context.WithValue(ctx, ctxYouTubeCommentCount, meta.CommentCount)
	return ctx
}
