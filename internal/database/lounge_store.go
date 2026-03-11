package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

// PostType represents the type of lounge post
type PostType string

const (
	PostTypeDailySpark PostType = "daily_spark"
	PostTypeChain      PostType = "chain"
	PostTypeHotTake    PostType = "hot_take"
	PostTypeVibe       PostType = "vibe"
	PostTypePoll       PostType = "poll"
)

// LoungePost represents a post in the writers' lounge
type LoungePost struct {
	ID              string         `json:"id"`
	Type            PostType       `json:"type"`
	AuthorCharacter string         `json:"author_character"`
	AuthorUserID    string         `json:"author_user_id"`
	AuthorUsername  string         `json:"author_username"`
	Content         string         `json:"content"`
	Timestamp       time.Time      `json:"timestamp"`
	Reactions       map[string]int `json:"reactions"`
	ReplyCount      int            `json:"reply_count"`
	IsPrompt        bool           `json:"is_prompt"`
	PromptDate      string         `json:"prompt_date,omitempty"`
}

// LoungeReply represents a reply to a post
type LoungeReply struct {
	ID              string         `json:"id"`
	PostID          string         `json:"post_id"`
	ParentReplyID   string         `json:"parent_reply_id,omitempty"`
	AuthorCharacter string         `json:"author_character"`
	AuthorUserID    string         `json:"author_user_id"`
	AuthorUsername  string         `json:"author_username"`
	Content         string         `json:"content"`
	Timestamp       time.Time      `json:"timestamp"`
	Reactions       map[string]int `json:"reactions"`
}

// LoungeReaction tracks who reacted to what
type LoungeReaction struct {
	UserID    string    `json:"user_id"`
	PostID    string    `json:"post_id,omitempty"`
	ReplyID   string    `json:"reply_id,omitempty"`
	Emoji     string    `json:"emoji"`
	Timestamp time.Time `json:"timestamp"`
}

// LoungeActivity represents recent activity for the activity feed
type LoungeActivity struct {
	Type      string    `json:"type"` // "post", "reply", "reaction"
	Username  string    `json:"username"`
	Character string    `json:"character"`
	Action    string    `json:"action"`
	TargetID  string    `json:"target_id"`
	Timestamp time.Time `json:"timestamp"`
}

// LoungeStore handles lounge data persistence in SQL
type LoungeStore struct {
	db *DB
}

// NewLoungeStore creates a new SQL-based lounge store
func NewLoungeStore(db *DB) *LoungeStore {
	return &LoungeStore{db: db}
}

