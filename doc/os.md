# 操作系统

## 进程

### 状态

* 就绪态
  * 在ready队列中等待系统调度
* 执行态
  * 正在执行
* 阻塞态
  * 因I/O阻塞
* 挂起态
  * 类似就绪态，但不会被系统调用
    * 用户行为
    * 父进程检查或修改子进程
    * 因负荷过大被调节
    * 操作系统检查资源

### 异常

* 孤儿
  * 父进程退出，但子进程还在运行，此时系统安排init进程领养孤儿进程，若为僵尸态则杀死它
* 僵尸
  * 子进程退出后，父进程没有waitpid或wait方式释放子进程，子进程未运行但仍占用系统资源
    * SIGCHLD信号表示有子进程停止或终止

### 调度算法

* 先来先服务（FCFS）
  * 对CPU繁忙型作业友好
* 短作业优先调度（SJF）
  * 优点
    * 增加了系统吞吐量
  * 缺点
    * 对长作业不利
    * 无优先级考虑
    * 主观的作业时常估算
* 最高优先权优先（FPF）
  * 静态优先权
    * 进程类型
    * 资源需求
    * 用户设定
  * 动态优先权
    * 随着等待时间增加而提高优先权
* 基于时间片轮转调度
  * 时间片轮转
    * 轮流给等待队列中的进程执行，无优先级区分
  * 多及反馈队列
    * N级队列，新作业放入顶级队列，执行一次后移动到下一级，队列从上到下执行，只有在一层执行结束后才会执行下一层极，最后一级则循环执行

* 实时调度
  * 抢占式
    * 当前作业会被更高优先级的作业打断
  * 非抢占式
    * 当前作业执行完之前不会被打断

### 死锁

* 原因
  * 非剥夺性资源争夺
  * 进程非法推进导致资源竞争
    * 因进程执行顺序导致保持等待
* 必要条件
  * 互斥
  * 请求保持
  * 不剥夺
  * 环路等待{A,B,C} => A >> B >> C >> A
* 处理方式
  * 预防死锁
    * 设置条件来打破上述必要条件
      * mutex
      * spin lock
      * event
      * signal
  * 避免死锁
    * 资源分配时预先设计规则防止进入死锁条件
  * 检测死锁
    * 允许死锁，但要及时发现并解除死锁
  * 解除死锁
    * 挂起或撤销进程，同检测死锁

## 内存

### 内存置换算法

* 先进先出（FIFO）
  * 淘汰指针永远指向最早调入内存的块
* 最久未使用（LRU）
  * 淘汰指针永远指向最久未使用的内存块
* Clock置换
  * LRU的模拟，增加访问位，用FIFO遍历，访问为为0的直接淘汰，不为0的置为0继续FIFO
* 最少使用（LFU）
  * 增加位移位寄存器，每访问一次块则置高位为1，每隔一段时间向右移一位，淘汰时找到数最小的

## 网络

### TCP

* 选项
  * MSS：最大分节大小，也就是愿意接受的最大数据量，限制发送端发送分节最大大小
  * 窗口规模：默认为65535既16bit
  * 时间戳：防止失而复得的包
* TIME_WAIT
  * 防止4次挥手的最后一个ACK丢失，导致服务端重发FIN，但收不到ACK的情况
  * 防止连接结束后，又在相同port上启用了新的连接，但这时候又收到了失而复得的上一次连接的数据包
    ![state](https://i.stack.imgur.com/sNTOn.gif)
* 拥塞控制
  * 防止发送太快导致接收方来不及处理
  * 防止大家过于自私占用太多公共带宽，最终导致网络不可用
* 时序图
    ![window](https://upload.wikimedia.org/wikipedia/commons/b/b1/Tcp_transport_example.gif)
    ![connect](https://cdncontribute.geeksforgeeks.org/wp-content/uploads/net.png)
    ![disconnect](https://cdncontribute.geeksforgeeks.org/wp-content/uploads/CN-2-1.png)

### UDP

* 无需建立任何连接，随用随发