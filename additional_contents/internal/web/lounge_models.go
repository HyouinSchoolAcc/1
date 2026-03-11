package web

import (
	"data_labler_ui_go/internal/database"
)

// Re-export types from database package for backward compatibility
type PostType = database.PostType

const (
	PostTypeDailySpark = database.PostTypeDailySpark
	PostTypeChain      = database.PostTypeChain
	PostTypeHotTake    = database.PostTypeHotTake
	PostTypeVibe       = database.PostTypeVibe
	PostTypePoll       = database.PostTypePoll
)

type LoungePost = database.LoungePost
type LoungeReply = database.LoungeReply
type LoungeReaction = database.LoungeReaction
type LoungeActivity = database.LoungeActivity
type LoungeStore = database.LoungeStore

