# lrishman

![GitHub Workflow Status](https://img.shields.io/github/workflow/status/ranzhendong/irishman/IrishManCI?label=GithubBuild&logo=github&style=plastic)
![Travis (.org) branch](https://img.shields.io/travis/ranzhendong/irishman?label=TravisBulid&logo=travis&style=plastic)
![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/ranzhendong/irishman?include_prereleases&style=plastic)
![GitHub last commit ](https://img.shields.io/github/last-commit/ranzhendong/irishman/master?style=plastic)
![GitHub](https://img.shields.io/github/license/ranzhendong/irishman?style=plastic)



</br></br>

&emsp;&emsp;通过go二次重构，以更简单的方式实现原来在**openresty**通过lua实现的功能，主要的重构的部分在**upstream**、**healthCheck**和**weightServerList**部分



## 1. upstream

&emsp;&emsp;开发标准restful接口来实现对upstream的增删改查。

&emsp;&emsp;下面的实例都将以请求json来进行示例。

</br>

### 1.1 GET

&emsp;&emsp;获取资源。

&emsp;&emsp;获取已有资源。

&emsp;&emsp;只会读取**upstreamName**关键词，即使数据结构传入其他参数。

#### 1.1.1 获取单个资源

```json
{
	"upstreamName": "vmims"
}
```



#### 1.1.2 获取所有资源

```json
{
	"upstreamName": "ALL"
}
```



</br>

### 1.2 PUT

&emsp;&emsp;更新资源。

&emsp;&emsp;只能更新，不能创建。

#### 1.2.1 更新资源

```json
{
	"upstreamName": "vmims",
	"algorithms": "round-robin",
	"pool": [
		{
			"ipPort": "192.168.101.59:8080",
			"status": "up",
			"weight": 2
		},
		{
			"ipPort": "192.168.101.61:8080",
			"status": "down",
			"weight": 3
		},
		{
			"ipPort": "192.168.101.77:8080",
			"status": "down",
			"weight": 6
		}
	]
}
```



&emsp;&emsp;更新相同的**upstreamName**，会全量覆盖，如果更新不存在资源则会报错。



#### 1.2.2 更新不存在资源报错

```json
{
    "Error": "No Key [UpstreamVmims]",
    "Message": "Upstream PUT: Etcd Get: Key Not Exist Error",
    "Code": 2102,
    "TimeStamp": "2020-01-10T17:41:55.8043597+08:00",
    "ExecutorTime": "7.9785ms"
}
```



</br>

### 1.3 POST

&emsp;&emsp;创建资源。

&emsp;&emsp;会做校验，只能创建不能更新。

#### 1.3.1 创建资源

```json
{
	"upstreamName": "vmims",
	"algorithms": "round-robin",
	"pool": [
		{
			"ipPort": "192.168.101.59:8080",
			"status": "up",
			"weight": 3
		},
		{
			"ipPort": "192.168.101.61:8080",
			"status": "down",
			"weight": 4
		}
	]
}
```



#### 1.3.2 创建已有资源

&emsp;&emsp;如果创建已经存在的资源，则会报错：

```json
{
    "Error": "Upstream POST: Etcd Get: Repeat Key Error",
    "Message": "Upstream POST: Etcd Get: Repeat Key Error",
    "Code": 3103,
    "TimeStamp": "2020-01-10T17:43:14.9657714+08:00",
    "ExecutorTime": "3.9582ms"
}
```



</br>

### 1.4 PATCH

&emsp;&emsp;更新部分已有资源。



#### 1.4.1 更新局部资源

```json
{
    "upstreamName": "vm",
    "algorithms": "ip-hex",
    "pool": [
        {
            "ipPort": "192.168.101.59:8080",
            "status": "nohc",
        },
        {
            "ipPort": "192.168.101.61:8080",
            "status": "down",
            "weight": 7
        },
        {
            "ipPort": "192.168.101.61:8085",
            "status": "up"
        },
        
		{
            "ipPort": "192.168.101.6:9999",
            "status": "up",
            "weight": 15
        }
    ]
}
```



#### 1.4.2 更新局部资源四种情况

&emsp;&emsp;上面的示例当中展示了四种不同情况：

- patch修改已经存在资源的pool，status或者weight其中一个。
- patch修改已经存在资源pool，status和weight需要全部，将新创建server
- patch修改不存在资源pool，不会有资源生效。
- patch单独修改algorithms资源，pool可以为空，或者同时修改。





</br>

### 1.5 DELETE

&emsp;&emsp;删除资源，无论是否存在。



#### 1.5.1 删除资源

```json
{
	"upstreamName": "vmims"
}
```



#### 1.5.2 删除资源server list

&emsp;&emsp;删除upstream的pool，一个或者多个。

```json
{
	"upstreamName": "vm",
	"algorithms": "round-robin",
	"pool": [
		{
			"ipPort": "192.168.101.59:8080",
			"status": "up",
			"weight": 2
		},
		{
			"ipPort": "192.168.101.61:8080",
			"status": "down",
			"weight": 3
		}
	]
}
```



#### 1.5.3 删除资源三种情况

&emsp;&emsp;delete删除资源pool有三种情况：

- 删除一个或者多个，当etcd存储的资源的pool已经只剩一个server list，就会返回删除失败，保证每个upstream至少有一个server list。即使准备删除的server list和已经存在server list不匹配，也会返回删除失败。
- 删除一个或者多个，当检测到将pool下面的所有server list全部删除，也会返回删除失败。删除之后至少保证，pool当中还有一个server list可以使用，不能全部删除。
- 成功删除一个或者多个server list。



</br></br>

# Copyright & License

BSD 2-Clause License

Copyright (c) 2019, Zhendong
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

- Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

- Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.