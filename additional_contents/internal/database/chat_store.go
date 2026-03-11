package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ChatChannel represents a chat channel/room
type ChatChannel struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description,omitempty"`
	CreatedBy       string    `json:"created_by"`
	CreatedFromQuote string   `json:"created_from_quote,omitempty"`
	IsDefault       bool      `json:"is_default"`
	CreatedAt       time.Time `json:"created_at"`
	LastMessageAt   *time.Time `json:"last_message_at,omitempty"`
	UnreadCount     int       `json:"unread_count,omitempty"`
}

// ChatMessage represents a message in a channel
type ChatMessage struct {
	ID         string         `json:"id"`
	ChannelID  string         `json:"channel_id"`
	UserID     string         `json:"user_id"`
	Username   string         `json:"username"`
	Content    string         `json:"content"`
	ReplyToID  string         `json:"reply_to_id,omitempty"`
	ReplyTo    *ChatMessage   `json:"reply_to,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	EditedAt   *time.Time     `json:"edited_at,omitempty"`
	Reactions  map[string]int `json:"reactions,omitempty"`
}

// Friend represents a friendship between users
type Friend struct {
	ID         string     `json:"id"`
	UserID     string     `json:"user_id"`
	FriendID   string     `json:"friend_id"`
	FriendName string     `json:"friend_name"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	AcceptedAt *time.Time `json:"accepted_at,omitempty"`
	IsOnline   bool       `json:"is_online"`
}

// DMConversation represents a direct message conversation
type DMConversation struct {
	ID            string     `json:"id"`
	OtherUserID   string     `json:"other_user_id"`
	OtherUsername string     `json:"other_username"`
	LastMessageAt *time.Time `json:"last_message_at,omitempty"`
	UnreadCount   int        `json:"unread_count"`
	LastMessage   string     `json:"last_message,omitempty"`
}

// DirectMessage represents a DM
type DirectMessage struct {
	ID             string     `json:"id"`
	ConversationID string     `json:"conversation_id"`
	SenderID       string     `json:"sender_id"`
	SenderUsername string     `json:"sender_username"`
	Content        string     `json:"content"`
	CreatedAt      time.Time  `json:"created_at"`
	ReadAt         *time.Time `json:"read_at,omitempty"`
}

// ChatStore handles chat data persistence
type ChatStore struct {
	db *DB
}

// NewChatStore creates a new chat store
func NewChatStore(db *DB) *ChatStore {
	return &ChatStore{db: db}
}

// ==================== CHANNEL METHODS ====================

// EnsureDefaultChannel creates the default #general channel if it doesn't exist
func (cs *ChatStore) EnsureDefaultChannel() error {
	var exists bool
	err := cs.db.QueryRow("SELECT EXISTS(SELECT 1 FROM chat_channels WHERE is_default = 1)").Scan(&exists)
	if err != nil {
		return err
	}
	
	if !exists {
		// Use NULL for created_by since "system" isn't a real user
		_, err = cs.db.Exec(`
			INSERT INTO chat_channels (id, name, description, created_by, is_default, created_at)
			VALUES (?, ?, ?, NULL, 1, ?)
		`, generateID(), "general", "通用讨论区 - General discussion for all writers", time.Now())
		return err
	}
	return nil
}

