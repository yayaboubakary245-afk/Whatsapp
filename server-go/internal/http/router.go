package httpapi

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/groupe12-supptic/whatsapp-campus/internal/http/handlers"
	"github.com/groupe12-supptic/whatsapp-campus/internal/security"
)

func NewRouter(h *handlers.Handler) http.Handler {
	r := mux.NewRouter()
	rl := security.NewRateLimiter(120, time.Minute)
	r.Use(CORSMiddleware, RateLimitMiddleware(rl))

	// Routes publiques
	r.HandleFunc("/api/health", h.Health).Methods(http.MethodGet)
	r.HandleFunc("/api/register", h.Register).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/api/login", h.Login).Methods(http.MethodPost, http.MethodOptions)

	// Routes protégées par JWT
	api := r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware(h.Cfg.JWTSecret))

	api.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetProfile(w, r, UserIDFromCtx(r))
			return
		}
		h.UpdateProfile(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet, http.MethodPut)

	api.HandleFunc("/users/search", func(w http.ResponseWriter, r *http.Request) {
		h.SearchUsers(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet)

	api.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListContacts(w, r, UserIDFromCtx(r))
			return
		}
		h.AddContact(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet, http.MethodPost)

	api.HandleFunc("/conversations", func(w http.ResponseWriter, r *http.Request) {
		h.ListConversations(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet)

	api.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.GetMessages(w, r, UserIDFromCtx(r))
			return
		}
		h.SendMessage(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet, http.MethodPost)

	api.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListGroups(w, r, UserIDFromCtx(r))
			return
		}
		h.CreateGroup(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet, http.MethodPost)

	api.HandleFunc("/attachments/upload", func(w http.ResponseWriter, r *http.Request) {
		h.UploadAttachment(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodPost)

	api.HandleFunc("/announcements", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ListAnnouncements(w, r)
			return
		}
		h.PublishAnnouncement(w, r, UserIDFromCtx(r), RoleFromCtx(r))
	}).Methods(http.MethodGet, http.MethodPost)

	api.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		h.ListNotifications(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodGet)

	// Signalisation WebRTC & contrôle d'appel
	api.HandleFunc("/call/start", func(w http.ResponseWriter, r *http.Request) {
		h.StartCall(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodPost)
	api.HandleFunc("/call/accept", func(w http.ResponseWriter, r *http.Request) {
		h.AcceptCall(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodPost)
	api.HandleFunc("/call/end", func(w http.ResponseWriter, r *http.Request) {
		h.EndCall(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodPost)

	// WebSocket de signalisation (gorilla/websocket)
	api.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.Signaling.HandleWebSocket(w, r, UserIDFromCtx(r))
	})

	// Espace Administration (RBAC admin)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(RequireAdmin)
	admin.HandleFunc("/users/{id}/suspend", func(w http.ResponseWriter, r *http.Request) {
		h.AdminSuspendUser(w, r, UserIDFromCtx(r))
	}).Methods(http.MethodPost)
	admin.HandleFunc("/audit-logs", h.AdminAuditLogs).Methods(http.MethodGet)

	return r
}
