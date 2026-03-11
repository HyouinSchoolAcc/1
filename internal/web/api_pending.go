package web

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// registerAPIPending stubs all JSON endpoints with a clear TODO response.
// This keeps the route surface identical while we fill implementations.
func registerAPIPending(r chi.Router) {
	stub := func(name string) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotImplemented)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error": "not implemented yet in Go port",
				"route": name,
			})
		}
	}

	// Mirror Flask routes
	r.Get("/load_structured_writer_files", stub("load_structured_writer_files"))
	r.Get("/get_categories", stub("get_categories"))
	r.Post("/load_writer_file_content", stub("load_writer_file_content"))
	r.Post("/save_writer_file_content", stub("save_writer_file_content"))
	r.Post("/update_qc_status", stub("update_qc_status"))
	r.Post("/update_character_turn", stub("update_character_turn"))
	r.Post("/update_kurisu_turn", stub("update_kurisu_turn"))
	r.Post("/save_conversation_data", stub("save_conversation_data"))
	r.Post("/create_new_version", stub("create_new_version"))
	r.Post("/delete_version", stub("delete_version"))
	r.Post("/add_character", stub("add_character"))
	r.Post("/delete_character", stub("delete_character"))
	r.Post("/get_deletable_characters", stub("get_deletable_characters"))
	r.Post("/update_real_name", stub("update_real_name"))
	r.Post("/add_day", stub("add_day"))
	r.Post("/delete_day", stub("delete_day"))
	r.Post("/update_schedule", stub("update_schedule"))
	r.Post("/update_day_category", stub("update_day_category"))
	r.Post("/move_to_legacy", stub("move_to_legacy"))
	r.Post("/auto_archive_old_items", stub("auto_archive_old_items"))
	r.Post("/update_inner_thought_annotation", stub("update_inner_thought_annotation"))
	r.Post("/delete_inner_thought_annotation", stub("delete_inner_thought_annotation"))
	r.Post("/save_checklist_data", stub("save_checklist_data"))
	r.Post("/load_checklist_data", stub("load_checklist_data"))
}
