package app

import (
	"aurora-adminui/infra/victoria"
	"aurora-adminui/internal/cache"
	"aurora-adminui/internal/config"
	"aurora-adminui/internal/repository"
	"aurora-adminui/internal/service"
	"aurora-adminui/internal/transport/http/handler"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Modules struct {
	AdminAuthHandler     *handler.AdminAuthHandler
	AdminSecurityHandler *handler.AdminSecurityHandler
	HypervisorHandler    *handler.HypervisorHandler
	K8sHandler           *handler.K8sHandler
	ZoneHandler          *handler.ZoneHandler
	UserHandler          *handler.UserHandler
	RoleHandler          *handler.RoleHandler
	PlanHandler          *handler.PlanHandler
}

func NewModules(db *pgxpool.Pool, redisClient *redis.Client, cfg *config.Config) (*Modules, error) {
	adminRepo := repository.NewAdminRepo(db)
	adminCache := cache.NewAdminTokenCache(redisClient)
	adminSvc := service.NewAdminService(adminRepo, adminCache, cfg)

	hypervisorRepo := repository.NewHypervisorRepo(db)
	hypervisorSvc := service.NewHypervisorService(hypervisorRepo, victoria.NewClient(cfg.Victoria.QueryBaseURL))
	k8sRepo := repository.NewK8sRepo(db)
	k8sSvc, err := service.NewK8sService(k8sRepo, cfg.K8s.KubeconfigEncryptionKey)
	if err != nil {
		return nil, err
	}
	zoneRepo := repository.NewZoneRepo(db)
	zoneSvc := service.NewZoneService(zoneRepo)
	userRepo := repository.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	roleRepo := repository.NewRoleRepo(db)
	roleSvc := service.NewRoleService(roleRepo)
	planRepo := repository.NewPlanRepo(db)
	planSvc := service.NewPlanService(planRepo)

	return &Modules{
		AdminAuthHandler:     handler.NewAdminAuthHandler(adminSvc),
		AdminSecurityHandler: handler.NewAdminSecurityHandler(adminSvc),
		HypervisorHandler:    handler.NewHypervisorHandler(hypervisorSvc, adminSvc),
		K8sHandler:           handler.NewK8sHandler(k8sSvc, adminSvc),
		ZoneHandler:          handler.NewZoneHandler(zoneSvc, adminSvc),
		UserHandler:          handler.NewUserHandler(userSvc, adminSvc),
		RoleHandler:          handler.NewRoleHandler(roleSvc, adminSvc),
		PlanHandler:          handler.NewPlanHandler(planSvc, adminSvc),
	}, nil
}
