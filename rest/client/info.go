package client

const LocaleEN = "en"

var DefaultLocale = LocaleEN

type Info struct {
	Subject   string
	Roles     []string
	Locale    string
	IPAddress string
	UserAgent string
	DeviceID  string
}
