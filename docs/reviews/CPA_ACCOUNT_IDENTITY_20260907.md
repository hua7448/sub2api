# CPA 账号身份对比与迁入记录（2026-09-07）

用户范围：参考最新 CPA，迁入与账号指纹、客户端身份、会话、请求头和 metadata 有关的改进。

## 对比基线

- 官方仓库：https://github.com/router-for-me/CLIProxyAPI
- 本次拉取和再次查询远端 HEAD 均为 `c76dfd4e0edabab9000628b1560ab8ab379eadb8`。
- 提交时间：`2026-09-06T15:27:15+08:00`；标题：`chore(codex): update codex user-agent to 0.153.3`。
- 本地只读参考目录：`/home/ubuntu/cpa-reference-20260907`。
- Sub2API 修改基线：`432c6e28b8c3a15a94827b87caa82e37f7a905af`，规范分支 `production/backend-current`。
- CPA 为 MIT 许可；本次根据其设计对 Sub2API 现有模块作独立实现，没有引入 CPA 运行时或依赖树。

固定版本参考：

- [HTTP 身份、会话缓存、嵌入 metadata](https://github.com/router-for-me/CLIProxyAPI/blob/c76dfd4e0edabab9000628b1560ab8ab379eadb8/internal/runtime/executor/codex_executor_request.go)
- [WebSocket 身份与请求头](https://github.com/router-for-me/CLIProxyAPI/blob/c76dfd4e0edabab9000628b1560ab8ab379eadb8/internal/runtime/executor/codex_websockets_request.go)
- [配置与身份开关](https://github.com/router-for-me/CLIProxyAPI/blob/c76dfd4e0edabab9000628b1560ab8ab379eadb8/internal/config/config_types.go)
- [MIT 许可](https://github.com/router-for-me/CLIProxyAPI/blob/c76dfd4e0edabab9000628b1560ab8ab379eadb8/LICENSE)

## 功能对比与结果

| CPA 相关能力 | Sub2API 状态及本次处理 |
|---|---|
| Codex UA、originator、客户端版本配对 | 现有共用身份解析器已覆盖 HTTP、WS、刷新、模型、额度和探测路径，保留现有账号 UA 优先级。 |
| CPA 当前 Mac/iTerm UA 配置 | 昨日版本已包含账号 UA 配置和测试。本次只读观察到运行时版本自动同步为 `0.153.4`；没有用 CPA 内置 `0.153.3` 覆盖运行时设置。 |
| 按账号隔离 installation/session/thread/turn/window | 已有持久随机种子与 off/device/session/full 四模式；没有替换现有账号种子或改变默认模式。 |
| 请求体会话作为会话头缺失时的后备 | 本次补齐 Responses HTTP、HTTP passthrough、WS 转发、WS native ingress、WS passthrough 首帧及 WS→HTTP bridge。按显式会话头、头内 metadata、体内会话、体内 metadata、显式 prompt_cache_key 的顺序提取。 |
| metadata 中的身份别名跟随账号身份 | 本次补齐 conversation_id/conversation-id，并修复可证明是会话别名的嵌入 prompt_cache_key。各载体先根据自己的原始会话归一化，再合并，避免跨载体后备值指向旧身份。 |
| 保留客户端提供的 opaque metadata | 在 CPA 保留原始字段的思路上补齐头/体缺失字段互补。各载体保留自己的显式值；sandbox、permissions、thread_source 等仅来自真实入站数据，不自动编造。使用 json.Number 保留大整数精度和原有 ASCII JSON 编码。 |
| WebSocket 每回合 metadata | 修复后续 WS 步骤用请求头整份覆盖 body metadata 的问题。每回合 body 快照独立，不能沿用上一回合的 body-only 字段。 |
| Beta features、Responses Lite、账号认证头 | 现有路径已有相应转发与身份校验，维持现有能力判定。 |
| 客户端会话、缓存与续链兼容 | 保留显式独立 root prompt_cache_key；保留 off/device 边界及 compact 不改写 body 的历史合约。 |
| 对响应文本进行整串身份替换 | 不迁入 CPA 的全局字符串替换方式：Sub2API 的身份处理限定在结构化身份字段，避免改写正常输出、工具参数或内容中的相同字符串。 |

## 验证与生产边界

新增回归覆盖：会话来源优先级；不同 body-only 会话得到不同且稳定的线程；HTTP/WS 真实请求构造；头体 metadata 互补；未知字段和大整数保真；默认缓存别名一致性；off/device/full 边界；不同回合 metadata 隔离。

历史测试未删除或降低断言。新测试通过新增生产合约条目执行。

这是源代码更新，尚未产生生产候选二进制或部署。UI、迁移、线上账号配置、重试策略均未修改。当前信任仓库继续记录已部署版本；后续发布必须按生产锁流程准备新候选及封存证据，不能直接从本次源码变更安装二进制。


验证记录：扩展身份回归通过 489 个顶层测试、1,012 个含子用例的测试项。生产后端合约通过 289 项检查、119 个指定 Go 回归。Race 检测尝试因环境中没有 C 编译器而无法运行；普通测试通过不等于已通过 race 检测。最终结果保存在 `/tmp/sub2api-cpa-identity-final-tests.jsonl`、`/tmp/sub2api-cpa-production-contract.json` 与 `/tmp/sub2api-cpa-production-contract.log`。
