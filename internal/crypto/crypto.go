package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
)

var Key = ""

func HashWithKey(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	hash := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(hash)
}

func VerifyHash(data []byte, key string, hash string) bool {
	checkHash := HashWithKey(data, key)
	return hash == checkHash
}

func HashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hash := r.Header.Get("HashSHA256")
		if Key != "" {

			body := &bytes.Buffer{}
			io.Copy(body, r.Body)

			if !VerifyHash(body.Bytes(), Key, hash) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			r.Body = io.NopCloser(body)
		}
		next.ServeHTTP(w, r)
	})
}
