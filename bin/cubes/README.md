# Cubes 命令行
## 编译
```sh
cd $GOPATH/src/gitlab.xibao100.com/skyline/skyline/cubes/bin/cubes
go build
```
## Web Server 模式
```sh
cubes.exe --mode=server [--port=9100] [--log=./cubes.log]
```
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

## 命令行数据分析模式(默认)
* 命令行格式：
```sh
cubes.exe [-mode=reports] [--log=./cubes.log] --jsoncfg=e:/reports.json [--output=./output.json]
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
  
* 返回结果(output.json)内容示例：
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
## 数据导入模式
* 命令行格式：
```sh
cubes.exe --mode=import --file=e:/aaa.csv --db=./test.db --table=aaa
```
* 用途：将csv/json文件分批导入sqlite数据库表中，用于后续数据分析

## 数据导出模式
* 命令行格式：
```sh
cubes.exe --mode=export --db=./test.db --table=aaa [--file=e:/aaa.json] 
```
* 用途：将sqlite数据库表数据导出到json文件中，不指定文件时默认输出文件：
```sh
./<table>.json
```


