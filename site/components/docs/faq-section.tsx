import { Section } from './section'
import { FaqAccordion } from './faq-accordion'
import { FAQ } from '@/lib/docs-data'

export function FaqSection() {
  return (
    <Section
      id="faq"
      eyebrow="FAQ"
      title="常见问题"
      description="使用过程中遇到的高频问题与排查方式。"
    >
      <FaqAccordion items={FAQ} />
    </Section>
  )
}
