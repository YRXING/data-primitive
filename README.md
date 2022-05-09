# 数据原生场景验证接口设计

## 场景描述

### 供应链金融

#### 背景

一般来说，供应链包括供应商->生产商->经销商->零售商->客户。

![Image:一个典型的供应链.jpg](https://tva1.sinaimg.cn/large/e6c9d24ely1h014o632mfj20dw0afwey.jpg)

为了保障企业的正常运作，必须确保资金的及时回收，否则企业就无法建立完善的经营体系。供应链资金流的方向是由消费者->零售商->经销商->生产商->供货商。



#### 业务流程

下面是供应链场景下数源体交互图。

根据供应链金融的三种融资模式：应收账款融资、未来货权融资、融通仓融资，构建以下数源体联动图

![供应商](https://tva1.sinaimg.cn/large/e6c9d24ely1h01hqcbq92j213g0t00uu.jpg)

- 数源体A：上游供货商
- 数源体B：金融机构（商业银行等）
- 数源体C：第三方物流
- 数源体D：下游购货商、经销商

### 绿色金融

#### 背景

​	**碳交易，即把二氧化碳排放权作为一种商品，**买方通过向卖方支付一定金额从而获得一定数量的二氧化碳排放权，从而形成了二氧化碳排放权的交易。**碳交易市场是由政府通过对能耗企业的控制排放而人为制造的市场。**通常情况下，政府确定一个碳排放总额，并根据一定规则将碳排放配额分配至企业。**如果未来企业排放高于配额，需要到市场上购买配额。**与此同时，部分企业通过采用节能减排技术，最终碳排放低于其获得的配额，则可以通过碳交易市场出售多余配额。双方一般通过碳排放交易所进行交易。**第一种情况，如果企业减排成本低于碳交易市场价时，**企业会选择减排，减排产生的份额可以卖出从而获得盈利；**第二种情况，当企业减排成本高于碳市场价时**，会选择在碳市场上向拥有配额的政府、企业、或其他市场主体进行购买，以完成政府下达的减排量目标。**若未足量购买配额以覆盖其实际排放量则面临高价罚款。**通过这一套设计，**碳交易市场将碳排放内化为企业经营成本的一部分**，而交易形成的碳排放价格则引导企业选择成本最优的减碳手段，**包括节能减排改造、碳配额购买、或碳捕捉等**，市场化的方式使得在产业结构从高耗能向低耗能转型的同时，全社会减排成本保持最优化。

​	目前我国碳交易市场有两类基础产品，**一类为政府分配给企业的碳排放配额，另一类为核证自愿减排量(CCER)。****CCER** **是指对我国境内可再生能源、林业碳汇、甲烷利用等项目的温室气体减排效果进行量化核证**，并在国家温室气体自愿减排交易注册登记系统中登记的温室气体减排量。**第一类，配额交易**，是政府为完成控排目标采用的一种政策手段，即在一定的空间和时间内，将该控排目标转化为碳排放配额并分配给下级政府和企业，**若企业实际碳排放量小于政府分配的配额，则企业可以通过交易多余碳配额**，来实现碳配额在不同企业的合理分配，最终以相对较低的成本实现控排目标。**第二类，作为补充，在配额市场之外引入自愿减排市场交易，即 CCER 交易**。CCER 交易指控排企业向实施“碳抵消”活动的企业购买可用于抵消自身碳排的核证量。**“****碳抵消”**是指用于减少温室气体排放源或增加温室气体吸收汇，用来实现补偿或抵消其他排放源产生温室气体排放的活动，即控排企业的碳排放可用非控排企业使用清洁能源减少温室气体排放或增加碳汇来抵消。**抵消信用由通过特定减排项目的实施得到减排量后进行签发，项目包括可再生能源项目、森林碳汇项目等。**碳市场按照 1:1 的比例给予 CCER 替代碳排放配额，**即 1 个 CCER 等同于 1 个配额，可以抵消 1 吨二氧化碳当量的排放**，《碳排放权交易管理办法（试行）》规定重点排放单位每年可以使用国家核证自愿减排量抵销碳排放配额的清缴，**抵消比例不得超过应清缴碳排放配额的 5%。**

#### 碳资产抵押融资

​	碳排放权抵质押融资，指控排企业将自身获得的碳排放权进行担保，通过抵押或者质押的方式获得金融机构融资的一种业务模式，是一种新型的[绿色信贷](https://news.bjx.com.cn/topics/lusexindai/)产品和融资贷款模式。（每万吨价值20-50万元）

​	好处：	

- 对于企业尤其是重点控排企业来说，其具有的碳排放权配额是他们的一大无形资产，若企业不想出售其碳配额，又想降低资金占用压力，将碳排放权作为担保向银行申请贷款是他们最好的选择，解决了一些中小型企业融资难的问题。
- 发展碳排放权抵质押融资业务有助于银行推动绿色信贷业务发展，但各家银行的风控标准及风险偏好存在差异，对目标客户画像有所不同，也造成各家银行的贷款审批流程也存在一些差异。

#### 业务流程

![image-20220509113704790](https://tva1.sinaimg.cn/large/e6c9d24ely1h2200pfp2jj20qi0jjmzs.jpg)

制造企业数原体1、制造业企业数原体2、金融机构数原体1

数据包括企业的生产经营数据、企业的碳排数据、监管机构的各类政策数据

针对不同的碳排数据可执行的流程：

- 企业数原体可以对自己多余的碳配额卖给其他企业进行资产变现
- 还可以将自己的低碳技术输出给同类型企业
- 金融机构数原体可以针对企业数原体的经营数据、碳配额等进行绿色信贷，或者根据监管机构的各类政策数据开发更多的产品服务客户

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



## 实验结果

![image-20220331162210914](https://tva1.sinaimg.cn/large/e6c9d24ely1h0t53c32j3j21iw0u0taf.jpg)

目前三个数源体运行过程中的日志如下

- distributor

  感知不同的订单数据

  ```bash
  2022/03/31 16:17:46 debug logging disabled
  2022/03/31 16:17:46 Initializing logging reporter
  2022/03/31 16:17:46 debug logging disabled
  2022/03/31 16:17:46 debug logging disabled
  2022/03/31 16:17:46 Initializing logging reporter
  2022/03/31 16:17:46 debug logging disabled
  INFO[0000] Starting *distributor.distributor gRPC server, listener on 127.0.0.1:8081 
  INFO[0003] start simulation process...                  
  INFO[0003] finding supplier....                         
  INFO[0006] supplier find: supplierA, establish connection successfully 
  INFO[0006] perceived new order:&{OrderType:normal OrderPrice:10 OrderCount:10 DistributorName:distributorA} 
  INFO[0006] sending data to supplierA: {"order_type":"normal","order_price":10,"order_count":10,"distributor_name":"distributorA"} 
  2022/03/31 16:17:52 Reporting span 2866de34dff70d8b:2866de34dff70d8b:0000000000000000:1
  INFO[0006] distributor get the products:  type:TRANSPORT sourceAddress:"127.0.0.1:8080" transport:{data:"{\"supplier_name\":\"supplierA\",\"order_state\":1,\"message\":\"Get products successful.\"}"} 
  INFO[0016] perceived new order:&{OrderType:account-receivable-order OrderPrice:10 OrderCount:10 DistributorName:distributorA} 
  INFO[0016] sending data to supplierA: {"order_type":"account-receivable-order","order_price":10,"order_count":10,"distributor_name":"distributorA"} 
  INFO[0022] get a payment promise:{DistributorName:distributorA SupplierName:supplierA Signatured:false} 
  INFO[0022] I promise to pay for products, signatured!   
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:0c1d04a08468d71c:2fc4e9248b21887f:1
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:2fc4e9248b21887f:38fda3de0eafd03d:1
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:68a992153ec3e093:0000000000000000:1
  INFO[0022] distributor get the products:  type:TRANSPORT sourceAddress:"127.0.0.1:8080" transport:{data:"{\"supplier_name\":\"supplierA\",\"order_state\":1,\"message\":\"Get products successful.\"}"} 
  INFO[0022] get bank information from order              
  INFO[0025] bank find: bankA, establish connection successfully 
  INFO[0025] prepare capital for products...              
  INFO[0025] sending capital to bankA: {"bank_name":"bankA","num":100} 
  2022/03/31 16:18:11 Reporting span 68a992153ec3e093:4df9f4d333e03ac7:2fc4e9248b21887f:1
  INFO[0025] the payment result: type:TRANSPORT sourceAddress:"127.0.0.1:8082" transport:{data:"true"} 
  INFO[0030] simulation process finished. 
  ```

- supplier

  ```bash
  2022/03/31 16:16:20 debug logging disabled
  2022/03/31 16:16:20 Initializing logging reporter
  2022/03/31 16:16:20 debug logging disabled
  INFO[0000] Starting *supplier.supplier gRPC server, listener on 127.0.0.1:8080 
  INFO[0092] get an order {normal 10 10 distributorA} start producing.... 
  INFO[0092] products ready,start transportation...       
  2022/03/31 16:17:52 Reporting span 2866de34dff70d8b:046de7dede28869c:1ddf448b16f3af92:1
  2022/03/31 16:17:52 Reporting span 2866de34dff70d8b:1ddf448b16f3af92:2866de34dff70d8b:1
  INFO[0102] get an order {account-receivable-order 10 10 distributorA} start producing.... 
  INFO[0102] insufficient funds,looking for a bank to make a loan... 
  INFO[0102] finding bank...                              
  INFO[0105] bank find: bankA, establish connection successfully 
  INFO[0105] generate form: &{Type:account-receivable-order SupplierName:supplierA DistributorName:distributorA LogisticsName: Num:10000} 
  INFO[0105] start sending data: {"Type":"account-receivable-order","supplier_name":"supplierA","distributor_name":"distributorA","logistics_name":"","num":10000} 
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:17abccf1adc6fea7:11862841d65c2c66:1
  INFO[0108] products ready,start transportation...       
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:11862841d65c2c66:0679001a424cecf0:1
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:0679001a424cecf0:68a992153ec3e093:1
  ```

- bank

  ```bash
  2022/03/31 16:16:17 debug logging disabled
  2022/03/31 16:16:17 Initializing logging reporter
  2022/03/31 16:16:17 debug logging disabled
  INFO[0000] Starting *bank.bank gRPC server, listener on 127.0.0.1:8082 
  INFO[0107] get a form {account-receivable-order supplierA distributorA  10000} start processing.... 
  INFO[0107] get distributor information from form...     
  INFO[0110] distributor find: distributorA, establish connection successfully 
  INFO[0110] generate payment promise:&{DistributorName:distributorA SupplierName:supplierA Signatured:false} 
  INFO[0110] start sending data {"distributor_name":"distributorA","supplier_name":"supplierA","signatured":false}.... 
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:38fda3de0eafd03d:67b30e883a0dc7da:1
  INFO[0110] verify whether the payment promise is signatured... 
  INFO[0110] get the payment promise from distributorA {"distributor_name":"distributorA","supplier_name":"supplierA","signatured":true} 
  INFO[0110] the loan is approved.                        
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:67b30e883a0dc7da:48b8f69f0c5fe3c7:1
  2022/03/31 16:18:08 Reporting span 68a992153ec3e093:48b8f69f0c5fe3c7:17abccf1adc6fea7:1
  INFO[0113] I have received capital:  {bankA 100}        
  2022/03/31 16:18:11 Reporting span 68a992153ec3e093:34b157ef398e542d:4df9f4d333e03ac7:1
  ```

  
