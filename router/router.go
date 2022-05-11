package router

import (
	"net/http"
)

type rateParserImpl interface {
	Start(w http.ResponseWriter, r *http.Request)
	Status(w http.ResponseWriter, r *http.Request)
	Stop(w http.ResponseWriter, r *http.Request)
	SharesByID(w http.ResponseWriter, r *http.Request)
	SharesAll(w http.ResponseWriter, r *http.Request)
	SharesPage(w http.ResponseWriter, r *http.Request)
	//	ServeFrontend(w http.ResponseWriter, r *http.Request)
}

//New returns router
func New(imp rateParserImpl) *http.ServeMux {
	//API
	mux := http.NewServeMux()
	mux.HandleFunc("/dhtcrawler/start", imp.Start)
	mux.HandleFunc("/dhtcrawler/status", imp.Status)
	mux.HandleFunc("/dhtcrawler/stop", imp.Stop)

	mux.HandleFunc("/dhtcrawler/display", imp.SharesByID)
	mux.HandleFunc("/dhtcrawler/displayall", imp.SharesAll)
	mux.HandleFunc("/dhtcrawler/displaypage", imp.SharesPage)

	prefixMux := http.NewServeMux()
	prefixMux.Handle("/api/v1/", http.StripPrefix("/api/v1", mux))

	//static
	static := http.FileServer(http.Dir("assets/public/"))
	prefixMux.Handle("/", static)

	return prefixMux
}
