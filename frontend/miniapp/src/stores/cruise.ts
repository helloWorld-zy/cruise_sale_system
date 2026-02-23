// miniapp/src/stores/cruise.ts — 小程序端邮轮状态管理
// 使用 Pinia Options API 风格管理邮轮数据

import { defineStore } from 'pinia'

/**
 * useCruiseStore 管理小程序端的邮轮浏览数据。
 * - list: 邮轮搜索结果列表
 * - detail: 当前查看的邮轮/舱房详情
 */
export const useCruiseStore = defineStore('cruise', {
    state: () => ({
        list: [] as any[],    // 邮轮列表数据
        detail: null as any   // 当前邮轮详情
    }),
    actions: {
        // setList 更新邮轮列表
        setList(list: any[]) { this.list = list },
        // setDetail 更新邮轮详情
        setDetail(detail: any) { this.detail = detail },
    },
})
