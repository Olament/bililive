# LiveVUP.com

[LiveVUP.com](LiveVUP.com) 源码托管

## 部署

准备好 MongoDB 或者与 Monogo 兼容的数据库，在```docker-compose.yml```中添加

```yaml
version: '3'
services:
  nginx:
    image: nginx:1.15-alpine
    depends_on:
      - backend
      - frontend
    ports:
      - "80:80"
    volumes:
      - ./nginx:/etc/nginx/conf.d
  backend:
    build: ./backend
    environment:
      - GIN_MODE=release
      - DB_STR=<你的数据库连接字符串>
  frontend:
    build: ./frontend
```

编译并启动项目

```bash
cd bililive
docker-compose build
docker-compose up
```

## API

### 当前在线

```bash
GET livevup.com/api/online
```

#### JSON 响应

| 字段  	| 类型   	| 内容         	| 备注 	|
|-------	|--------	|--------------	|------	|
| count 	| Number 	| 当前在线人数 	|      	|
| list  	| Array  	| 主播列表     	|      	|

#### ```list``` 详情

| 字段                   	| 类型   	| 内容               	| 备注                    	|
|------------------------	|--------	|--------------------	|-------------------------	|
| roomid                 	| Number 	| 房间 ID            	|                         	|
| uid                    	| Number 	|                    	|                         	|
| uname                  	| String 	| 名字               	|                         	|
| title                  	| String 	| 直播间标题         	|                         	|
| usercover              	| String 	| 直播间封面         	|                         	|
| keyframe               	| String 	| 直播间关键帧       	|                         	|
| popularity             	| Number 	| 直播间人气         	|                         	|
| maxPopularity          	| Number 	| 直播间最高人气     	|                         	|
| livetime               	| String 	| 开播时间           	| RFC 3339 标准格式的时间 	|
| participant            	| Number 	| 参与人数           	|                         	|
| participantDuring10Min 	| Number 	| 过去10分钟参与人数 	|                         	|
| goldCoin               	| Number 	| 金瓜子             	|                         	|
| goldUser               	| Number 	| 付费用户数         	|                         	|
| silverCoin             	| Number 	| 银瓜子             	|                         	|
| danmuCount             	| Number 	| 弹幕数             	|                         	|

<details>
<summary>查看响应示例：</summary>

```json
{
    "count": 1277,
    "list": [
        {
            "roomid": 21919321,
            "uid": 508963009,
            "uname": "HiiroVTuber",
            "title": "学中文+写名字！！！",
            "usercover": "https://i0.hdslb.com/bfs/live/new_room_cover/b6b551cfa6db7bc3d392a50adef02144eeacbc70.jpg",
            "keyframe": "https://i0.hdslb.com/bfs/live-key-frame/keyframe06292128000021919321s8ajde.jpg",
            "popularity": 673834,
            "maxPopularity": 677370,
            "livetime": "2021-06-29T12:01:17.896328161Z",
            "participant": 4322,
            "participantDuring10Min": 583,
            "goldCoin": 3503950,
            "goldUser": 1302,
            "silverCoin": 2097900,
            "danmuCount": 13140
        },
        ...
    ]
}
```
</details>

### 个人资料

```bash
GET livevup.com/api/broadcast
```

#### URL 参数

| 参数名     | 类型 | 内容         | 
| ---------- | ---- | ------------ | 
| uid |  Number  | Bilibili 用户 UID | 

#### JSON 响应

| 字段  	| 类型   	| 内容         	| 备注 	|
|-------	|--------	|--------------	|------	|
| - 	| Array | 主播近十场直播 	|   Array of broadcast   	|

#### ```broadcast``` 详情

| 字段             	| 类型   	| 内容           	| 备注                                           	|
|------------------	|--------	|----------------	|------------------------------------------------	|
| title            	| String 	| 直播间标题     	|                                                	|
| uname            	| String 	| 名字           	|                                                	|
| maxPopularity    	| Number 	| 直播间最高人气 	|                                                	|
| livetime         	| String 	| 开播时间       	| RFC 3339 标准格式的时间                        	|
| endTime          	| String 	| 下播时间       	| RFC 3339 标准格式的时间                        	|
| participant      	| Number 	| 参与人数       	|                                                	|
| goldCoin         	| Number 	| 金瓜子         	|                                                	|
| goldUser         	| Number 	| 付费用户数     	|                                                	|
| silverCoin       	| Number 	| 银瓜子         	|                                                	|
| danmuCount       	| Number 	| 弹幕数         	|                                                	|
| participantTrend 	| Array  	| 同接趋势       	| 过去十分钟参与弹幕互动的人数，每五分钟采集一次 	|
| goldTrend        	| Array  	| 营收趋势       	| 主播累计收到的金瓜子，每五分钟采集一次         	|
| danmuTrend       	| Array  	| 弹幕趋势       	| 主播累计收到的弹幕，每五分钟采集一次           	|

