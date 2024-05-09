package pixivweb

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/KJHJason/Cultured-Downloader-Logic/api"
	"github.com/KJHJason/Cultured-Downloader-Logic/configs"
	"github.com/KJHJason/Cultured-Downloader-Logic/constants"
	cdlerrors "github.com/KJHJason/Cultured-Downloader-Logic/errors"
	"github.com/KJHJason/Cultured-Downloader-Logic/notify"
	"github.com/KJHJason/Cultured-Downloader-Logic/progress"
)

var (
	ACCEPTED_SORT_ORDER = []string{
		"date", "date_d",
		"popular", "popular_d",
		"popular_male", "popular_male_d",
		"popular_female", "popular_female_d",
	}
	ACCEPTED_SEARCH_MODE = []string{
		"s_tag",
		"s_tag_full",
		"s_tc",
	}
	ACCEPTED_RATING_MODE = []string{
		"safe",
		"r18",
		"all",
	}
	ACCEPTED_ARTWORK_TYPE = []string{
		"illust_and_ugoira",
		"manga",
		"all",
	}
)

// PixivToDl is the struct that contains the arguments of Pixiv download options.
type PixivWebDlOptions struct {
	ctx    context.Context
	cancel context.CancelFunc

	// Sort order of the results. Can be "date_desc" or "date_asc".
	SortOrder    string
	SearchMode   string
	SearchAiType int // 1: filter AI works, 0: Display AI works
	RatingMode   string
	ArtworkType  string

	Configs *configs.Config

	SessionCookies  []*http.Cookie
	SessionCookieId string

	Notifier notify.Notifier

	// Progress indicators
	MainProgBar          progress.ProgressBar
	DownloadProgressBars *[]*progress.DownloadProgressBar
}

func (p *PixivWebDlOptions) GetContext() context.Context {
	return p.ctx
}

// CancelCtx releases the resources used and cancels the context of the PixivWebDlOptions struct.
func (p *PixivWebDlOptions) CancelCtx() {
	p.cancel()
}

func (p *PixivWebDlOptions) SetContext(ctx context.Context) {
	p.ctx, p.cancel = context.WithCancel(ctx)
}

func (p *PixivWebDlOptions) CtxIsActive() bool {
	return p.ctx.Err() == nil
}

// ValidateArgs validates the arguments of the Pixiv download options.
//
// Should be called after initialising the struct.
func (p *PixivWebDlOptions) ValidateArgs(userAgent string) error {
	if p.GetContext() == nil {
		p.SetContext(context.Background())
	}

	if len(p.SessionCookies) > 0 {
		if err := api.VerifyCookies(constants.PIXIV, userAgent, p.SessionCookies); err != nil {
			return err
		}
		p.SessionCookieId = ""
	} else if p.SessionCookieId != "" {
		if cookie, err := api.VerifyAndGetCookie(constants.PIXIV, p.SessionCookieId, userAgent); err != nil {
			return err
		} else {
			p.SessionCookies = []*http.Cookie{cookie}
		}
	}

	if p.MainProgBar == nil {
		return fmt.Errorf(
			"pixiv web error %d, main progress bar is nil",
			cdlerrors.DEV_ERROR,
		)
	}

	if p.Notifier == nil {
		return fmt.Errorf(
			"pixiv web error %d: Notifier cannot be nil",
			cdlerrors.DEV_ERROR,
		)
	}

	p.SortOrder = strings.ToLower(p.SortOrder)
	_, err := api.ValidateStrArgs(
		p.SortOrder,
		ACCEPTED_SORT_ORDER,
		[]string{
			fmt.Sprintf(
				"pixiv web error %d: Sort order %s is not allowed",
				cdlerrors.INPUT_ERROR,
				p.SortOrder,
			),
		},
	)
	if err != nil {
		return err
	}

	p.SearchMode = strings.ToLower(p.SearchMode)
	_, err = api.ValidateStrArgs(
		p.SearchMode,
		ACCEPTED_SEARCH_MODE,
		[]string{
			fmt.Sprintf(
				"pixiv web error %d: Search order %s is not allowed",
				cdlerrors.INPUT_ERROR,
				p.SearchMode,
			),
		},
	)
	if err != nil {
		return err
	}

	p.RatingMode = strings.ToLower(p.RatingMode)
	_, err = api.ValidateStrArgs(
		p.RatingMode,
		ACCEPTED_RATING_MODE,
		[]string{
			fmt.Sprintf(
				"pixiv web error %d: Rating order %s is not allowed",
				cdlerrors.INPUT_ERROR,
				p.RatingMode,
			),
		},
	)
	if err != nil {
		return err
	}

	p.ArtworkType = strings.ToLower(p.ArtworkType)
	_, err = api.ValidateStrArgs(
		p.ArtworkType,
		ACCEPTED_ARTWORK_TYPE,
		[]string{
			fmt.Sprintf(
				"pixiv web error %d: Artwork type %s is not allowed",
				cdlerrors.INPUT_ERROR,
				p.ArtworkType,
			),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
