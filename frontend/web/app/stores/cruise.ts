// web/app/stores/cruise.ts — 前台邮轮状态管理
// 使用 Pinia Options API 风格管理邮轮列表和详情数据

import { defineStore } from 'pinia'

/**
 * useCruiseStore 管理前台用户浏览的邮轮数据。
 * - list: 邮轮搜索结果列表
 * - detail: 当前查看的邮轮详情
 */
export const useCruiseStore = defineStore('cruise', {
    state: () => ({
        list: [] as any[],        // 邮轮列表数据
        detail: null as any       // 当前邮轮详情
    }),
    actions: {
        // setList 更新邮轮列表（搜索结果）
        setList(list: any[]) { this.list = list },
        // setDetail 更新当前邮轮详情
        setDetail(detail: any) { this.detail = detail },
    },
})
