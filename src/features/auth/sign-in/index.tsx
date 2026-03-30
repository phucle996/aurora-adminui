import { useState } from 'react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { useNavigate, useSearch } from '@tanstack/react-router'
import {
  Copy,
  Loader2,
  LogIn,
  ShieldCheck,
} from 'lucide-react'
import { toast } from 'sonner'
import { cn } from '@/lib/utils'
import { useAuthStore } from '@/stores/auth-store'
import { Button } from '@/components/ui/button'
import { ThemeSwitch } from '@/components/theme-switch'
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { loginWithAdminTokenInit, verifyAdminToken2FA } from './api'

const tokenFormSchema = z.object({
  token: z.string().trim().min(1, 'Please enter your admin token'),
})

const mfaFormSchema = z.object({
  code: z.string().trim().min(1, 'Please enter your verification code'),
})

interface SignInFormProps extends React.HTMLAttributes<HTMLDivElement> {
  redirectTo?: string
}

function SignInForm({
  className,
  redirectTo,
  ...props
}: SignInFormProps) {
  const [isLoading, setIsLoading] = useState(false)
  const [revealedToken, setRevealedToken] = useState('')
  const [preauthSession, setPreauthSession] = useState('')
  const navigate = useNavigate()
  const { auth } = useAuthStore()

  const tokenForm = useForm<z.infer<typeof tokenFormSchema>>({
    resolver: zodResolver(tokenFormSchema),
    defaultValues: { token: '' },
  })

  const mfaForm = useForm<z.infer<typeof mfaFormSchema>>({
    resolver: zodResolver(mfaFormSchema),
    defaultValues: { code: '' },
  })

  function completeLogin(
    tokenType: 'bootstrap' | 'apitoken',
    apiToken?: string
  ) {
    auth.setUser({
      label: 'Aurora Operator',
      tokenType,
      bootstrapExchanged: !!apiToken,
      lastLoginAt: Date.now(),
    })
    auth.setSessionMarker(
      apiToken ? 'bootstrap-exchanged' : 'admin-session-authenticated'
    )

    if (apiToken) {
      setRevealedToken(apiToken)
      navigator.clipboard
        .writeText(apiToken)
        .then(() => {
          toast.success('New API token copied to clipboard')
        })
        .catch(() => {
          toast.success('New API token generated')
        })
      return
    }

    toast.success('Admin session created')
    navigate({ to: redirectTo || '/', replace: true })
  }

  async function onSubmitToken(data: z.infer<typeof tokenFormSchema>) {
    setIsLoading(true)
    try {
      const result = await loginWithAdminTokenInit(data.token)
      if (result.mfa_required && result.preauth_session) {
        setPreauthSession(result.preauth_session)
        toast.success('Second factor required')
        return
      }
      completeLogin(result.token_type || 'apitoken', result.api_token)
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Login failed')
    } finally {
      setIsLoading(false)
    }
  }

  async function onSubmitMFA(data: z.infer<typeof mfaFormSchema>) {
    if (!preauthSession) {
      toast.error('Missing MFA session')
      return
    }

    setIsLoading(true)
    try {
      const result = await verifyAdminToken2FA(preauthSession, data.code)
      completeLogin(result.token_type, result.api_token)
    } catch (error) {
      toast.error(
        error instanceof Error
          ? error.message
          : 'Second-factor verification failed'
      )
    } finally {
      setIsLoading(false)
    }
  }

  async function handleCopyToken() {
    if (!revealedToken) return
    await navigator.clipboard.writeText(revealedToken)
    toast.success('API token copied to clipboard')
  }

  function continueToConsole() {
    navigate({ to: redirectTo || '/', replace: true })
  }

  function resetMFAFlow() {
    setPreauthSession('')
    mfaForm.reset({ code: '' })
  }

  return (
    <div className={cn('grid gap-4', className)} {...props}>
      {!preauthSession ? (
        <Form {...tokenForm}>
          <form
            onSubmit={tokenForm.handleSubmit(onSubmitToken)}
            className='grid gap-4'
          >
            <FormField
              control={tokenForm.control}
              name='token'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Admin token</FormLabel>
                  <FormControl>
                    <Input
                      placeholder='Paste bootstrap token or active API token'
                      className='font-mono text-[13px]'
                      autoComplete='off'
                      autoCapitalize='off'
                      autoCorrect='off'
                      spellCheck={false}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <Button className='mt-1' disabled={isLoading}>
              {isLoading ? <Loader2 className='animate-spin' /> : <LogIn />}
              Continue with token
            </Button>
          </form>
        </Form>
      ) : (
        <Form {...mfaForm}>
          <form
            onSubmit={mfaForm.handleSubmit(onSubmitMFA)}
            className='grid gap-4'
          >
            <div className='space-y-1'>
              <p className='text-sm font-semibold text-foreground'>
                Two-factor verification
              </p>
              <p className='text-sm leading-6 text-muted-foreground'>
                Enter the current authenticator code to finish this admin sign-in.
              </p>
            </div>

            <FormField
              control={mfaForm.control}
              name='code'
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Authenticator code</FormLabel>
                  <FormControl>
                    <Input
                      placeholder='Enter the 6-digit TOTP code'
                      className='font-mono text-[13px]'
                      autoComplete='one-time-code'
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className='flex flex-wrap gap-3'>
              <Button disabled={isLoading}>
                {isLoading ? <Loader2 className='animate-spin' /> : <ShieldCheck />}
                Verify and sign in
              </Button>
              <Button
                type='button'
                variant='outline'
                onClick={resetMFAFlow}
                disabled={isLoading}
              >
                Back
              </Button>
            </div>
          </form>
        </Form>
      )}

      {revealedToken ? (
        <div className='rounded-2xl border border-success/25 bg-success-soft px-4 py-4'>
          <div className='flex items-start gap-3'>
            <div className='w-full space-y-3'>
              <div>
                <p className='text-sm font-semibold text-foreground'>
                  New admin API token
                </p>
                <p className='text-sm leading-6 text-muted-foreground'>
                  This token is only shown once after bootstrap exchange. Save it
                  securely before you continue.
                </p>
              </div>
              <div className='rounded-xl border border-border/80 bg-background px-3 py-3 font-mono text-xs break-all text-foreground'>
                {revealedToken}
              </div>
              <div className='flex flex-wrap gap-3'>
                <Button type='button' variant='outline' onClick={handleCopyToken}>
                  <Copy className='size-4' />
                  Copy token
                </Button>
                <Button type='button' onClick={continueToConsole}>
                  Continue to console
                </Button>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  )
}

export function SignIn() {
  const { redirect } = useSearch({ from: '/(auth)/sign-in' })

  return (
    <div className='app-shell flex min-h-svh items-center justify-center px-6 py-12'>
      <div className='fixed right-6 bottom-6 z-20'>
        <ThemeSwitch />
      </div>
      <div className='mx-auto flex w-full max-w-[900px] items-center justify-center'>
        <div className='flex w-full max-w-[760px] flex-col justify-center'>
          <Card className='gap-4 rounded-[30px] px-3 py-4'>
            <CardHeader>
              <CardTitle className='text-xl tracking-tight'>
                Admin token access
              </CardTitle>
            </CardHeader>
            <CardContent className='pb-8'>
              <SignInForm redirectTo={redirect} />
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
}

export function SignIn2() {
  return (
    <div className='relative container grid h-svh flex-col items-center justify-center lg:max-w-none lg:grid-cols-2 lg:px-0'>
      <div className='lg:p-8'>
        <div className='mx-auto flex w-full flex-col justify-center space-y-2 py-8 sm:w-[480px] sm:p-8'>
        </div>
        <div className='mx-auto flex w-full max-w-sm flex-col justify-center space-y-2'>
          <div className='flex flex-col space-y-2 text-start'>
            <h2 className='text-lg font-semibold tracking-tight'>
              Admin token sign-in
            </h2>
            <p className='text-sm text-muted-foreground'>
              Use the environment admin token to open a session. If 2FA is
              enabled, you will confirm with your authenticator app next.
            </p>
          </div>
          <SignInForm />
        </div>
      </div>

      <div
        className={cn(
          'relative h-full overflow-hidden bg-muted max-lg:hidden',
          'before:absolute before:inset-0 before:bg-[radial-gradient(circle_at_top_left,hsl(var(--primary)/0.18),transparent_38%),linear-gradient(160deg,hsl(var(--card)),hsl(var(--muted)))]'
        )}
      >
      </div>
    </div>
  )
}
