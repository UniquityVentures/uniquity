package p_uniquity_video

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/lariv-in/lago/components"
	"github.com/lariv-in/lago/getters"
	. "maragu.dev/gomponents"
)

// YouTube video IDs are 11 characters from this set (YouTube's opaque id alphabet).
var ytVideoIDRe = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

// Order matters: more specific path patterns before generic v= in query.
var ytURLExtractors = []*regexp.Regexp{
	regexp.MustCompile(`youtu\.be/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`youtube\.com/embed/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`youtube\.com/shorts/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`youtube\.com/live/([a-zA-Z0-9_-]{11})`),
	regexp.MustCompile(`[?&]v=([a-zA-Z0-9_-]{11})`),
}

func isYouTubeURLish(s string) bool {
	s = strings.ToLower(s)
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") ||
		strings.Contains(s, "youtube.com") || strings.Contains(s, "youtu.be")
}

func isYouTubeHost(host string) bool {
	host = strings.ToLower(host)
	if i := strings.Index(host, ":"); i >= 0 {
		host = host[:i]
	}
	if host == "youtu.be" {
		return true
	}
	if host == "youtube.com" || strings.HasSuffix(host, ".youtube.com") {
		return true
	}
	return false
}

func parseYouTubeURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err == nil && u.Host != "" && isYouTubeHost(u.Host) {
		return u, nil
	}
	u2, err2 := url.Parse("https://" + strings.TrimPrefix(strings.TrimSpace(s), "/"))
	if err2 != nil {
		return nil, err2
	}
	if u2.Host == "" || !isYouTubeHost(u2.Host) {
		return nil, fmt.Errorf("not a YouTube URL")
	}
	return u2, nil
}

func cleanYouTubeVideoID(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return "", nil
	}
	if ytVideoIDRe.MatchString(s) {
		return s, nil
	}
	if !isYouTubeURLish(s) {
		return "", fmt.Errorf("invalid YouTube video ID: expected 11 characters [a-zA-Z0-9_-] or a YouTube URL")
	}
	u, err := parseYouTubeURL(s)
	if err != nil || u == nil {
		return "", fmt.Errorf("invalid YouTube video ID: expected 11 characters [a-zA-Z0-9_-] or a YouTube URL")
	}
	canonical := u.String()
	for _, re := range ytURLExtractors {
		if m := re.FindStringSubmatch(canonical); len(m) > 1 && ytVideoIDRe.MatchString(m[1]) {
			return m[1], nil
		}
	}
	// Paste without scheme: "youtube.com/watch?v=…" — try on original trimmed string too
	for _, re := range ytURLExtractors {
		if m := re.FindStringSubmatch(s); len(m) > 1 && ytVideoIDRe.MatchString(m[1]) {
			return m[1], nil
		}
	}
	return "", fmt.Errorf("could not parse a YouTube video id from the URL")
}

// YouTubeWatchURL returns a standard https://www.youtube.com/watch?v=… URL for a stored
// 11-character video id, or "" if the value is empty or not a valid id.
func YouTubeWatchURL(videoID string) string {
	s := strings.TrimSpace(videoID)
	if s == "" || !ytVideoIDRe.MatchString(s) {
		return ""
	}
	return "https://www.youtube.com/watch?v=" + s
}

// YouTubeStudioVideoURL returns the YouTube Studio edit URL for a video id, or "" if invalid.
// The link opens in Studio when the signed-in Google account owns or can manage the video.
func YouTubeStudioVideoURL(videoID string) string {
	s := strings.TrimSpace(videoID)
	if s == "" || !ytVideoIDRe.MatchString(s) {
		return ""
	}
	return "https://studio.youtube.com/video/" + s + "/edit"
}

// InputYouTubeVideoID is a text field like [components.InputText] that accepts either a bare
// 11-character YouTube video id or a watch / shorts / embed / youtu.be URL and stores the id.
type InputYouTubeVideoID struct {
	components.Page
	Label    string
	Name     string
	Getter   getters.Getter[string]
	Required bool
	Classes  string
}

func (e InputYouTubeVideoID) GetKey() string { return e.Key }

func (e InputYouTubeVideoID) GetRoles() []string { return e.Roles }

func (e InputYouTubeVideoID) Build(ctx context.Context) Node {
	return (components.InputText{
		Page:     e.Page,
		Label:    e.Label,
		Name:     e.Name,
		Getter:   e.Getter,
		Required: e.Required,
		Classes:  e.Classes,
	}).Build(ctx)
}

func (e InputYouTubeVideoID) Parse(v any, _ context.Context) (any, error) {
	vals, _ := v.([]string)
	raw := ""
	if len(vals) > 0 {
		raw = vals[0]
	}
	return cleanYouTubeVideoID(raw)
}

func (e InputYouTubeVideoID) GetName() string { return e.Name }
