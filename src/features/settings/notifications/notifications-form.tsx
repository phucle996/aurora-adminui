import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { showSubmittedData } from '@/lib/show-submitted-data'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'

const notificationsFormSchema = z.object({
  delivery: z.enum(['all', 'critical', 'disabled'], {
    error: (iss) =>
      iss.input === undefined
        ? 'Please select a Telegram delivery mode.'
        : undefined,
  }),
  bot_token: z.string(),
  chat_id: z.string(),
  security_alerts: z.boolean(),
  infrastructure_alerts: z.boolean(),
  platform_events: z.boolean(),
})

type NotificationsFormValues = z.infer<typeof notificationsFormSchema>

// This can come from your database or API.
const defaultValues: Partial<NotificationsFormValues> = {
  delivery: 'critical',
  bot_token: '',
  chat_id: '',
  security_alerts: true,
  infrastructure_alerts: true,
  platform_events: false,
}

export function NotificationsForm() {
  const form = useForm<NotificationsFormValues>({
    resolver: zodResolver(notificationsFormSchema),
    defaultValues,
  })

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit((data) => showSubmittedData(data))}
        className='space-y-8'
      >
        <FormField
          control={form.control}
          name='delivery'
          render={({ field }) => (
            <FormItem className='relative space-y-3'>
              <FormLabel>Telegram delivery mode</FormLabel>
              <FormControl>
                <RadioGroup
                  onValueChange={field.onChange}
                  defaultValue={field.value}
                  className='flex flex-col gap-2'
                >
                  <FormItem className='flex items-center'>
                    <FormControl>
                      <RadioGroupItem value='all' />
                    </FormControl>
                    <FormLabel className='font-normal'>
                      Deliver all platform notifications
                    </FormLabel>
                  </FormItem>
                  <FormItem className='flex items-center'>
                    <FormControl>
                      <RadioGroupItem value='critical' />
                    </FormControl>
                    <FormLabel className='font-normal'>
                      Deliver only critical alerts
                    </FormLabel>
                  </FormItem>
                  <FormItem className='flex items-center'>
                    <FormControl>
                      <RadioGroupItem value='disabled' />
                    </FormControl>
                    <FormLabel className='font-normal'>
                      Disable Telegram delivery
                    </FormLabel>
                  </FormItem>
                </RadioGroup>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <div className='grid gap-4 md:grid-cols-2'>
          <FormField
            control={form.control}
            name='bot_token'
            render={({ field }) => (
              <FormItem>
                <FormLabel>Telegram bot token</FormLabel>
                <FormControl>
                  <Input
                    placeholder='123456:AA...'
                    className='font-mono text-[13px]'
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Bot token used by the platform to send notifications into Telegram.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name='chat_id'
            render={({ field }) => (
              <FormItem>
                <FormLabel>Telegram chat ID</FormLabel>
                <FormControl>
                  <Input
                    placeholder='-1001234567890'
                    className='font-mono text-[13px]'
                    {...field}
                  />
                </FormControl>
                <FormDescription>
                  Destination chat, group, or channel ID for operator-facing alerts.
                </FormDescription>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <div className='relative'>
          <h3 className='mb-4 text-lg font-medium'>Telegram event categories</h3>
          <div className='space-y-4'>
            <FormField
              control={form.control}
              name='security_alerts'
              render={({ field }) => (
                <FormItem className='flex flex-row items-center justify-between rounded-lg border p-4'>
                  <div className='space-y-0.5'>
                    <FormLabel className='text-base'>Security alerts</FormLabel>
                    <FormDescription>
                      Send Telegram messages for admin auth failures, token rotation, and access anomalies.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='infrastructure_alerts'
              render={({ field }) => (
                <FormItem className='flex flex-row items-center justify-between rounded-lg border p-4'>
                  <div className='space-y-0.5'>
                    <FormLabel className='text-base'>
                      Infrastructure alerts
                    </FormLabel>
                    <FormDescription>
                      Send Telegram notifications for node health issues, cluster degradation, and platform substrate failures.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name='platform_events'
              render={({ field }) => (
                <FormItem className='flex flex-row items-center justify-between rounded-lg border p-4'>
                  <div className='space-y-0.5'>
                    <FormLabel className='text-base'>Platform events</FormLabel>
                    <FormDescription>
                      Send non-critical operator messages such as resource lifecycle events and provisioning completions.
                    </FormDescription>
                  </div>
                  <FormControl>
                    <Switch
                      checked={field.value}
                      onCheckedChange={field.onChange}
                    />
                  </FormControl>
                </FormItem>
              )}
            />
          </div>
        </div>
        <Button type='submit'>Update Telegram notifications</Button>
      </form>
    </Form>
  )
}
