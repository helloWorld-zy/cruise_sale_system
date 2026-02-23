import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import SearchPage from '../../app/pages/search/index.vue'

describe('Search Page', () => {
    it('renders filters', () => {
        const wrapper = mount(SearchPage)
        expect(wrapper.text()).toContain('Filter')
    })
})
