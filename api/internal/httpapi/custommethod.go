package httpapi

import "strings"

// parseCustomMethod splits an AIP-136 path segment such as "hello-world:like"
// at its first colon into the slug and the custom method name.
// ok is false if the segment has no colon or either part is empty.
func parseCustomMethod(segment string) (slug, method string, ok bool) {
	slug, method, ok = strings.Cut(segment, ":")
	if !ok || slug == "" || method == "" {
		return "", "", false
	}
	return slug, method, true
}
