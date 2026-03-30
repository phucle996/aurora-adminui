import { ContentSection } from '../components/content-section'
import { AppearanceForm } from './appearance-form'

export function SettingsAppearance() {
  return (
    <ContentSection
      title='Appearance'
      desc='Tune the visual presentation of the operator console.'
    >
      <AppearanceForm />
    </ContentSection>
  )
}
