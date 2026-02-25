import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/vue'
import Page from '../../../../pages/cabin/detail.vue'

describe('Cabin Detail', () => {
    it('renders view', () => {
        const { getByText } = render(Page)
        expect(getByText('Cabin Detail')).toBeTruthy()
    })
})