// CreateChannel creates a new channel
func (cs *ChatStore) CreateChannel(channel *ChatChannel) error {
	if channel.ID == "" {
		channel.ID = generateID()
	}
	if channel.CreatedAt.IsZero() {
		channel.CreatedAt = time.Now()
	}

	_, err := cs.db.Exec(`
		INSERT INTO chat_channels (id, name, description, created_by, created_from_quote, is_default, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, channel.ID, channel.Name, channel.Description, channel.CreatedBy, 
		nullString(channel.CreatedFromQuote), channel.IsDefault, channel.CreatedAt)
	return err
}

// GetChannels returns all channels
func (cs *ChatStore) GetChannels() []*ChatChannel {
	rows, err := cs.db.Query(`
		SELECT id, name, COALESCE(description, ''), COALESCE(created_by, ''), COALESCE(created_from_quote, ''),
		       is_default, created_at, last_message_at
		FROM chat_channels
		ORDER BY is_default DESC, last_message_at DESC NULLS LAST, created_at ASC
	`)
	if err != nil {
		return []*ChatChannel{}
	}
	defer rows.Close()

	channels := []*ChatChannel{}
	for rows.Next() {
		ch := &ChatChannel{}
		var createdAtStr string
		var lastMsgStr sql.NullString
		
		err := rows.Scan(&ch.ID, &ch.Name, &ch.Description, &ch.CreatedBy, 
			&ch.CreatedFromQuote, &ch.IsDefault, &createdAtStr, &lastMsgStr)
		if err != nil {
			continue
		}
		
		ch.CreatedAt = parseTimestamp(createdAtStr)
		if lastMsgStr.Valid {
			t := parseTimestamp(lastMsgStr.String)
			ch.LastMessageAt = &t
		}
		
		channels = append(channels, ch)
	}
	return channels
}

// GetChannel returns a channel by ID
func (cs *ChatStore) GetChannel(id string) (*ChatChannel, error) {
	ch := &ChatChannel{}
	var createdAtStr string
	var lastMsgStr sql.NullString
	
	err := cs.db.QueryRow(`
		SELECT id, name, COALESCE(description, ''), COALESCE(created_by, ''), COALESCE(created_from_quote, ''),
		       is_default, created_at, last_message_at
		FROM chat_channels WHERE id = ?
	`, id).Scan(&ch.ID, &ch.Name, &ch.Description, &ch.CreatedBy, 
		&ch.CreatedFromQuote, &ch.IsDefault, &createdAtStr, &lastMsgStr)
	
	if err == sql.ErrNoRows {
		return nil, errors.New("channel not found")
	}
	if err != nil {
		return nil, err
	}
	
	ch.CreatedAt = parseTimestamp(createdAtStr)
	if lastMsgStr.Valid {
		t := parseTimestamp(lastMsgStr.String)
		ch.LastMessageAt = &t
	}
	
	return ch, nil
}

// ==================== MESSAGE METHODS ====================

// CreateMessage creates a new message
func (cs *ChatStore) CreateMessage(msg *ChatMessage) error {
	if msg.ID == "" {
		msg.ID = generateID()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	tx, err := cs.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO chat_messages (id, channel_id, user_id, username, content, reply_to_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, msg.ID, msg.ChannelID, msg.UserID, msg.Username, msg.Content, 
		nullString(msg.ReplyToID), msg.CreatedAt)
	if err != nil {
		return err
	}

	// Update channel last message time
	_, err = tx.Exec(`
		UPDATE chat_channels SET last_message_at = ? WHERE id = ?
	`, msg.CreatedAt, msg.ChannelID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetMessages returns messages for a channel
func (cs *ChatStore) GetMessages(channelID string, limit int, before string) []*ChatMessage {
	query := `
		SELECT id, channel_id, user_id, username, content, COALESCE(reply_to_id, ''), created_at
		FROM chat_messages
		WHERE channel_id = ?
	`
	args := []interface{}{channelID}
	
	if before != "" {
		query += " AND created_at < (SELECT created_at FROM chat_messages WHERE id = ?)"
		args = append(args, before)
	}
	
	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := cs.db.Query(query, args...)
	if err != nil {
		return []*ChatMessage{}
	}
	defer rows.Close()

	messages := []*ChatMessage{}
	messageIDs := []string{}
	for rows.Next() {
		msg := &ChatMessage{}
		var createdAtStr string
		
		err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Username, 
			&msg.Content, &msg.ReplyToID, &createdAtStr)
		if err != nil {
			continue
		}
		
		msg.CreatedAt = parseTimestamp(createdAtStr)
		msg.Reactions = make(map[string]int) // Initialize empty, will be filled by batch query
		
		messages = append(messages, msg)
		messageIDs = append(messageIDs, msg.ID)
	}
	
	// Batch load all reactions for all messages in a single query (fixes N+1)
	if len(messageIDs) > 0 {
		reactionsMap := cs.batchGetMessageReactions(messageIDs)
		for _, msg := range messages {
			if reactions, ok := reactionsMap[msg.ID]; ok {
				msg.Reactions = reactions
			}
		}
	}
	
	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages
}

// batchGetMessageReactions loads reactions for multiple messages in a single query
// Returns a map of messageID -> map[emoji]count
func (cs *ChatStore) batchGetMessageReactions(messageIDs []string) map[string]map[string]int {
	result := make(map[string]map[string]int)
	if len(messageIDs) == 0 {
		return result
	}

	// Build placeholder string for IN clause
	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs))
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf(
		"SELECT message_id, emoji, COUNT(*) FROM chat_message_reactions WHERE message_id IN (%s) GROUP BY message_id, emoji",
		strings.Join(placeholders, ","),
	)

	rows, err := cs.db.Query(query, args...)
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var messageID, emoji string
		var count int
		if err := rows.Scan(&messageID, &emoji, &count); err == nil {
			if result[messageID] == nil {
				result[messageID] = make(map[string]int)
			}
			result[messageID][emoji] = count
		}
	}

	return result
}

// getMessageReactions returns reactions for a message
func (cs *ChatStore) getMessageReactions(messageID string) map[string]int {
	reactions := make(map[string]int)
	
	rows, err := cs.db.Query(`
		SELECT emoji, COUNT(*) FROM chat_message_reactions
		WHERE message_id = ? GROUP BY emoji
	`, messageID)
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

// AddMessageReaction adds a reaction to a message
func (cs *ChatStore) AddMessageReaction(messageID, userID, emoji string) error {
	_, err := cs.db.Exec(`
		INSERT OR IGNORE INTO chat_message_reactions (message_id, user_id, emoji, created_at)
		VALUES (?, ?, ?, ?)
	`, messageID, userID, emoji, time.Now())
	return err
}

// RemoveMessageReaction removes a reaction
func (cs *ChatStore) RemoveMessageReaction(messageID, userID, emoji string) error {
	_, err := cs.db.Exec(`
		DELETE FROM chat_message_reactions WHERE message_id = ? AND user_id = ? AND emoji = ?
	`, messageID, userID, emoji)
	return err
}

// ==================== FRIEND METHODS ====================

// SendFriendRequest sends a friend request
func (cs *ChatStore) SendFriendRequest(userID, friendID string) error {
	// Check if already friends or request pending
	var exists bool
	err := cs.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM friends 
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?))
	`, userID, friendID, friendID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("friend request already exists")
	}

	_, err = cs.db.Exec(`
		INSERT INTO friends (id, user_id, friend_id, status, created_at)
		VALUES (?, ?, ?, 'pending', ?)
	`, generateID(), userID, friendID, time.Now())
	return err
}

