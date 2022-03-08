# 数据原生场景验证接口设计

## 场景描述

### 背景

一般来说，供应链包括供应商->生产商->经销商->零售商->客户。

![Image:一个典型的供应链.jpg](https://tva1.sinaimg.cn/large/e6c9d24ely1h014o632mfj20dw0afwey.jpg)

为了保障企业的正常运作，必须确保资金的及时回收，否则企业就无法建立完善的经营体系。供应链资金流的方向是由消费者->零售商->经销商->生产商->供货商。



### 数源体交互图

下面是供应链场景下数源体交互图。

根据供应链金融的三种融资模式：应收账款融资、未来货权融资、融通仓融资，构建以下数源体联动图

![供应商](https://tva1.sinaimg.cn/large/e6c9d24ely1h01hqcbq92j213g0t00uu.jpg)

- 数源体A：上游供货商
- 数源体B：金融机构（商业银行等）
- 数源体C：第三方物流
- 数源体D：下游购货商、经销商

## 系统设计

### 感知

### 认知

根据数据确定需要的业务，然后给能提供服务的数源体发送联动报文，数源体之间交互需要数源体彼此可以互相发现。

相互了解的信息包括：数源体自身信息，包括资源、能力、数据。

**<font color=red>解决方案：</font>**

通过注册中心实现，带kv存储的。

经调研后，确定consul为注册中心：简单易用，不需要sdk集成，支持多数据中心、健康检查。

```json
{
   "services": [
       {
         "id":"server1",
         "name":"business_service1",
         "address": "10.10.102.1",
         "tags": [
             "webapi"
         ],
         "port":80,
       },
       {
         "id":"server2",
         "name":"business_service2",
         "address": "10.10.102.2",
         "tags": [
             "webapi2"
         ],
         "port":81,
       }
   ]
}
```

kv存储充当数源体的数据后端，key为数源体名字，value为数源体所有可能用到的信息

```json
{
  "name": "DSourceA",
	"services": [
    {
    "name": "business_service1",
    "args": ""
    }
	],
  
  /*
  *待补充
  */
}
```

当数源体感知到数据后，提炼数据，去存储中心查询服务。



### 联动

**<font color=red>解决方案：</font>**

通信采用gRPC，各个数源体对外定义成统一的服务格式和通信格式

每个数源体服务有两种类型：

- 业务服务：用来执行具体的业务 (主要被本体调用，直接在内部定义成普通方法，是否对外暴露？)
- 联动服务：用来和其他数源体交互（是否采用sidecar模式。把这部分功能独立出来）

通信格式定义：

```protobuf
syntax = "proto3"

service DataA {
	rpc interact(Packet) returns (string);
}

enum PacketType {
	DEPLOY = 0;
	TRANSPORT = 1;
	INVOKE = 2;
}

message Packet {
	PacketType type=1;
	string sourceAddress=2;
	string sendAddress;
	oneof payload{
		Invoke invoke=3;
		Transport transport=4;
	}
}

message Invoke {
	string funcName=1;
	string args=2;
}

message Transport {
	string data=1;
}
```



每个数源体的联动代码interact整体逻辑应该一致，区别是调用的业务服务不同,根据不同的数据包采取不同的动作。

目前主要实现INVOKE数据包，调用完服务后，下一步动作可以分为接着调用其他服务，或者返回调用方所需要的数据。

```go
func (s *supplier) Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error)  {
	switch p.Type {
	case agent.PacketType_INVOKE:
		res, err:= util.Call(s.funcs,p.GetInvoke().FuncName,p.GetInvoke().Args)
		//s.client = util.NewClient(p.SendAddress)
		if err != nil {
			// change []reflect.Value to []interface{}
			data := make([]interface{},0)
			for _, v := range res {
				data = append(data,v.Interface())
			}

			bytes,_ := json.Marshal(data)
			// make return packet
			pkt := &agent.Packet{
				Type: agent.PacketType_TRANSPORT,
				SourceAddress: s.address,
				Payload: &agent.Packet_Transport{
					Transport: &agent.Transport{
						Data: string(bytes),
					},
				},
			}
			return pkt,nil
		}
	case agent.PacketType_TRANSPORT:

	case agent.PacketType_DEPLOY:

	}
	return nil,nil
}
```

### 执行

## 问题及难点

- 如何根据感知的数据来确定联动其他数源体所需要的相关参数？

  先根据关键字进行服务匹配，根据原数据+简单逻辑判断-> 需要调用的服务

- 相同服务的数源体有多个时，如何却定联动哪一个？

  初期根据设定好的规则，不同数据走不同的路径

