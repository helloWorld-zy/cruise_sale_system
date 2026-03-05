import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

function collectVueFiles(rootDir: string): string[] {
  const result: string[] = []
  const entries = fs.readdirSync(rootDir, { withFileTypes: true })
  for (const entry of entries) {
    const fullPath = path.join(rootDir, entry.name)
    if (entry.isDirectory()) {
      result.push(...collectVueFiles(fullPath))
      continue
    }
    if (entry.isFile() && fullPath.endsWith('.vue')) {
      result.push(fullPath)
    }
  }
  return result
}

describe('Admin page action link guard', () => {
  it('does not use raw NuxtLink for 编辑/取消 actions', () => {
    const pagesRoot = path.resolve(process.cwd(), 'app/pages')
    const files = collectVueFiles(pagesRoot)
    const actionLinkPattern = /<NuxtLink\b[\s\S]*?>\s*(编辑|取消)\s*<\/NuxtLink>/g

    const violations: string[] = []
    for (const filePath of files) {
      const content = fs.readFileSync(filePath, 'utf-8')
      if (actionLinkPattern.test(content)) {
        violations.push(path.relative(process.cwd(), filePath))
      }
    }

    expect(
      violations,
      `Use AdminActionLink for 编辑/取消 actions. Violations: ${violations.join(', ')}`,
    ).toEqual([])
  })
})
