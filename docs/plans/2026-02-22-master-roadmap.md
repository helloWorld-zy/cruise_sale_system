# CruiseBooking 邮轮舱位预订平台 — 总路线图

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement plans referenced in this document.

**Goal:** 从零构建邮轮舱位预订全流程电商平台，覆盖邮轮介绍、舱位浏览、在线预订、在线支付、订单管理、售后全链路。

**技术栈:** Go (Gin+GORM+PostgreSQL) 后端 + Nuxt 4 管理后台 + Nuxt 4 Web前台 + uni-app 微信小程序

**测试策略:** 100% TDD — 每个功能先写失败测试，再写最小实现，再重构。全部代码覆盖率 100%。

---

## 总览表

| Phase | Sprint | 主题 | 核心交付物 | 依赖 |
|-------|--------|------|-----------|------|
| **Phase 1: MVP** | Sprint 1 | 基础设施 + 邮轮介绍模块 | Monorepo脚手架、Docker环境、邮轮/舱房类型/设施CRUD（后台+前台）、JWT+RBAC认证、CI/CD | 无 |
| | Sprint 2 | 舱位商品管理 | 航线管理、航次管理、舱位管理、定价矩阵、库存管理（后台CRUD + 前台浏览） | Sprint 1 |
| | Sprint 3 | 预订流程 + 用户系统 | 用户注册登录（微信/手机号）、舱位浏览筛选、在线下单、乘客信息管理 | Sprint 2 |
| | Sprint 4 | 订单支付 + 通知 + 统计 | 微信/支付宝支付、退改处理、财务对账、消息通知、数据看板 | Sprint 3 |
| **Phase 2: V1.0** | Sprint 5 | 智能发现与决策支持 | AI智能推荐、航线日历视图、价格趋势分析、舱位对比工具 | Sprint 4 |
| | Sprint 6 | 预订流程增强 | 舱位锁定倒计时、分期付款、多币种、历史乘客快填、证件OCR | Sprint 4 |
| | Sprint 7 | 团队预订 + 出行服务 | 批量预订入口、团队定价、出团通知书、出行倒计时 | Sprint 4 |
| | Sprint 8 | 实时信息 + 在线客服 | WebSocket实时舱位、价格变动提醒、预订进度追踪、AI客服、天气港口信息 | Sprint 4 |
| **Phase 3: V1.5** | Sprint 9 | 社交分享 | 行程海报生成、邀请同行优惠 | Sprint 4 |
| | Sprint 10 | 评价 + 游记社区 | 用户评价系统、UGC游记社区、拼团出行 | Sprint 4 |
| | Sprint 11 | 智能运营（上） | 动态定价引擎、渠道库存分配、智能预警系统 | Sprint 4 |
| | Sprint 12 | 智能运营（下） | 收益管理看板、客户生命周期管理、自动化营销引擎 | Sprint 11 |
| **Phase 4: V2.0** | Sprint 13 | 价格对比 + 多语言 | 客户价格对比工具、中英文/繁体多语言、翻译管理 | Sprint 4 |
| | Sprint 14 | 性能优化 + 压测 | 数据库优化、缓存策略、CDN、全链路压测 | Sprint 8 |
| | Sprint 15 | 安全审计 + 渗透测试 | OWASP安全检查、渗透测试修复、数据加密审计 | Sprint 14 |
| | Sprint 16 | 上线部署 + 监控 | K8s部署、Prometheus+Grafana监控、Loki日志、运维手册 | Sprint 15 |

---

## 文件索引

| Sprint | 计划文件 |
|--------|---------|
| Sprint 1 | `docs/plans/2026-02-22-sprint01.md` |
| Sprint 2 | `docs/plans/2026-02-22-sprint02.md` |
| Sprint 3 | `docs/plans/2026-02-22-sprint03.md` |
| Sprint 4 | `docs/plans/2026-02-22-sprint04.md` |
| Sprint 5 | `docs/plans/2026-02-22-sprint05.md` |
| Sprint 6 | `docs/plans/2026-02-22-sprint06.md` |
| Sprint 7 | `docs/plans/2026-02-22-sprint07.md` |
| Sprint 8 | `docs/plans/2026-02-22-sprint08.md` |
| Sprint 9 | `docs/plans/2026-02-22-sprint09.md` |
| Sprint 10 | `docs/plans/2026-02-22-sprint10.md` |
| Sprint 11 | `docs/plans/2026-02-22-sprint11.md` |
| Sprint 12 | `docs/plans/2026-02-22-sprint12.md` |
| Sprint 13 | `docs/plans/2026-02-22-sprint13.md` |
| Sprint 14 | `docs/plans/2026-02-22-sprint14.md` |
| Sprint 15 | `docs/plans/2026-02-22-sprint15.md` |
| Sprint 16 | `docs/plans/2026-02-22-sprint16.md` |

---

## Phase 详细说明

### Phase 1: MVP（Sprint 1-4，2个月）

核心目标：让用户能完成"浏览邮轮 → 选舱位 → 下单 → 支付"的完整闭环。

**Sprint 1 — 基础设施 + 邮轮介绍**
- Git仓库初始化、Monorepo脚手架
- Docker Compose开发环境（PostgreSQL 17 + Redis 7.4 + MinIO + Meilisearch + NATS）
- Go后端项目初始化（Gin + GORM + Viper + Zap）
- Nuxt 4管理后台项目初始化
- Nuxt 4 Web前台项目初始化
- uni-app小程序项目初始化
- 共享类型包
- JWT认证 + Casbin RBAC
- 邮轮公司CRUD
- 邮轮管理CRUD（含图片上传、状态管理）
- 舱房类型管理CRUD（含富文本、图片画廊）
- 设施分类 + 设施管理CRUD
- 前台展示：邮轮列表、邮轮详情
- GitHub Actions CI/CD
- 数据库迁移

