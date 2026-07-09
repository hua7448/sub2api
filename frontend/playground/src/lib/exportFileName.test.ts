import { describe, expect, it } from 'vitest'
import { getNumberedFileNameBase, sanitizeFileNamePart } from './exportFileName'

describe('exportFileName', () => {
  it('adds a padded index only when there are multiple files', () => {
    expect(getNumberedFileNameBase('task-abc', 0, 1)).toBe('task-abc')
    expect(getNumberedFileNameBase('task-abc', 0, 2)).toBe('task-abc-01')
    expect(getNumberedFileNameBase('task-abc', 1, 2)).toBe('task-abc-02')
  })

  it('sanitizes unsafe file name characters', () => {
    expect(sanitizeFileNamePart(' task:abc/01\n ')).toBe('task-abc-01')
  })
})
