import { describe, expect, test } from 'bun:test'
import { fetchGameResult } from './workflow'

const makeSendRequester = (body: unknown, statusCode = 200) =>
  ({
    sendRequest: () => ({
      result: () => ({
        statusCode,
        body: Buffer.from(JSON.stringify(body)),
      }),
    }),
  }) as any

const espnBody = () => ({
  status: { type: { completed: true } },
  competitions: [
    {
      competitors: [
        { homeAway: 'home', score: '123' },
        { homeAway: 'away', score: '107' },
      ],
    },
  ],
})

describe('fetchGameResult', () => {
  test('parses ESPN home and away scores', () => {
    const result = fetchGameResult(makeSendRequester(espnBody()), {
      url: 'https://espn.example/scoreboard',
      gameId: '401766123',
    })

    expect(result).toEqual({ homeScore: 123, awayScore: 107 })
  })

  test('throws when the sports API returns a non-200 response', () => {
    expect(() =>
      fetchGameResult(makeSendRequester({ error: 'not found' }, 404), {
        url: 'https://espn.example/scoreboard',
        gameId: '401766123',
      }),
    ).toThrow('ESPN returned HTTP 404')
  })

  test('throws when the ESPN competitor structure is missing', () => {
    expect(() =>
      fetchGameResult(makeSendRequester({ competitions: [] }), {
        url: 'https://espn.example/scoreboard',
        gameId: '401766123',
      }),
    ).toThrow('Unexpected competitor structure for game 401766123')
  })
})
