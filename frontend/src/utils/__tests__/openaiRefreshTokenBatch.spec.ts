import {
  parseOpenAIRefreshTokenBatchInput,
  parseOpenAIRefreshTokenBatchLine
} from '../openaiRefreshTokenBatch'

describe('openaiRefreshTokenBatch', () => {
  it('parses account-password-emailPassword-rt format', () => {
    expect(
      parseOpenAIRefreshTokenBatchLine(
        'SeanHale4619@outlook.com----key----sevr674456++----rt_xxx.yyy'
      )
    ).toEqual({
      raw: 'SeanHale4619@outlook.com----key----sevr674456++----rt_xxx.yyy',
      account: 'SeanHale4619@outlook.com',
      password: 'key',
      emailPassword: 'sevr674456++',
      refreshToken: 'rt_xxx.yyy'
    })
  })

  it('treats a plain line as a refresh token', () => {
    expect(parseOpenAIRefreshTokenBatchLine('rt_plain_token')).toEqual({
      raw: 'rt_plain_token',
      account: '',
      password: '',
      emailPassword: '',
      refreshToken: 'rt_plain_token'
    })
  })

  it('ignores blank lines when parsing batch input', () => {
    expect(
      parseOpenAIRefreshTokenBatchInput('\nrt_one\n\nuser@example.com----pass----mailpass----rt_two\n')
    ).toEqual([
      {
        raw: 'rt_one',
        account: '',
        password: '',
        emailPassword: '',
        refreshToken: 'rt_one'
      },
      {
        raw: 'user@example.com----pass----mailpass----rt_two',
        account: 'user@example.com',
        password: 'pass',
        emailPassword: 'mailpass',
        refreshToken: 'rt_two'
      }
    ])
  })
})
