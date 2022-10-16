package crawler

import (
	"context"
	"encoding/hex"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anacrolix/dht/v2"
	"github.com/anacrolix/torrent/bencode"
	"github.com/davecgh/go-spew/spew"

	"github.com/cenkalti/rain/torrent"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	"github.com/KekemonBS/dhtcrawler/crawler/models"
)

//DHTCrawler struct holds data that needed in process of capturing info_hashes
type DHTCrawler struct {
	ctx              context.Context
	logger           *log.Logger
	startupWaitgroup sync.WaitGroup
	SessionStorage   []string //will be modified by func that imports infohashes from csv
	dbImpl           DbImpl
}

//NewCrawler creates instance of DHTCrawler
//that will listen to DHT and process incoming get_peers queries
func NewCrawler(ctx context.Context, db DbImpl, logger *log.Logger) (*DHTCrawler, error) {
	return &DHTCrawler{
		ctx:              ctx,
		logger:           logger,
		startupWaitgroup: sync.WaitGroup{},
		SessionStorage:   make([]string, 10),
		dbImpl:           db,
	}, nil

}

//Crawl starts DHT node, captures get_peers queries
//, tries to search captured info_hash
func (dhtc *DHTCrawler) Crawl(sch chan struct{}) {

	ip, ifs := dhtc.getIP()

	dhtc.logger.Print(ip, " ", ifs)

	sc := dhtc.genConfig(ip, ":1337")
	ds, err := dht.NewServer(sc)
	if err != nil {
		dhtc.logger.Fatal(err)
	}

	go dhtc.readConn(ifs, ip)
	go ds.TableMaintainer()

	ticker := time.NewTicker(1 * time.Second)
	errChan := make(chan error)

	dhtc.startupWaitgroup.Add(1)
	var ses *torrent.Session
	go dhtc.createSession(&ses, ip, errChan)
	dhtc.startupWaitgroup.Wait()

	sch <- struct{}{}

	//addIH("ec7a402ff515d80f30f6244847b672ae9fbe5d7a")
	//addIH("74d89962b39acb3fc855d4b91c4eb234758a4ebc")

	for {
		select {
		case <-dhtc.ctx.Done():
			ds.Close()
			return
		case <-ticker.C:
			sl := len(dhtc.SessionStorage)
			var sfe string
			if sl != 0 {
				sfe = dhtc.SessionStorage[0]
				//spew.Dump(dhtc.SessionStorage)
				errchan := make(chan error, 0)

				go dhtc.addURI(ses, sfe, 30*time.Second, errchan, 4)
				dhtc.SessionStorage = dhtc.SessionStorage[1:]
			}
		case err = <-errChan:
			dhtc.logger.Fatal(err)
		}
	}
}

func (dhtc *DHTCrawler) genConfig(ip, port string) *dht.ServerConfig {
	sc := dht.NewDefaultServerConfig()
	conn, err := net.ListenPacket("udp", ip+port)
	if err != nil {
		panic(err)
	}
	sc.Conn = conn
	return sc
}

func (dhtc *DHTCrawler) getIP() (ip string, interfaceName string) {
	//Find IP
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		dhtc.logger.Fatal(err)
	}
	defer conn.Close()
	ipAddress := conn.LocalAddr().(*net.UDPAddr)
	resIP := ipAddress.IP.String()

	//Find interface
	ifs, err := net.Interfaces()
	var found net.Interface
	for _, v := range ifs {
		addrList, err := v.Addrs()
		if err != nil {
			dhtc.logger.Fatal(err)
		}
		for _, addr := range addrList {
			if strings.Split(addr.String(), "/")[0] == resIP {
				found = v
			}
		}
	}
	return resIP, found.Name
}

//Sniff right packets from interface
func (dhtc *DHTCrawler) readConn(ifs string, ip string) {
	//62144 default max packet size for tcpdump
	handle, err := pcap.OpenLive(ifs, 62144, true, pcap.BlockForever)
	if err != nil {
		dhtc.logger.Fatal(err)
	}
	defer handle.Close()
	pchan := gopacket.NewPacketSource(handle, handle.LinkType()).Packets()
	for {
		select {
		case <-dhtc.ctx.Done():
			handle.Close()
			return
		case packet := <-pchan:
			//Process packet
			var res interface{}
			if packet.ErrorLayer() == nil &&
				strings.Contains(string(packet.Data()), "get_peers") &&
				packet.NetworkLayer().NetworkFlow().Src().String() != ip {

				err = bencode.Unmarshal(packet.ApplicationLayer().LayerContents(), &res)
				if err != nil {
					dhtc.logger.Print(err)
				}
				if res != nil {
					vala, ok := res.(map[string]interface{})["a"]
					if !ok {
						dhtc.logger.Print("malformed packet")
					}
					if valHash, ok := vala.(map[string]interface{})["info_hash"]; ok {
						infoHash := valHash.(string)
						infoHashString := hex.EncodeToString([]byte(infoHash))
						if !dhtc.contains(dhtc.SessionStorage, infoHashString) {
							dhtc.SessionStorage = append(dhtc.SessionStorage, infoHashString)
							dhtc.logger.Println(infoHashString)
						}
					}
				}
			}
		}
	}
}

