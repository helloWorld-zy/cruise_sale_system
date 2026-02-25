import { describe, test, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BookingRow from '../../../app/components/BookingRow.vue'

describe('BookingRow', () => {
    test('renders booking summary text', () => {
        const wrapper = mount(BookingRow, {
            props: {
                booking: { id: 7, status: 'created', total: 12345 },
            },
        })

        expect(wrapper.text()).toContain('7')
        expect(wrapper.text()).toContain('created')
        expect(wrapper.text()).toContain('12345')
    })
})
