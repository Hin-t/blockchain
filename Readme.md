# Blockchain

### 1.实现区块结构以及与区块相关功能
1. 区块结构分析
2. 新建区块
3. 如何生成hash
4. 类型转换

### 2.实现区块链基本结构
1. 实现链表（通过切片进行缓存）-- 区块链的基本结构
2. 实现创世区块与区块链初始化功能
3. 实现上链功（区块添加到链上）

### 3.实现POW共识算法
1. pow结构分析
2. 设置目标难度值
3. hash碰撞
4. 数据准备

### 4.实现区块数据持久化
1. 安装：go get github.com/boltdb/bolt/...
2. 修改区块链结构
3. 定义数据和表名
4. bolt基本操作（insert，read）

### 5. 实现区块链遍历输出
1. 实现区块链的遍历输出函数PrintChain

### 6. go命令行
1. os.args
2. flags 

### 7. 获取blockchain对象实例

### 8. 比特币交易原理
1. 传统的web交易
   1. 账户
   2. 余额
   3. 参与者
   4. 货币
2. 基本概念
   1. 比特币系统中没有余额的概念，使用UTXO交易模型，在传统交易过程中所说的交易余额实际上指的是要给比特币钱包地址的UTXO集合
3. 交易组成
   1. 在比特币中，交易主要由输入，输出，ID（txHash），交易时间组成
4. UTXO交易模型
   1. 比特币专有的交易模型
   2. 在比特币中，交易实际上就是不断查找指定钱包地址的UTXO集合，然后进行修改的过程
   3. UTXO是比特币交易中最基本的单元，是不可拆分的，可以把UTXO理解成一个币，该币是拥有一个金额，一个拥有者
5. 交易过程
   1. 思考
      1. 如何保证用户只使用自己的比特币
      2. 如何保证一笔交易是有效的
   2. 流程图
   3. 交易分类
      1. coinbase：挖矿奖励的比特币，没有发送者，由系统提供，所以不包含input
      2. 普通转账：正常的转账交易，有发送者参与，包含input

### 9. 嵌入交易结构
1. 实现交易结构替换
2. 输入结构实现分析
3. 输出结构实现分析
4. 实现coinbase交易（没有输入）
   1. coinbase生成函数实现
   2. 交易哈希（序列化）的实现
   3. coinbase交易生成函数的调用
5. 实现CLI发起转账
   1. 添加命令行转账功能（封装）
   2. 实现JSON转为字符串切片功能
      1. 实现JSON2Array公共函数

### 10. 实现通过挖矿生成新的区块
1. 实现挖矿功能
2. 通过命令行send发起交易调用挖矿
3. 实现普通交易
   1. 实现生成普通交易函数
   2. 修改挖矿函数，调用NewSimpleTransaction（）
   3. 通过CLI实现普通转账交易调用
4. 实现余额查询与UTXO查询（封装）
   1. 实现余额查询cli端的封装
   2. 实现UTXO查找封装
   3. 实现先输入输出验证功能
5. UTXO查找的内部实现 
   1. 实现查找数据库指定地址所有已花费输出函数
   2. 实现coinbase交易判断函数
   3. 实现查找指定地址所有UTXO的函数
6. UTXO结构封装
   1.将所有与output相关的属性封装到一个结构
### 11. 转帐逻辑完善与UTXO查找优化
1. 实现查找可用UTXO的函数FindSpendableUTXO()
2. 多笔数据交易./src.exe send -from '[\"troytan\",\"Alice\"]' -to '[\"Bob\",\"troytan\"]' -amount  '[\"2\",\"1\"]'

### 12. base58编码
1. base58和base64编码区别
   1. 去掉'O'和'0','I'和'l'
   2. 去掉字符'+'和'/'，非字母数字作为账号难以接受
2. 实现go的base58
   1. 确定base58编码基数表
   2. 将传入的自负床转换为字节数组
   3. 将字节数组转换成big.int
   4. 对基数58求于，直到除数为0，将所有玉树对应作为索引在基数表中对应的自负床串联起来
   5. 实现base58编码函数、解码函数、切片反转

### 13. 钱包说明与实现
1. 比特币钱包，实际上管理的就是一个公钥-私钥的密钥对
2. 比特币地址特点
   1. version：版本前缀，1个字节，主要用来识别格式，前缀“1”代表公网上的比特币地址
   2. pubKey hash：20字节，代表公钥hash
   3. checksum：校验和，4个字节，表示添加到正在编码的数据的一段的校验字节，主要用于检测输入时产生的错误，该值通过pubKey得到
3. 实现
   1. 钱包结构
   2. 生成钱包对象
   3. 通过钱包生成密钥对

### 14. 通过钱包生成地址
1. 公钥哈希 sha256 - ripemd160
2. 添加version （在base58中已实现）
3. 添加checksum
4. checksum校验

### 15. 实现钱包集合
1. 实现集合基本结构
2. 实现集合对象生成
3. 实现创建钱包功能

### 16. 实例钱包模块集成
1. 将钱包创建功能加入命令行操作
2. 实现获取地址列表功能
3. 钱包功能持久化

### 17. 将钱包功能嵌入到区块链中
1. 添加命令行对钱包的操作
2. 实现获取地址列表功能
3. 钱包持久化功能

### 18. 将钱包与输入输出功能的结合
1. 实现输入结构与钱包功能相结合
2. 实现输出结构与钱包功能相结合
3. 调用方法相关修改

### 19. ecdse回顾&交易签名实现
1. 在什么时候对交易进行签名
   1. 在交易生成的时候进行签名
2. 在什么时候对交易进行验证
   1. 在交易被打包到区块之前进行验证
3. 是否需要对交易中所有属性进行交易和验证
   理论相关
      交易的基本概念
      交易实际就是解锁指定地址的output，然后重新分配他们的值，再重新 加锁到新的output中
   为了数据安全，必须加密数据如下
   1. 保存在已解锁的output的公钥哈希，代表交易的发送者
   2. 保存在新生成的output的公钥哈希， 代表交易的接收者
   3. 新生成的output包含的value

流程：
 1. 生成交易
2. 对交易进行签名
   1. 判断该交易是否是一笔coinbase，如果是，不签名
   2. 查找当前交易所引用的交易（输出所在的交易）
   3. 提取所需要的签名的属性
   4. 签名
3. 验证交易，验证交易的实现
4. 打包

### 20. 实现挖矿奖励
1. 默认谁发起交易，谁得到奖励（简化逻辑）
2. 修正生成交易hash的bug，添加时间戳作为哈希生成的标识

### 21.  实现UTXO查找优化
1. 分析当前区块链UTXO查找逻辑
2. 分析UTXO查找效率
3. 