**Sprint 2 — 舱位商品管理**
- 航线域模型 + CRUD
- 航次域模型 + CRUD
- 舱位域模型 + CRUD（含SKU属性）
- 舱位定价矩阵（日期/人数/舱型多维价格）
- 舱位库存管理（总库存、锁定、可用、预警）
- 库存变动日志
- Meilisearch 全文检索集成
- 前台舱位浏览（筛选、排序）
- 前台舱位详情 + 价格日历

**Sprint 3 — 预订流程 + 用户系统**
- 用户域模型（C端用户）
- 微信小程序登录（一键授权）
- Web端手机号+验证码登录
- 微信扫码登录
- 账号绑定机制
- 常用乘客信息管理
- 在线下单流程（选航次→选舱型→选舱位→填乘客→确认）
- 订单创建 + 舱位锁定（15分钟超时释放）
- 我的订单列表
- 后台订单列表

**Sprint 4 — 订单支付 + 通知 + 统计**
- 微信支付集成（小程序JSAPI + 网页Native）
- 支付宝支付集成
- 支付回调处理（验签 + 状态更新 + 库存确认）
- 订单状态流转完整生命周期
- 退改申请 + 后台审核 + 阶梯退款规则
- 退款处理 + 库存释放
- 财务对账报表
- 消息通知系统（微信订阅消息 + 短信 + 站内消息）
- 库存预警通知
- 数据看板（今日概览 + 趋势图 + 排行 + 库存概览）
- 员工账号管理 + 角色分配
- 店铺/品牌信息管理

### Phase 2: V1.0（Sprint 5-8，2个月）

**Sprint 5 — 智能发现与决策支持**
- 用户偏好收集（浏览记录、收藏）
- AI推荐算法（协同过滤 + 基于内容）
- 首页个性化推荐卡片
- 后台推荐权重配置
- 航线日历视图（按日期展示最低价）
- 价格趋势分析折线图
- 最佳预订时机预测
- 舱位并排对比工具（最多3个）

**Sprint 6 — 预订流程增强**
- 一键锁定舱位倒计时优化
- 分期付款（定金+尾款模式）
- 后台分期规则配置
- 多币种支付（CNY/USD/HKD/JPY）
- 汇率自动获取 + 缓存
- 历史乘客快速填充
- 证件OCR识别（护照/身份证）

**Sprint 7 — 团队预订 + 出行服务**
- 团队批量预订入口（前台+后台）
- Excel乘客名单批量导入
- 团队专属定价规则
- 团队订单独立管理
- 出团通知书（PDF上传 + 前台查看 + 推送通知）
- 出行倒计时 + 行前准备清单
- 后台行前模板配置

**Sprint 8 — 实时信息 + 在线客服**
- WebSocket实时舱位状态推送
- 实时浏览热度展示
- 价格变动关注 + 推送通知
- 预订进度可视化时间轴
- AI智能客服（FAQ自动回复）
- 人工客服转接
- 客服工作台
- 天气预报API集成
- 港口停靠信息展示

### Phase 3: V1.5（Sprint 9-12，2个月）

**Sprint 9 — 社交分享**
- 行程海报自动生成
- 海报模板管理（后台）
- 一键分享到微信朋友圈
- 邀请同行优惠机制
- 专属链接/二维码生成
- 后台邀请奖励规则配置

**Sprint 10 — 评价 + 游记社区**
- 出行后评价邀请
- 多维度评分（航线/舱位/服务）
- 图片+视频评价
- 评价审核机制
- UGC游记发布（图文）
- 游记点赞、收藏、评论
- 游记审核 + 首页推荐
- 拼团出行功能

**Sprint 11 — 智能运营（上）**
- 动态定价引擎（库存水位 + 时间 + 需求 + 竞品）
- 后台定价规则配置
- 渠道库存分配（直销/OTA/旅行社/分销）
- 超卖保护机制
- 多维预警系统

**Sprint 12 — 智能运营（下）**
- 收益管理看板（RevPAR/上座率/渠道贡献）
- 客户生命周期管理（CRM）
- 客户标签体系 + 画像
- 自动化营销引擎
- 行为触发规则配置

### Phase 4: V2.0（Sprint 13-16，3个月）

**Sprint 13 — 价格对比 + 多语言**
- 客户航次/舱位价格对比
- i18n框架集成
- 中文简体/繁体/英文支持
- 后台翻译管理
- 预留多语种扩展

**Sprint 14 — 性能优化 + 压测**
- 数据库慢查询优化
- Redis缓存策略完善
- CDN静态资源优化
- API响应时间优化
- 全链路压测脚本（k6/locust）
- 性能基准报告

**Sprint 15 — 安全审计 + 渗透测试**
- OWASP Top 10安全检查
- SQL注入/XSS防护验证
- 敏感数据加密审计
- API权限穿透测试
- 安全修复 + 回归测试
- 安全报告文档

**Sprint 16 — 上线部署 + 监控**
- Kubernetes部署配置
- Caddy反向代理 + HTTPS
- Prometheus指标采集
- Grafana监控面板
- Loki日志收集
- 告警规则配置
- 运维手册编写
- 灾备方案
