//go:build embed

package web

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	htmlpkg "html"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/gin-gonic/gin"
)

const (
	// NonceHTMLPlaceholder is the placeholder for nonce in HTML script tags
	NonceHTMLPlaceholder        = "__CSP_NONCE_VALUE__"
	frontendNoStore             = "no-store, no-cache, must-revalidate, max-age=0"
	frontendCacheRevision       = "registration-alias-toggle-20260724-v10"
	frontendCacheRevisionCookie = "sub2api_frontend_revision"
)

const staleFrontendAssetReloadModule = `
try {
  sessionStorage.removeItem("chunk_reload_attempted");
  sessionStorage.removeItem("sub2api_asset_reload");
  sessionStorage.removeItem("sub2api_asset_reload_at");
} catch (e) {}
try {
  var now = Date.now();
  var url = new URL(window.location.href);
  url.searchParams.set("__sub2api_refresh", String(now));
  window.location.replace(url.toString());
} catch (e) {
  window.location.reload();
}
throw new Error("Sub2API stale frontend asset");`

const frontendBootWatchdogScript = `
(function(){
  var retryKey="sub2api_boot_retry_at";
  function appHasRendered(){
    var app=document.getElementById("app");
    if(!app)return false;
    if(app.querySelector("[data-sub2api-ready],header,main,form,iframe"))return true;
    return (app.textContent||"").trim().length>20;
  }
  function showLoadError(){
    var app=document.getElementById("app");
    if(!app||appHasRendered())return;
    var shell=document.createElement("div");
    shell.style.cssText="min-height:100vh;display:flex;align-items:center;justify-content:center;background:#eef4fb;color:#0b1324;font-family:system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,sans-serif;padding:24px;text-align:center";
    var panel=document.createElement("div");
    var title=document.createElement("div");
    title.textContent="前端加载失败";
    title.style.cssText="font-size:18px;font-weight:700;margin-bottom:8px";
    var message=document.createElement("div");
    message.textContent="浏览器可能仍在使用旧缓存，请点击下方按钮重新加载。";
    message.style.cssText="font-size:14px;color:#475569;margin-bottom:16px";
    var button=document.createElement("button");
    button.type="button";
    button.textContent="重新加载";
    button.style.cssText="border:0;border-radius:8px;background:#2563eb;color:white;font-weight:600;padding:10px 16px;cursor:pointer";
    button.addEventListener("click",function(){
      try{sessionStorage.removeItem(retryKey);}catch(e){}
      location.reload();
    });
    panel.appendChild(title);
    panel.appendChild(message);
    panel.appendChild(button);
    shell.appendChild(panel);
    app.replaceChildren(shell);
  }
  setTimeout(function(){
    if(appHasRendered())return;
    try{
      var now=Date.now();
      var last=Number(sessionStorage.getItem(retryKey)||"0");
      if(!last||now-last>15000){
        sessionStorage.setItem(retryKey,String(now));
        sessionStorage.removeItem("chunk_reload_attempted");
        sessionStorage.removeItem("sub2api_asset_reload");
        sessionStorage.removeItem("sub2api_asset_reload_at");
        var url=new URL(location.href);
        url.searchParams.set("__sub2api_refresh",String(now));
        location.replace(url.toString());
        return;
      }
    }catch(e){}
    showLoadError();
  },2500);
})();`

//go:embed all:dist
var frontendFS embed.FS

// PublicSettingsProvider is an interface to fetch public settings
type PublicSettingsProvider interface {
	GetPublicSettingsForInjection(ctx context.Context) (any, error)
}

// FrontendServer serves the embedded frontend with settings injection
type FrontendServer struct {
	distFS      fs.FS
	fileServer  http.Handler
	baseHTML    []byte
	cache       *HTMLCache
	settings    PublicSettingsProvider
	overrideDir string // local file override directory
}

// NewFrontendServer creates a new frontend server with settings injection
func NewFrontendServer(settingsProvider PublicSettingsProvider) (*FrontendServer, error) {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		return nil, err
	}

	// Read base HTML once
	file, err := distFS.Open("index.html")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	baseHTML, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cache := NewHTMLCache()
	cache.SetBaseHTML(baseHTML)

	return &FrontendServer{
		distFS:      distFS,
		fileServer:  http.FileServer(http.FS(distFS)),
		baseHTML:    baseHTML,
		cache:       cache,
		settings:    settingsProvider,
		overrideDir: filepath.Join("data", "public"),
	}, nil
}

