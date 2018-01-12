# 报表分析配置文件格式
# 概览
* JSON格式
* 默认文件名：reports.json
* 示例：
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
# 格式描述
## reports
* 格式定义：
```json
{
  "reports": [
    "<report_1>":{
      ...
    },
    "<report_2>":{
      ...
    },
    ...
  ]
}
```
* 可以在同一个配置文件中定义一个或者多个report，名称不能重复
* 返回结果格式定义：
```json
{
  "<report_1>":[
    {
      <field1>:<val1>,
      <field2>:<val2>,
      ..
    },
    {
      <field1>:<val1_2>,
      <field2>:<val2_2>,
      ..
    },    
    ...
  ],
  "<report_2>":[
    {
      <field_1>:<val1>,
      <field_2>:<val2>,
      ..
    },
    {
      <field_1>:<val1_2>,
      <field_2>:<val2_2>,
      ..
    },    
    ...
  ],  
}
```
## report
### 格式定义：
* 数据来源：CSV / JSON
```json
"<report_1>":{
      "name": "<report_name>",
      "source": "csv|json,<csv/json_file_path>",
      "dimensions":"<field1>,<field2>,...<tag_name1>,<tag_name2>,...",
      "orderby":[
        "<field1>,DESC|ASC",
        "<field2>,DESC|ASC",
        ...
      ],
      "aggregates": [
        ["<func>","<field1>,<field1_alias>","<field2>,<field2_alias>,..."]
      ],
      "tags": {
        "<tag_name1>": [
          "<tag_name1_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name1_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
        "<tag_name2>": [
          "<tag_name2_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name2_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
      }
},
```
* 数据来源：mysql
```json
"<report_1>":{
      "name": "<report_name>",
      "source": "mysql,default",
      "dimensions":"<field1>,<field2>,...<tag_name1>,<tag_name2>,...",
      "orderby":[
        "<field1>,DESC|ASC",
        "<field2>,DESC|ASC",
        ...
      ],
      "aggregates": [
        ["<func>","<field1>,<field1_alias>","<field2>,<field2_alias>,..."]
      ],
      "tags": {
        "<tag_name1>": [
          "<tag_name1_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name1_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
        "<tag_name2>": [
          "<tag_name2_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name2_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
      }
},
```
* 数据来源：sqlite
```json
"<report_1>":{
      "name": "<report_name>",
      "source": "sqlite,<sqlite_db_path>",
      "dimensions":"<field1>,<field2>,...<tag_name1>,<tag_name2>,...",
      "orderby":[
        "<field1>,DESC|ASC",
        "<field2>,DESC|ASC",
        ...
      ],
      "aggregates": [
        ["<func>","<field1>,<field1_alias>","<field2>,<field2_alias>,..."]
      ],
      "tags": {
        "<tag_name1>": [
          "<tag_name1_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name1_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
        "<tag_name2>": [
          "<tag_name2_val1>, <fieldX>,REGEXP,<regular_expression>",
          "<tag_name2_val2>, <fieldX>,REGEXP,<regular_expression>",
          ...
        ],
      }
},
```