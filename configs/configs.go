package configs

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"
)

type Config struct {
	// DownloadPath will be used as the base path for all downloads
	DownloadPath   string

	// FfmpegPath is the path to the FFmpeg binary
	FfmpegPath     string

	// OverwriteFiles is a flag to overwrite existing files
	// If false, the download process will be skipped if the file already exists
	OverwriteFiles bool

	// Log any detected URLs of the post content that are being downloaded
	// Despite the variable name, it only logs URLs to any supported 
	// external file hosting providers such as MEGA, Google Drive, etc.
	LogUrls		   bool

	// UserAgent is the user agent to be used in the download process
	UserAgent      string
}

func ValidateFfmpegPathLogic(ctx context.Context, ffmpegPath string) error {
	_, ffmpegErr := exec.LookPath(ffmpegPath)
	if ffmpegErr != nil {
		return ffmpegErr
	}

	cmdCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// execute the ffmpeg binary to check if it's working
	cmd := exec.CommandContext(cmdCtx, ffmpegPath, "-version")
	stdout, ffmpegErr := cmd.Output()
	if ffmpegErr != nil {
		return ffmpegErr
	}

	if len(stdout) > 0 && strings.HasPrefix(string(stdout), "ffmpeg version") {
		return nil
	}
	return errors.New("unexpected output from ffmpeg binary, please ensure it is the correct ffmpeg binary")
}

func (c *Config) ValidateFfmpegPathLogic(ctx context.Context) error {
	return ValidateFfmpegPathLogic(ctx, c.FfmpegPath)
}
