package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	forward "github.com/axetroy/forward-cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return []byte("0.0.0.0")
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func printHelp() {
	println(`forward - A command line tool to quickly setup a reverse proxy server.

USAGE:
  forward [OPTIONS] [host]

OPTIONS:
  --help                              print help information
  --version                           show version information
  --address="<string>"                specify the address that the proxy server listens on. defaults: 0.0.0.0
  --port="<int>"                      specify the port that the proxy server listens on. defaults: 80
  --compress                          whether keep compress from target response, set to true may slowdown response. defaults: false
  --proxy-external                    whether to proxy external host. defaults: false
  --proxy-external-ignore=<host>      specify the external host without using a proxy. defaults: ""
  --req-header="key=value"            specify the request header attached to the request. defaults: ""
  --res-header="key=value"            specify the response headers. defaults: ""
  --cors                              whether enable cors. defaults: false
  --overwrite=<folder>                enable overwrite with folder. defaults: ""
  --no-cache                          disabled cache for response. defaults: true
  --tls-cert-file                     the cert file for enabled tls. defaults: ""
  --tls-key-file                      the key file for enabled tls. defaults: ""

EXAMPLES:
  forward http://example.com
  forward --port=80 http://example.com
  forward --req-header="foo=bar" http://example.com
  forward --cors --req-header="foo=bar" --req-header="hello=world" http://example.com
  forward --port=443 --tls-cert-file=/path/to/cert/file --tls-key-file=/path/to/key/file http://example.com`)
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
		showHelp             bool   = false
		showVersion          bool   = false
		address              string = "0.0.0.0"
		port                 string = "80"
		compress             bool   = false
		cors                 bool   = false
		noCache              bool   = true
		overwriteFolder      string = ""
		proxyExternal        bool   = false
		proxyExternalIgnores arrayFlags
		requestHeadersArray  arrayFlags
		responseHeadersArray arrayFlags
		certFilePath         string = ""
		keyFilePath          string = ""
	)

	if len(os.Getenv("PORT")) > 0 {
		PORT_FROM_ENV := os.Getenv("PORT")

		if PORT_FROM_ENV != "" {
			port = PORT_FROM_ENV
		}
	}

	flag.BoolVar(&showHelp, "help", showHelp, "")
	flag.BoolVar(&showVersion, "version", showVersion, "")
	flag.Var(&requestHeadersArray, "req-header", "")
	flag.Var(&responseHeadersArray, "res-header", "")
	flag.BoolVar(&compress, "compress", compress, "")
	flag.BoolVar(&cors, "cors", cors, "")
	flag.BoolVar(&noCache, "no-cache", noCache, "")
	flag.BoolVar(&proxyExternal, "proxy-external", proxyExternal, "")
	flag.Var(&proxyExternalIgnores, "proxy-external-ignore", "")
	flag.StringVar(&port, "port", port, "")
	flag.StringVar(&address, "address", address, "")
	flag.StringVar(&overwriteFolder, "overwrite", overwriteFolder, "")
	flag.StringVar(&certFilePath, "cert-file", certFilePath, "")
	flag.StringVar(&keyFilePath, "key-file", keyFilePath, "")

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

	if !filepath.IsAbs(overwriteFolder) {
		cwd, err := os.Getwd()

		if err != nil {
			log.Panic(err)
		}

		overwriteFolder = filepath.Join(cwd, overwriteFolder)
	}

	folder, err := os.Stat(overwriteFolder)

	if os.IsNotExist(err) {
		log.Panicln("the folder of '--overwrite=<folder>' not found in your system")
	}

	if err != nil {
		log.Panicln(err)
	}

	if !folder.IsDir() {
		log.Panicln("the flag '--overwrite=<folder>' must be a folder")
	}

	proxy := forward.NewProxyServer(&forward.ProxyServerOptions{
		ReqHeaders:           requestHeaders,
		ResHeaders:           responseHeaders,
		Compress:             compress,
		Cors:                 cors,
		ProxyExternal:        proxyExternal,
		ProxyExternalIgnores: proxyExternalIgnores,
		Target:               u,
		NoCache:              noCache,
		OverwriteFolder:      overwriteFolder,
	})

	http.HandleFunc("/", proxy.Handler())

	if address == "0.0.0.0" {
		log.Printf("Proxy 'http://%s:%s' to '%s'\n", getLocalIP(), port, target)
	} else {
		log.Printf("Proxy 'http://%s:%s' to '%s'\n", address, port, target)
	}

	if certFilePath != "" && keyFilePath != "" {
		log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%s:%s", address, port), certFilePath, keyFilePath, nil))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil))
	}
}
