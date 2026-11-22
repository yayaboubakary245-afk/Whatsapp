package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/groupe12-supptic/whatsapp-campus/internal/administration"
	"github.com/groupe12-supptic/whatsapp-campus/internal/announcements"
	"github.com/groupe12-supptic/whatsapp-campus/internal/attachments"
	"github.com/groupe12-supptic/whatsapp-campus/internal/auth"
	"github.com/groupe12-supptic/whatsapp-campus/internal/call"
	"github.com/groupe12-supptic/whatsapp-campus/internal/config"
	"github.com/groupe12-supptic/whatsapp-campus/internal/contacts"
	"github.com/groupe12-supptic/whatsapp-campus/internal/conversations"
	"github.com/groupe12-supptic/whatsapp-campus/internal/groups"
	"github.com/groupe12-supptic/whatsapp-campus/internal/messaging"
	"github.com/groupe12-supptic/whatsapp-campus/internal/network/tcp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/network/udp"
	"github.com/groupe12-supptic/whatsapp-campus/internal/notifications"
	"github.com/groupe12-supptic/whatsapp-campus/internal/profiles"
	"github.com/groupe12-supptic/whatsapp-campus/internal/protocol"
	"github.com/groupe12-supptic/whatsapp-campus/internal/security"
	"github.com/groupe12-supptic/whatsapp-campus/internal/users"
	"github.com/groupe12-supptic/whatsapp-campus/internal/websocket"
	"github.com/groupe12-supptic/whatsapp-campus/pkg/utils"
)

type Handler struct {
	Cfg       *config.Config
	DB        *pgxpool.Pool
	Hub       *websocket.Hub
	TCPServer *tcp.Server
	UDPServer *udp.Server
	Signaling *call.SignalingService
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	dbOK := h.DB.Ping(r.Context()) == nil
	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"service":     "whatsapp-campus-go",
		"status":      "ok",
		"postgres":    dbOK,
		"tcp_clients": h.TCPServer.Registry.Count(),
		"ws_clients":  h.Hub.ConnectedCount(),
		"udp":         h.UDPServer.Stats(),
		"timestamp":   time.Now().UTC().Format(time.RFC3339),
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Role      string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}
	if err := auth.ValidateCampusEmail(req.Email, h.Cfg.AllowedEmailDomain); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	roleName := "etudiant"
	if req.Role == "enseignant" {
		roleName = "enseignant"
	}
	codeNum, _ := rand.Int(rand.Reader, big.NewInt(900000))
	code := fmt.Sprintf("%06d", codeNum.Int64()+100000)

	tx, err := h.DB.Begin(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur base de données")
		return
	}
	defer tx.Rollback(r.Context())

	var userID string
	err = tx.QueryRow(r.Context(), `
		INSERT INTO users (email, password_hash, role_id, verification_code)
		VALUES ($1, $2, (SELECT id FROM roles WHERE name = $3), $4)
		RETURNING id::text
	`, strings.ToLower(strings.TrimSpace(req.Email)), hash, roleName, code).Scan(&userID)
	if err != nil {
		utils.WriteError(w, http.StatusConflict, "Cet e-mail universitaire est déjà enregistré")
		return
	}

	_, err = tx.Exec(r.Context(), `
		INSERT INTO profiles (user_id, first_name, last_name)
		VALUES ($1::uuid, $2, $3)
	`, userID, req.FirstName, req.LastName)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur création profil")
		return
	}
	if err := tx.Commit(r.Context()); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur validation transaction")
		return
	}

	token, _ := auth.GenerateJWT(h.Cfg.JWTSecret, userID, req.Email, roleName)
	security.LogAudit(r.Context(), h.DB, userID, "USER_REGISTER", "user", userID, r.RemoteAddr, map[string]any{"email": req.Email})

	utils.WriteJSON(w, http.StatusCreated, map[string]any{
		"user_id":           userID,
		"email":             req.Email,
		"role":              roleName,
		"token":             token,
		"verification_code": code,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}
	var userID, hash, role string
	var suspended bool
	err := h.DB.QueryRow(r.Context(), `
		SELECT u.id::text, u.password_hash, r.name, u.is_suspended
		FROM users u JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1 AND u.is_active = TRUE
	`, strings.ToLower(strings.TrimSpace(req.Email))).Scan(&userID, &hash, &role, &suspended)
	if err != nil || !auth.VerifyPassword(hash, req.Password) {
		utils.WriteError(w, http.StatusUnauthorized, "Identifiants invalides")
		return
	}
	if suspended {
		utils.WriteError(w, http.StatusForbidden, "Compte suspendu par l'administration")
		return
	}
	_, _ = h.DB.Exec(r.Context(), `UPDATE users SET is_online = TRUE, last_seen_at = NOW() WHERE id = $1::uuid`, userID)
	token, _ := auth.GenerateJWT(h.Cfg.JWTSecret, userID, req.Email, role)
	security.LogAudit(r.Context(), h.DB, userID, "USER_LOGIN", "user", userID, r.RemoteAddr, nil)

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"user_id": userID,
		"email":   req.Email,
		"role":    role,
		"token":   token,
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request, userID string) {
	p, err := profiles.GetByUserID(r.Context(), h.DB, userID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Profil introuvable")
		return
	}
	utils.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		Bio             string `json:"bio"`
		AvatarURL       string `json:"avatar_url"`
		DepartmentID    *int   `json:"department_id"`
		ProgramID       *int   `json:"program_id"`
		AcademicLevelID *int   `json:"academic_level_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}
	if err := profiles.Update(r.Context(), h.DB, userID, req.FirstName, req.LastName, req.Bio, req.AvatarURL, req.DepartmentID, req.ProgramID, req.AcademicLevelID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur mise à jour profil")
		return
	}
	p, _ := profiles.GetByUserID(r.Context(), h.DB, userID)
	utils.WriteJSON(w, http.StatusOK, p)
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request, userID string) {
	q := r.URL.Query().Get("q")
	list, err := users.Search(r.Context(), h.DB, q, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur recherche utilisateurs")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) ListContacts(w http.ResponseWriter, r *http.Request, userID string) {
	list, err := contacts.List(r.Context(), h.DB, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture contacts")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) AddContact(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		ContactUserID string `json:"contact_user_id"`
		Nickname      string `json:"nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !utils.IsValidUUID(req.ContactUserID) {
		utils.WriteError(w, http.StatusBadRequest, "Identifiant contact invalide")
		return
	}
	if err := contacts.Add(r.Context(), h.DB, userID, req.ContactUserID, req.Nickname); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur ajout contact")
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"status": "added"})
}

