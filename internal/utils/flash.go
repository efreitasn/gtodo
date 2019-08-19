package utils

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const cookieName = "flash"

// FlashMessage is a message used with cookies.
// The kind field of the message follows the linux commands' exit code semantics, i.e., 0 for success and any other number for failure.
type FlashMessage struct {
	Kind    int
	Content string
}

// AddFlashMessage adds a flash message to the provided response.
func AddFlashMessage(w http.ResponseWriter, fMsg *FlashMessage) {
	value := fmt.Sprintf(
		"%v:%v",
		strconv.Itoa(fMsg.Kind),
		base64.URLEncoding.EncodeToString([]byte(fMsg.Content)),
	)

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: value,
	})
}

// ReadFlashMessage returns the flash message present in the provided request.
func ReadFlashMessage(w http.ResponseWriter, r *http.Request) *FlashMessage {
	flashCookie, err := r.Cookie(cookieName)

	if err == http.ErrNoCookie {
		return nil
	}

	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	})

	parts := strings.Split(flashCookie.Value, ":")
	content, err := base64.URLEncoding.DecodeString(parts[1])

	if err != nil {
		return nil
	}

	kind, err := strconv.Atoi(parts[0])

	if err != nil {
		return nil
	}

	return &FlashMessage{
		Kind:    kind,
		Content: string(content),
	}
}
