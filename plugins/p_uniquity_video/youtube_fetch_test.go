package p_uniquity_video

import (
	"strings"
	"testing"
	"time"
)

func TestFormatUploadStatusLabel(t *testing.T) {
	if got := FormatUploadStatusLabel("processed"); got != "processed" {
		t.Fatalf("got %q", got)
	}
	if got := FormatUploadStatusLabel("upload_failed"); got != "upload failed" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatPublishedAtForTZ(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")
	s := FormatPublishedAtForTZ("2011-04-13T10:04:00Z", loc)
	if s == "" {
		t.Fatal("empty")
	}
	if !strings.Contains(s, "2011-04-13") {
		t.Fatalf("unexpected format %q", s)
	}
}
