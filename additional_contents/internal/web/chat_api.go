package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"data_labler_ui_go/internal/database"

	"github.com/go-chi/chi/v5"
)

// RegisterChatAPI wires all chat/Discord-like endpoints
func (a *App) RegisterChatAPI(r chi.Router) {
	// Channel endpoints
	r.Get("/api/chat/channels", a.handleGetChannels)
	r.Post("/api/chat/channels", a.handleCreateChannel)
	r.Get("/api/chat/channels/{channelID}", a.handleGetChannel)
	r.Get("/api/chat/channels/{channelID}/messages", a.handleGetMessages)
	r.Post("/api/chat/channels/{channelID}/messages", a.handleSendMessage)

	// Message reactions
	r.Post("/api/chat/messages/{messageID}/reactions", a.handleAddChatReaction)
	r.Delete("/api/chat/messages/{messageID}/reactions", a.handleRemoveChatReaction)

	// Friends endpoints
	r.Get("/api/chat/friends", a.handleGetFriends)
	r.Get("/api/chat/friends/pending", a.handleGetPendingFriendRequests)
	r.Post("/api/chat/friends/request", a.handleSendFriendRequest)
	r.Post("/api/chat/friends/accept", a.handleAcceptFriendRequest)
	r.Post("/api/chat/friends/decline", a.handleDeclineFriendRequest)
	r.Delete("/api/chat/friends/{friendID}", a.handleRemoveFriend)
	r.Get("/api/chat/users/search", a.handleSearchUsers)

	// DM endpoints
	r.Get("/api/chat/dms", a.handleGetDMConversations)
	r.Get("/api/chat/dms/{conversationID}", a.handleGetDMs)
	r.Post("/api/chat/dms", a.handleStartDM)
	r.Post("/api/chat/dms/{conversationID}/messages", a.handleSendDM)
	r.Post("/api/chat/dms/{conversationID}/read", a.handleMarkDMsRead)

	// Presence/stats
	r.Post("/api/chat/presence", a.handleUpdatePresence)
	r.Get("/api/chat/stats", a.handleGetChatStats)
}

// ==================== CHANNEL HANDLERS ====================

func (a *App) handleGetChannels(w http.ResponseWriter, r *http.Request) {
	channels := a.chatStore.GetChannels()
	writeJSON(w, http.StatusOK, map[string]any{
		"channels": channels,
	})
}

func (a *App) handleCreateChannel(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	name := stringField(body, "name", "")
	description := stringField(body, "description", "")
	quote := stringField(body, "quote", "")

	if name == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要频道名称")))
		return
	}

	channel := &database.ChatChannel{
		Name:             name,
		Description:      description,
		CreatedBy:        sessionData.UserID,
		CreatedFromQuote: quote,
		IsDefault:        false,
		CreatedAt:        time.Now(),
	}

	err = a.chatStore.CreateChannel(channel)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"channel": channel,
	})
}

func (a *App) handleGetChannel(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	
	channel, err := a.chatStore.GetChannel(channelID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, channel)
}

func (a *App) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	channelID := chi.URLParam(r, "channelID")
	before := r.URL.Query().Get("before")
	
	// Parse limit from query string (default 50, max 100)
	limit := 50
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if limit > 100 {
		limit = 100
	}

	messages := a.chatStore.GetMessages(channelID, limit, before)
	writeJSON(w, http.StatusOK, map[string]any{
		"messages": messages,
	})
}

func (a *App) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	channelID := chi.URLParam(r, "channelID")

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	content := stringField(body, "content", "")
	replyToID := stringField(body, "reply_to_id", "")

	if content == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("消息不能为空")))
		return
	}

	msg := &database.ChatMessage{
		ChannelID: channelID,
		UserID:    sessionData.UserID,
		Username:  sessionData.Username,
		Content:   content,
		ReplyToID: replyToID,
		CreatedAt: time.Now(),
	}

	err = a.chatStore.CreateMessage(msg)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	// Update presence
	_ = a.chatStore.UpdatePresence(sessionData.UserID, "online")

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": msg,
	})
}

// ==================== REACTION HANDLERS ====================

