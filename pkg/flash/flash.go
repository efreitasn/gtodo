package flash

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const cookieName = "flash"

// Message is a message used with cookies.
// The kind field of the message follows the linux commands' exit code semantics, i.e., 0 for success and any other number for failure.
type Message struct {
	Kind    int
	Content string
}

// Add adds a flash message to the provided response.
func Add(url string, w http.ResponseWriter, msg *Message) {
	value := fmt.Sprintf(
		"%v:%v",
		strconv.Itoa(msg.Kind),
		base64.URLEncoding.EncodeToString([]byte(msg.Content)),
	)

	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: value,
	})

	w.Header().Set("Location", url)
	w.WriteHeader(303)
}

// Read returns the flash message present in the provided request.
func Read(w http.ResponseWriter, r *http.Request) *Message {
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

	return &Message{
		Kind:    kind,
		Content: string(content),
	}
}
