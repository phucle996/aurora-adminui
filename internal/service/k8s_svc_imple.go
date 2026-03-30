package service

import (
	"context"
	"errors"
	"net"
	"net/url"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/security"

	"github.com/google/uuid"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sSvcImple struct {
	repo   domainrepo.K8sRepository
	cipher *security.SymmetricCipher
}

func NewK8sService(repo domainrepo.K8sRepository, rawKey string) (domainsvc.K8sService, error) {
	cipher, err := security.NewSymmetricCipher(rawKey)
	if err != nil {
		return nil, err
	}
	return &K8sSvcImple{repo: repo, cipher: cipher}, nil
}

func (s *K8sSvcImple) ListClusters(ctx context.Context) ([]entity.K8sCluster, error) {
	return s.repo.ListClusters(ctx)
}

func (s *K8sSvcImple) GetClusterDetail(ctx context.Context, id string) (*entity.K8sCluster, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.GetClusterByID(ctx, id)
}

func (s *K8sSvcImple) CreateCluster(ctx context.Context, input entity.K8sClusterCreateInput) (*entity.K8sCluster, error) {
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)
	if name == "" || len(input.Kubeconfig) == 0 {
		return nil, errorx.ErrInvalidArgument
	}

	cluster := &entity.K8sCluster{
		ID:                       uuid.New(),
		Name:                     name,
		Description:              description,
		ImportMode:               entity.K8sClusterImportModeKubeconfig,
		SupportsDBAAS:            input.SupportsDBAAS,
		SupportsServerless:       input.SupportsServerless,
		SupportsGenericWorkloads: input.SupportsGenericWorkloads,
		CreatedAt:                time.Now().UTC(),
	}

	if strings.TrimSpace(input.ZoneID) != "" {
		zoneID, err := uuid.Parse(strings.TrimSpace(input.ZoneID))
		if err != nil {
			return nil, errorx.ErrInvalidArgument
		}
		cluster.ZoneID = &zoneID
	}

	ciphertext, err := s.cipher.Encrypt(input.Kubeconfig)
	if err != nil {
		return nil, err
	}
	cluster.KubeconfigCiphertext = ciphertext

	validation, err := s.validateCluster(ctx, input.Kubeconfig)
	if err != nil {
		return nil, err
	}
	applyK8sValidation(cluster, validation)

	if err := s.repo.CreateCluster(ctx, cluster); err != nil {
		return nil, err
	}
	return s.repo.GetClusterByID(ctx, cluster.ID.String())
}

func (s *K8sSvcImple) RevalidateCluster(ctx context.Context, id string) (*entity.K8sCluster, error) {
	cluster, err := s.GetClusterDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	rawKubeconfig, err := s.cipher.Decrypt(cluster.KubeconfigCiphertext)
	if err != nil {
		return nil, err
	}

	validation, err := s.validateCluster(ctx, rawKubeconfig)
	if err != nil {
		return nil, err
	}
	applyK8sValidation(cluster, validation)

	if err := s.repo.UpdateClusterValidation(ctx, cluster); err != nil {
		return nil, err
	}
	return s.repo.GetClusterByID(ctx, cluster.ID.String())
}

func (s *K8sSvcImple) DeleteCluster(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errorx.ErrInvalidArgument
	}
	return s.repo.DeleteCluster(ctx, id)
}

type k8sValidationResult struct {
	APIServerURL        string
	CurrentContext      string
	KubernetesVersion   string
	ValidationStatus    entity.K8sClusterValidationStatus
	LastValidatedAt     time.Time
	LastValidationError string
}

func (s *K8sSvcImple) validateCluster(ctx context.Context, rawKubeconfig []byte) (*k8sValidationResult, error) {
	rawConfig, err := clientcmd.Load(rawKubeconfig)
	if err != nil {
		return nil, errorx.ErrInvalidArgument
	}

	currentContext := strings.TrimSpace(rawConfig.CurrentContext)
	if currentContext == "" {
		return nil, errorx.ErrInvalidArgument
	}
	contextConfig := rawConfig.Contexts[currentContext]
	if contextConfig == nil || strings.TrimSpace(contextConfig.Cluster) == "" {
		return nil, errorx.ErrInvalidArgument
	}
	clusterConfig := rawConfig.Clusters[contextConfig.Cluster]
	if clusterConfig == nil || strings.TrimSpace(clusterConfig.Server) == "" {
		return nil, errorx.ErrInvalidArgument
	}

	apiServerURL := strings.TrimSpace(clusterConfig.Server)
	if _, err := url.ParseRequestURI(apiServerURL); err != nil {
		return nil, errorx.ErrInvalidArgument
	}

	clientConfig, err := clientcmd.NewDefaultClientConfig(*rawConfig, nil).ClientConfig()
	if err != nil {
		return nil, errorx.ErrInvalidArgument
	}
	clientConfig.Timeout = 8 * time.Second

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, errorx.ErrInvalidArgument
	}

	result := &k8sValidationResult{
		APIServerURL:     apiServerURL,
		CurrentContext:   currentContext,
		ValidationStatus: entity.K8sClusterValidationStatusPending,
		LastValidatedAt:  time.Now().UTC(),
	}

	versionInfo, err := clientset.Discovery().ServerVersion()
	if err != nil {
		result.ValidationStatus = classifyK8sValidationError(err)
		result.LastValidationError = strings.TrimSpace(err.Error())
		return result, nil
	}
	result.KubernetesVersion = strings.TrimSpace(versionInfo.GitVersion)

	_, err = clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{Limit: 1})
	if err != nil {
		result.ValidationStatus = classifyK8sValidationError(err)
		result.LastValidationError = strings.TrimSpace(err.Error())
		return result, nil
	}

	result.ValidationStatus = entity.K8sClusterValidationStatusValid
	result.LastValidationError = ""
	return result, nil
}

func applyK8sValidation(cluster *entity.K8sCluster, validation *k8sValidationResult) {
	cluster.APIServerURL = validation.APIServerURL
	cluster.CurrentContext = validation.CurrentContext
	cluster.KubernetesVersion = validation.KubernetesVersion
	cluster.ValidationStatus = validation.ValidationStatus
	cluster.LastValidationError = validation.LastValidationError
	cluster.LastValidatedAt = &validation.LastValidatedAt
}

func classifyK8sValidationError(err error) entity.K8sClusterValidationStatus {
	if err == nil {
		return entity.K8sClusterValidationStatusValid
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return entity.K8sClusterValidationStatusUnreachable
	}
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(msg, "no such host"),
		strings.Contains(msg, "dial tcp"),
		strings.Contains(msg, "i/o timeout"),
		strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "tls handshake timeout"),
		strings.Contains(msg, "context deadline exceeded"),
		strings.Contains(msg, "temporarily unavailable"):
		return entity.K8sClusterValidationStatusUnreachable
	default:
		return entity.K8sClusterValidationStatusInvalid
	}
}
