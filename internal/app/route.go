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

	"aurora-adminui/internal/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, modules *Modules, controlPlaneURL *url.URL, distFS fs.FS) {
	static := http.FileServer(http.FS(distFS))
	apiProxy := newControlPlaneProxy(controlPlaneURL)
	wrapHTTP := func(next http.HandlerFunc) gin.HandlerFunc {
		return gin.WrapF(next)
	}
	requireAdminHTTP := func(next http.HandlerFunc) gin.HandlerFunc {
		return gin.WrapF(middleware.RequireAdminSession(modules.AdminService, next))
	}
	requireAdminGin := func(next gin.HandlerFunc) gin.HandlerFunc {
		return func(c *gin.Context) {
			middleware.RequireAdminSession(modules.AdminService, func(http.ResponseWriter, *http.Request) {
				next(c)
			})(c.Writer, c.Request)
		}
	}

	router.POST("/api/v1/admin/auth/token-login/init", wrapHTTP(modules.AdminAuthHandler.HandleTokenLoginInit))
	router.POST("/api/v1/admin/auth/token-login/verify-2fa", wrapHTTP(modules.AdminAuthHandler.HandleTokenLoginVerify2FA))
	router.POST("/api/v1/admin/auth/logout", requireAdminHTTP(modules.AdminAuthHandler.HandleLogout))
	router.GET("/api/v1/admin/auth/session", requireAdminHTTP(modules.AdminAuthHandler.HandleSessionStatus))
	router.POST("/api/v1/admin/auth/token/rotate", requireAdminHTTP(modules.AdminAuthHandler.HandleRotateToken))

	router.GET("/api/v1/admin/2fa/status", requireAdminHTTP(modules.AdminSecurityHandler.HandleTwoFactorStatus))
	router.POST("/api/v1/admin/2fa/totp/setup/begin", requireAdminHTTP(modules.AdminSecurityHandler.HandleBeginTOTPSetup))
	router.POST("/api/v1/admin/2fa/totp/setup/confirm", requireAdminHTTP(modules.AdminSecurityHandler.HandleConfirmTOTPSetup))
	router.POST("/api/v1/admin/2fa/disable", requireAdminHTTP(modules.AdminSecurityHandler.HandleDisableTwoFactor))

	router.GET("/api/v1/admin/hypervisor/metrics/ws", requireAdminHTTP(modules.HypervisorHandler.HandleMetricsStream))
	router.GET("/api/v1/admin/hypervisor/nodes", requireAdminHTTP(modules.HypervisorHandler.HandleListNodes))
	router.Any("/api/v1/admin/hypervisor/nodes/*rest", requireAdminHTTP(modules.HypervisorHandler.HandleNodeItem))

	router.GET("/api/v1/admin/k8s/clusters", requireAdminHTTP(modules.K8sHandler.HandleClustersCollection))
	router.POST("/api/v1/admin/k8s/clusters", requireAdminHTTP(modules.K8sHandler.HandleClustersCollection))
	router.Any("/api/v1/admin/k8s/clusters/*rest", requireAdminHTTP(modules.K8sHandler.HandleClusterItem))

	router.GET("/api/v1/admin/resource-definitions/template-options", requireAdminGin(modules.ResourceDefinitionHandler.ListRDTemplateOptions))
	router.GET("/api/v1/admin/resource-definitions", requireAdminGin(modules.ResourceDefinitionHandler.ListRD))
	router.POST("/api/v1/admin/resource-definitions", requireAdminGin(modules.ResourceDefinitionHandler.CreateRD))
	router.GET("/api/v1/admin/resource-definitions/:id/zones", requireAdminGin(modules.ResourceDefinitionHandler.ListRDZones))
	router.PUT("/api/v1/admin/resource-definitions/:id/zones", requireAdminGin(modules.ResourceDefinitionHandler.ReplaceRDZones))
	router.PATCH("/api/v1/admin/resource-definitions/:id", requireAdminGin(modules.ResourceDefinitionHandler.UpdateRDStatus))
	router.DELETE("/api/v1/admin/resource-definitions/:id", requireAdminGin(modules.ResourceDefinitionHandler.DeleteRD))

	router.GET("/api/v1/admin/marketplace/model-options", requireAdminHTTP(modules.MarketplaceHandler.HandleMarketplaceCollection))
	router.GET("/api/v1/admin/marketplace/template-options", requireAdminHTTP(modules.MarketplaceHandler.HandleMarketplaceCollection))
	router.GET("/api/v1/admin/marketplace", requireAdminHTTP(modules.MarketplaceHandler.HandleMarketplaceCollection))
	router.POST("/api/v1/admin/marketplace", requireAdminHTTP(modules.MarketplaceHandler.HandleMarketplaceCollection))
	router.Any("/api/v1/admin/marketplace/:id", requireAdminHTTP(modules.MarketplaceHandler.HandleMarketplaceItem))

	router.GET("/api/v1/admin/resource-templates/catalog", requireAdminHTTP(modules.TemplateRenderHandler.HandleListTemplateRenderCatalog))
	router.POST("/api/v1/admin/resource-templates", requireAdminHTTP(modules.TemplateRenderHandler.HandleTemplateRenderCollection))
	router.GET("/api/v1/admin/resource-templates/:id", requireAdminHTTP(modules.TemplateRenderHandler.HandleTemplateRenderItem))
	router.PATCH("/api/v1/admin/resource-templates/:id", requireAdminHTTP(modules.TemplateRenderHandler.HandleTemplateRenderItem))
	router.DELETE("/api/v1/admin/resource-templates/:id", requireAdminHTTP(modules.TemplateRenderHandler.HandleTemplateRenderItem))

	router.GET("/api/v1/admin/users", requireAdminHTTP(modules.UserHandler.HandleListUsers))
	router.GET("/api/v1/admin/plans", requireAdminHTTP(modules.PlanHandler.HandlePlansCollection))
	router.POST("/api/v1/admin/plans", requireAdminHTTP(modules.PlanHandler.HandlePlansCollection))
	router.GET("/api/v1/admin/roles/permissions", requireAdminHTTP(modules.RoleHandler.HandleListPermissions))
	router.GET("/api/v1/admin/roles", requireAdminHTTP(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			modules.RoleHandler.HandleListRoles(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}))
	router.POST("/api/v1/admin/roles", requireAdminHTTP(modules.RoleHandler.HandleCreateRole))
	router.GET("/api/v1/admin/zones", requireAdminHTTP(modules.ZoneHandler.HandleZonesCollection))
	router.POST("/api/v1/admin/zones", requireAdminHTTP(modules.ZoneHandler.HandleZonesCollection))
	router.DELETE("/api/v1/admin/zones/:id", requireAdminHTTP(modules.ZoneHandler.HandleZoneItem))

	router.GET("/healthz", func(c *gin.Context) {
		w := c.Writer
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.WriteString(w, "ok")
	})
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			apiProxy.ServeHTTP(c.Writer, c.Request)
			return
		}
		spaHandler(distFS, static)(c.Writer, c.Request)
	})
}

// newControlPlaneProxy forwards every unhandled /api request to controlplane.
func newControlPlaneProxy(target *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(w http.ResponseWriter, _ *http.Request, err error) {
		http.Error(w, "upstream controlplane unavailable", http.StatusBadGateway)
		log.Printf("proxy error: %v", err)
	}
	return proxy
}

// spaHandler serves built assets directly and falls back to index.html for client-side routes.
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

// fileExists checks whether a built frontend asset exists in the embedded dist filesystem.
func fileExists(distFS fs.FS, name string) bool {
	info, err := fs.Stat(distFS, name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
