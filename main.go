package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/miekg/dns"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
)

// DB represents the Bolt DB instance
type DB struct {
	*bolt.DB
}

var (
	db           *DB
	dnsServers   []*dns.Server
	httpServer   *http.Server
	dnsClient    = &dns.Client{}
	blacklistKey = []byte("blacklist")
	version      = "undefined"
	build        = "undefined"

	dbPath     = flag.String("db", "nogo.db", "Specify a file path for the database.")
	dnsAddr    = flag.String("dns-addr", ":53", "Specify an address for the DNS proxy server to listen on.")
	dnsNet     = flag.String("dns-net", "udp", "Specify the listener protocol(s) for the DNS proxy server to use (\"udp\", \"tcp\", or \"udp+tcp\").")
	dnsProxyTo = flag.String("dns-proxyto", "8.8.8.8:53,8.8.4.4:53", "Specify one or more (comma separated) upstream DNS server addresses to proxy allowed queries to.")
	blacklist  = flag.String("import", "", "Specify a file path to import records to block (traditional hosts file format, or simply one domain per line).")
	webAddr    = flag.String("web-addr", ":8080", "Specify an address for the control panel web server to listen on.")
	webOff     = flag.Bool("web-off", false, "Instruct nogo not to serve the web control panel/API.")
	webPasswd  = flag.String("web-password", "", "Instruct the web control panel/API to require basic auth, using the specified password and a username of \"admin\".")
	showVer    = flag.Bool("version", false, "Show version and exit.")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "nogo version %s+%s %s/%s\n", version, build, runtime.GOOS, runtime.GOARCH)
		fmt.Fprintln(os.Stderr, "Copyright (c) 2017 Seth Davis")
		fmt.Fprintf(os.Stderr, "http://nogo.curia.solutions/\n\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var wg sync.WaitGroup

	flag.Parse()

	if *showVer {
		fmt.Printf("nogo version %s+%s %s/%s\n", version, build, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	// Initialize the database
	bdb, err := bolt.Open(*dbPath, 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		log.Fatalf("bolt.Open() Error: %s\n", err)
	}
	db = &DB{bdb}
	defer db.Close()

	// Ensure blacklist bucket exists
	if err = db.Update(func(tx *bolt.Tx) error {
		// Get/create bucket
		_, err := tx.CreateBucketIfNotExists(blacklistKey)
		return err
	}); err != nil {
		log.Fatalf("CreateBucketIfNotExists(%s) Error: %s\n", blacklistKey, err)
	}

	// Import a blacklist, if specified
	if *blacklist != "" {
		db.NoSync = true

		fmt.Println("Importing blacklist file. Please wait...")
		if err := db.importBlacklist(*blacklist); err != nil {
			log.Fatalf("db.importBlacklist(%s) Error: %s\n", *blacklist, err)
		}

		if err := db.Sync(); err != nil {
			log.Fatalf("db.Sync() Error: %s\n", err)
		}

		db.NoSync = false
	}

	// Initialize the HTTP router
	r := chi.NewRouter()

	// Register HTTP middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	if *webPasswd != "" {
		r.Use(basicAuth(*webPasswd))
	}

	// Register HTTP routes/handlers
	r.Get("/", rootIndexHandler)
	r.Post("/", rootCreateHandler)
	r.Get("/:key", rootReadHandler)
	r.Get("/api/", apiIndexHandler)
	r.Get("/api/:key", apiReadHandler)
	r.Put("/api/:key", apiPutHandler)
	r.Delete("/api/:key", apiDeleteHandler)
	r.Get("/export/hosts.txt", exportHostsHandler)
	r.Get("/css/nogo.css", cssHandler)

	// Initialize/start the servers
	log.Println("Booting up nogo...")

	for _, n := range strings.Split(*dnsNet, "+") {
		dnsServers = append(dnsServers, &dns.Server{Addr: *dnsAddr, Net: n, Handler: dns.HandlerFunc(dnsHandler)})

		wg.Add(1)
		go func(s *dns.Server) {
			defer wg.Done()
			if err := s.ListenAndServe(); err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
				log.Fatal(err)
			}
		}(dnsServers[len(dnsServers)-1])
	}
	log.Printf("DNS proxy listening at: %s (%s)\n", *dnsAddr, *dnsNet)

	if *webOff != true {
		httpServer = &http.Server{Addr: *webAddr, Handler: r}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()
		log.Printf("Web control panel/API listening at: %s\n", *webAddr)
	}

	// Attempt to gracefully shut down when signaled
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	s := <-sig
	fmt.Printf("Signal (%s) received, shutting down... ", s)

	if *webOff != true {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		httpServer.Shutdown(ctx)
	}
	for _, s := range dnsServers {
		s.Shutdown()
	}

	wg.Wait() // Wait on goroutines
	fmt.Println("Done!")
}
