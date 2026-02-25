import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import Page from '../../../../app/pages/cabins/[id].vue'

describe('Cabin Detail', () => {
    it('renders detail view', () => {
        const wrapper = mount(Page)
        expect(wrapper.text()).toContain('Cabin Detail')
    })
})