//Taken from main.go of cenkalti/rain torrent client, modified
func (dhtc *DHTCrawler) createSession(resSession **torrent.Session, ip string, ech chan error) {
	cfg := torrent.DefaultConfig
	cfg.DHTHost = ip
	cfg.Database = "session.db"

	rs, err := torrent.NewSession(cfg)
	*resSession = rs

	if err != nil {
		ech <- err
	}
	defer rs.Close()
	dhtc.startupWaitgroup.Done()
	for {
		select {
		case s := <-dhtc.ctx.Done():
			dhtc.logger.Printf("received %s, stopping server", s)
			os.Remove(cfg.Database)
			return
		}
	}
}

func (dhtc *DHTCrawler) addURI(ses *torrent.Session, ih string, timeout time.Duration, ech chan error, retries int) {
	dhtc.logger.Print("looking up ----->", ih)
	ticker := time.NewTicker(timeout)
	opt := &torrent.AddTorrentOptions{
		StopAfterMetadata: true,
	}

	arg := dhtc.genMagnetLink(ih)
	t, err := ses.AddURI(arg, opt)
	if err != nil {
		ech <- err
	}

	metadataC := t.NotifyMetadata()

	for {
		select {
		case <-dhtc.ctx.Done():
			err = t.Stop()
			if err != nil {
				ech <- err
			}
			err := ses.RemoveTorrent(t.ID())
			if err != nil {
				ech <- err
			}
			return
		case <-metadataC:
			tbytes, _ := t.Torrent()
			var decodedTorrentFile interface{}
			bencode.Unmarshal(tbytes, &decodedTorrentFile)
			dhtc.writeResult(decodedTorrentFile.(map[string]interface{}), arg)
			err := ses.RemoveTorrent(t.ID())
			if err != nil {
				ech <- err
			}
			return
		case <-ticker.C:
			err := ses.RemoveTorrent(t.ID())
			if err != nil {
				ech <- err
			}
			dhtc.logger.Print("-----> metadata cannot be downloaded")

			//Recursively call itself, retries += -1
			if retries != 0 {
				dhtc.addURI(ses, ih, timeout, ech, retries-1)
			}
			return
		}
	}
}

func (dhtc *DHTCrawler) genMagnetLink(ih string) string {
	base := "magnet:?xt=urn:btih:"
	tr1 := "&tr=udp://tracker.openbittorrent.com:80"
	tr2 := "&tr=udp://tracker.opentrackr.org:1337/announce"
	return base + ih + tr1 + tr2
}

func (dhtc *DHTCrawler) contains(s []string, v string) bool {
	for _, vs := range s {
		if v == vs {
			return true
		}
	}
	return false
}

//AddIH adds info_hash to session for further lookup
func (dhtc *DHTCrawler) AddIH(ih string) {
	dhtc.SessionStorage = append(dhtc.SessionStorage, ih)
}

func (dhtc *DHTCrawler) writeResult(m map[string]interface{}, magnet string) {
	var res models.Share
	res.Name = m["info"].(map[string]interface{})["name"].(string)
	res.MagnetLink = magnet

	res.Size = 0
	res.FileTree = ""

	if m["info"].(map[string]interface{})["files"] != nil {
		files := m["info"].(map[string]interface{})["files"].([]interface{})
		for _, v := range files {
			res.Size += int(v.(map[string]interface{})["length"].(int64))

			pathMap := v.(map[string]interface{})["path"].([]interface{})
			path := "."
			for _, pv := range pathMap {
				path += "/" + pv.(string)
			}
			path += "\n"
			res.FileTree += path
		}
	} else {
		res.Size = int(m["info"].(map[string]interface{})["length"].(int64))
		res.FileTree = "./" + res.Name
	}

	//Write to db
	spew.Dump(res)
	dhtc.dbImpl.Create(dhtc.ctx, res)
}
