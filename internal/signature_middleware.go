package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

type SignatureMiddleware struct {
	secret string
	next   http.Handler
}

func NewSignatureMiddleware(secret string, next http.Handler) *SignatureMiddleware {
	return &SignatureMiddleware{secret, next}
}

func (h *SignatureMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	signature := r.URL.Query().Get(SignatureParam)
	removeParamFromQuery(r, SignatureParam) // remove signature from query so it's not forwarded (e.g. in the proxy handler)

	if h.secret == "" {
		h.next.ServeHTTP(w, r)
		return
	}

	url := currentRequestUrl(r)
	expectedSignature := computeHmac(url, h.secret)

	if !signaturesMatch(signature, expectedSignature) {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	h.next.ServeHTTP(w, r)
}

const SignatureParam = "s"

func computeHmac(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}

func signaturesMatch(submittedSignature string, expectedSignature string) bool {
	signatureBytes, _ := hex.DecodeString(submittedSignature)
	expectedBytes, _ := hex.DecodeString(expectedSignature)

	return hmac.Equal(signatureBytes, expectedBytes)
}
