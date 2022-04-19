# 可视化设计说明

交互图中一共有三个数源体：distributorA、supplierA、bankA

初始时，注册进来三个数源体，没有交互逻辑

![image-20220414135721834](https://tva1.sinaimg.cn/large/e6c9d24ely1h197kzcvm2j20jd0d9aay.jpg)

然后distributorA<font color=red>感知</font>到一个订单数据，开始寻找suppiler，寻找到合适的一个后，开始把订单发送给supplier。

我们模拟了两种订单，分别对应两种交互逻辑。一种是normal，代表供应商可以进行正常的生产，一种是account-receivable-order，代表供应商资金不足，需要进行贷款融资来生产产品。

```
INFO[0006] perceived new order:&{OrderType:normal OrderPrice:10 OrderCount:10 DistributorName:distributorA} 
INFO[0003] finding supplier....                         
INFO[0006] supplier find: supplierA, establish connection successfully 
INFO[0006] sending data to supplierA: {"order_type":"normal","order_price":10,"order_count":10,"distributor_name":"distributorA"} 
INFO[0006] distributor get the products:  type:TRANSPORT sourceAddress:"127.0.0.1:8080" transport:{data:"{\"supplier_name\":\"supplierA\",\"order_state\":1,\"message\":\"Get products successful.\"}"}
```

![image-20220414140732380](https://tva1.sinaimg.cn/large/e6c9d24ely1h197vj9o3sj20j50drgmk.jpg)

supplierA感知到一个发送过来的订单后，开始生产，因为是normal类型的订单，所以直接把产品返回给经销商即可。

```
INFO[0092] get an order {normal 10 10 distributorA} start producing.... 
INFO[0092] products ready,start transportation...  
```

![image-20220414140803793](https://tva1.sinaimg.cn/large/e6c9d24ely1h197w6dycfj20jq0d70tp.jpg)



当distributor再次感知到同种产品的订单后，不需要再去寻找供应商这个流程了，因为之前已经有了这个记录。可以直接建立连接发送过去，这个就是它的<font color=red>认知</font>能力。

```
INFO[0016] perceived new order:&{OrderType:account-receivable-order OrderPrice:10 OrderCount:10 DistributorName:distributorA} 
INFO[0016] sending data to supplierA: {"order_type":"account-receivable-order","order_price":10,"order_count":10,"distributor_name":"distributorA"} 
```

在这个流程中，supplierA发现资金不足，便会去寻找合适的，能提供贷款能力的数源体，建立连接，发送数据。

```
INFO[0102] get an order {account-receivable-order 10 10 distributorA} start producing.... 
INFO[0102] insufficient funds,looking for a bank to make a loan... 
INFO[0102] finding bank...                              
INFO[0105] bank find: bankA, establish connection successfully 
INFO[0105] generate form: &{Type:account-receivable-order SupplierName:supplierA DistributorName:distributorA LogisticsName: Num:10000} 
INFO[0105] start sending data: {"Type":"account-receivable-order","supplier_name":"supplierA","distributor_name":"distributorA","logistics_name":"","num":10000} 
```

![image-20220414143051687](https://tva1.sinaimg.cn/large/e6c9d24ely1h198jtoctxj20ho0ea3zk.jpg)

bankA感知到数据后，同样执行相应的流程。它从数据中认知到，需要去distributorA那里获得支付承诺，才可以放款给supplierA。于是和distributorA建立连接，发送相应数据。

```
INFO[0107] get a form {account-receivable-order supplierA distributorA  10000} start processing.... 
INFO[0107] get distributor information from form...     
INFO[0110] distributor find: distributorA, establish connection successfully 
INFO[0110] generate payment promise:&{DistributorName:distributorA SupplierName:supplierA Signatured:false} 
INFO[0110] start sending data {"distributor_name":"distributorA","supplier_name":"supplierA","signatured":false}.... 
```

![image-20220414142015555](https://tva1.sinaimg.cn/large/e6c9d24ely1h1988rxrdsj20iy0eqab8.jpg)

distributorA感知到paymentPromise这个数据后，进行相应的处理，这里它只需要验证消息的准确与否，然后进行签名，代表我承诺收到产品后会给你支付货款，你放心给它进行放款。

```
INFO[0022] get a payment promise:{DistributorName:distributorA SupplierName:supplierA Signatured:false} 
INFO[0022] I promise to pay for products, signatured! 
```

bankA收到签名后的支付承诺后，放款给申请者supplierA

```
INFO[0110] verify whether the payment promise is signatured... 
INFO[0110] get the payment promise from distributorA {"distributor_name":"distributorA","supplier_name":"supplierA","signatured":true} 
INFO[0110] the loan is approved. 
```

supplierA收到资金capital后，进行生产，然后返回给经销商货物。

```
INFO[0108] products ready,start transportation...  
```

![image-20220414142448612](https://tva1.sinaimg.cn/large/e6c9d24ely1h198dimdouj20hf0e4ab6.jpg)

distributorA收到货物后，给bankA支付货款。

```
INFO[0025] prepare capital for products...              
INFO[0025] sending capital to bankA: {"bank_name":"bankA","num":100} 
INFO[0025] the payment result: type:TRANSPORT sourceAddress:"127.0.0.1:8082" transport:{data:"true"} 
```

bankA

```
INFO[0113] I have received capital:  {bankA 100}  
```

![image-20220414142828914](https://tva1.sinaimg.cn/large/e6c9d24ely1h198hd9wonj20hp0f7aba.jpg)

至此，这个订单完成。

前端页面整体流程：

1. 页面有两个按钮，start和generate data
2. 点击start开始我们模拟过程，一共有5个数源体逐渐注册进来，需要有动态的效果展示
3. 然后点击generate data开始生成数据给distributorA，每个数源体感受到外界的数据后，做高亮放大闪烁一下。
4. 开始我们第一个normal订单的处理逻辑，发现、联动
5. 再次点击generate data开始第二个流程，此流程中supplierA会认知到需要贷款，旁边的日志做高亮展示，并短暂停留。然后开始发现能提供服务的数源体，开始建立连接，发送联动报文，依次类推。
6. 其他动态效果和之前一样。

![prototype](https://tva1.sinaimg.cn/large/e6c9d24ely1h1ezkldm0ug217o0o6mz0.gif)

整个交互的动态效果类似于下图：

![交互](https://tva1.sinaimg.cn/large/e6c9d24ely1h198xk51a6g20pl0esq51.gif)