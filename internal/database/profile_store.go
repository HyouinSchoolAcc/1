package database

import (
	"database/sql"
	"time"
)

// UserProfile stores extended profile info for a user
type UserProfile struct {
	UserID           string    `json:"user_id"`
	Bio              string    `json:"bio"`
	AvatarColor      string    `json:"avatar_color"`
	DisplayName      string    `json:"display_name"`
	IsPublic         bool      `json:"is_public"`
	XP               int       `json:"xp"`
	DaysLoggedIn     int       `json:"days_logged_in"`
	LastLoginDate    string    `json:"last_login_date"`
	ActiveBanner     string    `json:"active_banner"`
	ActiveBackground string    `json:"active_background"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// PointTransaction records a bonus point award or deduction (for Ink Tokens)
type PointTransaction struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Reason    string    `json:"reason"`
	AwardedBy string    `json:"awarded_by"`
	CreatedAt time.Time `json:"created_at"`
}

// UserCosmetic represents a cosmetic item owned by a user
type UserCosmetic struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	CosmeticType string    `json:"cosmetic_type"` // "banner" | "background"
	CosmeticID   string    `json:"cosmetic_id"`   // filename
	GrantedBy    string    `json:"granted_by"`
	GrantedAt    time.Time `json:"granted_at"`
}

// ProfileStore handles user profile, XP, cosmetic, and bonus point persistence
type ProfileStore struct {
	db *DB
}

// NewProfileStore creates a new profile store
func NewProfileStore(db *DB) *ProfileStore {
	return &ProfileStore{db: db}
}

// GetProfile retrieves a user profile, auto-creating a default if none exists
func (ps *ProfileStore) GetProfile(userID string) (*UserProfile, error) {
	profile := &UserProfile{}
	err := ps.db.QueryRow(
		`SELECT user_id, bio, avatar_color, display_name, is_public,
		        COALESCE(xp,0), COALESCE(days_logged_in,0), COALESCE(last_login_date,''),
		        COALESCE(active_banner,''), COALESCE(active_background,''), updated_at
		 FROM user_profiles WHERE user_id = ?`,
		userID,
	).Scan(&profile.UserID, &profile.Bio, &profile.AvatarColor, &profile.DisplayName,
		&profile.IsPublic, &profile.XP, &profile.DaysLoggedIn, &profile.LastLoginDate,
		&profile.ActiveBanner, &profile.ActiveBackground, &profile.UpdatedAt)

	if err == sql.ErrNoRows {
		profile = &UserProfile{
			UserID:      userID,
			Bio:         "",
			AvatarColor: "#6366f1",
			DisplayName: "",
			IsPublic:    true,
			UpdatedAt:   time.Now(),
		}
		_, err = ps.db.Exec(
			`INSERT INTO user_profiles (user_id, bio, avatar_color, display_name, is_public,
			 xp, days_logged_in, last_login_date, active_banner, active_background, updated_at)
			 VALUES (?, ?, ?, ?, ?, 0, 0, '', '', '', ?)`,
			profile.UserID, profile.Bio, profile.AvatarColor,
			profile.DisplayName, profile.IsPublic, profile.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		return profile, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateProfile upserts a user profile (does not overwrite XP/days/login tracking fields)
func (ps *ProfileStore) UpdateProfile(profile *UserProfile) error {
	profile.UpdatedAt = time.Now()
	_, err := ps.db.Exec(
		`INSERT INTO user_profiles
		 (user_id, bio, avatar_color, display_name, is_public, xp, days_logged_in,
		  last_login_date, active_banner, active_background, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(user_id) DO UPDATE SET
		   bio=excluded.bio, avatar_color=excluded.avatar_color,
		   display_name=excluded.display_name, is_public=excluded.is_public,
		   active_banner=excluded.active_banner, active_background=excluded.active_background,
		   updated_at=excluded.updated_at`,
		profile.UserID, profile.Bio, profile.AvatarColor, profile.DisplayName, profile.IsPublic,
		profile.XP, profile.DaysLoggedIn, profile.LastLoginDate,
		profile.ActiveBanner, profile.ActiveBackground, profile.UpdatedAt,
	)
	return err
}

// AddXP atomically adds XP to a user's profile
func (ps *ProfileStore) AddXP(userID string, amount int) error {
	_, err := ps.db.Exec(
		`INSERT INTO user_profiles (user_id, bio, avatar_color, display_name, is_public, xp,
		 days_logged_in, last_login_date, active_banner, active_background, updated_at)
		 VALUES (?, '', '#6366f1', '', 1, ?, 0, '', '', '', ?)
		 ON CONFLICT(user_id) DO UPDATE SET xp = xp + excluded.xp`,
		userID, amount, time.Now(),
	)
	return err
}

// RecordDailyLogin checks if the user has logged in today; if not, increments days_logged_in and adds XP.
func (ps *ProfileStore) RecordDailyLogin(userID string) (bool, error) {
	today := time.Now().Format("2006-01-02")
	profile, err := ps.GetProfile(userID)
	if err != nil {
		return false, err
	}
	if profile.LastLoginDate == today {
		return false, nil
	}
	_, err = ps.db.Exec(
		`UPDATE user_profiles SET days_logged_in = days_logged_in + 1, last_login_date = ?, xp = xp + 5
		 WHERE user_id = ?`,
		today, userID,
	)
	return true, err
}

// GrantCosmetic gives a cosmetic to a user (idempotent via IGNORE)
func (ps *ProfileStore) GrantCosmetic(userID, cosmeticType, cosmeticID, grantedBy string) error {
	id := generateID() + userID[max(0, len(userID)-4):]
	_, err := ps.db.Exec(
		`INSERT OR IGNORE INTO user_cosmetics (id, user_id, cosmetic_type, cosmetic_id, granted_by, granted_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, userID, cosmeticType, cosmeticID, grantedBy, time.Now(),
	)
	return err
}

// HasCosmetic checks if a user owns a specific cosmetic
func (ps *ProfileStore) HasCosmetic(userID, cosmeticType, cosmeticID string) (bool, error) {
	var count int
	err := ps.db.QueryRow(
		`SELECT COUNT(*) FROM user_cosmetics WHERE user_id=? AND cosmetic_type=? AND cosmetic_id=?`,
		userID, cosmeticType, cosmeticID,
	).Scan(&count)
	return count > 0, err
}

// GetUserCosmetics returns all cosmetics of a given type owned by a user
func (ps *ProfileStore) GetUserCosmetics(userID, cosmeticType string) ([]string, error) {
	rows, err := ps.db.Query(
		`SELECT cosmetic_id FROM user_cosmetics WHERE user_id=? AND cosmetic_type=? ORDER BY granted_at`,
		userID, cosmeticType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// AwardInk records an Ink Token transaction for a user
func (ps *ProfileStore) AwardInk(userID string, amount int, reason, awardedBy string) error {
	suffix := userID
	if len(suffix) > 4 {
		suffix = suffix[len(suffix)-4:]
	}
	id := generateID() + suffix
	_, err := ps.db.Exec(
		`INSERT INTO point_transactions (id, user_id, amount, reason, awarded_by, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		id, userID, amount, reason, awardedBy, time.Now(),
	)
	return err
}

// AwardPoints is an alias for AwardInk (backwards compatibility)
func (ps *ProfileStore) AwardPoints(userID string, amount int, reason, awardedBy string) error {
	return ps.AwardInk(userID, amount, reason, awardedBy)
}

// GetBonusInk returns total manually awarded Ink Tokens
func (ps *ProfileStore) GetBonusInk(userID string) (int, error) {
	var total int
	err := ps.db.QueryRow(
		`SELECT COALESCE(SUM(amount), 0) FROM point_transactions WHERE user_id = ?`,
		userID,
	).Scan(&total)
	return total, err
}

// GetBonusPoints is an alias for GetBonusInk
func (ps *ProfileStore) GetBonusPoints(userID string) (int, error) {
	return ps.GetBonusInk(userID)
}

// GetInkHistory returns recent Ink Token transactions
func (ps *ProfileStore) GetInkHistory(userID string, limit int) ([]PointTransaction, error) {
	rows, err := ps.db.Query(
		`SELECT id, user_id, amount, reason, awarded_by, created_at FROM point_transactions
		 WHERE user_id = ? ORDER BY created_at DESC LIMIT ?`,
		userID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txns []PointTransaction
	for rows.Next() {
		var t PointTransaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Reason, &t.AwardedBy, &t.CreatedAt); err != nil {
			return nil, err
		}
		txns = append(txns, t)
	}
	return txns, rows.Err()
}

// GetPointHistory is an alias for GetInkHistory
func (ps *ProfileStore) GetPointHistory(userID string, limit int) ([]PointTransaction, error) {
	return ps.GetInkHistory(userID, limit)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
