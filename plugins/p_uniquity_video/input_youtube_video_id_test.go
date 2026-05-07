package p_uniquity_video

import "testing"

func TestCleanYouTubeVideoID(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"  dQw4w9WgXcQ  ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://youtu.be/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://youtu.be/dQw4w9WgXcQ?t=12", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/watch?feature=share&v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/embed/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/shorts/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"https://www.youtube.com/live/dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"youtube.com/watch?v=dQw4w9WgXcQ", "dQw4w9WgXcQ", false},
		{"not-a-link", "", true},
		{"https://example.com/watch?v=dQw4w9WgXcQ", "", true},
		{"", "", false},
		{"short", "", true},
	}
	for _, tt := range tests {
		got, err := cleanYouTubeVideoID(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("cleanYouTubeVideoID(%q) err = nil, want error", tt.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("cleanYouTubeVideoID(%q) err = %v", tt.in, err)
			continue
		}
		if got != tt.want {
			t.Errorf("cleanYouTubeVideoID(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
