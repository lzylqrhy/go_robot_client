[TOC]

# go_robot_client

## 概述

这是一个人机器人客户端，它的作用：

- 服务端程序员自行验证服务端功能
- 数值人员跑数据
- 服务器压力测试

## 目录介绍

- ***core***：核心模块
  - ***mydb***：数据库模块，为各种数据库提供统一的操作接口
    - dbtypes：提供自定义数据类型，及类型转换逻辑
    - mymysql：由mysql实现操作接口
  - ***mynet***：提供统一的网络操作及消息解包，可以支持不同的协议，目前支持tcp和websocket
- common：公共功能和模块
- configs：各种配置文件
- games：具体的游戏模块，每个游戏一个子目录
  - fish：捕鱼
  - fruit：水果机
- global：全局数据及其操作逻辑
- protocols：游戏客户端和服务端之间通信的具体协议号及数据结构
- utile：工具函数

## 开发指南 

- 基于demo
  1. 在global/enum.go文件中，添加对应的游戏ID
  2. 在global/ini/目录下，添加新的配置
  3. 在games目录下，创建新的游戏目录，新建一个该游戏的client类，实现RobotClient接口，可从其他目录拷贝后修改
  4. 在Protocols目录下，创建该游戏的协议文件

- 自定义

  自己随意，参考demo，核心代码如下：

  ```go
  // 设置机器人数据
  games.SetRobotTestData(userList)
  // 启动机器人
  for i, user := range userList {
      ...
      // 网络连接
      d := mynet.NewConnect(cfg.NetProtocol, serverAddr)
      // 创建客户端
      c := games.NewClient(uint(i), user, d)
      // 开工
      core.DoWork(ctx, &wg, c, d, ini.GameCommonSetting.Frame)
  }
  ```

  

## 配置

### 主配置

- protocol：网络通信协议类型，目前支持tcp和websocket

- api_addr：机器人生成地址，根据不同的产品，修改对应的二级域名
- start：机器人ID起始位置，>= 2
- num：本次运行的机器人数量
- game_id：进入的游戏（非平台游戏ID，今后平台唯一后，可以保持一致）
- game_zone：平台设置的游戏类型ID

### 游戏配置

#### 通用

- frame：帧数，调用update()的频率

- server_addr：可以指定服务器地址，注释掉或空设置表示用平台下发的地址

- room_id：指定房间ID，不指定的话，会自动进入满足条件的房间
- [db_xx]：数据库配置

#### 捕鱼

- 是否开火

- 设置捕获指定类型的鱼，不指定表示所有鱼

- 是否攻击波塞冬

- 发射导弹方式

#### 水果机

- 必须指定room_id，否则无法进入房间

#### 阿拉丁

- 必须指定room_id，否则无法进入房间

