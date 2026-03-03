import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/vue'
import List from '../pages/cabin/list.vue'

describe('Cabin List', () => {
    it('shows title', () => {
        const { getByText } = render(List)
        expect(getByText('精选舱房')).toBeTruthy()
    })
})