// AcceptFriendRequest accepts a friend request
func (cs *ChatStore) AcceptFriendRequest(userID, friendID string) error {
	result, err := cs.db.Exec(`
		UPDATE friends SET status = 'accepted', accepted_at = ?
		WHERE user_id = ? AND friend_id = ? AND status = 'pending'
	`, time.Now(), friendID, userID)
	if err != nil {
		return err
	}
	
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("friend request not found")
	}
	return nil
}

// DeclineFriendRequest declines/removes a friend request
func (cs *ChatStore) DeclineFriendRequest(userID, friendID string) error {
	_, err := cs.db.Exec(`
		DELETE FROM friends WHERE user_id = ? AND friend_id = ? AND status = 'pending'
	`, friendID, userID)
	return err
}

// RemoveFriend removes a friendship
func (cs *ChatStore) RemoveFriend(userID, friendID string) error {
	_, err := cs.db.Exec(`
		DELETE FROM friends 
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)
	`, userID, friendID, friendID, userID)
	return err
}

// GetFriends returns all friends for a user
func (cs *ChatStore) GetFriends(userID string) []*Friend {
	rows, err := cs.db.Query(`
		SELECT f.id, f.user_id, f.friend_id, u.username, f.status, f.created_at, f.accepted_at,
		       COALESCE(p.status, 'offline') as online_status
		FROM friends f
		JOIN users u ON (CASE WHEN f.user_id = ? THEN f.friend_id ELSE f.user_id END) = u.id
		LEFT JOIN user_presence p ON u.id = p.user_id
		WHERE (f.user_id = ? OR f.friend_id = ?) AND f.status = 'accepted'
		ORDER BY u.username
	`, userID, userID, userID)
	if err != nil {
		return []*Friend{}
	}
	defer rows.Close()

	friends := []*Friend{}
	for rows.Next() {
		f := &Friend{}
		var createdAtStr string
		var acceptedAtStr sql.NullString
		var onlineStatus string
		
		err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.FriendName, &f.Status, 
			&createdAtStr, &acceptedAtStr, &onlineStatus)
		if err != nil {
			continue
		}
		
		f.CreatedAt = parseTimestamp(createdAtStr)
		if acceptedAtStr.Valid {
			t := parseTimestamp(acceptedAtStr.String)
			f.AcceptedAt = &t
		}
		f.IsOnline = onlineStatus == "online"
		
		// Normalize friend ID
		if f.UserID == userID {
			f.FriendID = f.FriendID
		} else {
			f.FriendID = f.UserID
		}
		
		friends = append(friends, f)
	}
	return friends
}

