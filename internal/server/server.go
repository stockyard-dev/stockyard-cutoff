package server

import (
	"encoding/json"
	"net/http"

	"github.com/stockyard-dev/stockyard-cutoff/internal/store"
)

type Server struct{ db *store.DB; mux *http.ServeMux; limits Limits }

func New(db *store.DB, limits Limits) *Server {
	s := &Server{db: db, mux: http.NewServeMux(), limits: limits}

	s.mux.HandleFunc("GET /api/rules", s.listRules)
	s.mux.HandleFunc("POST /api/rules", s.createRule)
	s.mux.HandleFunc("GET /api/rules/{id}", s.getRule)
	s.mux.HandleFunc("DELETE /api/rules/{id}", s.deleteRule)
	s.mux.HandleFunc("PATCH /api/rules/{id}/toggle", s.toggleRule)
	s.mux.HandleFunc("POST /api/check", s.check)
	s.mux.HandleFunc("GET /api/stats", s.stats)
	s.mux.HandleFunc("GET /api/health", s.health)
	s.mux.HandleFunc("GET /api/tier", func(w http.ResponseWriter, r *http.Request) {
		wj(w, 200, map[string]any{"tier": s.limits.Tier, "upgrade_url": "https://stockyard.dev/cutoff/"})
	})
	s.mux.HandleFunc("GET /ui", s.dashboard)
	s.mux.HandleFunc("GET /ui/", s.dashboard)
	s.mux.HandleFunc("GET /", s.root)
	return s
}
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) { s.mux.ServeHTTP(w, r) }
func wj(w http.ResponseWriter, c int, v any) { w.Header().Set("Content-Type", "application/json"); w.WriteHeader(c); json.NewEncoder(w).Encode(v) }
func we(w http.ResponseWriter, c int, m string) { wj(w, c, map[string]string{"error": m}) }
func (s *Server) root(w http.ResponseWriter, r *http.Request) { if r.URL.Path != "/" { http.NotFound(w, r); return }; http.Redirect(w, r, "/ui", 302) }
func or(r []store.Rule) []store.Rule { if r == nil { return []store.Rule{} }; return r }
func (s *Server) listRules(w http.ResponseWriter, r *http.Request) { wj(w, 200, map[string]any{"rules": or(s.db.ListRules())}) }
func (s *Server) createRule(w http.ResponseWriter, r *http.Request) {
	if s.limits.MaxItems > 0 && len(s.db.ListRules()) >= s.limits.MaxItems { we(w, 402, "Free tier limit reached"); return }
	var rule store.Rule; json.NewDecoder(r.Body).Decode(&rule); if rule.Name == "" { we(w, 400, "name required"); return }
	rule.Enabled = true; s.db.CreateRule(&rule); wj(w, 201, s.db.GetRule(rule.ID))
}
func (s *Server) getRule(w http.ResponseWriter, r *http.Request) { rule := s.db.GetRule(r.PathValue("id")); if rule == nil { we(w, 404, "not found"); return }; wj(w, 200, rule) }
func (s *Server) deleteRule(w http.ResponseWriter, r *http.Request) { s.db.DeleteRule(r.PathValue("id")); wj(w, 200, map[string]string{"status": "deleted"}) }
func (s *Server) toggleRule(w http.ResponseWriter, r *http.Request) { s.db.ToggleRule(r.PathValue("id")); wj(w, 200, s.db.GetRule(r.PathValue("id"))) }
func (s *Server) check(w http.ResponseWriter, r *http.Request) {
	var body struct { RuleID string `json:"rule_id"`; Key string `json:"key"` }
	json.NewDecoder(r.Body).Decode(&body)
	if body.RuleID == "" { we(w, 400, "rule_id required"); return }
	result := s.db.Check(body.RuleID, body.Key)
	status := 200; if !result.Allowed { status = 429 }
	wj(w, status, result)
}
func (s *Server) stats(w http.ResponseWriter, r *http.Request) { wj(w, 200, s.db.Stats()) }
func (s *Server) health(w http.ResponseWriter, r *http.Request) { st := s.db.Stats(); wj(w, 200, map[string]any{"service": "cutoff", "status": "ok", "rules": st.Rules}) }
