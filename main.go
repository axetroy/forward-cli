package main

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/andybalholm/brotli"
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
  --port="<int>"                      specify the port that the proxy server listens on. defaults: 8080
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

type ProxyServer struct {
	target            *url.URL
	reqHeaders        http.Header
	resHeaders        http.Header
	cors              bool
	corsAllowHeaders  string
	corsExposeHeaders string
	proxy             *httputil.ReverseProxy
}

type ProxyServerOptions struct {
	Target            *url.URL
	ReqHeaders        http.Header
	ResHeaders        http.Header
	Cors              bool
	CorsAllowHeaders  string
	CorsExposeHeaders string
}

func NewProxyServer(options ProxyServerOptions) *ProxyServer {
	proxy := httputil.NewSingleHostReverseProxy(options.Target)

	server := &ProxyServer{
		reqHeaders:        options.ReqHeaders,
		resHeaders:        options.ResHeaders,
		cors:              options.Cors,
		corsAllowHeaders:  options.CorsAllowHeaders,
		corsExposeHeaders: options.CorsExposeHeaders,
		target:            options.Target,
		proxy:             proxy,
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
	req.Header.Set("X-Origin-Host", req.Host)
	req.Host = p.target.Host
	req.URL.Host = p.target.Host
	req.URL.Scheme = p.target.Scheme
	req.Header.Set("Host", p.target.Host)
	req.Header.Set("Origin", fmt.Sprintf("%s://%s", p.target.Scheme, p.target.Host))
	req.Header.Set("Referrer", fmt.Sprintf("%s://%s%s", p.target.Scheme, p.target.Host, req.URL.RawPath))
	req.Header.Set("X-Real-IP", req.RemoteAddr)

	for k := range p.reqHeaders {
		req.Header.Add(k, p.reqHeaders.Get(k))
	}
}

func (p *ProxyServer) modifyContent(body []byte, host string) []byte {
	bodyStr := string(body)

	replaceMap := map[string]string{
		"https://": "http://",
		"wss://":   "ws://",
		"//":       "//",
	}

	var newBodyStr = bodyStr

	for oldScheme, newScheme := range replaceMap {
		originHost := oldScheme + p.target.Host
		proxyHost := newScheme + host
		// newBodyStr = strings.ReplaceAll(newBodyStr, originHost, proxyHost)
		newBodyStr = regexp.MustCompile("\\b"+strings.ReplaceAll(originHost, ".", "\\.")+"\\b").ReplaceAllString(newBodyStr, proxyHost)
	}

	// TODO: fix me
	newBodyStr = strings.ReplaceAll(newBodyStr, host+".", p.target.Host+".")
	newBodyStr = strings.ReplaceAll(newBodyStr, "."+host, "."+p.target.Host)

	return []byte(newBodyStr)
}

func (p *ProxyServer) modifyResponse(res *http.Response) error {
	host := res.Request.Header.Get("X-Origin-Host") // localhost:8080

	hostName, _, err := net.SplitHostPort(host)

	if err != nil {
		return err
	}

	res.Header.Set("X-Proxy-Client", "Forward-Cli")
	res.Header.Del("Expect-CT")

	// overwrite cookies
	{
		cookies := res.Cookies()
		res.Header.Del("Set-Cookie")

		for _, v := range cookies {
			v.Domain = hostName
			if v.Secure {
				v.Secure = false
			}

			res.Header.Add("Set-Cookie", v.String())
		}
	}

	// overrit 302 Location
	{
		for _, v := range res.Header["Location"] {
			// replace location
			newLocation := strings.Replace(v, fmt.Sprintf("%s://%s", p.target.Scheme, p.target.Host), "", -1)

			res.Header.Set("Location", newLocation)
		}
	}

	if p.cors {
		res.Header.Set("Access-Control-Allow-Origin", "*")
		res.Header.Set("Access-Control-Allow-Credentials", "true")
		if p.corsAllowHeaders != "" {
			res.Header.Add("Access-Control-Allow-Headers", p.corsAllowHeaders)
		}
		if p.corsExposeHeaders != "" {
			res.Header.Add("Access-Control-Expose-Headers", p.corsExposeHeaders)
		}
	}

	for k := range p.resHeaders {
		res.Header.Add(k, p.resHeaders.Get(k))
	}

	// replace HTML/css/javascript content
	{
		contentType := res.Header.Get("Content-Type")

		ext, err := mime.ExtensionsByType(contentType)

		if err != nil {
			return nil
		}

		replaceExtNames := []string{".html", ".htm", ".css", ".js", ".ts", ".txt", ".text"}
		isSupportReplace := false

		for _, v := range replaceExtNames {
			if contains(ext, v) {
				isSupportReplace = true
				break
			}
		}

		if !isSupportReplace {
			return nil
		}

		encoding := res.Header.Get("Content-Encoding")

		// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Content-Encoding
		switch encoding {
		case "gzip":
			reader, err := gzip.NewReader(res.Body)

			if err != nil {
				return err
			}

			defer reader.Close()

			body, err := ioutil.ReadAll(reader)

			if err != nil {
				return err
			}

			newBody := p.modifyContent(body, host)

			var b bytes.Buffer
			gz := gzip.NewWriter(&b)

			if _, err := gz.Write(newBody); err != nil {
				return err
			}

			if err := gz.Close(); err != nil {
				return err
			}

			bin := b.Bytes()
			res.Header.Set("Content-Length", fmt.Sprint(len(bin)))
			res.Body = io.NopCloser(bytes.NewReader(bin))
		case "compress":
			// Deprecated by most browsers
		case "deflate":
			reader, err := zlib.NewReader(res.Body)

			if err != nil {
				return err
			}

			body, err := ioutil.ReadAll(reader)

			defer reader.Close()

			if err != nil {
				return err
			}

			newBody := p.modifyContent(body, host)

			buf := &bytes.Buffer{}

			w := zlib.NewWriter(buf)

			if n, err := w.Write(newBody); err != nil {
				return err
			} else if n < len(newBody) {
				return fmt.Errorf("n too small: %d vs %d for %s", n, len(newBody), string(newBody))
			}

			if err := w.Close(); err != nil {
				return err
			}

			res.Header.Set("Content-Length", fmt.Sprint(buf.Len()))
			res.Body = io.NopCloser(buf)
		case "br":
			reader := brotli.NewReader(res.Body)

			body, err := ioutil.ReadAll(reader)

			if err != nil {
				return err
			}

			newBody := p.modifyContent(body, host)

			buf := &bytes.Buffer{}
			w := brotli.NewWriter(buf)
			if n, err := w.Write(newBody); err != nil {
				return err
			} else if n < len(newBody) {
				return fmt.Errorf("n too small: %d vs %d for %s", n, len(newBody), string(newBody))
			}

			if err := w.Close(); err != nil {
				return err
			}

			res.Header.Set("Content-Length", fmt.Sprint(buf.Len()))
			res.Body = io.NopCloser(buf)
		case "identity":
			fallthrough
		default:
			// origin response data without compress

			body, err := ioutil.ReadAll(res.Body)

			if err != nil {
				return err
			}

			defer res.Body.Close()

			newBody := p.modifyContent(body, host)

			if err != nil {
				return err
			}

			res.Header.Set("Content-Length", fmt.Sprint(len(newBody)))
			res.Body = io.NopCloser(bytes.NewReader(newBody))
		}

	}

	return nil
}

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		log.Panicln(err)
	}

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "custom array flag"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
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
		port                 string = "8080"
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

	proxy := NewProxyServer(ProxyServerOptions{
		ReqHeaders:        requestHeaders,
		ResHeaders:        responseHeaders,
		Cors:              cors,
		CorsAllowHeaders:  corsAllowHeaders,
		CorsExposeHeaders: corsExposeHeaders,
		Target:            u,
	})

	http.HandleFunc("/", proxy.Handler())

	ip := getLocalIP()

	log.Printf("Proxy 'http://%s:%s' to '%s'\n", ip, port, target)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
