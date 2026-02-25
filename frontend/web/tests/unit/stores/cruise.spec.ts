import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useCruiseStore } from '../../../app/stores/cruise'

describe('Cruise Store', () => {
    beforeEach(() => {
        setActivePinia(createPinia())
    })

    it('manages list and detail', () => {
        const store = useCruiseStore()
        expect(store.list).toEqual([])
        expect(store.detail).toBeNull()

        store.setList([{ id: 1 }])
        expect(store.list.length).toBe(1)

        store.setDetail({ id: 1 })
        expect(store.detail.id).toBe(1)
    })
})