func (h *Handler) ListConversations(w http.ResponseWriter, r *http.Request, userID string) {
	list, err := conversations.ListForUser(r.Context(), h.DB, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture conversations")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) GetMessages(w http.ResponseWriter, r *http.Request, userID string) {
	convID := r.URL.Query().Get("conversation_id")
	if !utils.IsValidUUID(convID) {
		utils.WriteError(w, http.StatusBadRequest, "conversation_id invalide")
		return
	}
	list, err := messaging.ListByConversation(r.Context(), h.DB, convID, userID, r.URL.Query().Get("q"))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture messages")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		ConversationID   string `json:"conversation_id"`
		RecipientID      string `json:"recipient_id"`
		ReplyToMessageID string `json:"reply_to_message_id"`
		MessageType      string `json:"message_type"`
		Content          string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		utils.WriteError(w, http.StatusBadRequest, "Contenu du message requis")
		return
	}
	convID := req.ConversationID
	if convID == "" && utils.IsValidUUID(req.RecipientID) {
		var err error
		convID, err = conversations.GetOrCreatePrivate(r.Context(), h.DB, userID, req.RecipientID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, "Erreur création conversation")
			return
		}
	}
	msg, err := messaging.Send(r.Context(), h.DB, convID, userID, req.ReplyToMessageID, req.MessageType, req.Content)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur envoi message")
		return
	}

	// Push temps réel via TCP Go + WebSocket
	rows, err := h.DB.Query(r.Context(), `SELECT user_id::text FROM conversation_members WHERE conversation_id = $1::uuid AND user_id <> $2::uuid`, convID, userID)
	if err == nil {
		defer rows.Close()
		env, _ := protocol.NewEnvelope(protocol.TypeSendMessage, userID, msg)
		for rows.Next() {
			var peerID string
			if rows.Scan(&peerID) == nil {
				_ = h.TCPServer.Registry.SendToUser(peerID, env)
				h.Hub.SendToUser(peerID, map[string]any{"type": "new_message", "message": msg})
				_ = notifications.Create(r.Context(), h.DB, peerID, "message", "Nouveau message", msg.Content, msg.ID)
			}
		}
	}
	utils.WriteJSON(w, http.StatusCreated, msg)
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request, userID string) {
	list, err := groups.ListForUser(r.Context(), h.DB, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture groupes")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) CreateGroup(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		Name              string   `json:"name"`
		Description       string   `json:"description"`
		GroupType         string   `json:"group_type"`
		OnlyAdminsCanPost bool     `json:"only_admins_can_post"`
		MemberIDs         []string `json:"member_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		utils.WriteError(w, http.StatusBadRequest, "Nom et type du groupe requis")
		return
	}
	g, err := groups.Create(r.Context(), h.DB, userID, req.Name, req.Description, req.GroupType, req.OnlyAdminsCanPost, req.MemberIDs)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, g)
}

func (h *Handler) UploadAttachment(w http.ResponseWriter, r *http.Request, userID string) {
	if err := r.ParseMultipartForm(h.Cfg.MaxUploadBytes); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Fichier trop volumineux ou invalide")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Champ 'file' manquant")
		return
	}
	defer file.Close()
	att, err := attachments.SaveFile(r.Context(), h.DB, h.Cfg.UploadDir, userID, r.FormValue("message_id"), header.Filename, header.Header.Get("Content-Type"), header.Size, file)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteJSON(w, http.StatusCreated, att)
}

func (h *Handler) ListAnnouncements(w http.ResponseWriter, r *http.Request) {
	list, err := announcements.List(r.Context(), h.DB)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture annonces")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) PublishAnnouncement(w http.ResponseWriter, r *http.Request, userID, role string) {
	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
		IsPinned bool   `json:"is_pinned"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Corps JSON invalide")
		return
	}
	ann, err := announcements.Publish(r.Context(), h.DB, userID, role, req.Title, req.Content, req.Category, req.IsPinned)
	if err != nil {
		utils.WriteError(w, http.StatusForbidden, err.Error())
		return
	}
	h.Hub.Broadcast(map[string]any{"type": "new_announcement", "announcement": ann})
	utils.WriteJSON(w, http.StatusCreated, ann)
}

