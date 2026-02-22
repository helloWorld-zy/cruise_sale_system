import { defineStore } from 'pinia'

export const useCruiseStore = defineStore('cruise', {
    state: () => ({ list: [] as any[], detail: null as any }),
    actions: {
        setList(list: any[]) { this.list = list },
        setDetail(detail: any) { this.detail = detail },
    },
})
