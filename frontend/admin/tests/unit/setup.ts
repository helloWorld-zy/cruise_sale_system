import { config } from '@vue/test-utils'

config.global.stubs = {
  ...(config.global.stubs || {}),
  NuxtLink: { props: ['to'], template: '<a :href="String(to)"><slot /></a>' },
  AdminActionLink: { props: ['to'], template: '<a :href="String(to)"><slot /></a>' },
  AdminPageHeader: { props: ['title', 'subtitle'], template: '<div>{{ title }} {{ subtitle }}<slot /><slot name="actions" /></div>' },
  AdminFilterBar: { template: '<div><slot /></div>' },
  AdminDataCard: { props: ['flush'], template: '<div><slot /></div>' },
  AdminFormCard: { props: ['title'], template: '<div><slot /></div>' },
  AdminActionBar: { template: '<div><slot /></div>' },
  AdminStatusTag: { props: ['text'], template: '<span>{{ text }}</span>' },
}
