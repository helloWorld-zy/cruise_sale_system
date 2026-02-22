import { describe, test, expect } from 'vitest'
import { buildUrl } from '../../src/utils/request'

describe('request', () => {
    test('builds url', () => {
        expect(buildUrl('/cruises')).toContain('/api/v1/cruises')
    })
})
