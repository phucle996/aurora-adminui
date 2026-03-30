#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="${SCRIPT_DIR}"

SERVICE_NAME="${SERVICE_NAME:-aurora-adminui}"
BIN_NAME="${BIN_NAME:-aurora-adminui}"

INSTALL_BIN_DIR="${INSTALL_BIN_DIR:-/usr/local/bin}"
INSTALL_ETC_DIR="${INSTALL_ETC_DIR:-/etc/aurora}"
INSTALL_STATE_DIR="${INSTALL_STATE_DIR:-/var/lib/aurora/adminui}"
SYSTEMD_UNIT_DIR="${SYSTEMD_UNIT_DIR:-/etc/systemd/system}"

ENV_EXAMPLE_PATH="${ENV_EXAMPLE_PATH:-${PROJECT_DIR}/packaging/systemd/adminui.env.example}"
SERVICE_TEMPLATE_PATH="${SERVICE_TEMPLATE_PATH:-${PROJECT_DIR}/packaging/systemd/aurora-adminui.service}"

INSTALL_ENV_PATH="${INSTALL_ETC_DIR}/adminui.env"
INSTALL_SERVICE_PATH="${SYSTEMD_UNIT_DIR}/${SERVICE_NAME}.service"
INSTALL_BIN_PATH="${INSTALL_BIN_DIR}/${BIN_NAME}"

SKIP_ENABLE="${SKIP_ENABLE:-false}"
SKIP_START="${SKIP_START:-false}"

log() {
  printf '[adminui-install] %s\n' "$*" >&2
}

fail() {
  printf '[adminui-install] error: %s\n' "$*" >&2
  exit 1
}

bool_true() {
  case "${1:-}" in
    1|true|TRUE|yes|YES|y|Y|on|ON) return 0 ;;
    *) return 1 ;;
  esac
}

usage() {
  cat >&2 <<'EOF'
Usage: ./install.sh [-e /path/to/envfile] [-h]

Options:
  -e, --env-file PATH   Use PATH as the source env file instead of the default example file.
  -h, --help            Show this help message.

Behavior:
  - The destination env file is always overwritten at /etc/aurora/adminui.env.
EOF
}

require_file() {
  local path="$1"
  [[ -f "$path" ]] || fail "required file not found: $path"
}

ensure_go() {
  command -v go >/dev/null 2>&1 || fail "go not found in PATH"
}

ensure_npm() {
  command -v npm >/dev/null 2>&1 || fail "npm not found in PATH"
}

ensure_sudo_access() {
  if [[ "$(id -u)" -eq 0 ]]; then
    return
  fi
  command -v sudo >/dev/null 2>&1 || fail "sudo is required for install steps"
  log "checking sudo access"
  sudo -v
}

sudo_cmd() {
  if [[ "$(id -u)" -eq 0 ]]; then
    "$@"
    return
  fi
  command -v sudo >/dev/null 2>&1 || fail "sudo is required for install steps"
  sudo "$@"
}

build_binary() {
  ensure_go
  ensure_npm

  require_file "${PROJECT_DIR}/package.json"
  require_file "${PROJECT_DIR}/go.mod"

  log "installing npm dependencies"
  npm ci --prefix "${PROJECT_DIR}" >&2

  log "building adminui dist"
  local build_env
  build_env="$(mktemp)"
  cp "${ENV_EXAMPLE_PATH}" "${build_env}"
  set -a
  # shellcheck disable=SC1090
  source "${build_env}"
  set +a
  rm -f "${build_env}"
  npm run build --prefix "${PROJECT_DIR}" >&2

  log "building ${BIN_NAME} binary"
  mkdir -p "${PROJECT_DIR}/bin"
  go build -o "${PROJECT_DIR}/bin/${BIN_NAME}" ./cmd/server >&2
  printf '%s\n' "${PROJECT_DIR}/bin/${BIN_NAME}"
}

install_files() {
  local built_binary="$1"

  require_file "${built_binary}"
  require_file "${ENV_EXAMPLE_PATH}"
  require_file "${SERVICE_TEMPLATE_PATH}"

  log "creating install directories"
  sudo_cmd install -d -m 0755 "${INSTALL_BIN_DIR}"
  sudo_cmd install -d -m 0755 "${INSTALL_ETC_DIR}"
  sudo_cmd install -d -m 0755 "${INSTALL_STATE_DIR}"
  sudo_cmd install -d -m 0755 "${SYSTEMD_UNIT_DIR}"

  log "installing binary to ${INSTALL_BIN_PATH}"
  sudo_cmd install -m 0755 "${built_binary}" "${INSTALL_BIN_PATH}"

  log "installing systemd unit to ${INSTALL_SERVICE_PATH}"
  sudo_cmd install -m 0644 "${SERVICE_TEMPLATE_PATH}" "${INSTALL_SERVICE_PATH}"

  log "installing env file to ${INSTALL_ENV_PATH}"
  sudo_cmd install -m 0644 "${ENV_EXAMPLE_PATH}" "${INSTALL_ENV_PATH}"

  sudo_cmd chown root:root "${INSTALL_BIN_PATH}" "${INSTALL_SERVICE_PATH}"
  sudo_cmd chown -R root:root "${INSTALL_ETC_DIR}" "${INSTALL_STATE_DIR}"
}

configure_systemd() {
  log "reloading systemd"
  sudo_cmd systemctl daemon-reload

  if ! bool_true "${SKIP_ENABLE}"; then
    log "enabling ${SERVICE_NAME}.service"
    sudo_cmd systemctl enable "${SERVICE_NAME}.service"
  fi

  if ! bool_true "${SKIP_START}"; then
    log "restarting ${SERVICE_NAME}.service"
    sudo_cmd systemctl restart "${SERVICE_NAME}.service"
    check_service_status
  else
    log "skip start requested"
  fi
}

check_service_status() {
  log "checking ${SERVICE_NAME}.service status"
  if sudo_cmd systemctl is-active --quiet "${SERVICE_NAME}.service"; then
    log "${SERVICE_NAME}.service is active"
    sudo_cmd systemctl --no-pager --full status "${SERVICE_NAME}.service" || true
    return
  fi

  log "${SERVICE_NAME}.service is not active"
  sudo_cmd systemctl --no-pager --full status "${SERVICE_NAME}.service" || true
  log "recent journal output"
  sudo_cmd journalctl -u "${SERVICE_NAME}.service" -n 50 --no-pager || true
  fail "${SERVICE_NAME}.service failed to become active"
}

main() {
  while [[ $# -gt 0 ]]; do
    case "$1" in
      -e|--env-file)
        [[ $# -ge 2 ]] || fail "missing value for $1"
        ENV_EXAMPLE_PATH="$2"
        shift 2
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        fail "unknown argument: $1"
        ;;
    esac
  done

  ensure_sudo_access

  local built_binary
  built_binary="$(build_binary)"
  install_files "${built_binary}"
  configure_systemd

  log "install completed"
  log "binary: ${INSTALL_BIN_PATH}"
  log "env: ${INSTALL_ENV_PATH}"
  log "unit: ${INSTALL_SERVICE_PATH}"

}

main "$@"
