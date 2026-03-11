# Language System Migration

## What Changed

The website has been migrated from the `/page/e` suffix pattern to a cleaner `/en/page` and `/zh/page` prefix pattern for internationalization.

### Before (Old System)
- Chinese: `/landing`, `/guide`, `/navigation`, etc.
- English: `/landing/e`, `/guide/e`, `/navigation/e`, etc.
- **Problem**: Required duplicate route definitions for every endpoint

### After (New System)
- Chinese: `/zh/landing`, `/zh/guide`, `/zh/navigation`, etc.
- English: `/en/landing`, `/en/guide`, `/en/navigation`, etc.
- Root: `/` shows language selection page
- **Benefit**: Single route definition, automatic redirects from old URLs

## Key Features

### 1. Language Selection Page
- Accessing `wl2.studio/` shows a beautiful language picker
- Users choose English or Chinese
- Choice is remembered via cookie

### 2. Automatic URL Redirects
Old URLs automatically redirect to new format:
- `/e` → `/en/landing` (301 permanent redirect)
- `/guide/e` → `/en/guide` (301 permanent redirect)
- `/guide` → `/zh/guide` (301 permanent redirect)

This means all your existing links will continue to work!

### 3. Language Persistence
- Language choice stored in cookie for 1 year
- Users stay in their chosen language across sessions

### 4. Language Switcher
- Both landing pages now have a language switcher in the header
- One click to switch between English and Chinese

## Files Modified

### New Files
1. `templates/language_select.html` - Beautiful language selection page
2. `internal/web/language_middleware.go` - Middleware for language handling and redirects

### Modified Files
1. `internal/web/router.go`
   - Added middleware
   - Removed duplicate route definitions (cut routes in half!)
   - Added language selection page handler

2. `internal/web/handlers.go`
   - Updated `getTemplateData()` to use middleware language detection
   - Added `handleLandingMultiLang()` for template selection
   - Added `langURL()` template function for generating language-aware URLs
   - Updated character routes to include language prefix

3. `templates/landing.html` & `templates/landing_en.html`
   - Updated all navigation links to use `langURL` function
   - Added language switcher component

## Testing

### 1. Build the Application
```bash
cd /home/exx/Desktop/fine-tune/data_labler_UI_production
go build -o server_sql ./cmd/server
```

### 2. Start the Server
```bash
./server_sql
# Or use your existing run script:
# ./run_go3.sh
```

### 3. Test These URLs

#### Root and Language Selection
- `http://localhost:5002/` - Should show language selection page

#### New URL Format (These should work)
- `http://localhost:5002/en/landing` - English landing page
- `http://localhost:5002/zh/landing` - Chinese landing page
- `http://localhost:5002/en/guide` - English guide
- `http://localhost:5002/zh/guide` - Chinese guide

#### Old URL Redirects (Should redirect automatically)
- `http://localhost:5002/e` → redirects to `/en/landing`
- `http://localhost:5002/guide/e` → redirects to `/en/guide`
- `http://localhost:5002/guide` → redirects to `/zh/guide`

#### Navigation
- Click links in the header menu - should maintain language
- Click language switcher - should switch languages
- Refresh page - should remember your language choice

## Next Steps (Optional Improvements)

### 1. Merge Duplicate Templates
Currently we still have separate `landing.html` and `landing_en.html`. You could:
- Create a translation key system
- Use a single template with translation strings
- Reduce template duplication

### 2. Add More Language Switchers
Add the language switcher component to other pages:
- FAQ page
- Guide page
- Writing pages
- etc.

### 3. Translation Management
Consider using a translation file (JSON/YAML) instead of hardcoded strings:
```json
{
  "zh": {
    "nav.guide": "写手指南",
    "nav.writing": "写作"
  },
  "en": {
    "nav.guide": "Writer's Guide",
    "nav.writing": "Writing"
  }
}
```

### 4. Update API Endpoints
Some API endpoints might need language context. Consider:
- Passing language in request headers
- Adding language to API responses
- Internationalizing error messages

## Benefits Summary

✅ **Reduced Code**: Cut route definitions in half  
✅ **SEO Friendly**: Standard `/en/` and `/zh/` prefixes  
✅ **Backward Compatible**: Old URLs redirect automatically  
✅ **User Friendly**: Language selection page + switcher  
✅ **Maintainable**: Single source of truth for routes  
✅ **Scalable**: Easy to add new languages (French, Spanish, etc.)  

## Rollback (If Needed)

If you need to rollback, the old code is still in git history. However, both old and new URL patterns work, so there's no risk of breaking existing links.

---

**Status**: ✅ Implementation Complete  
**Build**: ✅ Successful  
**Tests**: Ready for manual testing  
**Documentation**: Complete  