// GetPendingFriendRequests returns pending friend requests for a user
func (cs *ChatStore) GetPendingFriendRequests(userID string) []*Friend {
	rows, err := cs.db.Query(`
		SELECT f.id, f.user_id, f.friend_id, u.username, f.status, f.created_at
		FROM friends f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`, userID)
	if err != nil {
		return []*Friend{}
	}
	defer rows.Close()

	requests := []*Friend{}
	for rows.Next() {
		f := &Friend{}
		var createdAtStr string
		
		err := rows.Scan(&f.ID, &f.UserID, &f.FriendID, &f.FriendName, &f.Status, &createdAtStr)
		if err != nil {
			continue
		}
		f.CreatedAt = parseTimestamp(createdAtStr)
		requests = append(requests, f)
	}
	return requests
}

// SearchUsers searches for users to add as friends
func (cs *ChatStore) SearchUsers(query string, excludeUserID string) []map[string]string {
	rows, err := cs.db.Query(`
		SELECT id, username FROM users
		WHERE id != ? AND username LIKE ?
		LIMIT 20
	`, excludeUserID, "%"+query+"%")
	if err != nil {
		return []map[string]string{}
	}
	defer rows.Close()

	users := []map[string]string{}
	for rows.Next() {
		var id, username string
		if err := rows.Scan(&id, &username); err == nil {
			users = append(users, map[string]string{"id": id, "username": username})
		}
	}
	return users
}

// ==================== DM METHODS ====================

