package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	forward "github.com/axetroy/forward-cli"
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
  --help                              print help information
  --version                           show version information
  --port="<int>"                      specify the port that the proxy server listens on. defaults: 80
  --req-header="key=value"            specify the request header attached to the request. defaults: ""
  --res-header="key=value"            specify the response headers. defaults: ""
  --cors                              enable cors. defaults: false
  --cors-allow-headers="<string>"     allow send headers from client when cors enabled. defaults: ""
  --cors-expose-headers="<string>"    expose response headers from server when cors enabled. defaults: ""

EXAMPLES:
  forward http://example.com
  forward --port=80 http://example.com
  forward --req-header="foo=bar" http://example.com
  forward --cors --cors-allow-headers="Auth, Token" http://example.com`)
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "custom array flag"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var (
		showHelp             bool
		showVersion          bool
		cors                 bool
		corsAllowHeaders     string
		corsExposeHeaders    string
		requestHeadersArray  arrayFlags
		responseHeadersArray arrayFlags
		port                 string = "80"
	)

	if len(os.Getenv("PORT")) > 0 {
		PORT_FROM_ENV := os.Getenv("PORT")

		if PORT_FROM_ENV != "" {
			port = PORT_FROM_ENV
		}
	}

	flag.BoolVar(&showHelp, "help", false, "")
	flag.BoolVar(&showVersion, "version", false, "")
	flag.Var(&requestHeadersArray, "req-header", "")
	flag.Var(&responseHeadersArray, "res-header", "")
	flag.BoolVar(&cors, "cors", false, "")
	flag.StringVar(&corsAllowHeaders, "cors-allow-headers", corsAllowHeaders, "")
	flag.StringVar(&corsExposeHeaders, "cors-expose-headers", corsExposeHeaders, "")
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

	if server == "" {
		fmt.Printf("ERR: proxy server is required\n\n")
		printHelp()
		os.Exit(1)
	}

	u, err := url.Parse(server)

	if err != nil {
		panic("invalid host")
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		panic("invalid proxy target")
	}

	target := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	requestHeaders := http.Header{}
	responseHeaders := http.Header{}

	for _, paren := range requestHeadersArray {
		arr := strings.Split(paren, "=")
		requestHeaders.Set(arr[0], strings.Join(arr[1:], "="))
	}

	for _, paren := range responseHeadersArray {
		arr := strings.Split(paren, "=")
		responseHeaders.Set(arr[0], strings.Join(arr[1:], "="))
	}

	proxy := forward.NewProxyServer(forward.ProxyServerOptions{
		ReqHeaders:        requestHeaders,
		ResHeaders:        responseHeaders,
		Cors:              cors,
		CorsAllowHeaders:  corsAllowHeaders,
		CorsExposeHeaders: corsExposeHeaders,
		Target:            u,
	})

	http.HandleFunc("/", proxy.Handler())

	ip, err := forward.GetLocalIP()

	if err != nil {
		log.Panicln(err)
	}

	log.Printf("Proxy 'http://%s:%s' to '%s'\n", ip, port, target)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