// CreatePost creates a new post
func (ls *LoungeStore) CreatePost(post *LoungePost) error {
	if post.ID == "" {
		post.ID = generateID()
	}
	if post.Timestamp.IsZero() {
		post.Timestamp = time.Now()
	}
	if post.Reactions == nil {
		post.Reactions = make(map[string]int)
	}

	_, err := ls.db.Exec(
		`INSERT INTO lounge_posts (id, type, author_character, author_user_id, author_username, 
		 content, timestamp, reply_count, is_prompt, prompt_date) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.ID, string(post.Type), post.AuthorCharacter, post.AuthorUserID, post.AuthorUsername,
		post.Content, post.Timestamp, post.ReplyCount, post.IsPrompt, post.PromptDate,
	)
	return err
}

// GetPost retrieves a post by ID
func (ls *LoungeStore) GetPost(postID string) (*LoungePost, error) {
	post := &LoungePost{}
	var typeStr string
	var timestampStr string

	err := ls.db.QueryRow(
		`SELECT id, type, author_character, author_user_id, author_username, content, 
		 timestamp, reply_count, is_prompt, COALESCE(prompt_date, '') 
		 FROM lounge_posts WHERE id = ?`,
		postID,
	).Scan(
		&post.ID, &typeStr, &post.AuthorCharacter, &post.AuthorUserID, &post.AuthorUsername,
		&post.Content, &timestampStr, &post.ReplyCount, &post.IsPrompt, &post.PromptDate,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("post not found")
	}
	if err != nil {
		return nil, err
	}

	post.Timestamp = parseTimestamp(timestampStr)
	post.Type = PostType(typeStr)
	post.Reactions = ls.getPostReactions(postID)

	return post, nil
}

// GetPosts retrieves all posts, sorted by timestamp (newest first)
func (ls *LoungeStore) GetPosts() []*LoungePost {
	rows, err := ls.db.Query(
		`SELECT id, type, author_character, author_user_id, author_username, content, 
		 timestamp, reply_count, is_prompt, COALESCE(prompt_date, '') 
		 FROM lounge_posts ORDER BY timestamp DESC`,
	)
	if err != nil {
		return []*LoungePost{}
	}
	defer rows.Close()

	posts := []*LoungePost{}
	for rows.Next() {
		post := &LoungePost{}
		var typeStr string
		var timestampStr string

		err := rows.Scan(
			&post.ID, &typeStr, &post.AuthorCharacter, &post.AuthorUserID, &post.AuthorUsername,
			&post.Content, &timestampStr, &post.ReplyCount, &post.IsPrompt, &post.PromptDate,
		)
		if err != nil {
			continue
		}

		// Parse timestamp from string
		post.Timestamp = parseTimestamp(timestampStr)
		post.Type = PostType(typeStr)
		post.Reactions = ls.getPostReactions(post.ID)

		posts = append(posts, post)
	}

	return posts
}

// GetDailySparkPrompt gets today's daily spark prompt
func (ls *LoungeStore) GetDailySparkPrompt() (*LoungePost, error) {
	today := time.Now().Format("2006-01-02")

	post := &LoungePost{}
	var typeStr string
	var timestampStr string

	err := ls.db.QueryRow(
		`SELECT id, type, author_character, author_user_id, author_username, content, 
		 timestamp, reply_count, is_prompt, COALESCE(prompt_date, '') 
		 FROM lounge_posts WHERE type = ? AND is_prompt = 1 AND prompt_date = ?`,
		string(PostTypeDailySpark), today,
	).Scan(
		&post.ID, &typeStr, &post.AuthorCharacter, &post.AuthorUserID, &post.AuthorUsername,
		&post.Content, &timestampStr, &post.ReplyCount, &post.IsPrompt, &post.PromptDate,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("no daily spark prompt for today")
	}
	if err != nil {
		return nil, err
	}

	post.Timestamp = parseTimestamp(timestampStr)
	post.Type = PostType(typeStr)
	post.Reactions = ls.getPostReactions(post.ID)

	return post, nil
}

// CreateReply creates a new reply
func (ls *LoungeStore) CreateReply(reply *LoungeReply) error {
	if reply.ID == "" {
		reply.ID = generateID()
	}
	if reply.Timestamp.IsZero() {
		reply.Timestamp = time.Now()
	}
	if reply.Reactions == nil {
		reply.Reactions = make(map[string]int)
	}

	tx, err := ls.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert reply
	_, err = tx.Exec(
		`INSERT INTO lounge_replies (id, post_id, parent_reply_id, author_character, 
		 author_user_id, author_username, content, timestamp) 
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		reply.ID, reply.PostID, nullString(reply.ParentReplyID), reply.AuthorCharacter,
		reply.AuthorUserID, reply.AuthorUsername, reply.Content, reply.Timestamp,
	)
	if err != nil {
		return err
	}

	// Update post reply count
	_, err = tx.Exec(
		"UPDATE lounge_posts SET reply_count = reply_count + 1 WHERE id = ?",
		reply.PostID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetReplies retrieves all replies for a post
func (ls *LoungeStore) GetReplies(postID string) []*LoungeReply {
	rows, err := ls.db.Query(
		`SELECT id, post_id, COALESCE(parent_reply_id, ''), author_character, 
		 author_user_id, author_username, content, timestamp 
		 FROM lounge_replies WHERE post_id = ? ORDER BY timestamp ASC`,
		postID,
	)
	if err != nil {
		return []*LoungeReply{}
	}
	defer rows.Close()

	replies := []*LoungeReply{}
	for rows.Next() {
		reply := &LoungeReply{}
		var timestampStr string
		err := rows.Scan(
			&reply.ID, &reply.PostID, &reply.ParentReplyID, &reply.AuthorCharacter,
			&reply.AuthorUserID, &reply.AuthorUsername, &reply.Content, &timestampStr,
		)
		if err != nil {
			continue
		}

		reply.Timestamp = parseTimestamp(timestampStr)
		reply.Reactions = ls.getReplyReactions(reply.ID)
		replies = append(replies, reply)
	}

	return replies
}

// AddReaction adds a reaction to a post or reply
func (ls *LoungeStore) AddReaction(userID, targetID, emoji string, isPost bool) error {
	tx, err := ls.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check if user already reacted with this emoji
	var exists bool
	if isPost {
		err = tx.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM lounge_reactions WHERE user_id = ? AND post_id = ? AND emoji = ?)",
			userID, targetID, emoji,
		).Scan(&exists)
	} else {
		err = tx.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM lounge_reactions WHERE user_id = ? AND reply_id = ? AND emoji = ?)",
			userID, targetID, emoji,
		).Scan(&exists)
	}
	if err != nil {
		return err
	}
	if exists {
		return errors.New("already reacted")
	}

	// Insert reaction
	if isPost {
		_, err = tx.Exec(
			"INSERT INTO lounge_reactions (user_id, post_id, emoji, timestamp) VALUES (?, ?, ?, ?)",
			userID, targetID, emoji, time.Now(),
		)
	} else {
		_, err = tx.Exec(
			"INSERT INTO lounge_reactions (user_id, reply_id, emoji, timestamp) VALUES (?, ?, ?, ?)",
			userID, targetID, emoji, time.Now(),
		)
	}
	if err != nil {
		return err
	}

	return tx.Commit()
}

// RemoveReaction removes a reaction
func (ls *LoungeStore) RemoveReaction(userID, targetID, emoji string, isPost bool) error {
	var result sql.Result
	var err error

	if isPost {
		result, err = ls.db.Exec(
			"DELETE FROM lounge_reactions WHERE user_id = ? AND post_id = ? AND emoji = ?",
			userID, targetID, emoji,
		)
	} else {
		result, err = ls.db.Exec(
			"DELETE FROM lounge_reactions WHERE user_id = ? AND reply_id = ? AND emoji = ?",
			userID, targetID, emoji,
		)
	}

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("reaction not found")
	}

	return nil
}