// GetOrCreateDMConversation gets or creates a DM conversation
func (cs *ChatStore) GetOrCreateDMConversation(user1ID, user2ID string) (*DMConversation, error) {
	// Normalize order
	if user1ID > user2ID {
		user1ID, user2ID = user2ID, user1ID
	}

	conv := &DMConversation{}
	var lastMsgStr sql.NullString
	
	err := cs.db.QueryRow(`
		SELECT c.id, c.last_message_at
		FROM dm_conversations c
		WHERE c.user1_id = ? AND c.user2_id = ?
	`, user1ID, user2ID).Scan(&conv.ID, &lastMsgStr)
	
	if err == sql.ErrNoRows {
		// Create new conversation
		conv.ID = generateID()
		_, err = cs.db.Exec(`
			INSERT INTO dm_conversations (id, user1_id, user2_id, created_at)
			VALUES (?, ?, ?, ?)
		`, conv.ID, user1ID, user2ID, time.Now())
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	
	if lastMsgStr.Valid {
		t := parseTimestamp(lastMsgStr.String)
		conv.LastMessageAt = &t
	}
	
	return conv, nil
}

// GetDMConversations returns all DM conversations for a user
func (cs *ChatStore) GetDMConversations(userID string) []*DMConversation {
	rows, err := cs.db.Query(`
		SELECT c.id, 
		       CASE WHEN c.user1_id = ? THEN c.user2_id ELSE c.user1_id END as other_id,
		       u.username as other_username,
		       c.last_message_at,
		       (SELECT COUNT(*) FROM direct_messages dm 
		        WHERE dm.conversation_id = c.id AND dm.sender_id != ? AND dm.read_at IS NULL) as unread,
		       (SELECT content FROM direct_messages dm 
		        WHERE dm.conversation_id = c.id ORDER BY created_at DESC LIMIT 1) as last_msg
		FROM dm_conversations c
		JOIN users u ON (CASE WHEN c.user1_id = ? THEN c.user2_id ELSE c.user1_id END) = u.id
		WHERE c.user1_id = ? OR c.user2_id = ?
		ORDER BY c.last_message_at DESC NULLS LAST
	`, userID, userID, userID, userID, userID)
	if err != nil {
		return []*DMConversation{}
	}
	defer rows.Close()

	convs := []*DMConversation{}
	for rows.Next() {
		conv := &DMConversation{}
		var lastMsgStr sql.NullString
		var lastMsgContent sql.NullString
		
		err := rows.Scan(&conv.ID, &conv.OtherUserID, &conv.OtherUsername, 
			&lastMsgStr, &conv.UnreadCount, &lastMsgContent)
		if err != nil {
			continue
		}
		
		if lastMsgStr.Valid {
			t := parseTimestamp(lastMsgStr.String)
			conv.LastMessageAt = &t
		}
		if lastMsgContent.Valid {
			conv.LastMessage = lastMsgContent.String
		}
		
		convs = append(convs, conv)
	}
	return convs
}

// SendDM sends a direct message
func (cs *ChatStore) SendDM(msg *DirectMessage) error {
	if msg.ID == "" {
		msg.ID = generateID()
	}
	if msg.CreatedAt.IsZero() {
		msg.CreatedAt = time.Now()
	}

	tx, err := cs.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO direct_messages (id, conversation_id, sender_id, sender_username, content, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, msg.ID, msg.ConversationID, msg.SenderID, msg.SenderUsername, msg.Content, msg.CreatedAt)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE dm_conversations SET last_message_at = ? WHERE id = ?
	`, msg.CreatedAt, msg.ConversationID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetDMs returns messages for a DM conversation
func (cs *ChatStore) GetDMs(conversationID string, limit int, before string) []*DirectMessage {
	query := `
		SELECT id, conversation_id, sender_id, sender_username, content, created_at, read_at
		FROM direct_messages
		WHERE conversation_id = ?
	`
	args := []interface{}{conversationID}
	
	if before != "" {
		query += " AND created_at < (SELECT created_at FROM direct_messages WHERE id = ?)"
		args = append(args, before)
	}
	
	query += " ORDER BY created_at DESC LIMIT ?"
	args = append(args, limit)

	rows, err := cs.db.Query(query, args...)
	if err != nil {
		return []*DirectMessage{}
	}
	defer rows.Close()

	messages := []*DirectMessage{}
	for rows.Next() {
		msg := &DirectMessage{}
		var createdAtStr string
		var readAtStr sql.NullString
		
		err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, 
			&msg.SenderUsername, &msg.Content, &createdAtStr, &readAtStr)
		if err != nil {
			continue
		}
		
		msg.CreatedAt = parseTimestamp(createdAtStr)
		if readAtStr.Valid {
			t := parseTimestamp(readAtStr.String)
			msg.ReadAt = &t
		}
		
		messages = append(messages, msg)
	}
	
	// Reverse for chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	
	return messages
}

// MarkDMsAsRead marks DMs as read
func (cs *ChatStore) MarkDMsAsRead(conversationID, userID string) error {
	_, err := cs.db.Exec(`
		UPDATE direct_messages SET read_at = ?
		WHERE conversation_id = ? AND sender_id != ? AND read_at IS NULL
	`, time.Now(), conversationID, userID)
	return err
}

// ==================== PRESENCE METHODS ====================

// UpdatePresence updates a user's online status
func (cs *ChatStore) UpdatePresence(userID, status string) error {
	_, err := cs.db.Exec(`
		INSERT INTO user_presence (user_id, last_seen, status)
		VALUES (?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET last_seen = ?, status = ?
	`, userID, time.Now(), status, time.Now(), status)
	return err
}

// GetOnlineUsers returns count of online users
func (cs *ChatStore) GetOnlineUsers() int {
	var count int
	cs.db.QueryRow(`
		SELECT COUNT(*) FROM user_presence 
		WHERE status = 'online' AND last_seen > datetime('now', '-5 minutes')
	`).Scan(&count)
	return count
}

// GetTodayMessageCount returns today's message count
func (cs *ChatStore) GetTodayMessageCount() int {
	var count int
	cs.db.QueryRow(`
		SELECT COUNT(*) FROM chat_messages 
		WHERE date(created_at) = date('now')
	`).Scan(&count)
	return count
}

