package errorx

import "errors"

var (
	ErrInvalidArgument         = errors.New("invalid argument")
	ErrTokenInvalid            = errors.New("token invalid")
	ErrTokenExpired            = errors.New("token expired")
	ErrMFACodeInvalid          = errors.New("mfa code invalid")
	ErrMFAMethodNotFound       = errors.New("mfa method not found")
	ErrMFAMethodAlreadyEnabled = errors.New("mfa method already enabled")
	ErrAPITokenNotFound        = errors.New("api token not found")
	ErrAdminSecurityNotFound   = errors.New("admin security state not found")
	ErrAdminSessionNotFound    = errors.New("admin session not found")
	ErrHypervisorNodeNotFound  = errors.New("hypervisor node not found")
	ErrZoneNotFound            = errors.New("zone not found")
	ErrZoneAlreadyExists       = errors.New("zone already exists")
	ErrZoneHasResources        = errors.New("zone has resources")
	ErrRoleAlreadyExists       = errors.New("role already exists")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPlanAlreadyExists       = errors.New("plan already exists")
	ErrK8sClusterNotFound      = errors.New("k8s cluster not found")
	ErrK8sClusterAlreadyExists = errors.New("k8s cluster already exists")
)
