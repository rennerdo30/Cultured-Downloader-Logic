package httpfuncs

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"runtime"

	"github.com/KJHJason/Cultured-Downloader-Logic/constants"
	cdlerrors "github.com/KJHJason/Cultured-Downloader-Logic/errors"
)

var DEFAULT_USER_AGENT string

func init() {
	// https://www.whatismybrowser.com/guides/the-latest-user-agent/chrome
	var userAgent = map[string]string{
		"linux":   "X11; Linux x86_64",
		"darwin":  "Macintosh; Intel Mac OS X 10_15_7",
		"windows": "Windows NT 10.0; Win64; x64",
	}
	userAgentOS, ok := userAgent[runtime.GOOS]
	if !ok {
		panic(
			fmt.Errorf(
				"error %d: Failed to get user agent OS as your OS, %q, is not supported",
				cdlerrors.OS_ERROR,
				runtime.GOOS,
			),
		)
	}
	DEFAULT_USER_AGENT = fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36", userAgentOS)
}

// Returns a boolean value indicating whether the specified site supports HTTP/3
//
// Usually, the API endpoints of a site do not support HTTP/3, so the isApi parameter must be provided.
func IsHttp3Supported(site string, isApi bool) bool {
	switch site {
	case constants.FANTIA:
		return !isApi
	case constants.PIXIV_FANBOX:
		return false
	case constants.PIXIV:
		return !isApi
	case constants.PIXIV_MOBILE:
		return true
	case constants.KEMONO:
		return false
	default:
		panic(
			fmt.Errorf(
				"error %d, invalid site, %q in IsHttp3Supported",
				cdlerrors.DEV_ERROR,
				site,
			),
		)
	}
}

// Returns the last part of the given URL string
func GetLastPartOfUrl(url string) string {
	removedParams := strings.SplitN(url, "?", 2)
	splittedUrl := strings.Split(removedParams[0], "/")
	return splittedUrl[len(splittedUrl)-1]
}

// Converts a map of string back to a string
func ParamsToString(params map[string]string) string {
	paramsStr := ""
	for key, value := range params {
		paramsStr += fmt.Sprintf("%s=%s&", key, url.QueryEscape(value))
	}
	return paramsStr[:len(paramsStr)-1] // remove the last &
}

// Reads and returns the response body in bytes and closes it
func ReadResBody(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf(
			"error %d: failed to read response body from %s due to %w",
			cdlerrors.RESPONSE_ERROR,
			res.Request.URL.String(),
			err,
		)
	}
	return body, nil
}
