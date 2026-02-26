import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../pages/dashboard/index.vue'

describe('Dashboard', () => {
  it('shows summary title', () => {
    const wrapper = mount(Page)
    expect(wrapper.text()).toContain('Dashboard')
  })
})
