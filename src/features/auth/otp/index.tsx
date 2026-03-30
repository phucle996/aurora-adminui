import { Link } from '@tanstack/react-router'
import { Logo } from '@/assets/logo'
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { OtpForm } from './components/otp-form'

export function Otp() {
  return (
    <div className='app-shell min-h-svh'>
      <div className='container grid min-h-svh max-w-none items-center py-8 lg:grid-cols-[1.15fr_0.85fr] lg:gap-12'>
        <div className='mx-auto flex w-full max-w-[540px] flex-col justify-center space-y-4 py-8'>
          <div className='mb-2 flex items-center gap-3'>
            <span className='flex size-12 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-[var(--shadow-card)]'>
              <Logo className='size-6' />
            </span>
          </div>
          <Card className='gap-4'>
            <CardHeader>
              <CardTitle className='text-base tracking-tight'>
                Two-factor Authentication
              </CardTitle>
              <CardDescription>
                Please enter the authentication code. <br /> We have sent the
                authentication code to your email.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <OtpForm />
            </CardContent>
            <CardFooter>
              <p className='px-8 text-center text-sm text-muted-foreground'>
                Haven&apos;t received it?{' '}
                <Link
                  to='/sign-in'
                  className='underline underline-offset-4 hover:text-primary'
                >
                  Resend a new code.
                </Link>
                .
              </p>
            </CardFooter>
          </Card>
        </div>
      </div>
    </div>
  )
}
