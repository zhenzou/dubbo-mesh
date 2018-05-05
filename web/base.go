package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rcrowley/go-metrics"
)

const (
	// @formatter:off
	MinPage          = 1
	MinLimit         = 10
	MaxLimit         = 50

	CodeOk           = 0   // 正常
	CodeError        = 402 // 未知错误
	CodeUnauthorized = 401 // 未授权

	ParamLimit       = "limit"
	ParamPage        = "page"
	ParamKeyword     = "kw"
	ParamName        = "name"

	Debug            = "debug"
	DoDebug          = "yyyyy"
	// @formatter:on
)

// Headers
const (
	HeaderAccept                        = "Accept"
	HeaderAcceptEncoding                = "Accept-Encoding"
	HeaderAuthorization                 = "Authorization"
	HeaderContentDisposition            = "Content-Disposition"
	HeaderContentEncoding               = "Content-Encoding"
	HeaderContentLength                 = "Content-Length"
	HeaderContentType                   = "Content-Type"
	HeaderContentDescription            = "Content-Description"
	HeaderContentTransferEncoding       = "Content-Transfer-Encoding"
	HeaderCookie                        = "Cookie"
	HeaderSetCookie                     = "Set-Cookie"
	HeaderIfModifiedSince               = "If-Modified-Since"
	HeaderLastModified                  = "Last-Modified"
	HeaderLocation                      = "Location"
	HeaderReferer                       = "Referer"
	HeaderUserAgent                     = "User-Agent"
	HeaderUpgrade                       = "Upgrade"
	HeaderVary                          = "Vary"
	HeaderWWWAuthenticate               = "WWW-Authenticate"
	HeaderXForwardedProto               = "X-Forwarded-Proto"
	HeaderXHTTPMethodOverride           = "X-HTTP-Method-Override"
	HeaderXForwardedFor                 = "X-Forwarded-For"
	HeaderXRealIP                       = "X-Real-IP"
	HeaderXRequestedWith                = "X-Requested-With"
	HeaderServer                        = "Server"
	HeaderOrigin                        = "Origin"
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"
	HeaderExpires                       = "Expires"
	HeaderCacheControl                  = "Cache-Control"
	HeaderPragma                        = "Pragma"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json" + "; " + charsetUTF8
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=utf-8"
)

var (
	EmptyResponse interface{}
)

type ResponseCode int

type Response struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
	Data interface{}  `json:"data"`
}

type ListData struct {
	Total    int64       `json:"total"`
	Page     int         `json:"Page"`
	PageSize int         `json:"page_size"`
	Data     interface{} `json:"data"`
}

func Ok(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusOK, &Response{CodeOk, "", result})
}

func Error(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, &Response{CodeError, err.Error(), EmptyResponse})
}

func Result(ctx *gin.Context, result interface{}, err error) {
	if err != nil {
		Error(ctx, err)
	} else {
		Ok(ctx, result)
	}
}

func PageOk(ctx *gin.Context, result interface{}, total int64, page, size int) {
	ctx.JSON(http.StatusOK, &Response{CodeOk, "", ListData{total, page, size, result}})
}

func Offset(ctx *gin.Context) int {
	return (Page(ctx) - 1) * Limit(ctx)
}

func Keyword(ctx *gin.Context) string {
	return StrQuery(ctx, ParamKeyword)
}

func Limit(ctx *gin.Context) int {
	ret := IntQuery(ctx, ParamLimit)
	if ret < MinLimit {
		ret = MinLimit
	}
	if ret > MaxLimit {
		ret = MaxLimit
	}
	return ret
}

func Page(ctx *gin.Context) int {
	ret := IntQuery(ctx, ParamPage)
	if ret < MinPage {
		ret = MinPage
	}
	return ret
}

func IntQuery(ctx *gin.Context, key string) int {
	ret := StrQuery(ctx, key)
	res, _ := strconv.Atoi(ret)
	return res
}

func Int64Query(ctx *gin.Context, key string) int64 {
	ret := StrQuery(ctx, key)
	res, _ := strconv.ParseInt(ret, 10, 64)
	return res
}

func Uint64Query(ctx *gin.Context, key string) uint64 {
	ret := StrQuery(ctx, key)
	res, _ := strconv.ParseUint(ret, 10, 64)
	return res
}

func StrQuery(ctx *gin.Context, key string) string {
	if strings.Index(key, ":") == 0 {
		return ctx.Param(key[1:])
	}
	return ctx.Query(key)
}

func IsDebug(ctx *gin.Context) bool {
	debug := StrQuery(ctx, Debug)
	return debug == DoDebug
}

func NewGin() *gin.Engine {
	g := gin.New()
	g.Use(Logger(), Recovery())
	return g
}

func WriteContentType(w http.ResponseWriter, value string) {
	header := w.Header()
	if val := header.Get("Content-Type"); val == "" {
		w.Header().Set("Content-Type", value)
	}
}

func EnableMetrics(g *gin.Engine, path string) {
	g.GET(path, Metrics)
}

func Metrics(ctx *gin.Context) {
	name := StrQuery(ctx, ParamName)
	m := metrics.DefaultRegistry.GetAll()
	if name != "" {
		Ok(ctx, m[name])
	} else {
		Ok(ctx, m)
	}
}

type SingleHandler struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (s *SingleHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	s.handler(writer, req)
}

func HandlerFunc(handler func(http.ResponseWriter, *http.Request)) http.Handler {
	return &SingleHandler{handler}
}