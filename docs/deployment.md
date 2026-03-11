# Deployment Guide

## Prerequisites

- Go 1.21+
- SQLite3
- ngrok account with custom domain configured
- Python 3 with PIL (for dialogue image generation)

## Production Setup

### 1. Configure ngrok

Ensure ngrok is authenticated and the custom domain is configured:

```bash
ngrok config add-authtoken YOUR_TOKEN
```

The script uses `wl2.studio` as the domain. To change this, edit `run_go3.sh`:

```bash
NGROK_HOST="your-domain.ngrok.io"
```

### 2. Environment Configuration

Edit `run_go3.sh` to customize:

```bash
PORT="5002"           # Server port
NGROK_HOST="wl2.studio"  # Public domain
```

### 3. Database Initialization

The database is automatically created on first run. To manually initialize:

```bash
cd data_labler_UI_production
go run migrate_to_sql.go data
```

### 4. Start Services

```bash
./run_go3.sh
```

This script:
1. Kills any existing processes on ports 5002 and 5003
2. Kills existing ngrok processes
3. Checks database integrity
4. Builds the server binary
5. Starts the server with debug mode
6. Starts ngrok tunnel

### 5. Verify Deployment

```bash
# Check server is running
curl http://localhost:5002/api/current_user

# Check ngrok tunnel
curl https://wl2.studio/api/current_user
```

## Stopping Services

The PIDs are displayed when `run_go3.sh` starts:

```bash
kill $BACKEND_PID   # Stop server
kill $NGROK_PID     # Stop ngrok
```

Or forcefully:

```bash
lsof -ti:5002 | xargs kill -9
pkill -9 ngrok
```

## Backup Procedures

### Database Backup

```bash
sqlite3 data/app.db ".backup data/backup_$(date +%Y%m%d).db"
```

### Full Backup

```bash
tar -czvf backup_$(date +%Y%m%d).tar.gz \
    data/app.db \
    presets/
```

## Monitoring

### View Logs

```bash
# Server logs
tail -f data_labler_UI_production/server.log

# ngrok logs
tail -f data_labler_UI_production/ngrok.log
```

### Check Database Size

```bash
ls -lh data/app.db
```

### Database Integrity Check

```bash
sqlite3 data/app.db "PRAGMA integrity_check;"
```

## Troubleshooting

### Server Won't Start

1. **Port already in use:**
   ```bash
   lsof -ti:5002 | xargs kill -9
   ```

2. **Database corrupted:**
   ```bash
   # Backup and recreate
   mv data/app.db data/app.db.corrupted
   # Server will create new database on next start
   ```

3. **Missing dependencies:**
   ```bash
   go mod download
   ```

### ngrok Issues

1. **Domain not working:**
   - Check ngrok dashboard for errors
   - Verify authtoken is set
   - Confirm domain ownership

2. **Tunnel keeps disconnecting:**
   - Check internet connection
   - Review ngrok logs for errors

### Database Issues

1. **"Database locked" error:**
   - Only one server instance should run
   - Kill all server processes: `pkill -f server_sql`

2. **Migration fails:**
   - Check existing JSON files for valid format
   - Review migration script output for specific errors

## Security Considerations

### Production Checklist

- [ ] Change default test passwords
- [ ] Update session key in `internal/web/auth.go`
- [ ] Enable HTTPS-only cookies
- [ ] Configure proper CORS if needed
- [ ] Set up database encryption at rest
- [ ] Regular backup schedule

### Session Security

The default session key should be changed for production:

```go
// internal/web/auth.go
SessionKey = "your-secure-random-key-here"
```

Generate a secure key:

```bash
openssl rand -hex 32
```

## Performance Tuning

### Database Optimization

```bash
# Run periodically
sqlite3 data/app.db "VACUUM;"
sqlite3 data/app.db "ANALYZE;"
```

### Connection Limits

SQLite handles concurrent reads well but serializes writes. For high-write workloads, consider migrating to PostgreSQL.

## Scaling

### Current Limitations

- Single server instance
- SQLite write serialization
- In-memory session storage

### Future Improvements

For higher scale:
1. Migrate to PostgreSQL
2. Add Redis for session storage
3. Deploy behind load balancer
4. Separate static asset serving