// InvalidateCache invalidates the HTML cache (call when settings change)
func (s *FrontendServer) InvalidateCache() {
	if s != nil && s.cache != nil {
		s.cache.Invalidate()
	}
}

// Middleware returns the Gin middleware handler
func (s *FrontendServer) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// Skip API routes
		if shouldBypassEmbeddedFrontend(path) {
			c.Next()
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// 目录请求（以 / 结尾，如内嵌的 AI 画图子应用 /image-studio/）应解析到该目录下的
		// index.html 并按静态文件原样返回；否则 fs.Open 对带尾斜杠的路径返回 invalid argument，
		// 使 fileExists 判定为 false，从而错误地把主应用的 index.html 注入到 iframe 中，
		// 导致子应用白屏 / 一直加载。
		if strings.HasSuffix(cleanPath, "/") {
			indexPath := cleanPath + "index.html"
			if s.fileExists(indexPath) {
				if s.tryServeOverride(c, indexPath) {
					return
				}
				// Do not pass an index.html path to http.FileServer. FileServer redirects
				// index.html to "./", which resolves back to the same directory URL and
				// creates an infinite 301 loop for nested apps such as /image-studio/.
				if serveEmbeddedHTMLFile(c, s.distFS, indexPath) {
					return
				}
			}
		}

		// For index.html or SPA routes, serve with injected settings.
		// Missing static assets must not fall back to HTML; browsers will otherwise
		// try to execute index.html as JS/CSS after a deployment changes chunk hashes.
		if cleanPath == "index.html" {
			s.serveIndexHTML(c)
			return
		}
		if !s.fileExists(cleanPath) {
			if isFrontendStaticAssetRequest(cleanPath) {
				serveMissingFrontendAsset(c, cleanPath)
				return
			}
			s.serveIndexHTML(c)
			return
		}

		// Try local override first
		if s.tryServeOverride(c, cleanPath) {
			return
		}

		// Serve static files normally (hashed assets get long-lived cache headers)
		applyStaticAssetCacheHeaders(c.Writer.Header(), cleanPath)
		s.fileServer.ServeHTTP(c.Writer, c.Request)
		c.Abort()
	}
}

func (s *FrontendServer) fileExists(path string) bool {
	file, err := s.distFS.Open(path)
	if err != nil {
		return false
	}
	_ = file.Close()
	return true
}

// tryServeOverride checks if a local override file exists and serves it.
// Files in overrideDir take precedence over embedded files.
func (s *FrontendServer) tryServeOverride(c *gin.Context, cleanPath string) bool {
	if s.overrideDir == "" {
		return false
	}
	filePath := filepath.Join(s.overrideDir, filepath.Clean("/"+cleanPath))
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	c.File(filePath)
	c.Abort()
	return true
}

func (s *FrontendServer) serveIndexHTML(c *gin.Context) {
	// Get nonce from context (generated by SecurityHeaders middleware)
	nonce := middleware.GetNonceFromContext(c)
	setFrontendNoStoreHeaders(c)
	applyFrontendCacheRevision(c)

	// Check cache first
	cached := s.cache.Get()
	if cached != nil {
		// Replace nonce placeholder with actual nonce before serving
		content := replaceNoncePlaceholder(cached.Content, nonce)

		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
		c.Abort()
		return
	}

	// Cache miss - fetch settings and render
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	settings, err := s.settings.GetPublicSettingsForInjection(ctx)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", s.baseHTML)
		c.Abort()
		return
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		// Fallback: serve without injection
		c.Data(http.StatusOK, "text/html; charset=utf-8", s.baseHTML)
		c.Abort()
		return
	}

	rendered := s.injectSettings(settingsJSON)
	s.cache.Set(rendered, settingsJSON)

	// Replace nonce placeholder with actual nonce before serving
	content := replaceNoncePlaceholder(rendered, nonce)

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

