import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import RouteTable from '../../../app/components/RouteTable.vue'

describe('RouteTable', () => {
  it('renders table headers and rows', () => {
    const wrapper = mount(RouteTable, {
      props: {
        items: [
          { id: 1, code: 'R-001', name: '上游线' },
          { id: 2, code: 'R-002', name: '下游线' },
        ],
      },
    })

    expect(wrapper.find('table.route-table').exists()).toBe(true)
    expect(wrapper.text()).toContain('编码')
    expect(wrapper.text()).toContain('R-001')
    expect(wrapper.text()).toContain('上游线')
    expect(wrapper.findAll('tbody tr')).toHaveLength(2)
  })

  it('renders empty tbody when no items', () => {
    const wrapper = mount(RouteTable, { props: { items: [] } })
    expect(wrapper.findAll('tbody tr')).toHaveLength(0)
  })
})
