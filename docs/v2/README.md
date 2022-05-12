# 绿色金融可视化设计

![image-20220509152548575](https://tva1.sinaimg.cn/large/e6c9d24ely1h226mpdq79j20r90luwh5.jpg)

下方数据与信息展示主要展示感知的数据、认知结果和业务流程。数据展示以json形式。

联动数据通过点击连线显示在下方信息栏。

点击数原体跳转到单独页面，可以查看数原体内部的业务数据（日志信息）以及自身的业务架构

点击认知场景需求，下方展示数原体基于目前数据挖掘出来的需求

点击定义业务流程，展示需要完成这些需求的伪代码

点击联动并产生业务，开始执行业务流程

点击构建场景数原体，展示数原体的构建过程：拉镜像、启动

## 数据结构及业务流程

![未命名文件](https://tva1.sinaimg.cn/large/e6c9d24ely1h25enu1pahj20i40gmq3i.jpg)

### 流程一

制造企业数原体A感知到场景数据，生产运营资金不足以进行后续生产：

```json
{
	"message": "insufficient funds",
  "need": "10000000",
  "surplus": "3762342",
  "timestamp": "2022-05-10 18:00:48.533563 +0800 CST",
}
```

制造企业数原体A经营数据：

```json
{
  "capital": 3762342,
  "carbon_quota": 269,
}
```

然后资金不足，产生需求：卖炭、贷款。

如果选择卖炭，需要寻找碳交易机构数原体，发送数据

```json
{
  "enterprise": "ManufacturingA",
  "carbon_quota": 20,
  "money_per_ton": 350000,
  "timestamp": "2022-05-11 18:00:48.533563 +0800 CST",
  "bank_account": "1234 5678 9101",
}
```

碳交易机构数原体收到之后，会给一个反馈结果数据：

```json
{
  "enterprise": "CarbonTradingAgency",
  "money_per_ton": 350000,
  "result": "ok",
  "timestamp": "2022-05-11 19:00:40.533563 +0800 CST"
}
```

也可以向金融机构以碳资产作为抵押进行贷款，找到合适的金融机构后，发送贷款合同数据

```json
{
  "enterprise": "ManufacturingA",
  "carbon_quota": 20,
  "timestamp": "2022-05-11 18:00:48.533563 +0800 CST",
  "bank_account": "1234 5678 9101",
  "relevant_certificates": "https://www.manufaturingA.com/files",
}
```

金融机构核查相关资料后，进行放款：

```json
{
  "enterprise": "bank",
  "result": "ok",
  "timestamp": "2022-05-11 18:00:48.533563 +0800 CST",
  "contract": "https://www.bank.com/contract100021",
}
```





### 流程二

制造企业数原体B感知到场景数据：碳超标

```json
{
  "message": "high carbon emissions",
  "carbon_emission_per_day": "0.8324",
  "carbon_need": "80",
  "carbon_suplus": "125",
  "timestamp": "2022-05-11 18:00:48.533563 +0800 CST",
}
```

企业经营数据

```json
{
  "capital": 543704582,
  "carbon_quota": 125,
}
```

认知需求：买碳，买技术优化碳排数据

如果是买碳，向碳交易机构数原体发送数据：

```json
{
  "enterprise": "ManufacturingB",
  "carbon_need": 80,
  "timestamp": "2022-05-11 18:00:48.533563 +0800 CST",
}
```

碳交易机构数原体收到请求后，核实后发放碳配额

```json
{
  "enterprise": "CarbonTradingAgency",
  "money_per_ton": 350000,
  "result": "ok",
  "timestamp": "2022-05-11 19:00:40.533563 +0800 CST",
  "bank_account": "2343 3214 8098",
}
```

如果发现同类型企业有好的技术，可以去购买先进技术优化碳排

```json
{
  "enterprise": "ManufacturingB",
  "capital": 10000000,
  "timestamp": "2022-05-13 15:00:48.533563 +0800 CST",
}
```

制造业企业数原体A收到资金后，给出反馈数据：

```json
{
  "enterprise": "ManufacturingA",
  "result": "ok",
  "timestamp": "2022-05-13 15:00:48.533563 +0800 CST",
}
```





## 初始数原体模版

```go
type object struct {
  address        string
	name           string
	funcs          map[string]interface{}
	parentSpan     opentracing.Span
	receivedPacket *agent.Packet
}

Methods:
Run() error
Interact(ctx context.Context, p *agent.Packet) (*agent.Packet, error)
```

具备基础的联动能力。
