package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/edubarbieri/proxy/config"
	"github.com/edubarbieri/proxy/middleware"
)

type AdminMiddleware struct {
	statsMiddleware *middleware.StatsMiddleware
	configManager   *config.ConfigManager
}

func NewAdminMiddleware(statsMiddleware *middleware.StatsMiddleware, configManager *config.ConfigManager) AdminMiddleware {
	return AdminMiddleware{
		statsMiddleware: statsMiddleware,
		configManager:   configManager,
	}
}

func (md *AdminMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.RequestURI, "/proxy-admin/stats") {
			if req.Method == "DELETE" {
				md.handlerResetStats(res)
				return
			}
			md.handlerStats(res)
			return
		}
		if strings.HasPrefix(req.RequestURI, "/proxy-admin/config") {
			if req.Method == "PUT" {
				md.handlerUpdateConfig(res, req)
				return
			}
			md.handlerReadConfig(res)
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
	var config config.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		log.Printf("error parsing json config %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	adm.configManager.UpdateConfig(config)
	file, _ := json.MarshalIndent(config, "", " ")
	ioutil.WriteFile(os.Getenv("CONFIG_PATH"), file, 0644)

	adm.handlerReadConfig(w)

}
func (adm *AdminMiddleware) handlerReadConfig(w http.ResponseWriter) {
	config := adm.configManager.GetCurrentConfig()
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(config)
	w.Write(b)
}
