import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import OAuthAuthorizationFlow from '../OAuthAuthorizationFlow.vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copied: { value: false },
    copyToClipboard: vi.fn()
  })
}))

describe('OAuthAuthorizationFlow', () => {
  it('允许只填写 OpenAI 批量账号信息输入框时触发校验事件', async () => {
    const wrapper = mount(OAuthAuthorizationFlow, {
      props: {
        addMethod: 'oauth',
        platform: 'openai',
        showCookieOption: false,
        showRefreshTokenOption: true,
        allowMultiple: true
      },
      global: {
        stubs: {
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    await wrapper.get('input[value="refresh_token"]').setValue(true)

    const batchValue = 'StephanieMason5060@outlook.com----xVcgygWFCrqq----hfsqan87362----rt_RxtqD2_oTT1o5s4EjWjRp_XiuD799cJotZQ7uynLMw0.rTDMxyY5Qs_tEJOz1Tk7yli-ruihip3MTr8rQHGP5ao'
    await wrapper.get('textarea[placeholder="账号----密码----邮箱密码----RT"]').setValue(batchValue)
    await wrapper.get('button.btn.btn-primary').trigger('click')

    expect(wrapper.emitted('validate-refresh-token')).toEqual([[ '', batchValue ]])
  })
})
