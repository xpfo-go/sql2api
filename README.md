# sql2api:

一个简单的sql生成api的程序，初步为数据分析场景设计，此程序路由未做优化，后续逐步实现todo的内容。

# quick start:



# used:

### 获取api列表
```http request
GET http://localhost:8000/api/list
```

### 创建api
```http request
POST http://localhost:8000/api/create
Content-Type: application/json

{
    "method": "GET",
    "url": "t1",
    "sql": "select * from user"
}
```

### 删除api
```http request
POST http://localhost:8000/api/delete
Content-Type: application/json

{
    "url": "t1"
}
```

# todo list:

1. 支持多个数据库连接 ✅
2. 支持更多数据库类型 pgsql clickhouse redis...
3. 接口支持分页? 待定
4. 直接生成到文件， gin、spring、php...
5. 路由持久化, 数据库链接持久化， 引入sqlite

## todo 3个模块
1. 添加数据库
    1.1 mysql.....
2. 生成api
3. 生成文件
