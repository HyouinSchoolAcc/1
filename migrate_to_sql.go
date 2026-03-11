package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"data_labler_ui_go/internal/database"
)

// Migration script to move data from JSON files to SQL database

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate_to_sql.go <data_directory>")
	}

	dataDir := os.Args[1]
	log.Printf("Starting migration from %s", dataDir)

	// Initialize database
	db, err := database.New(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	userStore := database.NewUserStore(db)
	loungeStore := database.NewLoungeStore(db)

	// Migrate users
	if err := migrateUsers(dataDir, userStore, db); err != nil {
		log.Printf("Warning: Failed to migrate users: %v", err)
	} else {
		log.Println("✓ Users migrated successfully")
	}

	// Migrate lounge data
	loungeDataDir := filepath.Join(dataDir, "lounge")
	if err := migrateLoungePosts(loungeDataDir, loungeStore); err != nil {
		log.Printf("Warning: Failed to migrate lounge posts: %v", err)
	} else {
		log.Println("✓ Lounge posts migrated successfully")
	}

	if err := migrateLoungeReplies(loungeDataDir, loungeStore); err != nil {
		log.Printf("Warning: Failed to migrate lounge replies: %v", err)
	} else {
		log.Println("✓ Lounge replies migrated successfully")
	}

	if err := migrateLoungeReactions(loungeDataDir, loungeStore, db); err != nil {
		log.Printf("Warning: Failed to migrate lounge reactions: %v", err)
	} else {
		log.Println("✓ Lounge reactions migrated successfully")
	}

	log.Println("Migration completed!")
	log.Println("\nBackup your JSON files before deleting them.")
	log.Println("Once you've verified the migration, you can remove:")
	log.Printf("  - %s/users.json\n", dataDir)
	log.Printf("  - %s/lounge/posts.json\n", dataDir)
	log.Printf("  - %s/lounge/replies.json\n", dataDir)
	log.Printf("  - %s/lounge/reactions.json\n", dataDir)
}

func migrateUsers(dataDir string, userStore *database.UserStore, db *database.DB) error {
	usersFile := filepath.Join(dataDir, "users.json")
	
	// Check if file exists
	if _, err := os.Stat(usersFile); os.IsNotExist(err) {
		log.Println("No users.json file found, skipping user migration")
		return nil
	}

	data, err := os.ReadFile(usersFile)
	if err != nil {
		return fmt.Errorf("failed to read users file: %w", err)
	}

	type OldUser struct {
		ID        string    `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		Role      string    `json:"role"`
		CreatedAt time.Time `json:"created_at"`
	}

	var users []OldUser
	if err := json.Unmarshal(data, &users); err != nil {
		return fmt.Errorf("failed to parse users JSON: %w", err)
	}

	log.Printf("Migrating %d users...", len(users))

	for _, u := range users {
		// Insert directly into database (password is already hashed)
		_, exists := userStore.GetUser(u.Username)
		if exists {
			log.Printf("  - User %s already exists, skipping", u.Username)
			continue
		}

		// Direct SQL insert to preserve hashed password
		_, err = db.Exec(
			"INSERT INTO users (id, username, email, password, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
			u.ID, u.Username, u.Email, u.Password, u.Role, u.CreatedAt,
		)
		if err != nil {
			log.Printf("  - Failed to migrate user %s: %v", u.Username, err)
		} else {
			log.Printf("  - Migrated user: %s (%s)", u.Username, u.Role)
		}
	}

	return nil
}

func migrateLoungePosts(loungeDir string, loungeStore *database.LoungeStore) error {
	postsFile := filepath.Join(loungeDir, "posts.json")
	
	// Check if file exists
	if _, err := os.Stat(postsFile); os.IsNotExist(err) {
		log.Println("No posts.json file found, skipping posts migration")
		return nil
	}

	data, err := os.ReadFile(postsFile)
	if err != nil {
		return fmt.Errorf("failed to read posts file: %w", err)
	}

	var posts []*database.LoungePost
	if err := json.Unmarshal(data, &posts); err != nil {
		return fmt.Errorf("failed to parse posts JSON: %w", err)
	}

	log.Printf("Migrating %d posts...", len(posts))

	for _, post := range posts {
		if err := loungeStore.CreatePost(post); err != nil {
			log.Printf("  - Failed to migrate post %s: %v", post.ID, err)
		} else {
			log.Printf("  - Migrated post: %s (by %s)", post.ID, post.AuthorUsername)
		}
	}

	return nil
}

func migrateLoungeReplies(loungeDir string, loungeStore *database.LoungeStore) error {
	repliesFile := filepath.Join(loungeDir, "replies.json")
	
	// Check if file exists
	if _, err := os.Stat(repliesFile); os.IsNotExist(err) {
		log.Println("No replies.json file found, skipping replies migration")
		return nil
	}

	data, err := os.ReadFile(repliesFile)
	if err != nil {
		return fmt.Errorf("failed to read replies file: %w", err)
	}

	var replies []*database.LoungeReply
	if err := json.Unmarshal(data, &replies); err != nil {
		return fmt.Errorf("failed to parse replies JSON: %w", err)
	}

	log.Printf("Migrating %d replies...", len(replies))

	for _, reply := range replies {
		if err := loungeStore.CreateReply(reply); err != nil {
			log.Printf("  - Failed to migrate reply %s: %v", reply.ID, err)
		} else {
			log.Printf("  - Migrated reply: %s (by %s)", reply.ID, reply.AuthorUsername)
		}
	}

	return nil
}

func migrateLoungeReactions(loungeDir string, loungeStore *database.LoungeStore, db *database.DB) error {
	reactionsFile := filepath.Join(loungeDir, "reactions.json")
	
	// Check if file exists
	if _, err := os.Stat(reactionsFile); os.IsNotExist(err) {
		log.Println("No reactions.json file found, skipping reactions migration")
		return nil
	}

	data, err := os.ReadFile(reactionsFile)
	if err != nil {
		return fmt.Errorf("failed to read reactions file: %w", err)
	}

	var reactions []*database.LoungeReaction
	if err := json.Unmarshal(data, &reactions); err != nil {
		return fmt.Errorf("failed to parse reactions JSON: %w", err)
	}

	log.Printf("Migrating %d reactions...", len(reactions))

	for _, reaction := range reactions {
		isPost := reaction.PostID != ""
		targetID := reaction.PostID
		if !isPost {
			targetID = reaction.ReplyID
		}

		if err := loungeStore.AddReaction(reaction.UserID, targetID, reaction.Emoji, isPost); err != nil {
			// Ignore "already reacted" errors
			if err.Error() != "already reacted" {
				log.Printf("  - Failed to migrate reaction: %v", err)
			}
		}
	}

	log.Printf("  - Migrated reactions")
	return nil
}

