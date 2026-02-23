import { describe, it, expect } from 'vitest'
import { render } from '@testing-library/vue'
import Login from '../pages/login/login.vue'

describe('Miniapp Login', () => {
    it('shows login title', () => {
        const { getByText } = render(Login)
        expect(getByText('Login')).toBeTruthy()
    })
})
