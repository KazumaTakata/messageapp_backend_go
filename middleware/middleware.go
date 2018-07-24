package middleware

import (
	"context"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

func NewContextWithRequestID(ctx context.Context, req *http.Request) context.Context {

	tokenstring := req.Header["X-Access-Token"][0]
	tokenData := Token{}
	_, err := jwt.ParseWithClaims(tokenstring, &tokenData, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {

	}

	return context.WithValue(ctx, "id", tokenData.ID)
}

func RequestIDFromContext(ctx context.Context) string {
	return ctx.Value("id").(string)
}

func Middleware(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := NewContextWithRequestID(req.Context(), req)
		next(rw, req.WithContext(ctx))
	})
}

type Token struct {
	ID string `json:"id"`
	jwt.StandardClaims
}
