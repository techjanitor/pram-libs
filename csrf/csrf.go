package csrf

import (
	e "github.com/eirka/eirka-libs/errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

const (
	// the name of CSRF header
	HeaderName = "X-XSRF-TOKEN"
	// the name of the form field
	FormFieldName = "csrf_token"
	// the name of CSRF cookie
	CookieName = "csrf_token"
	// the name of the session cookie for angularjs
	SessionName = "XSRF-TOKEN"
)

var ()

// generates two cookies: a long term csrf token for a user, and a masked session token to verify against
func Cookie() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Vary", "Cookie")

		// the token from the users cookie
		var csrfToken []byte

		// get the token from the cookie
		tokenCookie, err := c.Request.Cookie(CookieName)
		if err == nil {
			csrfToken = b64decode(tokenCookie.Value)
		}

		// if the user doesnt have a csrf token create one
		if len(csrfToken) != tokenLength {
			// creates a 32 bit token
			csrfToken = generateToken()

			// set the users csrf token cookie
			csrfCookie := &http.Cookie{
				Name:     CookieName,
				Value:    b64encode(csrfToken),
				Expires:  time.Now().Add(356 * 24 * time.Hour),
				Path:     "/",
				HttpOnly: true,
			}

			// set the csrf token cookie
			http.SetCookie(c.Writer, csrfCookie)

		}

		// set the users csrf token tookie
		sessionCookie := &http.Cookie{
			Name:  SessionName,
			Value: b64encode(maskToken(csrfToken)),
			Path:  "/",
		}

		// set the session cookie
		http.SetCookie(c.Writer, sessionCookie)

		// pass token to controllers
		c.Set("csrf_token", string(sessionCookie))

		c.Next()

	}
}

// verify the sent csrf token
func Verify() gin.HandlerFunc {
	return func(c *gin.Context) {

		// the token from the users cookie
		var csrfToken []byte

		// get the token from the cookie
		tokenCookie, err := c.Request.Cookie(CookieName)
		if err == nil {
			csrfToken = b64decode(tokenCookie.Value)
		}

		var sentToken string

		// Prefer the header over form value
		sentToken = c.Header.Get(HeaderName)

		// Then POST values
		if len(sentToken) == 0 {
			sentToken = r.PostFormValue(FormFieldName)
		}

		// error if there was no csrf token or it isnt verified
		if csrfToken == nil || !verifyToken(csrfToken, b64decode(sentToken)) {
			c.JSON(e.ErrorMessage(e.ErrUnauthorized))
			c.Error(e.ErrCsrfNotValid)
			c.Abort()
			return
		}

		c.Next()

	}
}