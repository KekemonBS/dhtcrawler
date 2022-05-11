package crawler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/KekemonBS/dhtcrawler/crawler/models"
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

//DHTCrawler struct that holds general info needed to crawl DHT
type DHTCrawler struct {
	ctx     context.Context
	logger  *log.Logger
	threads int
	dbImpl  DbImpl
}

//DbImpl interface describes needed database operations
type DbImpl interface {
	Create(ctx context.Context, rate models.Share) error
	DeleteByID(ctx context.Context, id string) error
	ReadByID(ctx context.Context, id string) (models.Share, error)
	ReadAll(ctx context.Context) ([]models.Share, error)
	ReadPage(ctx context.Context, size, n int) ([]models.Share, error)
}

//New returs crawler instance
func New(ctx context.Context, logger *log.Logger, dbImpl DbImpl, threads int) (*DHTCrawler, error) {
	return &DHTCrawler{
		ctx:     ctx,
		logger:  logger,
		threads: threads,
		dbImpl:  dbImpl,
	}, nil
}

//Start starts crawler
func (dc *DHTCrawler) Start(w http.ResponseWriter, r *http.Request) {
	resChan := make(chan models.Share)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(dc.threads)
	//for i, v := range dc.currencies {
	//go ParseRate(dc.ctx, resChan, errChan, wg, v, time.Second*time.Duration(i), dc.dbImpl)
	//}

	go func() {
		wg.Wait()
		close(resChan)
		close(errChan)
	}()

	for {
		select {
		case val := <-errChan:
			{
				dc.logger.Println(val)
				return
			}
		case val, ok := <-resChan:
			{
				if !ok {
					return
				}

				n, err := w.Write([]byte(fmt.Sprint(val) + "\n"))
				if err != nil || n == 0 {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}
	}
}

//Status checks crawler progress
func (dc *DHTCrawler) Status(w http.ResponseWriter, r *http.Request) {
	n, err := w.Write([]byte("status\n"))
	if err != nil || n == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//Stop stops crawler
func (dc *DHTCrawler) Stop(w http.ResponseWriter, r *http.Request) {
	n, err := w.Write([]byte("stop\n"))
	if err != nil || n == 0 {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//SharesByID calls to db and writes specific share of shares table
func (dc *DHTCrawler) SharesByID(w http.ResponseWriter, r *http.Request) {
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
func (dc *DHTCrawler) SharesAll(w http.ResponseWriter, r *http.Request) {
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
func (dc *DHTCrawler) SharesPage(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	size, ok := query["size"]
	if !ok {
		dc.logger.Print("wrong query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageNum, ok := query["n"]
	if !ok {
		dc.logger.Print("wrong query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sizeConv, err := strconv.ParseInt(size[0], 10, 32)
	if err != nil {
		dc.logger.Print("malformed query param")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pageNumConv, err := strconv.ParseInt(pageNum[0], 10, 32)
	if err != nil {
		dc.logger.Print("malformed query param")
		w.WriteHeader(http.StatusBadRequest)
		return

	}
	rates, err := dc.dbImpl.ReadPage(dc.ctx, int(sizeConv), int(pageNumConv))
	if err != nil {
		dc.logger.Fatal(err)
	}
	_, err = w.Write([]byte(fmt.Sprint(rates) + "\n"))
	if err != nil {
		dc.logger.Fatal(err)
	}
}
