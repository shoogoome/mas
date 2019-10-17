# mas
由helm管理的轻量级分布式对象存储系统
已实现功能: 分布式部署、数据校验去重、数据冗余和即时修复、断点续存、断点下载、数据压缩

# 业务流程

1. 通过系统token在后端客户端发起请求获取上传或下载令牌，
2. 前端通过令牌进行上传或下载
3. 对于分片上传需访问初始化接口进行初始化，若得到已持久化信息则无需再上传
4. token时效问题，token在创建之后2分钟内必须进行使用。在每次使用后添加10s有效使用时间。过完时间则token无效

# 安装说明

当前仍处于开发阶段，暂不开放安装说明。可参考`mas1.0`仓库
<!-- **前置条件: 系统配置好Kubernetes分布式环境**

- 本地版  
1. clone下当前资源仓库
2. 创建nfs服务并自行修改values.yaml配置
3. 执行helm install ./liuma 自行按需添加其他参数

- 线上版  
1. helm repo add '自定义仓库名' https://docker.hub.shoogoome.com/chartrepo/liuma
2. helm repo update
3. 创建nfs服务
3. helm install '自定义仓库名'/liuma 自行按需添加其他参数（按照values文件格式修改nfs配置）

**启动系统后需手动初始化mongo环境，数据库: 'liuma'** -->

# 接口说明
**ps: 需附带systemToken的接口应为后端访问接口，前端与系统交互应使用后端获取的临时token**
```
/api/server/signal [get] 获取活跃信号
headers: systemToken 附带系统token
return status 'ip列表'
```
```
/api/token/upload [get] 生成上传令牌
headers: systemToken 附带系统token
url参数: hash 文件hash
return token token
```
```
/api/token/download [get] 生成下载令牌
headers: systemToken 附带系统token
url参数: hash 文件hash
return token token
```
```
/api/file/info [get] 获取文件信息
headers: systemToken 附带系统token
url参数: hash 文件hash
return {
    "status": 状态,
    "code": 状态码,
    "msg": "",
    "data": {
        "size": 文件大小,
        "name": 文件名称,
        "hash": 文件hash,
        "persistence": 是否持续化,
        "create_time": 创建时间,
        "server_ip": null
    }
}
```
```
/api/file/upload/single [post] 单文件上传
headers: token 上传临时token
form-data: file 文件
return {
    "status": 状态,
    "code": 状态码,
    "msg": "",
    "data": {
        "size": 文件大小,
        "name": 文件名称,
        "hash": 文件hash,
        "persistence": 是否持续化,
        "create_time": 创建时间,
        "server_ip": null
    }
}
```
```
/api/file/upload/init [get] 初始化文件信息（当分片上传需要访问）
headers: token 上传临时token
name: 文件名称
return {
    "status": 状态,
    "code": 状态码,
    "msg": "",
    "data": {
        "size": 文件大小,
        "name": 文件名称,
        "hash": 文件hash,
        "persistence": 是否持续化,
        "create_time": 创建时间,
        "server_ip": null
    }
}
```
```
/api/file/upload/chuck [post] 文件分块上传
headers: token 上传临时token
url参数: chuck 当前分片数
form-data: file 分片文件
return {
    "status": 状态,
    "code": 状态码,
    "msg": "",
    "data": {
        "speed": 当前上传块数,
        "status": "success"
    }
}
```
```
/api/file/upload/finish [get] 完成上传(仅断点续传模式需要)
headers: token 上传临时token  
         systemToken 附带系统token
return {
    "status": 状态,
    "code": 状态码,
    "msg": "",
    "data": {
        "size": 文件大小,
        "name": 文件名称,
        "hash": 文件hash,
        "persistence": 是否持续化,
        "create_time": 创建时间,
        "server_ip": null
    }
}
```
```
/api/file/upload/download [get] 下载文件
headers: token 下载临时token
url参数: seek 获取起始游标(可选，用于断点下载)
return file(文件数据)
```

# 使用实践
个人在使用过程中，由于系统将存储服务独立抽象，不依赖业务系统，所以难免出现存储权限问题。  
目前我的解决方案是: 开放获取上传token接口，在业务后端添加完成上传接口，此接口无论单文件上传还是分片上传，前端在上传后均需要访问，以便后端进行权限控制。
此外，存储服务的完成上传接口不由前端访问，将由后端的完成上传接口进行触发，同样为了兼容权限控制。
使用中可根据业务需求让业务提供存储凭证以达到权限控制的目的; 而提供下载的token由业务在特定情况获取并传递给前端，并不直接暴露获取下载token功能

# PS
1.0版本存在许多问题，将在后续版本持续更新改善  
这个ps就当作2.0改版需求文档..
1. 2.0将原本把服务间信息写入环境变量的方式修改为使用rabbitmq进行广播身份信息，以达到动态服务发现。适配kubernetes的statefulset扩张服务，但这样设计宜扩不宜缩，因为减少的服务可能存在文件数据分片，此时只能依赖数据修复功能进行修复减少的数据分片，如果一次性减少的服务过多导致无法修复成功的话，将会有部分文件损失。
2. 当前使用的中间件 redis、mongo仍处于单机形式，后续将部署为集群可扩容形式
3. 补上了1.0的断点下载，其实就是加了个游标。。
4. 计划将加上后台管理服务
