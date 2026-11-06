package account

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

// loginRequest je tijelo zahtjeva za prijavu koje šalje Android aplikacija.
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Handler prima HTTP zahtjeve za korisničke račune i vraća JSON odgovore.
type Handler struct {
	auth *AuthService
}

func NewHandler(auth *AuthService) *Handler {
	return &Handler{auth: auth}
}

// RegisterRoutes registrira rute account servisa. Trenutno samo POST /api/auth/login.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", h.login)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "username and password are required")
		return
	}

	user, err := h.auth.login(r.Context(), req.Username, req.Password)
	if errors.Is(err, errInvalidCredentials) {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if err != nil {
		log.Printf("login: %v", err)
		writeError(w, http.StatusInternalServerError, "login failed")
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// writeJSON šalje podatke kao JSON odgovor sa zadanim statusom.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("encode response: %v", err)
	}
}

// writeError šalje grešku u obliku {"error": "poruka"}.
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
