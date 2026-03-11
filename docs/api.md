# API Reference

All API endpoints are prefixed with `/api/` unless otherwise noted.

## Authentication

### Get Current User
```
GET /api/current_user
```

**Response:**
```json
{
    "logged_in": true,
    "username": "testwriter",
    "role": "writer",
    "user_id": "12345"
}
```

### Login
```
POST /login
Content-Type: application/x-www-form-urlencoded

username=testwriter&password=password123
```

### Register
```
POST /register
Content-Type: application/x-www-form-urlencoded

username=newuser&email=user@example.com&password=password123
```

## File Operations

### List Writer Files
```
GET /load_structured_writer_files?preset_set=presets_kurisu_CN
```

**Response:**
```json
{
    "users": {
        "user_1": {
            "days": {
                "Day1": {
                    "files": [
                        {
                            "name": "user_1_Day1_dup_0_simplified.json",
                            "qc_status": "pending"
                        }
                    ]
                }
            }
        }
    }
}
```

### Load File Content
```
POST /load_writer_file_content
Content-Type: application/json

{
    "filename": "user_1_Day1_dup_0_simplified.json",
    "preset_set": "presets_kurisu_CN"
}
```

### Save File Content
```
POST /save_writer_file_content
Content-Type: application/json

{
    "filename": "user_1_Day1_dup_0_simplified.json",
    "content": { ... },
    "preset_set": "presets_kurisu_CN"
}
```

### Update QC Status
```
POST /update_qc_status
Content-Type: application/json

{
    "filename": "user_1_Day1_dup_0_simplified.json",
    "qc_status": "approved",
    "preset_set": "presets_kurisu_CN"
}
```

## Writers' Lounge

### Get Posts
```
GET /api/lounge/posts?type=daily_spark&limit=20
```

### Create Post
```
POST /api/lounge/posts
Content-Type: application/json

{
    "type": "vibe",
    "content": "Anyone else struggling with dialogue pacing?",
    "author_character": "Writer"
}
```

### Reply to Post
```
POST /api/lounge/posts/{postId}/replies
Content-Type: application/json

{
    "content": "Yes! Try reading it aloud.",
    "author_character": "Writer"
}
```

### React to Post
```
POST /api/lounge/posts/{postId}/reactions
Content-Type: application/json

{
    "emoji": "🔥"
}
```

## Discord-Style Chat

### Get Channels
```
GET /api/chat/channels
```

### Create Channel
```
POST /api/chat/channels
Content-Type: application/json

{
    "name": "character-discussion",
    "description": "Discuss character development"
}
```

### Get Messages
```
GET /api/chat/channels/{channelId}/messages?limit=50
```

### Send Message
```
POST /api/chat/channels/{channelId}/messages
Content-Type: application/json

{
    "content": "Hello everyone!"
}
```

### React to Message
```
POST /api/chat/messages/{messageId}/reactions
Content-Type: application/json

{
    "emoji": "👍"
}
```

### Direct Messages
```
GET /api/chat/dm/{userId}/messages
POST /api/chat/dm/{userId}/messages
```

## Tutorial System

### Get Progress
```
GET /api/tutorial/progress
```

### Update Progress
```
POST /api/tutorial/progress
Content-Type: application/json

{
    "step": 3,
    "completed": true
}
```

### Get Certifications
```
GET /api/tutorial/certifications
```

## Payment API

### Get Earnings
```
GET /api/payment/earnings
```

### Get Payment History
```
GET /api/payment/history
```

## Character API

### Get Characters
```
GET /api/characters
```

### Vote for Character
```
POST /api/characters/{characterId}/vote
```

## LLM Integration

### Check Status
```
GET /api/llm/status
```

### Chat with Model
```
POST /api/llm/chat
Content-Type: application/json

{
    "messages": [
        {"role": "user", "content": "你好"}
    ],
    "use_lora": true,
    "max_tokens": 512
}
```

### Toggle Server (Editor only)
```
POST /api/llm/toggle
Content-Type: application/json

{
    "action": "start"
}
```

## Error Responses

All errors return appropriate HTTP status codes with JSON body:

```json
{
    "error": "Description of the error"
}
```

Common status codes:
- `400` — Bad request (invalid input)
- `401` — Unauthorized (not logged in)
- `403` — Forbidden (insufficient permissions)
- `404` — Not found
- `500` — Internal server error
