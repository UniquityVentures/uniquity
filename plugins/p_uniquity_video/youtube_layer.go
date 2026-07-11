package p_uniquity_video

import (
	"net/http"

	"github.com/UniquityVentures/lamu/views"
)

type youtubePublishedMetaLayer struct{}

func (youtubePublishedMetaLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if v, err := views.GetValueFromContext[string, PublishedVideo](ctx, "publishedVideo"); err == nil {
			ctx = attachYouTubeMetaToContext(ctx, v)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
