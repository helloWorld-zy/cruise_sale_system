import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import RoutesPage from '../../app/pages/routes/index.vue'

vi.stubGlobal('useApi', () => ({ request: vi.fn().mockResolvedValue({ data: [] }) }))

describe('Admin Routes List', () => {
    it('renders title', () => {
        const wrapper = mount(RoutesPage)
        expect(wrapper.text()).toContain('Routes')
    })
})
