## StoreKey

### 1. ordersStoreKey 
| key                           | value           |  number(keys)                 | value detail                                                                                                     | value size | clean up              | 备注                                            |
|-------------------------------|-----------------|----------------------------|------------------------------------------------------------------------------------------------------------------|------------|-----------------------|-------------------------------------------------|
| orderNum:block({blockHeight}) | int64           |  区块高度                  |                                                                                                                  | <1k        | 每区块删除3天前的数据 | 某一区块的order数量                             |
| ID{0-blockHeight}-${Num}      | order.Order     |  所有订单数量              |                                                                                                                  | <1k        | 每区块删除3天前的数据 | 某一区块更新过的订单id列表                      |
|** updatedAt({blockHeight})    | []string        |  区块高度                  | 数组长度取决于该区块内更新的订单数，包括取消、过期、成交<br> 平均值6000，峰值无上限                              | >1k        | 每区块删除3天前的数据 | 某一区块更新过的订单id列表 <br> 仅供backend查询 |
|** depthbook:{product}         | []DepthBookItem |  币对数量                  | 假设某币对价格精度为当前价格的万分之一，<br>正常挂单都在当前价格+-5%以内，<br>则一个币对深度表中含有1000个表项   | >1k        |                       | 某一币对当前的深度表 <br>DepthBookItem数组      |
|** {product}-{price}-{side}    | []string        |  币对数量*价格可能取值数量 | 数组长度取决于某币对某价格的买/卖单数量<br> 平均值不好预估，峰值无上限                                           | >1k        |                       | 某一币对在某一价位的所有买单或卖单的订单id列表  |
|   lastprice:{product}         | sdk.Dec         |  币对数量                  |                                                                                                                  | <1k        |                       | 某一币对的最近成交价                            |

### 2. tradesStoreKey 
| key                            | value                  | detail                     | number(keys) | value detail                                                     | value size | clean up              | 备注                |
|--------------------------------|------------------------|----------------------------|-----------|------------------------------------------------------------------|------------|-----------------------|---------------------|
|** blockMatchResult:{blockHeight} | order.BlockMatchResult | 某一区块的集合竞价撮合结果 | 区块高度  | 撮合结果取决于该区块内成交的订单数目，平均数量为6000，峰值无上限 | >1k        | 每区块删除3天前的数据 | 仅供backend查询使用 |

## Http api

| url              | method | 读key                         | 写key                                                                                                                                                       |
|------------------|--------|-------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| /order/new       | POST   | orderNum:block({blockHeight}) | orderNum:block({blockHeight})<br>ID{0-blockHeight}-${Num}<br>depthbook:{product}<br>{product}-{price}-{side}<br>lastprice:{product}<br>blockMatchResult:{blockHeight} |
| /order/cancel    | POST   | ID{0-blockHeight}-${Num}      | ID{0-blockHeight}-${Num}<br>depthbook:{product}<br>{product}-{price}-{side}                                                                                     |
| /order/depthbook | GET    | depthbook:{product}           |                                                                                                                                                             |
| /order/{orderId} | GET    | ID{0-blockHeight}-${Num}      |                                                                                                                                                             |
