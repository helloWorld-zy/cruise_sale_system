import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import Page from '../../../../app/pages/cruises/[id].vue'

const mockRequest = vi.fn()
vi.stubGlobal('useApi', () => ({ request: mockRequest }))
vi.stubGlobal('useRoute', () => ({ params: { id: '9' } }))

beforeEach(() => {
  mockRequest.mockReset()
  mockRequest.mockImplementation((url: string) => {
    if (url === '/cruises/9') return Promise.resolve({ data: { id: 9, name: 'Atlantic Dream', tonnage: 140000, passenger_capacity: 3200, length: 315, deck_count: 16, images: [{ url: 'https://img/1.jpg' }, { url: 'https://img/2.jpg' }] } })
    if (url === '/cabin-types') return Promise.resolve({ data: { list: [{ id: 1, name: '阳台房', area_min: 24, max_capacity: 3, min_price_cents: 188000, amenities: '阳台,浴缸', description: '宽敞舒适' }] } })
    if (url === '/facilities') return Promise.resolve({ data: [{ id: 3, category_id: 1, name: '海上剧院', open_hours: '10:00-22:00', extra_charge: false }, { id: 4, category_id: 2, name: '付费SPA', open_hours: '09:00-21:00', extra_charge: true }] })
    if (url === '/facility-categories') return Promise.resolve({ data: [{ id: 1, name: '娱乐' }, { id: 2, name: '康养' }] })
    if (url === '/routes') return Promise.resolve({ data: [{ id: 5, name: '地中海之旅', departure_date: '2026-06-01', min_price_cents: 256000 }] })
    return Promise.resolve({ data: [] })
  })
})

describe('Web cruise detail page', () => {
  it('renders detail sections', async () => {
    const wrapper = mount(Page)
    await flushPromises()
    expect(wrapper.text()).toContain('Atlantic Dream')
    expect(wrapper.text()).toContain('舱房类型')
    expect(wrapper.text()).toContain('设施导览')
    expect(wrapper.text()).toContain('关联航线')
  })

  it('expands cabin type details when toggled', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const toggleBtn = wrapper.findAll('button').find((b) => b.text().includes('阳台房'))
    await toggleBtn.trigger('click')

    expect(wrapper.text()).toContain('宽敞舒适')
    expect(wrapper.text()).toContain('阳台')
  })

  it('filters facilities by selected category', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    expect(wrapper.text()).toContain('海上剧院')
    expect(wrapper.text()).toContain('付费SPA')

    const categoryBtn = wrapper.findAll('button').find((b) => b.text() === '康养')
    await categoryBtn!.trigger('click')

    expect(wrapper.text()).not.toContain('海上剧院')
    expect(wrapper.text()).toContain('付费SPA')
  })

  it('switches gallery slides with prev/next controls', async () => {
    const wrapper = mount(Page)
    await flushPromises()

    const nextBtn = wrapper.findAll('button').find((b) => b.text() === '›')
    const prevBtn = wrapper.findAll('button').find((b) => b.text() === '‹')
    const heroImg = wrapper.find('section img')

    expect(heroImg.attributes('src')).toContain('img/1.jpg')
    await nextBtn!.trigger('click')
    expect(heroImg.attributes('src')).toContain('img/2.jpg')
    await prevBtn!.trigger('click')
    expect(heroImg.attributes('src')).toContain('img/1.jpg')
  })
})
