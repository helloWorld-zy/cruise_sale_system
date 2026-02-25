import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'

import CabinsIndexPage from '../../../app/pages/cabins/index.vue'
import CabinsInventoryPage from '../../../app/pages/cabins/inventory.vue'
import RouteNewPage from '../../../app/pages/routes/new.vue'
import VoyagesIndexPage from '../../../app/pages/voyages/index.vue'
import VoyagesNewPage from '../../../app/pages/voyages/new.vue'

describe('Sprint2 admin pages', () => {
    test('renders cabin pages', () => {
        const cabins = mount(CabinsIndexPage)
        expect(cabins.text()).toContain('Cabins')

        const inventory = mount(CabinsInventoryPage)
        expect(inventory.text()).toContain('Inventory')
    })

    test('renders route new form', async () => {
        const wrapper = mount(RouteNewPage)
        expect(wrapper.text()).toContain('New Route')

        const inputs = wrapper.findAll('input')
        expect(inputs).toHaveLength(2)
        await inputs[0]!.setValue('R-001')
        await inputs[1]!.setValue('Asia Route')
        expect((inputs[0]!.element as HTMLInputElement).value).toBe('R-001')
        expect((inputs[1]!.element as HTMLInputElement).value).toBe('Asia Route')
    })

    test('renders voyage pages', async () => {
        const list = mount(VoyagesIndexPage)
        expect(list.text()).toContain('Voyages')
        expect(list.findAll('tbody tr').length).toBeGreaterThan(0)

        const create = mount(VoyagesNewPage)
        expect(create.text()).toContain('New Voyage')
        const inputs = create.findAll('input')
        await inputs[0]!.setValue('V-001')
        await inputs[1]!.setValue('Voyage 1')
        expect((inputs[0]!.element as HTMLInputElement).value).toBe('V-001')
        expect((inputs[1]!.element as HTMLInputElement).value).toBe('Voyage 1')
    })
})
