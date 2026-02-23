import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../pages/bookings/index.vue'

describe('Admin Bookings', () => {
    it('renders bookings title', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Bookings')
    })
})
