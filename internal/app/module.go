package app

import (
	"aurora-adminui/infra/victoria"
	"aurora-adminui/internal/cache"
	"aurora-adminui/internal/config"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/repository"
	"aurora-adminui/internal/service"
	controlplanegrpc "aurora-adminui/internal/transport/grpc"
	"aurora-adminui/internal/transport/http/handler"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Modules struct {
	AdminService              domainsvc.AdminService
	AdminAuthHandler          *handler.AdminAuthHandler
	AdminSecurityHandler      *handler.AdminSecurityHandler
	HypervisorHandler         *handler.HypervisorHandler
	K8sHandler                *handler.K8sHandler
	ResourceDefinitionHandler *handler.RDHandler
	TemplateRenderHandler     *handler.TemplateRenderHandler
	MarketplaceHandler        *handler.MarketplaceHandler
	ZoneHandler               *handler.ZoneHandler
	UserHandler               *handler.UserHandler
	RoleHandler               *handler.RoleHandler
	PlanHandler               *handler.PlanHandler
}

// NewModules builds all repositories, services, and HTTP handlers used by adminui.
func NewModules(db *pgxpool.Pool, redisClient *redis.Client, cfg *config.Config) (*Modules, error) {
	adminRepo := repository.NewAdminRepo(db)
	adminCache := cache.NewAdminTokenCache(redisClient)
	adminSvc := service.NewAdminService(adminRepo, adminCache, cfg)

	hypervisorRepo := repository.NewHypervisorRepo(db)
	hypervisorSvc := service.NewHypervisorService(hypervisorRepo, victoria.NewClient(cfg.Victoria.QueryBaseURL))
	controlPlaneClient := controlplanegrpc.New(cfg.ControlPlane.GRPCAddr)
	k8sRepo := repository.NewK8sRepo(controlPlaneClient)
	k8sSvc := service.NewK8sService(k8sRepo)
	resourceDefinitionRepo := repository.NewResourceDefinitionRepo(controlPlaneClient)
	resourceDefinitionSvc := service.NewResourceDefinitionService(resourceDefinitionRepo)
	zoneRepo := repository.NewZoneRepo(db)
	zoneSvc := service.NewZoneService(zoneRepo)
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	roleRepo := repository.NewRoleRepo(db)
	roleSvc := service.NewRoleService(roleRepo)
	planRepo := repository.NewPlanRepo(db)
	planSvc := service.NewPlanService(planRepo)
	templateRenderRepo := repository.NewTemplateRenderRepo(controlPlaneClient)
	templateRenderSvc := service.NewTemplateRenderService(templateRenderRepo)
	marketplaceRepo := repository.NewMarketplaceRepo(db)
	marketplaceSvc := service.NewMarketplaceService(marketplaceRepo)

	return &Modules{
		AdminService:              adminSvc,
		AdminAuthHandler:          handler.NewAdminAuthHandler(adminSvc),
		AdminSecurityHandler:      handler.NewAdminSecurityHandler(adminSvc),
		HypervisorHandler:         handler.NewHypervisorHandler(hypervisorSvc),
		K8sHandler:                handler.NewK8sHandler(k8sSvc, zoneSvc),
		ResourceDefinitionHandler: handler.NewRDHandler(resourceDefinitionSvc),
		TemplateRenderHandler:     handler.NewTemplateRenderHandler(templateRenderSvc),
		MarketplaceHandler:        handler.NewMarketplaceHandler(marketplaceSvc),
		ZoneHandler:               handler.NewZoneHandler(zoneSvc),
		UserHandler:               handler.NewUserHandler(userSvc),
		RoleHandler:               handler.NewRoleHandler(roleSvc),
		PlanHandler:               handler.NewPlanHandler(planSvc),
	}, nil
}
