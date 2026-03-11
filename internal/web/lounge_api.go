package web

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// RegisterLoungeAPI wires all lounge JSON endpoints
func (a *App) RegisterLoungeAPI(r chi.Router) {
	// Public endpoints (can view)
	r.Get("/api/lounge/posts", a.handleGetLoungePosts)
	r.Get("/api/lounge/post/{postID}", a.handleGetLoungePost)
	r.Get("/api/lounge/post/{postID}/replies", a.handleGetLoungeReplies)
	r.Get("/api/lounge/daily-spark", a.handleGetDailySpark)
	r.Get("/api/lounge/activity", a.handleGetLoungeActivity)
	r.Get("/api/lounge/stats", a.handleGetLoungeStats)

	// Protected endpoints (require authentication)
	r.Post("/api/lounge/post", a.handleCreateLoungePost)
	r.Post("/api/lounge/reply", a.handleCreateLoungeReply)
	r.Post("/api/lounge/react", a.handleAddReaction)
	r.Delete("/api/lounge/react", a.handleRemoveReaction)
}

// handleGetLoungePosts retrieves all lounge posts
func (a *App) handleGetLoungePosts(w http.ResponseWriter, r *http.Request) {
	posts := a.loungeStore.GetPosts()

	// Filter by type if specified
	postType := r.URL.Query().Get("type")
	if postType != "" {
		filtered := make([]*LoungePost, 0)
		for _, post := range posts {
			if string(post.Type) == postType {
				filtered = append(filtered, post)
			}
		}
		posts = filtered
	}

	// Limit results
	limit := 50
	if len(posts) > limit {
		posts = posts[:limit]
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"posts": posts,
		"total": len(posts),
	})
}

// handleGetLoungePost retrieves a single post
func (a *App) handleGetLoungePost(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	post, err := a.loungeStore.GetPost(postID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, post)
}

// handleGetLoungeReplies retrieves replies for a post
func (a *App) handleGetLoungeReplies(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "postID")

	replies := a.loungeStore.GetReplies(postID)

	writeJSON(w, http.StatusOK, map[string]any{
		"replies": replies,
		"total":   len(replies),
	})
}

// handleGetDailySpark gets today's daily spark prompt and replies
func (a *App) handleGetDailySpark(w http.ResponseWriter, r *http.Request) {
	prompt, err := a.loungeStore.GetDailySparkPrompt()
	if err != nil {
		// If no prompt exists, create one
		prompt = a.createDailySparkPrompt()
	}

	replies := a.loungeStore.GetReplies(prompt.ID)

	writeJSON(w, http.StatusOK, map[string]any{
		"prompt":  prompt,
		"replies": replies,
	})
}

// handleGetLoungeActivity gets recent activity feed
func (a *App) handleGetLoungeActivity(w http.ResponseWriter, r *http.Request) {
	activities := a.loungeStore.GetRecentActivity()

	writeJSON(w, http.StatusOK, map[string]any{
		"activities": activities,
	})
}

// handleGetLoungeStats gets lounge statistics
func (a *App) handleGetLoungeStats(w http.ResponseWriter, r *http.Request) {
	activeUsers := a.loungeStore.GetActiveUsers()
	todayMessages := a.loungeStore.GetTodayMessageCount()

	writeJSON(w, http.StatusOK, map[string]any{
		"active_users":   activeUsers,
		"today_messages": todayMessages,
	})
}

// handleCreateLoungePost creates a new post
func (a *App) handleCreateLoungePost(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	// Only writers and editors can post
	if sessionData.Role != RoleWriter && sessionData.Role != RoleEditor {
		writeJSON(w, http.StatusForbidden, errJSON(errors.New("仅限写手")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	content := stringField(body, "content", "")
	postType := stringField(body, "type", string(PostTypeVibe))
	character := stringField(body, "character", "")

	if content == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要内容")))
		return
	}

	if character == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要选择角色")))
		return
	}

	post := &LoungePost{
		Type:            PostType(postType),
		AuthorCharacter: character,
		AuthorUserID:    sessionData.UserID,
		AuthorUsername:  sessionData.Username,
		Content:         content,
		Timestamp:       time.Now(),
		Reactions:       make(map[string]int),
		ReplyCount:      0,
	}

	err = a.loungeStore.CreatePost(post)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	// Award XP for posting (writers only — editors don't earn XP)
	if sessionData.Role != RoleEditor {
		_ = a.profileStore.AddXP(sessionData.UserID, 10)
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"post":    post,
	})
}