func (h *Handler) ListNotifications(w http.ResponseWriter, r *http.Request, userID string) {
	list, err := notifications.ListForUser(r.Context(), h.DB, userID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur notifications")
		return
	}
	utils.WriteJSON(w, http.StatusOK, list)
}

func (h *Handler) StartCall(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		TargetUserID   string `json:"target_user_id"`
		ConversationID string `json:"conversation_id"`
		CallType       string `json:"call_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !utils.IsValidUUID(req.TargetUserID) {
		utils.WriteError(w, http.StatusBadRequest, "target_user_id requis")
		return
	}
	if req.CallType == "" {
		req.CallType = "audio"
	}
	callID, err := h.Signaling.StartCall(r.Context(), userID, req.TargetUserID, req.ConversationID, req.CallType)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Impossible de démarrer l'appel")
		return
	}
	utils.WriteJSON(w, http.StatusCreated, map[string]string{"call_id": callID})
}

func (h *Handler) AcceptCall(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		CallID string `json:"call_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !utils.IsValidUUID(req.CallID) {
		utils.WriteError(w, http.StatusBadRequest, "call_id requis")
		return
	}
	if err := h.Signaling.AcceptCall(r.Context(), req.CallID, userID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur acceptation appel")
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "active"})
}

func (h *Handler) EndCall(w http.ResponseWriter, r *http.Request, userID string) {
	var req struct {
		CallID string `json:"call_id"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !utils.IsValidUUID(req.CallID) {
		utils.WriteError(w, http.StatusBadRequest, "call_id requis")
		return
	}
	if err := h.Signaling.EndCall(r.Context(), req.CallID, userID, req.Status); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur clôture appel")
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]string{"status": "ended"})
}

func (h *Handler) AdminSuspendUser(w http.ResponseWriter, r *http.Request, adminID string) {
	targetID := mux.Vars(r)["id"]
	var req struct {
		Suspended bool `json:"suspended"`
	}
	_ = json.NewDecoder(r.Body).Decode(&req)
	if err := administration.SetUserSuspension(r.Context(), h.DB, targetID, req.Suspended); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur suspension")
		return
	}
	security.LogAudit(r.Context(), h.DB, adminID, "ADMIN_SUSPEND_USER", "user", targetID, r.RemoteAddr, map[string]any{"suspended": req.Suspended})
	utils.WriteJSON(w, http.StatusOK, map[string]any{"user_id": targetID, "suspended": req.Suspended})
}

func (h *Handler) AdminAuditLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := administration.ListAuditLogs(r.Context(), h.DB)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Erreur lecture audit_logs")
		return
	}
	utils.WriteJSON(w, http.StatusOK, logs)
}
