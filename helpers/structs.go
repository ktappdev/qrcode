package helpers

// Define FormData struct
type FormDataStruct struct {
	OriginalLink     string `json:"originalLink"`
	Opacity          string `json:"opacity"`
	BackgroundColour string `json:"backgroundColour"`
	QRCodeColour     string `json:"qrCodeColour"`
	Name             string `json:"name"`
	UseDots          string `json:"useDots"`
	OverlayOurLogo   string `json:"overlayOurLogo"`
}