// handleCreateLoungeReply creates a reply to a post
func (a *App) handleCreateLoungeReply(w http.ResponseWriter, r *http.Request) {
	// Require authentication
	sessionData := a.getCurrentUser(r)
	if sessionData == nil {
		writeJSON(w, http.StatusUnauthorized, errJSON(errors.New("需要登录")))
		return
	}

	// Only writers and editors can reply
	if sessionData.Role != RoleWriter && sessionData.Role != RoleEditor {
		writeJSON(w, http.StatusForbidden, errJSON(errors.New("仅限写手")))
		return
	}

	body, err := decodeBody(r.Body)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	postID := stringField(body, "post_id", "")
	content := stringField(body, "content", "")
	character := stringField(body, "character", "")
	parentReplyID := stringField(body, "parent_reply_id", "")

	if postID == "" || content == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要帖子ID和内容")))
		return
	}

	if character == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要选择角色")))
		return
	}

	// Verify post exists
	_, err = a.loungeStore.GetPost(postID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, errJSON(errors.New("帖子不存在")))
		return
	}

	reply := &LoungeReply{
		PostID:          postID,
		ParentReplyID:   parentReplyID,
		AuthorCharacter: character,
		AuthorUserID:    sessionData.UserID,
		AuthorUsername:  sessionData.Username,
		Content:         content,
		Timestamp:       time.Now(),
		Reactions:       make(map[string]int),
	}

	err = a.loungeStore.CreateReply(reply)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errJSON(err))
		return
	}

	// Award XP for replying (writers only — editors don't earn XP)
	if sessionData.Role != RoleEditor {
		_ = a.profileStore.AddXP(sessionData.UserID, 5)
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"success": true,
		"reply":   reply,
	})
}

// handleAddReaction adds a reaction to a post or reply
func (a *App) handleAddReaction(w http.ResponseWriter, r *http.Request) {
	// Require authentication
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

	targetID := stringField(body, "target_id", "")
	emoji := stringField(body, "emoji", "")
	isPost := boolField(body, "is_post", true)

	if targetID == "" || emoji == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要目标ID和表情")))
		return
	}

	// Validate emoji (only allowed emojis)
	allowedEmojis := map[string]bool{
		"😂": true, "🔥": true, "💀": true, "😭": true,
		"👑": true, "🍕": true, "✨": true, "🤝": true,
	}
	if !allowedEmojis[emoji] {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("无效的表情")))
		return
	}

	err = a.loungeStore.AddReaction(sessionData.UserID, targetID, emoji, isPost)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// handleRemoveReaction removes a reaction
func (a *App) handleRemoveReaction(w http.ResponseWriter, r *http.Request) {
	// Require authentication
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

	targetID := stringField(body, "target_id", "")
	emoji := stringField(body, "emoji", "")
	isPost := boolField(body, "is_post", true)

	if targetID == "" || emoji == "" {
		writeJSON(w, http.StatusBadRequest, errJSON(errors.New("需要目标ID和表情")))
		return
	}

	err = a.loungeStore.RemoveReaction(sessionData.UserID, targetID, emoji, isPost)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errJSON(err))
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"success": true,
	})
}

// Daily spark prompts pool
var dailySparkPrompts = []string{
	"What's your character doing RIGHT NOW?",
	"If your character had a superpower for 5 minutes...",
	"Your character's last 3 Google searches?",
	"Rate today's writing session: 🔥 / 😅 / 💀",
	"Your character just won the lottery. First thing they do?",
	"What's your character's go-to comfort food?",
	"Your character's most embarrassing moment in 10 words",
	"If your character could time travel once, when/where?",
	"Your character's current mood in emojis",
	"What's playing in your character's headphones right now?",
}

// createDailySparkPrompt creates today's daily spark prompt
func (a *App) createDailySparkPrompt() *LoungePost {
	today := time.Now().Format("2006-01-02")

	// Use day of year to pick prompt (deterministic but changes daily)
	dayOfYear := time.Now().YearDay()
	promptIndex := dayOfYear % len(dailySparkPrompts)

	prompt := &LoungePost{
		Type:            PostTypeDailySpark,
		Content:         dailySparkPrompts[promptIndex],
		Timestamp:       time.Now(),
		Reactions:       make(map[string]int),
		ReplyCount:      0,
		IsPrompt:        true,
		PromptDate:      today,
		AuthorCharacter: "✨ Daily Spark",
		AuthorUsername:  "system",
		AuthorUserID:    "system",
	}

	err := a.loungeStore.CreatePost(prompt)
	if err != nil {
		log.Printf("Error creating daily spark prompt: %v", err)
	}

	return prompt
}
