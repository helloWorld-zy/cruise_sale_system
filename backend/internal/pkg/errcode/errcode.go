package errcode

// 业务错误码 —— 通过 Response.Code 字段返回给客户端。
// HTTP 状态码在传输层（handler）单独设置。
const (
	// 成功
	OK = 0

	// 通用客户端错误（4xx 范围）
	ErrBadRequest   = 40000 // 请求参数错误
	ErrUnauthorized = 40001 // 未认证（未登录或令牌无效）
	ErrForbidden    = 40003 // 无权限（RBAC 拒绝）
	ErrNotFound     = 40004 // 资源不存在
	ErrConflict     = 40009 // 资源冲突

	// 参数验证错误
	ErrValidation = 42200 // 请求参数验证失败

	// 业务逻辑错误
	ErrCruiseHasCabins   = 42201 // 邮轮下仍有舱房类型，无法删除
	ErrCompanyHasCruises = 42202 // 公司下仍有邮轮，无法删除
	ErrPasswordMismatch  = 42203 // 密码不匹配

	// 服务器内部错误（5xx 范围）
	ErrInternal = 50000 // 服务器内部错误
)
