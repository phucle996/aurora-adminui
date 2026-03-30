import { ContentSection } from '../components/content-section'
import { AccountForm } from './account-form'

export function SettingsAccount() {
  return (
    <ContentSection
      title='Account'
      desc='Manage the single admin login surface and its authenticator-based protection.'
    >
      <AccountForm />
    </ContentSection>
  )
}
