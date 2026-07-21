import { describe, expect, test } from 'bun:test'
import { configSchema } from './workflow'

describe('configSchema', () => {
  test('requires at least one EVM configuration', () => {
    const result = configSchema.safeParse({
      evms: [],
      dataSourceUrls: ['https://api.example1.com/v1', 'https://api.example2.com/v1'],
      aggregationMode: 'majority',
    })

    expect(result.success).toBe(false)
  })
})
