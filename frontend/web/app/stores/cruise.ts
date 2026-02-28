// web/app/stores/cruise.ts — 前台邮轮状态管理
// 使用 Pinia Composition API 风格管理邮轮列表和详情数据

import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * useCruiseStore 管理前台用户浏览的邮轮数据。
 * - list: 邮轮搜索结果列表
 * - detail: 当前查看的邮轮详情
 */
export const useCruiseStore = defineStore('cruise', () => {
    const list = ref<any[]>([]) // 邮轮列表数据
    const detail = ref<any>(null) // 当前邮轮详情

    // setList 更新邮轮列表（搜索结果）。
    function setList(newList: any[]) {
        list.value = newList
    }

    // setDetail 更新当前邮轮详情。
    function setDetail(newDetail: any) {
        detail.value = newDetail
    }

    return { list, detail, setList, setDetail }
})
