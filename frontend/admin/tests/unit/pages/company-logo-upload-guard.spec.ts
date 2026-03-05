import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

describe('Company logo upload guard', () => {
  it('uses local file input instead of Logo URL text input on company pages', () => {
    const files = [
      path.resolve(process.cwd(), 'app/pages/companies/index.vue'),
      path.resolve(process.cwd(), 'app/pages/companies/[id].vue'),
    ]

    for (const filePath of files) {
      const content = fs.readFileSync(filePath, 'utf-8')
      expect(content).toContain('type="file"')
      expect(content).not.toMatch(/Logo URL/)
    }
  })
})
