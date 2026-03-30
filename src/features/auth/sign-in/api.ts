type ApiResponse<T> = {
  message: string
  data: T
}

export type AdminTokenLoginInitData = {
  mfa_required?: boolean
  preauth_session?: string
  preauth_ttl_seconds?: number
  token_type?: 'bootstrap' | 'apitoken'
  api_token?: string
}

export type AdminTokenLoginVerifyData = {
  token_type: 'bootstrap' | 'apitoken'
  api_token?: string
}

export type AdminSessionStatusData = {
  authenticated: boolean
  session_id: string
  session_expires_at: string
  two_factor_enabled: boolean
}

export type AdminTwoFactorStatusData = {
  two_factor_enabled: boolean
  totp_enabled_at?: string | null
}

export type AdminTOTPSetupBeginData = {
  setup_session: string
  secret: string
  otpauth_url: string
  setup_ttl_seconds: number
}

async function parseResponse<T>(response: Response): Promise<T> {
  const payload = (await response.json().catch(() => null)) as
    | ApiResponse<T>
    | null
    | { message?: string; error?: string }

  if (!response.ok) {
    const message =
      payload && typeof payload.message === 'string'
        ? payload.message
        : 'Request failed'
    throw new Error(message)
  }

  if (!payload || !('data' in payload) || !payload.data) {
    throw new Error('Invalid server response')
  }

  return payload.data
}

export async function loginWithAdminTokenInit(
  token: string
): Promise<AdminTokenLoginInitData> {
  const response = await fetch('/api/v1/admin/auth/token-login/init', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ token }),
  })
  return parseResponse<AdminTokenLoginInitData>(response)
}

export async function verifyAdminToken2FA(
  preauthSession: string,
  code: string
): Promise<AdminTokenLoginVerifyData> {
  const response = await fetch('/api/v1/admin/auth/token-login/verify-2fa', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      preauth_session: preauthSession,
      code,
    }),
  })
  return parseResponse<AdminTokenLoginVerifyData>(response)
}

export async function getAdminSessionStatus(): Promise<AdminSessionStatusData> {
  const response = await fetch('/api/v1/admin/auth/session', {
    credentials: 'include',
  })
  return parseResponse<AdminSessionStatusData>(response)
}

export async function logoutAdminSession(): Promise<void> {
  const response = await fetch('/api/v1/admin/auth/logout', {
    method: 'POST',
    credentials: 'include',
  })
  if (!response.ok) {
    throw new Error('Failed to sign out')
  }
}

export async function getAdminTwoFactorStatus(): Promise<AdminTwoFactorStatusData> {
  const response = await fetch('/api/v1/admin/2fa/status', {
    credentials: 'include',
  })
  return parseResponse<AdminTwoFactorStatusData>(response)
}

export async function beginAdminTOTPSetup(): Promise<AdminTOTPSetupBeginData> {
  const response = await fetch('/api/v1/admin/2fa/totp/setup/begin', {
    method: 'POST',
    credentials: 'include',
  })
  return parseResponse<AdminTOTPSetupBeginData>(response)
}

export async function confirmAdminTOTPSetup(
  setupSession: string,
  code: string
): Promise<void> {
  const response = await fetch('/api/v1/admin/2fa/totp/setup/confirm', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({
      setup_session: setupSession,
      code,
    }),
  })
  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as
      | { message?: string }
      | null
    throw new Error(payload?.message || 'Failed to enable 2FA')
  }
}

export async function disableAdminTwoFactor(code: string): Promise<void> {
  const response = await fetch('/api/v1/admin/2fa/disable', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include',
    body: JSON.stringify({ code }),
  })
  if (!response.ok) {
    const payload = (await response.json().catch(() => null)) as
      | { message?: string }
      | null
    throw new Error(payload?.message || 'Failed to disable 2FA')
  }
}
