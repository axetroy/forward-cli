package forward

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/andybalholm/brotli"
)

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
		rw.WriteHeader(http.StatusInternalServerError)
		_, _ = rw.Write([]byte(err.Error()))
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
		"https://":     "http://",
		"https:\\/\\/": "http://",
		"http://":      "http://",
		"http:\\/\\/":  "http://",
		"wss://":       "ws://",
		"//":           "//",
		"":             "",
	}

	var newBodyStr = bodyStr

	for oldScheme, newScheme := range replaceMap {
		originHost := oldScheme + p.target.Host
		proxyHost := newScheme + host
		newBodyStr = strings.ReplaceAll(newBodyStr, originHost, proxyHost)
		// newBodyStr = regexp.MustCompile("\\b"+strings.ReplaceAll(originHost, ".", "\\.")+"\\b").ReplaceAllString(newBodyStr, proxyHost)
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

		replaceExtNames := []string{".html", ".htm", ".xhtml", ".xml", ".yml", ".yaml", ".css", ".js", ".txt", ".text", ".json"}
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
