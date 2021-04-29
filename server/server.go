package server

import (
	"github.com/desotech-it/whoami/app"
	"github.com/desotech-it/whoami/server/util"
	"github.com/desotech-it/whoami/view"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	v := view.IndexView{
		Title:   "WhoAmI",
		Info:    app.Info,
		Request: util.GetRequestAsString(r),
	}
	if util.IsFromCurl(r) {
		v.WriteAsText(w)
	} else {
		v.Write(w)
	}
}

func cpustressHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)

	switch r.Method {
	case "GET":
		if app.IsOSDarwin() {
			v := view.UnsupportedView{
				Title:   "Unsupported Operation",
				Message: "Retrieving the CPU info on MacOS is not currently implemented.",
				Link:    "https://github.com/shirou/gopsutil/issues/1000",
			}
			v.Write(w)
			return
		}
		stats := app.CPUInfo()
		v := view.CPUStressView{
			Title: "CPU Load",
			Stats: stats,
		}
		v.Write(w)
	case "POST":
		r.ParseForm()
		form := r.Form
		magnitude := form.Get("magnitude")
		unit := form.Get("unit")
		duration, err := time.ParseDuration(magnitude + unit)
		if err == nil {
			go util.GenerateCPULoadFor(duration)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func memstressHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)

	switch r.Method {
	case "GET":
		stats := app.MemInfo()
		v := view.MemStressView{
			Title: "Memory Usage",
			Stats: stats,
		}
		v.Write(w)
	case "POST":
		r.ParseForm()
		form := r.Form
		magnitude := form.Get("magnitude")
		unit := form.Get("unit")
		duration, err := time.ParseDuration(magnitude + unit)
		if err == nil {
			go util.GenerateHighMemoryUsageFor(duration)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	health := app.GetHealth()
	var status string
	if health == 200 {
		status = "UP"
	} else {
		status = "ERRORED"
	}
	v := view.HealthView{
		Status: status,
	}
	w.WriteHeader(health)
	v.Write(w)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	readiness := app.GetReadiness()
	var status string
	if readiness == 200 {
		status = "UP"
	} else {
		status = "STARTING"
	}
	v := view.ReadinessView{
		Status: status,
	}
	w.WriteHeader(readiness)
	v.Write(w)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	promhttp.Handler().ServeHTTP(w, r)
}

func goldieHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	v := view.ImageView{
		Title:    "Goldie",
		Filename: "goldie.png",
		AltText:  "goldie",
		Info:     app.Info,
		Request:  util.GetRequestAsString(r),
	}
	if util.IsFromCurl(r) {
		v.WriteAsPlainText(w)
	} else {
		v.Write(w)
	}
}

func zeeHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	v := view.ImageView{
		Title:    "Zee",
		Filename: "zee.png",
		AltText:  "zee",
		Info:     app.Info,
		Request:  util.GetRequestAsString(r),
	}
	if util.IsFromCurl(r) {
		v.WriteAsPlainText(w)
	} else {
		v.Write(w)
	}
}

func captainkubeHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	v := view.ImageView{
		Title:    "Captain Kube",
		Filename: "captain-kube.png",
		AltText:  "captain kube",
		Info:     app.Info,
		Request:  util.GetRequestAsString(r),
	}
	if util.IsFromCurl(r) {
		v.WriteAsPlainText(w)
	} else {
		v.Write(w)
	}
}

func phippyHandler(w http.ResponseWriter, r *http.Request) {
	app.LogRequest(r)
	v := view.ImageView{
		Title:    "Phippy",
		Filename: "phippy.png",
		AltText:  "phippy",
		Info:     app.Info,
		Request:  util.GetRequestAsString(r),
	}
	if util.IsFromCurl(r) {
		v.WriteAsPlainText(w)
	} else {
		v.Write(w)
	}
}

func cpuusageHandler(w http.ResponseWriter, _ *http.Request) {
	cpusUsage := app.CPULoad()
	util.WriteJSONResponse(w, cpusUsage)
}

func memusageHandler(w http.ResponseWriter, _ *http.Request) {
	memUsage := app.MemInfo()
	util.WriteJSONResponse(w, memUsage)
}

type Server struct {
	Port uint64
}

func assignRootHandler() {
	targetImage := os.Getenv("NAME_APPLICATION")
	var targetHandler func(http.ResponseWriter, *http.Request)

	switch targetImage {
	case "goldie":
		targetHandler = goldieHandler
	case "zee":
		targetHandler = zeeHandler
	case "captainkube":
		targetHandler = captainkubeHandler
	case "phippy":
		targetHandler = phippyHandler
	default:
		targetHandler = rootHandler
	}

	http.HandleFunc("/", targetHandler)
}

func (s *Server) Start() {
	assignRootHandler()
	http.HandleFunc("/cpustress", cpustressHandler)
	http.HandleFunc("/memstress", memstressHandler)
	http.HandleFunc("/readiness", readinessHandler)
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.HandleFunc("/goldie", goldieHandler)
	http.HandleFunc("/zee", zeeHandler)
	http.HandleFunc("/captainkube", captainkubeHandler)
	http.HandleFunc("/phippy", phippyHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/cpuusage", cpuusageHandler)
	http.HandleFunc("/memusage", memusageHandler)

	address := fmt.Sprintf(":%d", s.Port)
	log.Fatal(http.ListenAndServe(address, nil))
}
