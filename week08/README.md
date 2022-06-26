# 极客时间 Go训练营作业

## 作业
1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。
```
root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 10 -q
SET: 163934.42 requests per second, p50=0.159 msec
GET: 163132.14 requests per second, p50=0.159 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 20 -q
SET: 159744.41 requests per second, p50=0.159 msec
GET: 167224.08 requests per second, p50=0.151 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 50 -q
SET: 167504.19 requests per second, p50=0.151 msec
GET: 165562.92 requests per second, p50=0.151 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 100 -q
SET: 166666.66 requests per second, p50=0.151 msec
GET: 163132.14 requests per second, p50=0.159 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 200 -q
SET: 175131.36 requests per second, p50=0.143 msec
GET: 177304.97 requests per second, p50=0.143 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 1024 -q
SET: 159235.66 requests per second, p50=0.159 msec
GET: 210526.31 requests per second, p50=0.127 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 5120 -q
SET: 153609.83 requests per second, p50=0.167 msec
GET: 167785.23 requests per second, p50=0.159 msec

root@06dcde9914d7:/data# redis-benchmark -t set,get -n 100000 -d 51200 -q
SET: 130548.30 requests per second, p50=0.199 msec
GET: 96993.21 requests per second, p50=0.167 msec
```
<p>redis基准测试文档: https://redis.io/docs/reference/optimization/benchmarks/</p>


2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。
```
go run main.go
--------------------插入:10000个10字节--------------------
插入前used_memory:880176(byte)
插入后used_memory:1877056(byte)
平均每个 key 的占用内存空间:99(byte)
--------------------插入:500000个10字节--------------------
插入前used_memory:880448(byte)
插入后used_memory:48994752(byte)
平均每个 key 的占用内存空间:96(byte)
--------------------插入:10000个200字节--------------------
插入前used_memory:880600(byte)
插入后used_memory:1957208(byte)
平均每个 key 的占用内存空间:107(byte)
--------------------插入:500000个20字节--------------------
插入前used_memory:880600(byte)
插入后used_memory:52994904(byte)
平均每个 key 的占用内存空间:104(byte)
--------------------插入:10000个1k字节--------------------
插入前used_memory:880752(byte)
插入后used_memory:14517360(byte)
平均每个 key 的占用内存空间:1363(byte)
--------------------插入:500000个1k字节--------------------
插入前used_memory:880752(byte)
插入后used_memory:680995056(byte)
平均每个 key 的占用内存空间:1360(byte)
--------------------插入:10000个5k字节--------------------
插入前used_memory:880904(byte)
插入后used_memory:63157512(byte)
平均每个 key 的占用内存空间:6227(byte)
--------------------插入:500000个5k字节--------------------
插入前used_memory:880904(byte)
插入后used_memory:3112995208(byte)
平均每个 key 的占用内存空间:6224(byte)
```

Redis_Memory_Info 信息含义

|指标 |含义| 
|------|:------| 
|used_memory|由 Redis 分配器分配的内存总量，以字节（byte）为单位，即当前redis使用内存大小。|
|used_memory_human|已更直观的单位展示分配的内存总量。|
|used_memory_rss|向操作系统申请的内存大小，即redis使用的物理内存大小。|
|used_memory_rss_human|已更直观的单位展示向操作系统申请的内存大小。| 
|used_memory_peak|redis的内存消耗峰值，以字节为单位，即历史使用记录中redis使用内存峰值。|
|used_memory_peak_human|以更直观的格式返回redis的内存消耗峰值| 
|used_memory_peak_perc|使用内存达到峰值内存的百分比|
|used_memory_overhead|Redis为了维护数据集的内部机制所需的内存开销，包括所有客户端输出缓冲区、查询缓冲区、AOF重写缓冲区和主从复制的backlog。|
|used_memory_startup|Redis服务器启动时消耗的内存|
|used_memory_dataset|数据实际占用的内存大小，即 used_memory-used_memory_overhead|
|used_memory_dataset_perc|数据占用的内存大小的百分比|
|total_system_memory|整个系统内存| 
|total_system_memory_human|以更直观的格式显示整个系统内存|
|used_memory_lua|Lua脚本存储占用的内存|
|used_memory_lua_human|以更直观的格式显示Lua脚本存储占用的内存| 
|maxmemory|Redis实例的最大内存配置| 
|maxmemory_human|以更直观的格式显示Redis实例的最大内存配置|
|maxmemory_policy|当达到maxmemory时的淘汰策略|
|mem_fragmentation_ratio|碎片率，used_memory_rss/ used_memory。ratio指数>1表明有内存碎片，越大表明越多，<1表明正在使用虚拟内存，虚拟内存其实就是硬盘，性能比内存低得多，这是应该增强机器的内存以提高性能。一般来说，mem_fragmentation_ratio的数值在1 ~ 1.5之间是比较健康的。|
|mem_allocator|内存分配器| 
|active_defrag_running|表示没有活动的defrag(表示内存碎片整理)任务正在运行，1表示有活动的defrag任务正在运行|
|lazyfree_pending_objects|0表示不存在延迟释放的挂起对象|