func (s *FrontendServer) injectSettings(settingsJSON []byte) []byte {
	// Create the script tag to inject with nonce placeholder
	// The placeholder will be replaced with actual nonce at request time
	script := []byte(`<script nonce="` + NonceHTMLPlaceholder + `">window.__APP_CONFIG__=` + string(settingsJSON) + `;` + frontendBootWatchdogScript + `</script>`)

	// Inject before </head>
	headClose := []byte("</head>")
	result := bytes.Replace(s.baseHTML, headClose, append(script, headClose...), 1)

	// Apply custom branding before the browser paints the static defaults.
	result = injectSiteTitle(result, settingsJSON)
	result = injectSiteFavicon(result, settingsJSON)

	return result
}

// injectSiteFavicon replaces the static favicon with a configured, browser-safe image URL.
func injectSiteFavicon(html, settingsJSON []byte) []byte {
	var cfg struct {
		SiteLogo string `json:"site_logo"`
	}
	if err := json.Unmarshal(settingsJSON, &cfg); err != nil {
		return html
	}

	logoURL := safeImageURL(cfg.SiteLogo)
	if logoURL == "" {
		return html
	}

	linkStart := bytes.Index(html, []byte(`<link rel="icon"`))
	if linkStart == -1 {
		return html
	}
	linkEndOffset := bytes.IndexByte(html[linkStart:], '>')
	if linkEndOffset == -1 {
		return html
	}
	linkEnd := linkStart + linkEndOffset + 1
	replacement := []byte(`<link rel="icon" href="` + htmlpkg.EscapeString(logoURL) + `" />`)

	var buf bytes.Buffer
	buf.Write(html[:linkStart])
	buf.Write(replacement)
	buf.Write(html[linkEnd:])
	return buf.Bytes()
}

func safeImageURL(value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	if strings.HasPrefix(trimmed, "/") && !strings.HasPrefix(trimmed, "//") {
		return trimmed
	}
	if strings.HasPrefix(strings.ToLower(trimmed), "data:image/") {
		return trimmed
	}

	parsed, err := url.Parse(trimmed)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return ""
	}
	return trimmed
}

// injectSiteTitle replaces the static <title> in HTML with the configured site name.
// This ensures the browser tab shows the correct title before JS executes.
func injectSiteTitle(html, settingsJSON []byte) []byte {
	var cfg struct {
		SiteName string `json:"site_name"`
	}
	if err := json.Unmarshal(settingsJSON, &cfg); err != nil || cfg.SiteName == "" {
		return html
	}

	// Find and replace the existing <title>...</title>
	titleStart := bytes.Index(html, []byte("<title>"))
	titleEnd := bytes.Index(html, []byte("</title>"))
	if titleStart == -1 || titleEnd == -1 || titleEnd <= titleStart {
		return html
	}

	newTitle := []byte("<title>" + htmlpkg.EscapeString(cfg.SiteName) + " - AI API Gateway</title>")
	var buf bytes.Buffer
	buf.Write(html[:titleStart])
	buf.Write(newTitle)
	buf.Write(html[titleEnd+len("</title>"):])
	return buf.Bytes()
}

// replaceNoncePlaceholder replaces the nonce placeholder with actual nonce value
func replaceNoncePlaceholder(html []byte, nonce string) []byte {
	return bytes.ReplaceAll(html, []byte(NonceHTMLPlaceholder), []byte(nonce))
}

// ServeEmbeddedFrontend returns a middleware for serving embedded frontend
// This is the legacy function for backward compatibility when no settings provider is available
func ServeEmbeddedFrontend() gin.HandlerFunc {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic("failed to get dist subdirectory: " + err.Error())
	}
	fileServer := http.FileServer(http.FS(distFS))
	overrideDir := filepath.Join("data", "public")

	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if shouldBypassEmbeddedFrontend(path) {
			c.Next()
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		if cleanPath == "" {
			cleanPath = "index.html"
		}

		// 目录请求（尾斜杠，如 /image-studio/）解析到该目录的 index.html，避免回退到主应用 SPA。
		if strings.HasSuffix(cleanPath, "/") {
			indexPath := cleanPath + "index.html"
			if file, err := distFS.Open(indexPath); err == nil {
				_ = file.Close()
				if tryServeOverrideFile(c, overrideDir, indexPath) {
					return
				}
				if serveEmbeddedHTMLFile(c, distFS, indexPath) {
					return
				}
			}
		}

		if file, err := distFS.Open(cleanPath); err == nil {
			_ = file.Close()
			// Try local override first
			if tryServeOverrideFile(c, overrideDir, cleanPath) {
				return
			}
			applyStaticAssetCacheHeaders(c.Writer.Header(), cleanPath)
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}

		if isFrontendStaticAssetRequest(cleanPath) {
			serveMissingFrontendAsset(c, cleanPath)
			return
		}

		serveIndexHTML(c, distFS)
	}
}

