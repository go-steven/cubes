# Cubes Web Server for Data Reports
## 编译
```sh
cd $GOPATH/src/gitlab.xibao100.com/skyline/skyline/cubes/bin/win_service
go build -o cubes_service.exe
```

## Cubes Web Server Windows服务
以管理员身份运行以下命令：
* 安装服务(默认端口：9100)：
```sh
cubes_service.exe [--port=9100] install
```
* 启动服务：
```sh
cubes_service.exe start
```
* 停止服务：
```sh
cubes_service.exe stop
```
* 卸载服务：
```sh
cubes_service.exe uninstall
```
## Cubes Web Server
### API:
1. http://localhost:9100/hello
* HTTP GET
* 返回：
 ```json
  {"res":"Cubes web server works."}
 ```
* 用途：用于测试Cubes Web Server是否正常工作

2. http://localhost:9100/rpts
* HTTP POST:
* 示例：
```json
{
  "json_cfg":"{\"reports\":[{\"name\":\"keyword_report\",\"source\":\"csv,e:/aaa.csv\",\"dimensions\":\"tag1, tag2\",\"aggregates\":[[\"SUM\",\"f8,TotalClick\",\"f9,TotalCost\"]],\"tags\":{\"tag1\":[\"EC, f2,REGEXP,.*EC.*\",\"FC, f2,REGEXP,.*FC.*\"],\"tag2\":[\"取暖器, f3,REGEXP,.*取暖器.*\",\"暖风机, f3,REGEXP,.*暖风机.*\"]}}]}"
}
```
* reports.json示例
```json
{
  "reports": [
    {
      "name": "keyword_report",
      "source": "csv,e:/auto_keyword_report_2017-03-05_to_2017-03-11.csv",
      "dimensions":"tag1,tag2",
      "aggregates": [
        ["SUM","f8,TotalClick","f9,TotalCost"]
      ],
      "tags": {
        "tag1": [
          "EC, f2,REGEXP,.*EC.*",
          "FC, f2,REGEXP,.*FC.*"
        ],
        "tag2": [
          "取暖器, f3,REGEXP,.*取暖器.*",
          "暖风机, f3,REGEXP,.*暖风机.*"
        ]
      }
    }
  ]
}
```
  
* 返回示例：
```json
{
       "keyword_report": [
           {
               "TotalClick": "2310",
               "TotalCost": "354342",
               "tag1": "",
               "tag2": ""
           },
           ...
       ]      
   }
```
* 用途：数据分析，返回分析的结果


   







