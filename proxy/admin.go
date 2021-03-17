package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"proxy/middleware"
	"strings"
)

type AdminMiddleware struct {
	statsMiddleware *middleware.StatsMiddleware
}

func NewAdminMiddleware(statsMiddleware *middleware.StatsMiddleware) AdminMiddleware {
	return AdminMiddleware{
		statsMiddleware: statsMiddleware,
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

func (md *AdminMiddleware) handlerStats(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	stats := md.statsMiddleware.GenerateStats()
	b, _ := json.Marshal(stats)
	w.Write(b)
}
func (md *AdminMiddleware) handlerResetStats(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	stats := md.statsMiddleware.GenerateStatsAndReset()
	b, _ := json.Marshal(stats)
	w.Write(b)
}
func (md *AdminMiddleware) handlerUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var config Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		log.Printf("error parsing json config %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updateConfig(config)
	file, _ := json.MarshalIndent(config, "", " ")
	ioutil.WriteFile(os.Getenv("CONFIG_PATH"), file, 0644)

	md.handlerReadConfig(w)

}
func (md *AdminMiddleware) handlerReadConfig(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	config := ReadConfigJson()
	b, _ := json.Marshal(config)
	w.Write(b)
}
