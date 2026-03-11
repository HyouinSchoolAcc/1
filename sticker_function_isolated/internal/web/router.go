package web

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

// RegisterRoutes wires all HTTP handlers.
func RegisterRoutes(r chi.Router, cfg Config) {
	app := NewApp(cfg)
	
	// Start background cleanup routine for temp characters
	app.StartTempCleanup()

	// Apply language middleware globally
	r.Use(RedirectLegacyLanguageURLs) // Redirect old /e URLs first
	r.Use(LanguageMiddleware)          // Then handle language detection

	// Load Go templates for landing and writing pages - fail fast if they don't load
	log.Printf("Loading Go templates from: %s", cfg.TemplatesDir)
	landingTpl, err := loadGoTemplate(cfg.TemplatesDir, "landing.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load landing.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded landing.html template")

	landingEnTpl, err := loadGoTemplate(cfg.TemplatesDir, "landing_en.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load landing_en.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded landing_en.html template")

	writingTpl, err := loadGoTemplate(cfg.TemplatesDir, "writing.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load writing.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded writing.html template")

	// Load additional page templates
	faqTpl, err := loadGoTemplate(cfg.TemplatesDir, "faq.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load faq.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded faq.html template")

	characterTpl, err := loadGoTemplate(cfg.TemplatesDir, "character.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load character.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded character.html template")

	paymentTpl, err := loadGoTemplate(cfg.TemplatesDir, "payment.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load payment.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded payment.html template")

	loginTpl, err := loadGoTemplate(cfg.TemplatesDir, "login.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load login.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded login.html template")

	registerTpl, err := loadGoTemplate(cfg.TemplatesDir, "register.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load register.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded register.html template")

	guideTpl, err := loadGoTemplate(cfg.TemplatesDir, "guide.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load guide.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded guide.html template")

	mainTpl, err := loadGoTemplate(cfg.TemplatesDir, "main.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load main.html template: %v", err)
	}
	conversationTpl, err := loadGoTemplate(cfg.TemplatesDir, "conversation.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load conversation.html template: %v", err)
	}
	writingMainTpl, err := loadGoTemplate(cfg.TemplatesDir, "writing_main.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load writing_main.html template: %v", err)
	}
	newCharacterTpl, err := loadGoTemplate(cfg.TemplatesDir, "new_character.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load new_character.html template: %v", err)
	}

	loungeTpl, err := loadGoTemplate(cfg.TemplatesDir, "lounge.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load lounge.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded lounge.html template")

	tutorialTpl, err := loadGoTemplate(cfg.TemplatesDir, "tutorial.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load tutorial.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded tutorial.html template")

	quizTpl, err := loadGoTemplate(cfg.TemplatesDir, "quiz.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load quiz.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded quiz.html template")

	llmTpl, err := loadGoTemplate(cfg.TemplatesDir, "llm.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load llm.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded llm.html template")

	storyboardTpl, err := loadGoTemplate(cfg.TemplatesDir, "storyboard.html")
	if err != nil {
		log.Fatalf("CRITICAL: failed to load storyboard.html template: %v", err)
	}
	log.Printf("✓ Successfully loaded storyboard.html template")

	// Initialize auth middleware
	authMiddleware := NewAuthMiddleware(app.authService)

	// Root route - landing page (Chinese by default, English with /e suffix)
	r.Get("/", app.handleLandingMultiLang(landingTpl, landingEnTpl))

	// Landing page alias
	r.Get("/landing", app.handleLandingMultiLang(landingTpl, landingEnTpl))

	// Authentication routes (public)
	r.Get("/login", app.handleLogin(loginTpl))
	r.Post("/login", app.handleLogin(loginTpl))
	r.Post("/signup", app.handleSignup)
	r.Get("/logout", app.handleLogout)  // Support GET for href links
	r.Post("/logout", app.handleLogout) // Support POST for form submissions
	
	// Password recovery (public)
	r.Post("/api/recover-password", app.handleRecoverPassword)

	// Registration routes (public)
	r.Get("/register", app.handleRegisterPage(registerTpl))
	r.Post("/register", app.handleRegisterPage(registerTpl))
	r.Get("/confirm-email", app.handleConfirmEmail)
	r.Post("/confirm-email", app.handleConfirmEmail)

	r.Get("/main", app.handleTemplate(mainTpl))
	r.Get("/conversation", app.handleTemplate(conversationTpl))
	r.Post("/generate_dialogue_image", app.handleGenerateDialogueImage())

	// Public routes (accessible to all users including new users)
	r.Get("/faq", app.handleFAQ(faqTpl))
	r.Get("/descriptions", app.handleCharacter(characterTpl))
	r.Get("/character", app.handleCharacter(characterTpl)) // Backward compatibility
	r.Get("/guide", app.handleGuide(guideTpl))             // Writer's guide/onboarding
	r.Get("/tutorial", app.handleTutorial(tutorialTpl))    // Interactive Cao Cao tutorial
	r.Get("/quiz", app.handleTutorial(quizTpl))            // Character certification quiz

	// Character-specific routes (public - new users can view these)
	r.Get("/kurisu", app.handleTemplate(writingMainTpl))
	r.Get("/kafka", app.handleTemplate(writingMainTpl))
	r.Get("/linlu", app.handleTemplate(writingMainTpl))
	r.Get("/newcharacter_1", app.handleTemplate(newCharacterTpl))

	// Navigation page (character selection) - now viewable without login (UI will gate edits)
	r.Get("/navigation", app.handleWriting(writingTpl))
	r.Get("/writing", app.handleWriting(writingTpl)) // Backward compatibility redirect
	r.Get("/payment", authMiddleware.RequireWriterOrEditor(app.handlePayment(paymentTpl)))

	// Writers' Lounge (writers and editors only)
	r.Get("/lounge", authMiddleware.RequireWriterOrEditor(app.handleLounge(loungeTpl)))

	// LLM Chat page (public)
	r.Get("/llm", app.handleLLM(llmTpl))

	// Storyboard planning page (editors only)
	r.Get("/storyboard", authMiddleware.RequireWriterOrEditor(app.handleStoryboard(storyboardTpl)))

	// Static assets
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.StaticDir))))

	// Sticker files
	r.Handle("/static/stickers/*", http.StripPrefix("/static/stickers/", http.FileServer(http.Dir(filepath.Join(cfg.RootDir, "stickers")))))

	// Preset files passthrough
	r.Get("/presets/{preset}/{file:.+}", func(w http.ResponseWriter, r *http.Request) {
		preset := chi.URLParam(r, "preset")
		file := chi.URLParam(r, "file")
		target := filepath.Join(cfg.PresetBaseDir, preset, file)
		http.ServeFile(w, r, target)
	})

	// Full API parity implementation
	app.RegisterAPI(r)

	// Writers' Lounge API (legacy)
	app.RegisterLoungeAPI(r)

	// Discord-style Chat API
	app.RegisterChatAPI(r)

	// Tutorial and certification API
	app.RegisterTutorialAPI(r)

	// Payment and earnings API
	app.RegisterPaymentAPI(r)

	// LLM API
	app.RegisterLLMAPI(r)

	log.Printf("Template base: %s", cfg.TemplatesDir)
	log.Printf("Static base:   %s", cfg.StaticDir)
	log.Printf("Preset base:   %s", cfg.PresetBaseDir)
}
