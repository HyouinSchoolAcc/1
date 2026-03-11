-- Database schema for data labeler UI
-- Using SQLite for simplicity and portability

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Lounge posts table
CREATE TABLE IF NOT EXISTS lounge_posts (
    id TEXT PRIMARY KEY,
    type TEXT NOT NULL,
    author_character TEXT NOT NULL,
    author_user_id TEXT NOT NULL,
    author_username TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    reply_count INTEGER NOT NULL DEFAULT 0,
    is_prompt BOOLEAN NOT NULL DEFAULT 0,
    prompt_date TEXT,
    FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_lounge_posts_timestamp ON lounge_posts(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_lounge_posts_type ON lounge_posts(type);
CREATE INDEX IF NOT EXISTS idx_lounge_posts_author ON lounge_posts(author_user_id);
CREATE INDEX IF NOT EXISTS idx_lounge_posts_prompt_date ON lounge_posts(prompt_date);

-- Lounge replies table
CREATE TABLE IF NOT EXISTS lounge_replies (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    parent_reply_id TEXT,
    author_character TEXT NOT NULL,
    author_user_id TEXT NOT NULL,
    author_username TEXT NOT NULL,
    content TEXT NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES lounge_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_reply_id) REFERENCES lounge_replies(id) ON DELETE CASCADE,
    FOREIGN KEY (author_user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_lounge_replies_post_id ON lounge_replies(post_id);
CREATE INDEX IF NOT EXISTS idx_lounge_replies_timestamp ON lounge_replies(timestamp);
CREATE INDEX IF NOT EXISTS idx_lounge_replies_author ON lounge_replies(author_user_id);

-- Lounge reactions table
CREATE TABLE IF NOT EXISTS lounge_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    post_id TEXT,
    reply_id TEXT,
    emoji TEXT NOT NULL,
    timestamp DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES lounge_posts(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_id) REFERENCES lounge_replies(id) ON DELETE CASCADE,
    UNIQUE(user_id, post_id, emoji),
    UNIQUE(user_id, reply_id, emoji),
    CHECK ((post_id IS NOT NULL AND reply_id IS NULL) OR (post_id IS NULL AND reply_id IS NOT NULL))
);

CREATE INDEX IF NOT EXISTS idx_lounge_reactions_post ON lounge_reactions(post_id);
CREATE INDEX IF NOT EXISTS idx_lounge_reactions_reply ON lounge_reactions(reply_id);
CREATE INDEX IF NOT EXISTS idx_lounge_reactions_user ON lounge_reactions(user_id);

-- Sessions table (optional - for persistent sessions)
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON sessions(expires_at);

-- User characters table
CREATE TABLE IF NOT EXISTS user_characters (
    id TEXT PRIMARY KEY,
    creator_id TEXT NOT NULL,
    creator_name TEXT NOT NULL,
    name TEXT NOT NULL,
    character_values TEXT NOT NULL,
    experiences TEXT NOT NULL,
    judgements TEXT NOT NULL,
    abilities TEXT NOT NULL,
    story TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (creator_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_characters_creator ON user_characters(creator_id);
CREATE INDEX IF NOT EXISTS idx_user_characters_created_at ON user_characters(created_at DESC);

-- New character discussion messages table
CREATE TABLE IF NOT EXISTS newcharacter_messages (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    section TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_newcharacter_messages_section ON newcharacter_messages(section, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_newcharacter_messages_user ON newcharacter_messages(user_id);

-- Character votes table
CREATE TABLE IF NOT EXISTS character_votes (
    id TEXT PRIMARY KEY,
    character_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (character_id) REFERENCES user_characters(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(character_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_character_votes_character_id ON character_votes(character_id);
CREATE INDEX IF NOT EXISTS idx_character_votes_user_id ON character_votes(user_id);

-- Email confirmation tokens table
CREATE TABLE IF NOT EXISTS email_confirmations (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL,
    username TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    confirmed BOOLEAN NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_email_confirmations_token ON email_confirmations(token);
CREATE INDEX IF NOT EXISTS idx_email_confirmations_email ON email_confirmations(email);
CREATE INDEX IF NOT EXISTS idx_email_confirmations_expires ON email_confirmations(expires_at);

-- Writer tutorial progress table
-- Tracks user's progress through the writer onboarding tutorial
CREATE TABLE IF NOT EXISTS tutorial_progress (
    id TEXT PRIMARY KEY,
    user_id TEXT UNIQUE NOT NULL,
    current_step INTEGER NOT NULL DEFAULT 0,
    completed_steps TEXT NOT NULL DEFAULT '[]',  -- JSON array of completed step IDs
    started_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    is_completed BOOLEAN NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_tutorial_progress_user_id ON tutorial_progress(user_id);

-- Character certifications table
-- Tracks which characters a user is certified to write for
CREATE TABLE IF NOT EXISTS character_certifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    character_id TEXT NOT NULL,
    certified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    certified_by TEXT,  -- editor who approved certification, null if self-certified through tutorial
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, character_id)
);

CREATE INDEX IF NOT EXISTS idx_character_certifications_user ON character_certifications(user_id);
CREATE INDEX IF NOT EXISTS idx_character_certifications_character ON character_certifications(character_id);

-- Demo dialogues table
-- Stores curated demo dialogues to show on the landing page
CREATE TABLE IF NOT EXISTS demo_dialogues (
    id TEXT PRIMARY KEY,
    character_id TEXT NOT NULL,
    dialogue_data TEXT NOT NULL,  -- JSON data of the dialogue
    display_order INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT 1,
    added_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    added_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_demo_dialogues_character ON demo_dialogues(character_id);
CREATE INDEX IF NOT EXISTS idx_demo_dialogues_order ON demo_dialogues(display_order);

-- ============================================
-- DISCORD-STYLE CHAT SYSTEM TABLES
-- ============================================

-- Chat channels (rooms) table
CREATE TABLE IF NOT EXISTS chat_channels (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    created_by TEXT,  -- NULL for system-created channels
    created_from_quote TEXT,  -- If created from selected text
    is_default BOOLEAN NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_channels_created_at ON chat_channels(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_channels_last_message ON chat_channels(last_message_at DESC);

-- Chat messages table
CREATE TABLE IF NOT EXISTS chat_messages (
    id TEXT PRIMARY KEY,
    channel_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    username TEXT NOT NULL,
    content TEXT NOT NULL,
    reply_to_id TEXT,  -- For threaded replies
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    edited_at DATETIME,
    FOREIGN KEY (channel_id) REFERENCES chat_channels(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (reply_to_id) REFERENCES chat_messages(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_channel ON chat_messages(channel_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_chat_messages_user ON chat_messages(user_id);

-- Message reactions table
CREATE TABLE IF NOT EXISTS chat_message_reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    message_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    emoji TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (message_id) REFERENCES chat_messages(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(message_id, user_id, emoji)
);

CREATE INDEX IF NOT EXISTS idx_chat_reactions_message ON chat_message_reactions(message_id);

-- Friends table
CREATE TABLE IF NOT EXISTS friends (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    friend_id TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',  -- pending, accepted, blocked
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    accepted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, friend_id)
);

CREATE INDEX IF NOT EXISTS idx_friends_user ON friends(user_id, status);
CREATE INDEX IF NOT EXISTS idx_friends_friend ON friends(friend_id, status);

-- Direct message conversations table
CREATE TABLE IF NOT EXISTS dm_conversations (
    id TEXT PRIMARY KEY,
    user1_id TEXT NOT NULL,
    user2_id TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_message_at DATETIME,
    FOREIGN KEY (user1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (user2_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user1_id, user2_id)
);

CREATE INDEX IF NOT EXISTS idx_dm_conversations_users ON dm_conversations(user1_id, user2_id);

-- Direct messages table
CREATE TABLE IF NOT EXISTS direct_messages (
    id TEXT PRIMARY KEY,
    conversation_id TEXT NOT NULL,
    sender_id TEXT NOT NULL,
    sender_username TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    read_at DATETIME,
    FOREIGN KEY (conversation_id) REFERENCES dm_conversations(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_direct_messages_conversation ON direct_messages(conversation_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_direct_messages_unread ON direct_messages(conversation_id, read_at);

-- User online status tracking
CREATE TABLE IF NOT EXISTS user_presence (
    user_id TEXT PRIMARY KEY,
    last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'offline',  -- online, away, offline
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- USER PROFILES & POINTS SYSTEM
-- ============================================

-- Extended profile info (bio, avatar, display settings)
CREATE TABLE IF NOT EXISTS user_profiles (
    user_id TEXT PRIMARY KEY,
    bio TEXT NOT NULL DEFAULT '',
    avatar_color TEXT NOT NULL DEFAULT '#6366f1',
    display_name TEXT NOT NULL DEFAULT '',
    is_public BOOLEAN NOT NULL DEFAULT 1,
    xp INTEGER NOT NULL DEFAULT 0,
    days_logged_in INTEGER NOT NULL DEFAULT 0,
    last_login_date TEXT NOT NULL DEFAULT '',
    equipped_profile_background TEXT NOT NULL DEFAULT '',
    equipped_banner TEXT NOT NULL DEFAULT '',
    equipped_profile_border TEXT NOT NULL DEFAULT '',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Point transactions: manual bonuses and event-driven awards
CREATE TABLE IF NOT EXISTS point_transactions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    amount INTEGER NOT NULL,
    reason TEXT NOT NULL DEFAULT '',
    awarded_by TEXT NOT NULL DEFAULT 'system',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_point_transactions_user ON point_transactions(user_id, created_at DESC);

-- ============================================
-- STORE: user inventory and equipped cosmetics
-- ============================================

-- Items the user has purchased (one row per user per item)
CREATE TABLE IF NOT EXISTS user_store_inventory (
    user_id TEXT NOT NULL,
    item_id TEXT NOT NULL,
    purchased_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, item_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_store_inventory_user ON user_store_inventory(user_id);

-- Cosmetics owned by users (banners, backgrounds)
CREATE TABLE IF NOT EXISTS user_cosmetics (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    cosmetic_type TEXT NOT NULL,
    cosmetic_id TEXT NOT NULL,
    granted_by TEXT NOT NULL DEFAULT 'system',
    granted_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, cosmetic_type, cosmetic_id)
);

CREATE INDEX IF NOT EXISTS idx_user_cosmetics_user ON user_cosmetics(user_id, cosmetic_type);

