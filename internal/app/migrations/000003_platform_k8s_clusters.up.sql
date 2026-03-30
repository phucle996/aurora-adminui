CREATE SCHEMA IF NOT EXISTS platform;

CREATE TABLE IF NOT EXISTS platform.k8s_clusters (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  import_mode TEXT NOT NULL DEFAULT 'kubeconfig' CHECK (import_mode = 'kubeconfig'),
  kubeconfig_ciphertext TEXT NOT NULL,
  api_server_url TEXT NOT NULL,
  current_context TEXT NOT NULL,
  kubernetes_version TEXT NOT NULL DEFAULT '',
  validation_status TEXT NOT NULL CHECK (validation_status IN ('pending', 'valid', 'invalid', 'unreachable')),
  last_validated_at TIMESTAMPTZ NULL,
  last_validation_error TEXT NOT NULL DEFAULT '',
  supports_dbaas BOOLEAN NOT NULL DEFAULT FALSE,
  supports_serverless BOOLEAN NOT NULL DEFAULT FALSE,
  supports_generic_workloads BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_platform_k8s_clusters_created_at
  ON platform.k8s_clusters (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_platform_k8s_clusters_validation_status
  ON platform.k8s_clusters (validation_status);
