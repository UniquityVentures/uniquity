package p_uniquity_video

import (
	"net/http"

	"github.com/UniquityVentures/lamu/views"
)

type youtubePublishedMetaLayer struct{}

func (youtubePublishedMetaLayer) Next(_ views.View, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		switch v := ctx.Value("publishedVideo").(type) {
		case PublishedVideo:
			ctx = attachYouTubeMetaToContext(ctx, v)
		case *PublishedVideo:
			if v != nil {
				ctx = attachYouTubeMetaToContext(ctx, *v)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
