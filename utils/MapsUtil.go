package utils

import (
	"fmt"
	"net/url"
)

func CreateGoogleMapsLocationUrl(location string) string {
	l := url.PathEscape(location)
	return fmt.Sprintf("https://maps.google.com/maps?q=%s&z=15&ie=UTF8", l)
}

func CreateGoogleMapsNavLink(saddr, daddr string) string {
	src := url.PathEscape(saddr)
	dest := url.PathEscape(daddr)
	return fmt.Sprintf("https://maps.google.com/maps?ie=UTF8&saddr=%s&daddr=%s&dirflg=d", src, dest)
}
