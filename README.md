# ruomu-uc
用户中心

提供基础用户、角色、资源、权限的管理

依赖于 ruomu-core

# 注入信息
## 模块名称
uc

## 模块注入点
| 代码                        | 注入类型        | 说明         |
|---------------------------|-------------|------------|
| authorizationInfoByUserId | CODE        | 获取用户角色ID列表 |
| authorizationInfoByRoleId | CODE        | 获取角色资源列表   |
| /user/login               | HTTP-POST   | 用户登录       |
| /user/add                 | HTTP-POST   | 新增用户       |
| /user/update              | HTTP-PUT    | 修改用户       |
| /user/delete              | HTTP-DELETE | 删除用户       |
| /user/instance            | HTTP-GET    | 用户详情       |
| /user/list                | HTTP-GET    | 用户列表       |
| /user/password            | HTTP-PUT    | 用户密码修改     |


## 注入请求
POST /module/add
```json
{
    "name": "用户中心",
    "code": "uc",
    "cmd": "./plugins/ruomu_uc.exe",
    "status": 1,
    "needDb": true,
    "needUserTokenExpire": true,
    "dependencies": [],
    "injects": [
        {
            "name": "获取用户角色ID列表",
            "type": 51,
            "injectCode": "authorizationInfoByUserId",
            "authorizationCode": "inner",
        },
        {
            "name": "获取角色资源列表",
            "type": 51,
            "injectCode": "authorizationInfoByRoleId",
            "authorizationCode": "inner",
        },
        {
            "name": "登录",
            "type": 2,
            "injectCode": "/user/login",
            "authorizationCode": "anno"
        },
        {
            "name": "新增用户",
            "type": 2,
            "injectCode": "/user/add",
            "authorizationCode": "user:add"
        },
        {
            "name": "修改用户",
            "type": 3,
            "injectCode": "/user/update",
            "authorizationCode": "user:update"
        },
        {
            "name": "删除用户",
            "type": 4,
            "injectCode": "/user/delete",
            "authorizationCode": "user:delete"
        },
        {
            "name": "用户详情",
            "type": 1,
            "injectCode": "/user/instance",
            "authorizationCode": "user:instance"
        },
        {
            "name": "用户列表",
            "type": 1,
            "injectCode": "/user/list",
            "authorizationCode": "user:list"
        },
        {
            "name": "修改用户密码",
            "type": 3,
            "injectCode": "/user/password",
            "authorizationCode": "user:password"
        }
    ]
}
```

[//]: # (goreleaser release --skip-publish --rm-dist --snapshot)