<details>
<summary>查看响应示例：</summary>

```json 
[
    {
        "title": "【秋舍夜话】夏日天气物语",
        "uname": "秋凛子Rinco",
        "maxPopularity": 89239,
        "livetime": "2021-06-28T13:59:18.995Z",
        "endTime": "2021-06-28T16:49:52.537Z",
        "participant": 686,
        "goldCoin": 1213040,
        "goldUser": 311,
        "silverCoin": 523400,
        "danmuCount": 5739,
        "participantTrend": [
            73,
            121,
            121,
            104,
            96,
            135,
            138,
            129,
            97,
            88,
            100,
            97,
            112,
            107,
            72,
            83,
            104,
            102,
            96,
            105,
            96,
            88,
            85,
            134,
            157,
            136,
            122,
            102,
            98,
            101,
            108,
            113,
            96,
            168
        ],
        "goldTrend": [
            1500,
            3950,
            15350,
            26250,
            34250,
            40250,
            342250,
            381250,
            381250,
            381310,
            381340,
            388840,
            388840,
            388840,
            388840,
            399540,
            399540,
            399540,
            399540,
            399540,
            399540,
            399540,
            399640,
            414740,
            1008840,
            1206840,
            1207840,
            1207840,
            1207840,
            1207940,
            1209040,
            1210040,
            1210040,
            1213040
        ],
        "danmuTrend": [
            123,
            277,
            454,
            585,
            744,
            997,
            1172,
            1351,
            1470,
            1686,
            1843,
            1973,
            2167,
            2264,
            2372,
            2551,
            2718,
            2860,
            3007,
            3189,
            3342,
            3508,
            3625,
            4002,
            4206,
            4369,
            4505,
            4622,
            4729,
            4874,
            5038,
            5205,
            5424,
            5729
        ]
    },
    ...
]
```
</details>

### 排名

```bash
GET livevup.com/api/rank
```

#### URL 参数

| 参数名     | 类型 | 内容         | 详情 |
| ---------- | ---- | ------------ |  ---- |
| sortBy |  String  | 排名依据 | ```income```\|```viewership```\|```paid```\|```duration```|

#### JSON 响应

| 字段  	| 类型   	| 内容         	| 备注 	|
|-------	|--------	|--------------	|------	|
| - 	| Array | VUP 数据列表 	|   Array of VUP   	|

#### ```VUP``` 详情

| 字段           	| 类型   	| 内容         	| 备注                 	|
|----------------	|--------	|--------------	|----------------------	|
| uid            	| Number 	|              	|                      	|
| uname          	| String 	| 名字         	|                      	|
| face           	| String 	| 主播头像     	|                      	|
| duration       	| Number 	| 播出时长     	| 以小时为单位，浮点数 	|
| income         	| Number 	| 收入         	| 以元为单位，浮点数   	|
| danmuCount     	| Number 	| 弹幕总数     	|                      	|
| avgPaidUser    	| Number 	| 场均付费用户 	|                      	|
| avgParticipant 	| Number 	| 场均参与人数 	|                      	|
| avgViewership  	| Number 	| 场均同接     	|                      	|

<details>
<summary>查看响应示例：</summary>

```json
[
    {
        "uid": 33081544,
        "uname": "菜小仙_Channel",
        "face": "https://i1.hdslb.com/bfs/face/19aa214cf5bca096e65dca7ac557c5572d1dc0dc.jpg",
        "duration": 7.11578861111111,
        "income": 483475.4,
        "danmuCount": 2525,
        "avgPaidUser": 25.5,
        "avgParticipant": 122.5,
        "avgViewership": 16.226190476190474
    },
    ...
]
```
</details>