// GetRecentActivity gets recent activity (last 50 actions)
func (ls *LoungeStore) GetRecentActivity() []*LoungeActivity {
	// Create activity feed from posts and replies
	activities := []*LoungeActivity{}

	// Get recent posts
	posts := ls.GetPosts()
	for i, post := range posts {
		if i >= 25 { // Limit to 25 posts
			break
		}
		if !post.IsPrompt {
			activities = append(activities, &LoungeActivity{
				Type:      "post",
				Username:  post.AuthorUsername,
				Character: post.AuthorCharacter,
				Action:    "posted",
				TargetID:  post.ID,
				Timestamp: post.Timestamp,
			})
		}
	}

	// Get recent replies across all posts
	rows, err := ls.db.Query(
		`SELECT id, post_id, author_username, author_character, timestamp 
		 FROM lounge_replies ORDER BY timestamp DESC LIMIT 25`,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var id, postID, username, character string
			var timestampStr string
			if err := rows.Scan(&id, &postID, &username, &character, &timestampStr); err == nil {
				activities = append(activities, &LoungeActivity{
					Type:      "reply",
					Username:  username,
					Character: character,
					Action:    "replied",
					TargetID:  postID,
					Timestamp: parseTimestamp(timestampStr),
				})
			}
		}
	}

	// Sort by timestamp (newest first) and limit to 50
	sortActivitiesByTimestamp(activities)
	if len(activities) > 50 {
		activities = activities[:50]
	}

	return activities
}

// GetActiveUsers returns count of users active in last 15 minutes
func (ls *LoungeStore) GetActiveUsers() int {
	cutoff := time.Now().Add(-15 * time.Minute)

	var count int
	err := ls.db.QueryRow(
		`SELECT COUNT(DISTINCT author_user_id) FROM (
			SELECT author_user_id FROM lounge_posts WHERE timestamp > ?
			UNION
			SELECT author_user_id FROM lounge_replies WHERE timestamp > ?
		)`,
		cutoff, cutoff,
	).Scan(&count)

	if err != nil {
		return 0
	}
	return count
}

// GetTodayMessageCount returns count of posts + replies today
func (ls *LoungeStore) GetTodayMessageCount() int {
	today := time.Now().Truncate(24 * time.Hour)

	var count int
	err := ls.db.QueryRow(
		`SELECT 
			(SELECT COUNT(*) FROM lounge_posts WHERE timestamp > ? AND is_prompt = 0) +
			(SELECT COUNT(*) FROM lounge_replies WHERE timestamp > ?)`,
		today, today,
	).Scan(&count)

	if err != nil {
		return 0
	}
	return count
}

// Helper functions

func (ls *LoungeStore) getPostReactions(postID string) map[string]int {
	reactions := make(map[string]int)

	rows, err := ls.db.Query(
		"SELECT emoji, COUNT(*) FROM lounge_reactions WHERE post_id = ? GROUP BY emoji",
		postID,
	)
	if err != nil {
		return reactions
	}
	defer rows.Close()

	for rows.Next() {
		var emoji string
		var count int
		if err := rows.Scan(&emoji, &count); err == nil {
			reactions[emoji] = count
		}
	}

	return reactions
}

func (ls *LoungeStore) getReplyReactions(replyID string) map[string]int {
	reactions := make(map[string]int)

	rows, err := ls.db.Query(
		"SELECT emoji, COUNT(*) FROM lounge_reactions WHERE reply_id = ? GROUP BY emoji",
		replyID,
	)
	if err != nil {
		return reactions
	}
	defer rows.Close()

	for rows.Next() {
		var emoji string
		var count int
		if err := rows.Scan(&emoji, &count); err == nil {
			reactions[emoji] = count
		}
	}

	return reactions
}

func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func sortActivitiesByTimestamp(activities []*LoungeActivity) {
	// Simple bubble sort for timestamp (newest first)
	for i := 0; i < len(activities)-1; i++ {
		for j := 0; j < len(activities)-i-1; j++ {
			if activities[j].Timestamp.Before(activities[j+1].Timestamp) {
				activities[j], activities[j+1] = activities[j+1], activities[j]
			}
		}
	}
}

// parseTimestamp parses a timestamp string from SQLite into time.Time
// Handles multiple formats that SQLite might store
func parseTimestamp(s string) time.Time {
	// Try common formats
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02 15:04:05.999999999-07:00",
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.999999999-07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t
		}
	}

	// If all parsing fails, return current time as fallback
	return time.Now()
}

// SaveData is a no-op for SQL store (data is already persisted)
func (ls *LoungeStore) SaveData() error {
	return nil
}

// LoadData is a no-op for SQL store (data is loaded on demand)
func (ls *LoungeStore) LoadData() error {
	return nil
}

// MarshalJSON custom marshaller for proper JSON output
func (p *LoungePost) MarshalJSON() ([]byte, error) {
	type Alias LoungePost
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