func (a *App) handleAddChatReaction(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	messageID := chi.URLParam(r, "messageID")

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	emoji := stringField(body, "emoji", "")
	if emoji == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要表情")))
		return
	}

	err = a.chatStore.AddMessageReaction(messageID, sessionData.UserID, emoji)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleRemoveChatReaction(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	messageID := chi.URLParam(r, "messageID")

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	emoji := stringField(body, "emoji", "")
	if emoji == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要表情")))
		return
	}

	err = a.chatStore.RemoveMessageReaction(messageID, sessionData.UserID, emoji)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

// ==================== FRIENDS HANDLERS ====================

func (a *App) handleGetFriends(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	friends := a.chatStore.GetFriends(sessionData.UserID)
	writeJSON(w, http.StatusOK, map[string]any{
		"friends": friends,
	})
}

func (a *App) handleGetPendingFriendRequests(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	requests := a.chatStore.GetPendingFriendRequests(sessionData.UserID)
	writeJSON(w, http.StatusOK, map[string]any{
		"requests": requests,
	})
}

func (a *App) handleSendFriendRequest(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	friendID := stringField(body, "friend_id", "")
	if friendID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要用户ID")))
		return
	}

	if friendID == sessionData.UserID {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("不能添加自己为好友")))
		return
	}

	err = a.chatStore.SendFriendRequest(sessionData.UserID, friendID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleAcceptFriendRequest(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	friendID := stringField(body, "friend_id", "")
	if friendID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要用户ID")))
		return
	}

	err = a.chatStore.AcceptFriendRequest(sessionData.UserID, friendID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleDeclineFriendRequest(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	friendID := stringField(body, "friend_id", "")
	if friendID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要用户ID")))
		return
	}

	err = a.chatStore.DeclineFriendRequest(sessionData.UserID, friendID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleRemoveFriend(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	friendID := chi.URLParam(r, "friendID")
	
	err := a.chatStore.RemoveFriend(sessionData.UserID, friendID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleSearchUsers(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		writeJSON(w, http.StatusOK, map[string]any{"users": []any{}})
		return
	}

	users := a.chatStore.SearchUsers(query, sessionData.UserID)
	writeJSON(w, http.StatusOK, map[string]any{
		"users": users,
	})
}

// ==================== DM HANDLERS ====================

func (a *App) handleGetDMConversations(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	conversations := a.chatStore.GetDMConversations(sessionData.UserID)
	writeJSON(w, http.StatusOK, map[string]any{
		"conversations": conversations,
	})
}

func (a *App) handleStartDM(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	otherUserID := stringField(body, "user_id", "")
	if otherUserID == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要用户ID")))
		return
	}

	conv, err := a.chatStore.GetOrCreateDMConversation(sessionData.UserID, otherUserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"conversation": conv,
	})
}

func (a *App) handleGetDMs(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")
	before := r.URL.Query().Get("before")
	limit := numberField(map[string]any{"limit": r.URL.Query().Get("limit")}, "limit", 50)
	if limit > 100 {
		limit = 100
	}

	messages := a.chatStore.GetDMs(conversationID, limit, before)
	writeJSON(w, http.StatusOK, map[string]any{
		"messages": messages,
	})
}

func (a *App) handleSendDM(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	content := stringField(body, "content", "")
	if content == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("消息不能为空")))
		return
	}

	msg := &database.DirectMessage{
		ConversationID: conversationID,
		SenderID:       sessionData.UserID,
		SenderUsername: sessionData.Username,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	err = a.chatStore.SendDM(msg)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"message": msg,
	})
}

func (a *App) handleMarkDMsRead(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	conversationID := chi.URLParam(r, "conversationID")
	
	err := a.chatStore.MarkDMsAsRead(conversationID, sessionData.UserID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

// ==================== PRESENCE/STATS HANDLERS ====================

func (a *App) handleUpdatePresence(w http.ResponseWriter, r *http.Request) {
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	status := stringField(body, "status", "online")
	if status != "online" && status != "away" && status != "offline" {
		status = "online"
	}

	err = a.chatStore.UpdatePresence(sessionData.UserID, status)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}

func (a *App) handleGetChatStats(w http.ResponseWriter, r *http.Request) {
	onlineUsers := a.chatStore.GetOnlineUsers()
	todayMessages := a.chatStore.GetTodayMessageCount()

	writeJSON(w, http.StatusOK, map[string]any{
		"online_users":   onlineUsers,
		"today_messages": todayMessages,
	})
}

