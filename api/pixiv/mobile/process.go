package pixivmobile

import (
	"strconv"
	"path/filepath"

	"github.com/KJHJason/Cultured-Downloader-Logic/api/pixiv/ugoira"
	"github.com/KJHJason/Cultured-Downloader-Logic/constants"
	"github.com/KJHJason/Cultured-Downloader-Logic/httpfuncs"
	"github.com/KJHJason/Cultured-Downloader-Logic/iofuncs"
)

// Process the artwork JSON and returns a slice of map that contains the urls of the images and the file path
func (pixiv *PixivMobile) processArtworkJson(artworkJson *IllustJson, downloadPath string) ([]*httpfuncs.ToDownload, *ugoira.Ugoira, error) {
	if artworkJson == nil {
		return nil, nil, nil
	}

	artworkId := strconv.Itoa(artworkJson.Id)
	artworkTitle := artworkJson.Title
	artworkType := artworkJson.Type
	artistName := artworkJson.User.Name
	artworkFolderPath := iofuncs.GetPostFolder(
		filepath.Join(downloadPath, constants.PIXIV_TITLE), artistName, artworkId, artworkTitle,
	)

	if artworkType == "ugoira" {
		ugoiraInfo, err := pixiv.getUgoiraMetadata(artworkId, artworkFolderPath)
		if err != nil {
			return nil, nil, err
		}
		return nil, ugoiraInfo, nil
	}

	var artworksToDownload []*httpfuncs.ToDownload
	singlePageImageUrl := artworkJson.MetaSinglePage.OriginalImageUrl
	if singlePageImageUrl != "" {
		artworksToDownload = append(artworksToDownload, &httpfuncs.ToDownload{
			Url:      singlePageImageUrl,
			FilePath: artworkFolderPath,
		})
	} else {
		for _, image := range artworkJson.MetaPages {
			imageUrl := image.ImageUrls.Original
			artworksToDownload = append(artworksToDownload, &httpfuncs.ToDownload{
				Url:      imageUrl,
				FilePath: artworkFolderPath,
			})
		}
	}
	return artworksToDownload, nil, nil
}

// The same as the processArtworkJson function but for mutliple JSONs at once
// (Those with the "illusts" key which holds a slice of maps containing the artwork JSON)
func (pixiv *PixivMobile) processMultipleArtworkJson(resJson *ArtworksJson, downloadPath string) ([]*httpfuncs.ToDownload, []*ugoira.Ugoira, []error) {
	if resJson == nil {
		return nil, nil, nil
	}

	artworksMaps := resJson.Illusts
	if len(artworksMaps) == 0 {
		return nil, nil, nil
	}

	var errSlice []error
	var ugoiraToDl []*ugoira.Ugoira
	var artworksToDl []*httpfuncs.ToDownload
	for _, artwork := range artworksMaps {
		artworks, ugoira, err := pixiv.processArtworkJson(artwork, downloadPath)
		if err != nil {
			errSlice = append(errSlice, err)
			continue
		}
		if ugoira != nil {
			ugoiraToDl = append(ugoiraToDl, ugoira)
			continue
		}
		artworksToDl = append(artworksToDl, artworks...)
	}
	return artworksToDl, ugoiraToDl, errSlice
}