// tryServeOverrideFile is a standalone version of tryServeOverride for legacy usage.
func tryServeOverrideFile(c *gin.Context, overrideDir, cleanPath string) bool {
	if overrideDir == "" {
		return false
	}
	filePath := filepath.Join(overrideDir, filepath.Clean("/"+cleanPath))
	info, err := os.Stat(filePath)
	if err != nil || info.IsDir() {
		return false
	}
	c.File(filePath)
	c.Abort()
	return true
}

func setFrontendNoStoreHeaders(c *gin.Context) {
	c.Header("Cache-Control", frontendNoStore)
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
}

func applyFrontendCacheRevision(c *gin.Context) {
	if revision, err := c.Cookie(frontendCacheRevisionCookie); err == nil && revision == frontendCacheRevision {
		return
	}

	c.Header("Clear-Site-Data", `"cache"`)
	secure := c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https")
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     frontendCacheRevisionCookie,
		Value:    frontendCacheRevision,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func serveEmbeddedHTMLFile(c *gin.Context, distFS fs.FS, cleanPath string) bool {
	content, err := fs.ReadFile(distFS, cleanPath)
	if err != nil {
		return false
	}
	setFrontendNoStoreHeaders(c)
	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
	return true
}

func shouldBypassEmbeddedFrontend(path string) bool {
	trimmed := strings.TrimSpace(path)
	return strings.HasPrefix(trimmed, "/api/") ||
		strings.HasPrefix(trimmed, "/v1/") ||
		strings.HasPrefix(trimmed, "/v1beta/") ||
		strings.HasPrefix(trimmed, "/backend-api/") ||
		strings.HasPrefix(trimmed, "/antigravity/") ||
		strings.HasPrefix(trimmed, "/setup/") ||
		trimmed == "/health" ||
		trimmed == "/models" ||
		trimmed == "/responses" ||
		strings.HasPrefix(trimmed, "/responses/") ||
		trimmed == "/alpha/search" ||
		strings.HasPrefix(trimmed, "/images/") ||
		strings.HasPrefix(trimmed, "/videos/")
}

func isFrontendStaticAssetRequest(cleanPath string) bool {
	cleanPath = strings.TrimSpace(strings.TrimPrefix(cleanPath, "/"))
	if cleanPath == "" {
		return false
	}
	if strings.HasPrefix(cleanPath, "assets/") {
		return true
	}
	if idx := strings.LastIndexByte(cleanPath, '/'); idx >= 0 {
		cleanPath = cleanPath[idx+1:]
	}
	return strings.Contains(cleanPath, ".")
}

func serveMissingFrontendAsset(c *gin.Context, cleanPath string) {
	ext := strings.ToLower(filepath.Ext(cleanPath))
	setFrontendNoStoreHeaders(c)
	c.Header("Clear-Site-Data", `"cache"`)

	switch ext {
	case ".js", ".mjs":
		// Old cached HTML can still point at removed Vite chunks. Force a current
		// HTML fetch instead of letting the app stay on a blank page.
		c.Data(http.StatusOK, "text/javascript; charset=utf-8", []byte(staleFrontendAssetReloadModule))
	case ".css":
		// Missing CSS should not block JS execution or be served as HTML.
		c.Data(http.StatusOK, "text/css; charset=utf-8", nil)
	default:
		c.String(http.StatusNotFound, "Not found")
	}
	c.Abort()
}

func serveIndexHTML(c *gin.Context, fsys fs.FS) {
	setFrontendNoStoreHeaders(c)
	applyFrontendCacheRevision(c)

	file, err := fsys.Open("index.html")
	if err != nil {
		c.String(http.StatusNotFound, "Frontend not found")
		c.Abort()
		return
	}
	defer func() { _ = file.Close() }()

	content, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read index.html")
		c.Abort()
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	c.Abort()
}

func HasEmbeddedFrontend() bool {
	_, err := frontendFS.ReadFile("dist/index.html")
	return err == nil
}
