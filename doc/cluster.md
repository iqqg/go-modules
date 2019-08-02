# Redis

## Pipelining

* 本地redis内存执行效率高
* socket的read/write的代价高于redis本身效率
* 将多条指令合并成一条发送，一次read或write减少I/O的成本

## EXPIRE {key}

* EXPIRE特性
  * 指定key的生存期外自动消除
  * 覆盖和删除(DEL,SET,GETSET)key会清除原有的expire属性
  * 修改(INCR,LPUSH,HSET)不会去除原有的expire属性
  * PERSIST后expire属性将被清除
  * REAME后key会携带原有expire属性到新key中
  * REAME后覆盖了已存在的key，expire属性和覆盖方一致
  * 执行EXPIRE重置想要的时间

* EXPIRE的检测
  * 没秒检测10次，每次随机20个拥有expire属性的key
  * 删除所有过期的key
  * 如果超过25%（4个key）的key过期，则立即重新执行一次

## WAIT {numsalve} {timeout}

* return特性
  * timeout被触发
  * replicas满足时被触发
  * timeout给定0，则忽略超时
  * 永远返回replicas的个数
  * 检测replicas个数是否大于等于numsalve，判断执行成功与否
* master宕机时，redis只会尽力replicas，不保证强一致性

## Redis LRU cache

* 设置内存限制
  * 通过redis.conf或执行CONFIG SET修改
  * maxmemory默认64bit为无限制（0MB），32bit为3G
  * LRU（概率模拟LRU）
    * maxmemory-samples 5 // 默认为5，10最好可增加释放的准确度，增加CPU代价
  * LFU（Morris counter）
    * 对命中率友好
    * lfu-log-factor 10 // 1百万命中，满分255，因子越小对访问量越小的环境越好
    * lfu-decay-time 1 // 1分钟后衰减，0每次LFU扫描的时候都会衰减
      | factor | 100 hits | 1000 hits | 100K hits | 1M hits | 10M hits |
      | --- | --- | --- | --- | --- | --- |
      | 0 | 104 | 255 | 255 | 255 | 255 |
      | 1 | 18 | 49 | 255 | 255 | 255 |
      | 10 | 10 | 18 | 142 | 255 | 255 |
      | 100 | 8 | 11 | 49 | 143 | 255 |

* LRU模式（<4.0）
  * noeviction: 超过限制后返回错误，不淘汰任何key
  * allkeys-lru: 超过限制后，从所有key中淘汰最不常用的key，然后写入新数据
  * volatile-lru: 超过限制后，从所有包含expire的key中淘汰最不经常使用的，然后写入新数据
  * allkeys-random: 超过限制后，随机删除key以释放内存，然后写入新数据
  * volatile-random: 超过限制后，从所有包含expire的key中随机选择淘汰，然后写入新数据
  * volatile-ttl: 超过限制后，从所有包含expire的key中选择TTL最少的删除，然后写入新数据
  * volatile-lru; volatile-random; volatile-ttl模式下没有key可以删除，其行为和noeviction一样
* LFU模式（>=4.0）
  * volatile-lfu 只对包含expire属性的key做LFU
  * allkeys-lfu 对所有key做LFU

## 数据持久化

* AOF
  * 只将写入指令保存在文件中
  * 允许将aof文件压缩，类似将inr x; inr x; 改为set x+2
  * 一般每秒保存一次
* RDB
  * 一般达到一个写入频率时触发备份操作
  * 通过fork创建子进程，然后写入本地文件，在此期间父进程发生的改变会通知给子进程
  * 是一种内存快照

## 事务

* 指令
  * MULTI 开启指令缓存
  * EXEC 执行缓存中的指令
  * DISCARD 放弃缓存的指令
* 特性
  * 缓存中出现语法错误，则整条指令都不会执行
  * 缓存中指令执行失败，则继续执行下一个指令

## Pub/Sub

* 指令
  * SUBSCRIBE {channel RegExp}
  * PUBLISH {channel RegExp} {message}
* 特性
  * 一对多

## Redis cluster

* Ports
  * client端口6379，Redis bus端口固定为client端口+1000
  * docker使用了port mapping技术，所以为了使上述规则成立，需要continer使用host networking模式

* Hash slot
  * 非一致性哈希，全局有16384个hash slot
  * key进行CRC16后模16384得到所在的hash slot
  * hash slot可以自由移动，如：

    ``` bash
    # 期初
    Node A 0~5500.
    Node B 5501~11000.
    Node C 11001~16383.

    # 加入Node D：分别从A/B/C中取slot
    Node A 2751~5500
    Node B 8251~11000
    Node C 13751~16383
    Node D 0~2750; 5501~8251; 11001~13751;

    # 删除Node A：删除前需将其保有的slot移动到其他node上
    # 将NodeA的slot清空后，便可以将其下线，因为slot的移动添加删除不会导致系统暂停
    Node B 2751~4126; 5501~11000
    Node C 4127~5500; 11001~16383
    Node D 0~2750; 5501~8251; 11001~13751;
    ```

* Hash tags
  * 允许将多个key映射到一个slot中
  * {XXX}中间的字符串为计算时slot是要用到的key如：`this{foo}key 和 another{foo}key`

* Consistency guarantees
  * 非强一致性，client向master写入后则立即返回，不会等所有replicas写入成功
  * client可以使用WAIT同步等待，但无法保证完全的强一致性，如：

    ```base
    # init:
    #   master A,B,C
    #   salve A1,B1,C1; client Z1;
    # partition:
    #   majority A,C,A1,B1,C1;
    #   minority B,Z1;
    ```

    * 给定最大超时时间内，Z1可以向B写入任何数据，超过maximum window时间后，则任何master在minority侧都停止接受write，此称为node timeout
    * node timeout前写入或分区结束，无任何问题
    * 发生node timeout后master进入失败状态，若无法发现majority中其他的master，则进停止接受write
    * 进入失败状态的master，可以被它之前的replicas方覆盖

* redis.conf
  * cluster-enabled \<yes/no>
    * 是否开启cluster
  * cluster-config-file \<filename>
    * node内部沟通时的状态记录，非人为操作
  * cluster-node-timeout \<milliseconds>
    * node最大的失联时间窗口
  * cluster-slave-validity-factor \<factor>
    * 等于0，slave则一直尝试failover master，忽略master-slave失联的总时间
    * 大于>0，slave则在factor * node-timeout之后将停止尝试failover master。且若没有可用的slave，则有可能导致整个cluster不可用，直到master重新加入。
  * cluster-migration-barrier \<count>
    * 每个master需要slave最少的数目
  * cluster-require-full-coverage \<yes/no>
    * yes，一旦key槽中出现不可达，整个cluster停止服务
    * no，只有key槽中可达的请求才会被处理

* 使用
  
  ```bash

  ```
