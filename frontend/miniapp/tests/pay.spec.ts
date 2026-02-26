import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/vue'
import Pay from '../pages/pay/pay.vue'

describe('Miniapp Pay', () => {
    it('shows pay label', () => {
        const { getByText } = render(Pay)
        expect(getByText('Pay Now')).toBeTruthy()
    })
})
