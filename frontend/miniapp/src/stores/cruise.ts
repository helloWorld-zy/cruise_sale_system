// miniapp/src/stores/cruise.ts — 小程序端邮轮状态管理
// 使用 Pinia Setup API (Composition API) 风格管理邮轮数据

import { defineStore } from 'pinia'
import { ref } from 'vue'

/**
 * useCruiseStore 管理小程序端的邮轮浏览数据。
 * - list: 邮轮搜索结果列表
 * - detail: 当前查看的邮轮/舱房详情
 */
export const useCruiseStore = defineStore('cruise', () => {
    // 邮轮列表数据
    const list = ref<any[]>([])
    // 当前邮轮详情
    const detail = ref<any>(null)

    // setList 更新邮轮列表
    function setList(newList: any[]) {
        list.value = newList
    }

    // setDetail 更新邮轮详情
    function setDetail(newDetail: any) {
        detail.value = newDetail
    }

    return {
        list,
        detail,
        setList,
        setDetail
    }
})

