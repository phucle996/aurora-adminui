import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Shield, ShieldCheck, ShieldOff } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import {
  beginAdminTOTPSetup,
  confirmAdminTOTPSetup,
  disableAdminTwoFactor,
  getAdminTwoFactorStatus,
} from '@/features/auth/sign-in/api'

export function AccountForm() {
  const queryClient = useQueryClient()
  const [setupSession, setSetupSession] = useState('')
  const [setupSecret, setSetupSecret] = useState('')
  const [setupOTPAuthURL, setSetupOTPAuthURL] = useState('')
  const [setupCode, setSetupCode] = useState('')
  const [verificationCode, setVerificationCode] = useState('')

  const statusQuery = useQuery({
    queryKey: ['admin-2fa-status'],
    queryFn: getAdminTwoFactorStatus,
  })

  const beginSetupMutation = useMutation({
    mutationFn: beginAdminTOTPSetup,
    onSuccess: (data) => {
      setSetupSession(data.setup_session)
      setSetupSecret(data.secret)
      setSetupOTPAuthURL(data.otpauth_url)
      toast.success('Authenticator setup ready')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to begin setup')
    },
  })

  const confirmSetupMutation = useMutation({
    mutationFn: () => confirmAdminTOTPSetup(setupSession, setupCode),
    onSuccess: () => {
      setSetupSession('')
      setSetupSecret('')
      setSetupOTPAuthURL('')
      setSetupCode('')
      queryClient.invalidateQueries({ queryKey: ['admin-2fa-status'] })
      toast.success('Admin 2FA enabled')
    },
    onError: (error) => {
      toast.error(
        error instanceof Error ? error.message : 'Failed to confirm authenticator setup'
      )
    },
  })

  const disableMutation = useMutation({
    mutationFn: () => disableAdminTwoFactor(verificationCode),
    onSuccess: () => {
      setVerificationCode('')
      queryClient.invalidateQueries({ queryKey: ['admin-2fa-status'] })
      toast.success('Admin 2FA disabled')
    },
    onError: (error) => {
      toast.error(error instanceof Error ? error.message : 'Failed to disable 2FA')
    },
  })

  const isEnabled = !!statusQuery.data?.two_factor_enabled

  return (
    <div className='space-y-6'>
      <div className='rounded-2xl border border-border/80 bg-card px-5 py-5 shadow-xs'>
        <div className='flex flex-wrap items-start justify-between gap-4'>
          <div className='space-y-2'>
            <div className='flex items-center gap-3'>
              <span className='flex size-11 items-center justify-center rounded-2xl bg-accent text-accent-foreground'>
                {isEnabled ? (
                  <ShieldCheck className='size-5' />
                ) : (
                  <Shield className='size-5' />
                )}
              </span>
              <div>
                <h3 className='text-base font-semibold text-foreground'>
                  Admin two-factor authentication
                </h3>
                <p className='text-sm leading-6 text-muted-foreground'>
                  Protect the single admin token login with an authenticator
                  app.
                </p>
              </div>
            </div>
          </div>
          <Badge variant={isEnabled ? 'default' : 'outline'}>
            {isEnabled ? '2FA enabled' : '2FA disabled'}
          </Badge>
        </div>
      </div>

      {!isEnabled ? (
        <div className='rounded-2xl border border-border/80 bg-card px-5 py-5 shadow-xs'>
          <div className='space-y-4'>
            <div>
              <h4 className='text-sm font-semibold text-foreground'>
                Enable admin 2FA
              </h4>
              <p className='mt-1 text-sm leading-6 text-muted-foreground'>
                Start by generating a TOTP secret, then verify the first code
                from your authenticator app.
              </p>
            </div>

            {!setupSession ? (
              <Button
                type='button'
                onClick={() => beginSetupMutation.mutate()}
                disabled={beginSetupMutation.isPending}
              >
                Begin authenticator setup
              </Button>
            ) : (
              <div className='space-y-4'>
                <div className='rounded-xl border border-border/80 bg-muted/40 px-4 py-4'>
                  <p className='text-sm font-medium text-foreground'>Secret</p>
                  <p className='mt-2 font-mono text-sm break-all text-foreground'>
                    {setupSecret}
                  </p>
                  <p className='mt-4 text-sm font-medium text-foreground'>
                    OTPAuth URL
                  </p>
                  <p className='mt-2 font-mono text-xs break-all text-muted-foreground'>
                    {setupOTPAuthURL}
                  </p>
                </div>
                <div className='space-y-2'>
                  <label className='text-sm font-medium text-foreground'>
                    First authenticator code
                  </label>
                  <Input
                    value={setupCode}
                    onChange={(event) => setSetupCode(event.target.value)}
                    placeholder='Enter the 6-digit code'
                    className='font-mono text-[13px]'
                  />
                </div>
                <div className='flex flex-wrap gap-3'>
                  <Button
                    type='button'
                    onClick={() => confirmSetupMutation.mutate()}
                    disabled={confirmSetupMutation.isPending || !setupCode.trim()}
                  >
                    Confirm and enable 2FA
                  </Button>
                  <Button
                    type='button'
                    variant='outline'
                    onClick={() => {
                      setSetupSession('')
                      setSetupSecret('')
                      setSetupOTPAuthURL('')
                      setSetupCode('')
                    }}
                  >
                    Cancel
                  </Button>
                </div>
              </div>
            )}
          </div>
        </div>
      ) : (
        <div className='rounded-2xl border border-border/80 bg-card px-5 py-5 shadow-xs'>
          <div className='space-y-4'>
            <div>
              <h4 className='text-sm font-semibold text-foreground'>
                Manage admin 2FA
              </h4>
              <p className='mt-1 text-sm leading-6 text-muted-foreground'>
                Use a current TOTP code from your authenticator app to disable
                2FA.
              </p>
            </div>

            <div className='grid gap-4'>
              <Input
                value={verificationCode}
                onChange={(event) => setVerificationCode(event.target.value)}
                placeholder='Enter the 6-digit authenticator code'
                className='font-mono text-[13px]'
              />
            </div>

            <div className='flex flex-wrap gap-3'>
              <Button
                type='button'
                variant='destructive'
                onClick={() => disableMutation.mutate()}
                disabled={disableMutation.isPending || !verificationCode.trim()}
              >
                <ShieldOff className='size-4' />
                Disable 2FA
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
