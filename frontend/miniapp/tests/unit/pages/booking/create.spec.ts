import { afterEach, beforeEach, describe, it, expect, vi } from 'vitest'
import { cleanup, fireEvent, render } from '@testing-library/vue'
import Page from '../../../../pages/booking/create.vue'

const mockRequest = vi.fn().mockResolvedValue({ data: { status: 'created' } })

vi.mock('../../../../src/utils/request', () => ({
    request: (...args: any[]) => mockRequest(...args)
}))

beforeEach(() => {
    mockRequest.mockClear()
})

afterEach(() => {
    cleanup()
})

describe('Booking Create', () => {
    it('renders view', () => {
        const { getByText } = render(Page)
        expect(getByText('Create Booking')).toBeTruthy()
    })

    it('submits booking when form is valid', async () => {
        const { getByPlaceholderText, getByText } = render(Page)

        await fireEvent.update(getByPlaceholderText('Voyage ID'), '1001')
        await fireEvent.update(getByPlaceholderText('Cabin SKU ID'), '2002')
        await fireEvent.update(getByPlaceholderText('Guests'), '2')
        await fireEvent.click(getByText('Submit Booking'))

        expect(mockRequest).toHaveBeenCalledWith('/bookings', expect.objectContaining({
            method: 'POST',
            data: { voyage_id: 1001, cabin_sku_id: 2002, guests: 2 },
        }))
    })

    it('does not submit when form is invalid', async () => {
        const { getByText } = render(Page)
        await fireEvent.click(getByText('Submit Booking'))
        expect(mockRequest).not.toHaveBeenCalled()
    })

    it('shows error message when request fails', async () => {
        mockRequest.mockRejectedValueOnce(new Error('network failed'))
        const { getByPlaceholderText, getByText, findByText } = render(Page)

        await fireEvent.update(getByPlaceholderText('Voyage ID'), '1001')
        await fireEvent.update(getByPlaceholderText('Cabin SKU ID'), '2002')
        await fireEvent.update(getByPlaceholderText('Guests'), '2')
        await fireEvent.click(getByText('Submit Booking'))

        expect(await findByText('network failed')).toBeTruthy()
    })
})
