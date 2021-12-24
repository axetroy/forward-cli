package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printHelp() {
	println(`forward - A command line tool to quickly setup a reverse proxy server.

USAGE:
  forward [OPTIONS] [host]

OPTIONS:
  --help              print help information
  --version           show version information
  --cors              enable cors. defaults: false
  --port=<int>        Specify the port that the proxy server listens on. defaults: 8080

EXAMPLE:
  forward http://example.com
  forward --port=80 http://example.com`)
}

type ProxyServer struct {
	target *url.URL
	cors   bool
	proxy  *httputil.ReverseProxy
}

type ProxyServerOptions struct {
	Target *url.URL
	Cors   bool
}

func NewProxyServer(options ProxyServerOptions) *ProxyServer {
	proxy := httputil.NewSingleHostReverseProxy(options.Target)

	server := &ProxyServer{
		cors:   options.Cors,
		target: options.Target,
		proxy:  proxy,
	}

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		server.modifyRequest(req)
	}

	proxy.ModifyResponse = server.modifyResponse
	proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("Got error while modifying response: %v \n", err)
		rw.Header().Set("X-Proxy-Error", err.Error())
	}

	return server
}

func (p *ProxyServer) Handler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		p.proxy.ServeHTTP(w, r)
	}
}

func (p *ProxyServer) modifyRequest(req *http.Request) {
	req.Header.Set("Host", p.target.Host)
	req.Header.Set("Referrer", fmt.Sprintf("%s://%s%s", p.target.Scheme, p.target.Host, req.URL.RawPath))
	req.Header.Set("X-Real-IP", req.RemoteAddr)
	req.Header.Set("X-Forwarded-For", fmt.Sprintf("%s://%s", p.target.Scheme, p.target.Host))
	req.Header.Set("X-Proxy", "Forward-Cli")
}

func (p *ProxyServer) modifyResponse(res *http.Response) error {
	res.Header.Set("X-Proxy", "Forward-Cli")
	return nil
}

func main() {
	var (
		showHelp    bool
		showVersion bool
		cors        bool
		port        string = "8080"
	)

	if len(os.Getenv("PORT")) > 0 {
		PORT_FROM_ENV := os.Getenv("PORT")

		if PORT_FROM_ENV != "" {
			port = PORT_FROM_ENV
		}
	}

	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.BoolVar(&cors, "cors", false, "")
	flag.StringVar(&port, "port", port, "")

	flag.Usage = printHelp

	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	if showVersion {
		println(fmt.Sprintf("%s %s %s", version, commit, date))
		return
	}

	server := flag.Arg(0)

	u, err := url.Parse(server)

	if err != nil {
		panic("invalid host")
	}

	target := fmt.Sprintf("%s://%s", u.Scheme, u.Host)

	proxy := NewProxyServer(ProxyServerOptions{
		Cors:   cors,
		Target: u,
	})

	http.HandleFunc("/", proxy.Handler())

	log.Printf("Proxy 'http://0.0.0.0:%s' to '%s'\n", port, target)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
