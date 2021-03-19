package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/edubarbieri/proxy/config"
	"github.com/edubarbieri/proxy/middleware"
)

type AdminMiddleware struct {
	statsMiddleware *middleware.StatsMiddleware
	configManager   *config.Manager
}

func NewAdminMiddleware(statsMiddleware *middleware.StatsMiddleware, configManager *config.Manager) AdminMiddleware {
	return AdminMiddleware{
		statsMiddleware: statsMiddleware,
		configManager:   configManager,
	}
}

func (adm *AdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/proxy-admin/stats") {
			if req.Method == "DELETE" {
				adm.handlerResetStats(res)
				return
			}
			adm.handlerStats(res)
			return
		}
		if strings.HasPrefix(req.RequestURI, "/proxy-admin/config") {
			if req.Method == "PUT" {
				adm.handlerUpdateConfig(res, req)
				return
			}
			adm.handlerReadConfig(res)
			return
		}

		next.ServeHTTP(res, req)

	})
}

func (adm *AdminMiddleware) handlerStats(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	stats := adm.statsMiddleware.GenerateStats()
	b, _ := json.Marshal(stats)
	w.Write(b)
}
func (adm *AdminMiddleware) handlerResetStats(w http.ResponseWriter) {
	stats := adm.statsMiddleware.GenerateStatsAndReset()
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(stats)
	w.Write(b)
}
func (adm *AdminMiddleware) handlerUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var configReq config.Config
	err := json.NewDecoder(r.Body).Decode(&configReq)
	if err != nil {
		log.Printf("error parsing json config %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	adm.configManager.UpdateConfig(configReq)
	adm.handlerReadConfig(w)

}
func (adm *AdminMiddleware) handlerReadConfig(w http.ResponseWriter) {
	configError := adm.configManager.GetCurrentConfig()
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(configError)
	w.Write(b)
}
