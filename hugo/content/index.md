---
weight: 10
title: API Reference
---



# 介绍

ETH订单管理微服务，提供以下功能：

1. 用户钱包注册，用于钱包交易推送（无钱包管理功能）；
2. 订单管理，包括：获取订单列表,创建订单，查询订单状态以及状态推送；

所有接口通过 HTTP RESTful 方式提供

# 

## 注册用户钱包

### HTTP Request

`POST http://xxxxx.com/wallet/:userid/:address` 

#### 请求参数


Parameter | Type | Description
--------- | ------- | -----------
userid|string|阿里云推送账号ID
address|string|ETH钱包地址

## 删除用户钱包

### HTTP Request

`DELETE http://xxxxx.com/wallet/:userid/:address` 

#### 请求参数


Parameter | Type | Description
--------- | ------- | -----------
userid|string|阿里云推送账号ID
address|string|ETH钱包地址

## 创建订单

### HTTP Request

`POST http://xxxxx.com/order` 

#### 请求参数


Parameter | Type | Description
--------- | ------- | -----------
tx|string|订单ID
from|string|转账来源钱包地址
to|string|转账目标钱包地址
asset|string|转账资产类型ID
value|string|订单转账金额
context|json|订单上下文数据，json字符串（非JSON对象）

> 请求参数

```json
{
    "tx":"",
    "from":"",
    "to":"",
    "value":"",
    "context":""
}
```
## 获取订单状态

### HTTP Request

`GET http://xxxxx.com/order/:tx` 

#### 请求参数


Parameter | Type | Description
--------- | ------- | -----------
tx|string|订单ID

> 响应参数

```json
{
    "status":false
}
```

## 获取订单列表

### HTTP Request

`GET http://xxxxx.com/orders/:address/:asset/:offset/:size` 

#### 请求参数


Parameter | Type | Description
--------- | ------- | -----------
address|string|钱包地址
asset|string|资产类型
offset|number|分页
size|number|分页大小

> 响应参数

```json
[
{
"tx": "0x67905b068cde98d0450168bf6f8feac5eac390073a97cb660dadb056fa31ca11",
"from": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"to": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"asset": "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
"value": "1",
"createTime": "2017-11-26T22:38:16.133121Z",
"confirmTime": "2017-11-26T22:38:50.41296Z"
},
{
"tx": "0x526c5d94b828a35ac1a165008a0777ef052be3192e194c134a78f34fedab7e36",
"from": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"to": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"asset": "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
"value": "1",
"createTime": "2017-11-26T22:41:44.013348Z",
"confirmTime": "2017-11-26T22:42:05.859609Z"
},
{
"tx": "0x1e1cda7e791cf896f321efe5524d78ebf5aacb874b9f17999bd79dd445b7dac3",
"from": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"to": "AMpupnF6QweQXLfCtF4dR45FDdKbTXkLsr",
"asset": "0xc56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b",
"value": "1",
"createTime": "2017-11-26T22:57:50.941295Z",
"confirmTime": "2017-11-26T22:58:24.039529Z"
}
]
```


