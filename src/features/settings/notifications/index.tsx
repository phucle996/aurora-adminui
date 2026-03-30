import { ContentSection } from '../components/content-section'
import { NotificationsForm } from './notifications-form'

export function SettingsNotifications() {
  return (
    <ContentSection
      title='Notifications'
      desc='Configure Telegram delivery for platform alerts and operator notifications.'
    >
      <NotificationsForm />
    </ContentSection>
  )
}
