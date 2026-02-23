import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import RoutesPage from '../../app/pages/routes/index.vue'

describe('Admin Routes List', () => {
    it('renders title', () => {
        const wrapper = mount(RoutesPage)
        expect(wrapper.text()).toContain('Routes')
    })
})
