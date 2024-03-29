//
// types returned by giphy API as generated by `https://github.com/ChimeraCoder/gojson`
//

package GTCHA

type searchResult struct {
	Data []*Image `json:"data"`

	*Meta `json:"meta"`

	Pagination *struct {
		Count      int `json:"count"`
		Offset     int `json:"offset"`
		TotalCount int `json:"total_count"`
	} `json:"pagination"`
}

type Meta struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

// Image represents an image returned from the giphy API.
type Image struct {
	BitlyGifURL      string     `json:"bitly_gif_url"`
	BitlyURL         string     `json:"bitly_url"`
	Caption          string     `json:"caption"`
	ContentURL       string     `json:"content_url"`
	EmbedURL         string     `json:"embed_url"`
	ID               string     `json:"id"`
	Images           *imageURLs `json:"images"`
	ImportDatetime   string     `json:"import_datetime"`
	Rating           string     `json:"rating"`
	Source           string     `json:"source"`
	TrendingDatetime string     `json:"trending_datetime"`
	Type             string     `json:"type"`
	URL              string     `json:"url"`
	Username         string     `json:"username"`
}

type imageURLs struct {
	Downsized              *imageInfo `json:"downsized"`
	DownsizedLarge         *imageInfo `json:"downsized_large"`
	DownsizedStill         *imageInfo `json:"downsized_still"`
	FixedHeight            *imageInfo `json:"fixed_height"`
	FixedHeightDownsampled *imageInfo `json:"fixed_height_downsampled"`
	FixedHeightSmall       *imageInfo `json:"fixed_height_small"`
	FixedHeightSmallStill  *imageInfo `json:"fixed_height_small_still"`
	FixedHeightStill       *imageInfo `json:"fixed_height_still"`
	FixedWidth             *imageInfo `json:"fixed_width"`
	FixedWidthDownsampled  *imageInfo `json:"fixed_width_downsampled"`
	FixedWidthSmall        *imageInfo `json:"fixed_width_small"`
	FixedWidthSmallStill   *imageInfo `json:"fixed_width_small_still"`
	FixedWidthStill        *imageInfo `json:"fixed_width_still"`
	Original               *imageInfo `json:"original"`
	OriginalStill          *imageInfo `json:"original_still"`
}

type imageInfo struct {
	Height   string `json:"height"`
	Mp4      string `json:"mp4"`
	Mp4Size  string `json:"mp4_size"`
	Size     string `json:"size"`
	URL      string `json:"url"`
	Webp     string `json:"webp"`
	WebpSize string `json:"webp_size"`
	Width    string `json:"width"`
}

type tagResult struct {
	Data string `json:"data"`
	*Meta
}
