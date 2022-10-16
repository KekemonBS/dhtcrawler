package crawler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/KekemonBS/dhtcrawler/crawler/models"
	"github.com/KekemonBS/dhtcrawler/infrastructure/env"
	"github.com/KekemonBS/dhtcrawler/storage/csvimport"

	"github.com/google/uuid"
)

////go:generate cp -r ../assets/public/ served/
//var (
//	//go:embed served
//	frontend embed.FS
//	pages    = map[string]string{
//		"/": "served/public/index.html",
//	}
//)

//DHTCrawlerHandlers struct that holds general info needed to crawl DHT
type DHTCrawlerHandlers struct {
	ctx    context.Context
	logger *log.Logger
	dbImpl DbImpl
	cfg    *env.Config
}

//DbImpl interface describes needed database operations
type DbImpl interface {
	Create(ctx context.Context, rate models.Share) error
	DeleteByID(ctx context.Context, id string) error
	ReadByID(ctx context.Context, id string) (models.Share, error)
	ReadAll(ctx context.Context) ([]models.Share, error)
	ReadPage(ctx context.Context, size, n int, queryShares string) (models.SharesPage, error)
}

//New returs crawler instance
func New(ctx context.Context, logger *log.Logger, dbImpl DbImpl, cfg *env.Config) (*DHTCrawlerHandlers, error) {
	return &DHTCrawlerHandlers{
		ctx:    ctx,
		logger: logger,
		dbImpl: dbImpl,
		cfg:    cfg,
	}, nil
}

//Start starts crawler
func (dc *DHTCrawlerHandlers) Start(w http.ResponseWriter, r *http.Request) {
	sesChan := make(chan struct{})

	dhtcses, err := NewCrawler(dc.ctx, dc.dbImpl, dc.logger)
	if err != nil {
		dc.logger.Fatal(err)
	}
	go dhtcses.Crawl(sesChan)

	<-sesChan
	if dc.cfg.ImportCSV {
		err := csvimport.ImportCSV("./storage/csvimport/mine", dhtcses)
		if err != nil {
			dc.logger.Fatal(err)
		}
	}
}

//Stop stops crawler
func (dc *DHTCrawlerHandlers) Stop(w http.ResponseWriter, r *http.Request) {
	n, err := w.Write([]byte("stop\n"))
	if err != nil || n == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//SharesByID calls to db and writes specific share of shares table
func (dc *DHTCrawlerHandlers) SharesByID(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id, ok := query["id"]
	if !ok {
		dc.logger.Print("wrong query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uuidString, err := uuid.Parse(id[0])
	if err != nil {
		dc.logger.Fatal(err)
	}
	rate, err := dc.dbImpl.ReadByID(dc.ctx, uuidString.String())
	if err != nil {
		dc.logger.Fatal(err)
	}
	_, err = w.Write([]byte(fmt.Sprint(rate) + "\n"))
	if err != nil {
		dc.logger.Fatal(err)
	}
}

//SharesAll calls to db and writes all shares of shares table
func (dc *DHTCrawlerHandlers) SharesAll(w http.ResponseWriter, r *http.Request) {
	rates, err := dc.dbImpl.ReadAll(dc.ctx)
	if err != nil {
		dc.logger.Fatal(err)
	}
	_, err = w.Write([]byte(fmt.Sprint(rates) + "\n"))
	if err != nil {
		dc.logger.Fatal(err)
	}

}

//SharesPage calls to db and writes page of shares table
func (dc *DHTCrawlerHandlers) SharesPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	size, ok := query["size"]
	if !ok {
		dc.logger.Print("wrong query param: size")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageNum, ok := query["n"]
	if !ok {
		dc.logger.Print("wrong query param: n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	queryShares, ok := query["query"]
	if !ok {
		dc.logger.Print("wrong query param: query")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sizeConv, err := strconv.ParseInt(size[0], 10, 32)
	if err != nil {
		dc.logger.Print("malformed query param: size")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pageNumConv, err := strconv.ParseInt(pageNum[0], 10, 32)
	if err != nil {
		dc.logger.Print("malformed query param: n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sharesPage, err := dc.dbImpl.ReadPage(dc.ctx, int(sizeConv), int(pageNumConv), queryShares[0])
	if err != nil {
		dc.logger.Fatal(err)
	}

	res, err := json.Marshal(sharesPage)
	_, err = w.Write(res)
	if err != nil {
		dc.logger.Fatal(err)
	}
}
