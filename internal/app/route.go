package app

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, modules *Modules, controlPlaneURL *url.URL, distFS fs.FS) {
	static := http.FileServer(http.FS(distFS))
	apiProxy := newControlPlaneProxy(controlPlaneURL)

	mux.HandleFunc("/api/v1/admin/auth/token-login/init", modules.AdminAuthHandler.HandleTokenLoginInit)
	mux.HandleFunc("/api/v1/admin/auth/token-login/verify-2fa", modules.AdminAuthHandler.HandleTokenLoginVerify2FA)
	mux.HandleFunc("/api/v1/admin/auth/logout", modules.AdminAuthHandler.RequireAdminSession(modules.AdminAuthHandler.HandleLogout))
	mux.HandleFunc("/api/v1/admin/auth/session", modules.AdminAuthHandler.RequireAdminSession(modules.AdminAuthHandler.HandleSessionStatus))
	mux.HandleFunc("/api/v1/admin/auth/token/rotate", modules.AdminAuthHandler.RequireAdminSession(modules.AdminAuthHandler.HandleRotateToken))

	mux.HandleFunc("/api/v1/admin/2fa/status", modules.AdminSecurityHandler.RequireAdminSession(modules.AdminSecurityHandler.HandleTwoFactorStatus))
	mux.HandleFunc("/api/v1/admin/2fa/totp/setup/begin", modules.AdminSecurityHandler.RequireAdminSession(modules.AdminSecurityHandler.HandleBeginTOTPSetup))
	mux.HandleFunc("/api/v1/admin/2fa/totp/setup/confirm", modules.AdminSecurityHandler.RequireAdminSession(modules.AdminSecurityHandler.HandleConfirmTOTPSetup))
	mux.HandleFunc("/api/v1/admin/2fa/disable", modules.AdminSecurityHandler.RequireAdminSession(modules.AdminSecurityHandler.HandleDisableTwoFactor))

	mux.HandleFunc("/api/v1/admin/hypervisor/metrics/ws", modules.HypervisorHandler.RequireAdminSession(modules.HypervisorHandler.HandleMetricsStream))
	mux.HandleFunc("/api/v1/admin/hypervisor/nodes", modules.HypervisorHandler.RequireAdminSession(modules.HypervisorHandler.HandleListNodes))
	mux.HandleFunc("/api/v1/admin/hypervisor/nodes/", modules.HypervisorHandler.RequireAdminSession(modules.HypervisorHandler.HandleNodeItem))
	mux.HandleFunc("/api/v1/admin/k8s/clusters", modules.K8sHandler.RequireAdminSession(modules.K8sHandler.HandleClustersCollection))
	mux.HandleFunc("/api/v1/admin/k8s/clusters/", modules.K8sHandler.RequireAdminSession(modules.K8sHandler.HandleClusterItem))
	mux.HandleFunc("/api/v1/admin/users", modules.UserHandler.RequireAdminSession(modules.UserHandler.HandleListUsers))
	mux.HandleFunc("/api/v1/admin/plans", modules.PlanHandler.RequireAdminSession(modules.PlanHandler.HandlePlansCollection))
	mux.HandleFunc("/api/v1/admin/roles/permissions", modules.RoleHandler.RequireAdminSession(modules.RoleHandler.HandleListPermissions))
	mux.HandleFunc("/api/v1/admin/roles", modules.RoleHandler.RequireAdminSession(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			modules.RoleHandler.HandleListRoles(w, r)
		case http.MethodPost:
			modules.RoleHandler.HandleCreateRole(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	mux.HandleFunc("/api/v1/admin/zones", modules.ZoneHandler.RequireAdminSession(modules.ZoneHandler.HandleZonesCollection))
	mux.HandleFunc("/api/v1/admin/zones/", modules.ZoneHandler.RequireAdminSession(modules.ZoneHandler.HandleZoneItem))

	mux.Handle("/api/", apiProxy)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "ok")
	})
	mux.Handle("/", spaHandler(distFS, static))
}

func newControlPlaneProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		http.Error(w, "upstream controlplane unavailable", http.StatusBadGateway)
		log.Printf("proxy error: %v", err)
	}
	return proxy
}

func spaHandler(distFS fs.FS, static http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cleanPath := path.Clean("/" + strings.TrimPrefix(r.URL.Path, "/"))
		if cleanPath != "/" {
			trimmed := strings.TrimPrefix(cleanPath, "/")
			if fileExists(distFS, trimmed) {
				static.ServeHTTP(w, r)
				return
			}
		}

		indexFile, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}
		defer indexFile.Close()

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.Copy(w, indexFile)
	}
}

func fileExists(distFS fs.FS, name string) bool {
	info, err := fs.Stat(distFS, name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
