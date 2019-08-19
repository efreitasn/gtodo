package utils

import "net/http"

const cookieName = "flash"

// AddFlashMessage adds a flash message to the provided response.
func AddFlashMessage(w http.ResponseWriter, content string) {
	http.SetCookie(w, &http.Cookie{
		Name:  cookieName,
		Value: content,
	})
}

// ReadFlashMessage returns the flash message present in the provided request.
func ReadFlashMessage(w http.ResponseWriter, r *http.Request) string {
	flashCookie, err := r.Cookie(cookieName)

	if err == http.ErrNoCookie {
		return ""
	}

	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	})

	return flashCookie.Value
}
