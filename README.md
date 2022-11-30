# ruomu-uc
用户中心

提供基础用户、角色、资源、权限的管理

依赖于 ruomu-core

# 注入信息
## 模块名称
uc

## 模块注入点
| 代码                        | 注入类型      | 说明         |
|---------------------------|-----------|------------|
| authorizationInfoByUserId | CODE      | 获取用户角色ID列表 |
| authorizationInfoByRoleId | CODE      | 获取角色资源列表   |
| user:add                  | HTTP-POST | 新增用户       |
| user:login                | HTTP-POST | 用户登录       |




[//]: # (goreleaser release --skip-publish --rm-dist --snapshot)
