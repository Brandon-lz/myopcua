// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "email": "advanced_to@163.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/v1/": {
            "get": {
                "description": "根路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "summary": "根路由",
                "responses": {
                    "200": {
                        "description": "欢迎使用OPC-UA OpenAPI",
                        "schema": {
                            "$ref": "#/definitions/routers.ApiResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/opc-node/add-node-to-read": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "AddNodeToRead 路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nodes"
                ],
                "summary": "AddNodeToRead 路由",
                "parameters": [
                    {
                        "description": "见下方JSON",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/noderouters.AddNodeToReadRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "节点添加成功",
                        "schema": {
                            "$ref": "#/definitions/noderouters.AddNodeToReadResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/opc-node/delete-node/{id}": {
            "delete": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "DeleteNode 路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nodes"
                ],
                "summary": "DeleteNode 路由",
                "parameters": [
                    {
                        "type": "string",
                        "description": "节点ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "节点删除成功",
                        "schema": {
                            "$ref": "#/definitions/noderouters.DeleteNodeResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/opc-node/get-node/{id}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "GetNode 路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "nodes"
                ],
                "summary": "GetNode 路由",
                "parameters": [
                    {
                        "type": "string",
                        "description": "节点ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "节点信息",
                        "schema": {
                            "$ref": "#/definitions/noderouters.GetNodeResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/ping": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "description": "ping 路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "summary": "ping 路由",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/api/v1/webhook/condition": {
            "post": {
                "description": "# 创建触发条件\n## 请求参数\n| 参数名称 | 类型 | 必填 | 描述 |\n| --- | --- | --- | --- |\n| and | []Condition | 否 | 规则列表，逻辑与 |\n| or | []Condition | 否 | 规则列表，逻辑或 |\n| rule | Rule | 否 | 规则 |\n*注意：Condition是嵌套类型，Condition包含and，or，rule，所以and里面可以嵌套and。。。无限嵌套*",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhook"
                ],
                "summary": "创建触发条件",
                "parameters": [
                    {
                        "description": "创建条件",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.CreateConditionRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.CreateConditionResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/webhook/config": {
            "post": {
                "description": "# 配置一条新的webhook\n\n## 例1：当节点node1值等于123时，发送通知到http://localhost:8080/api/v1/webhook/example\n` + "`" + `` + "`" + `` + "`" + `json\n{\n\"active\": true,\n\"name\":\"webhook1\",\n\"url\": \"http://localhost:8080/api/v1/webhook/example\",\n\"when\": {\n\"rule\": {\n\"node_name\": \"node1\",\n\"type\": \"eq\",\n\"value\": \"123\"\n}\n}\n}\n` + "`" + `` + "`" + `` + "`" + `\n使用when字段会创建新的条件condition，并将其配置在这个webhook上\n## 例2：使用已经配置好的条件condition\n` + "`" + `` + "`" + `` + "`" + `json\n{\n\"active\": true,\n\"url\": \"http://localhost:8080/api/v1/webhook/example\",\n\"condition_id\": 10\n}\n` + "`" + `` + "`" + `` + "`" + `\n## 常见异常\n- \"code\": 2007 代表数据重复，不能创建重复的webhook，具体重复了哪个字段，请看ConstraintName最后一个下划线后面的字段名\n- \"code\": 400 \"json: cannot unmarshal string into Go struct field AddWebhookConfigRequest.condition_id of type int64\"  ： 看下body参数，数字类型传成了字符串",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhook"
                ],
                "summary": "配置一条新的webhook",
                "parameters": [
                    {
                        "description": "Webhook configuration",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.AddWebhookConfigRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.AddWebhookConfigResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/webhook/config-by-name/{name}": {
            "get": {
                "description": "根据名称获取webhook配置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhook"
                ],
                "summary": "根据名称获取webhook配置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "webhook名称",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.GetWebhookConfigByNameResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/webhook/config/{id}": {
            "get": {
                "description": "根据id获取webhook配置",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhook"
                ],
                "summary": "根据id获取webhook配置",
                "parameters": [
                    {
                        "type": "string",
                        "description": "webhook id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.GetWebhookConfigByIdResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/webhook/example": {
            "post": {
                "description": "webhook示例",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Webhook"
                ],
                "summary": "webhook示例",
                "parameters": [
                    {
                        "description": "Webhook example",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.WebHookExampleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/webhookrouters.WebHookExampleResponse"
                        }
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "healthCheck 路由",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "default"
                ],
                "summary": "healthCheck 路由",
                "responses": {
                    "200": {
                        "description": "pong",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "globaldata.Condition": {
            "type": "object",
            "properties": {
                "and": {
                    "description": "规则列表，逻辑与",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/globaldata.Condition"
                    }
                },
                "or": {
                    "description": "规则列表，逻辑或",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/globaldata.Condition"
                    }
                },
                "rule": {
                    "description": "规则",
                    "allOf": [
                        {
                            "$ref": "#/definitions/globaldata.Rule"
                        }
                    ]
                }
            }
        },
        "globaldata.Rule": {
            "type": "object",
            "required": [
                "node_name",
                "type"
            ],
            "properties": {
                "node_name": {
                    "description": "节点名称",
                    "type": "string",
                    "example": "MyVariable"
                },
                "type": {
                    "description": "规则类型 eq, ne, gt, lt, all-time : 相等, 不相等, 大于, 小于, 全时间",
                    "type": "string",
                    "enum": [
                        "eq",
                        "ne",
                        "gt",
                        "lt",
                        "all-time"
                    ],
                    "example": "eq"
                },
                "value": {
                    "description": "规则value",
                    "type": "string",
                    "example": "123"
                }
            }
        },
        "noderouters.AddNodeToReadRequest": {
            "type": "object",
            "required": [
                "name",
                "node-id"
            ],
            "properties": {
                "data-type": {
                    "type": "string",
                    "example": "Int32"
                },
                "name": {
                    "type": "string",
                    "example": "MyVariable"
                },
                "node-id": {
                    "type": "string",
                    "example": "ns=2;i=2"
                }
            }
        },
        "noderouters.AddNodeToReadResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/noderouters.OpcNodeOutput"
                },
                "message": {
                    "type": "string",
                    "example": "节点添加成功"
                }
            }
        },
        "noderouters.DeleteNodeResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "type": "string"
                },
                "message": {
                    "type": "string",
                    "example": "节点删除成功"
                }
            }
        },
        "noderouters.GetNodeResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/noderouters.OpcNodeWithDataOutput"
                },
                "message": {
                    "type": "string",
                    "example": "节点信息"
                }
            }
        },
        "noderouters.OpcNodeOutput": {
            "type": "object",
            "properties": {
                "data-type": {
                    "type": "string",
                    "example": "Int32"
                },
                "name": {
                    "type": "string",
                    "example": "MyVariable"
                },
                "node-id": {
                    "type": "string",
                    "example": "ns=2;s=MyVariable"
                }
            }
        },
        "noderouters.OpcNodeWithDataOutput": {
            "type": "object",
            "properties": {
                "data-type": {
                    "type": "string",
                    "example": "Int32"
                },
                "name": {
                    "type": "string",
                    "example": "MyVariable"
                },
                "node-id": {
                    "type": "string",
                    "example": "ns=2;s=MyVariable"
                },
                "value": {
                    "type": "string",
                    "example": "123"
                }
            }
        },
        "routers.ApiResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "欢迎使用OPC-UA OpenAPI"
                }
            }
        },
        "webhookrouters.AddWebhookConfigRequest": {
            "type": "object",
            "required": [
                "url"
            ],
            "properties": {
                "active": {
                    "description": "是否激活，不传的话默认true",
                    "type": "boolean",
                    "example": true
                },
                "condition_id": {
                    "description": "条件id，不传的话默认新增条件",
                    "type": "integer",
                    "example": 1
                },
                "name": {
                    "description": "webhook名称，可以为空",
                    "type": "string",
                    "example": "webhook1"
                },
                "url": {
                    "description": "webhook地址",
                    "type": "string",
                    "example": "http://192.168.1.1:8800/notify"
                },
                "when": {
                    "description": "触发条件，为空时相当于通知所有数据变化",
                    "allOf": [
                        {
                            "$ref": "#/definitions/globaldata.Condition"
                        }
                    ]
                }
            }
        },
        "webhookrouters.AddWebhookConfigResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/webhookrouters.WebHookConfigRead"
                },
                "message": {
                    "type": "string",
                    "example": "节点添加成功"
                }
            }
        },
        "webhookrouters.CreateConditionRequest": {
            "type": "object",
            "properties": {
                "and": {
                    "description": "规则列表，逻辑与",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/globaldata.Condition"
                    }
                },
                "or": {
                    "description": "规则列表，逻辑或",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/globaldata.Condition"
                    }
                },
                "rule": {
                    "description": "规则",
                    "allOf": [
                        {
                            "$ref": "#/definitions/globaldata.Rule"
                        }
                    ]
                }
            }
        },
        "webhookrouters.CreateConditionResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/webhookrouters.WebHookConditionRead"
                },
                "message": {
                    "type": "string",
                    "example": "Condition created successfully"
                }
            }
        },
        "webhookrouters.GetWebhookConfigByIdResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/webhookrouters.WebHookConfigRead"
                },
                "message": {
                    "type": "string",
                    "example": "Webhook configuration get successfully"
                }
            }
        },
        "webhookrouters.GetWebhookConfigByNameResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "$ref": "#/definitions/webhookrouters.WebHookConfigRead"
                },
                "message": {
                    "type": "string",
                    "example": "Webhook configuration get successfully"
                }
            }
        },
        "webhookrouters.WebHookConditionRead": {
            "type": "object",
            "required": [
                "condition",
                "created_at",
                "id",
                "updated_at"
            ],
            "properties": {
                "condition": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "webhookrouters.WebHookConfigRead": {
            "type": "object",
            "required": [
                "active",
                "created_at",
                "id",
                "name",
                "updated_at",
                "url"
            ],
            "properties": {
                "active": {
                    "type": "boolean"
                },
                "condition_id": {
                    "type": "integer"
                },
                "created_at": {
                    "type": "string"
                },
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "url": {
                    "type": "string"
                },
                "when": {
                    "type": "string"
                }
            }
        },
        "webhookrouters.WebHookExampleRequest": {
            "type": "object",
            "properties": {
                "node_id": {
                    "description": "节点id",
                    "type": "string",
                    "example": "ns=1;s=MyVariable"
                },
                "node_name": {
                    "description": "节点名称",
                    "type": "string",
                    "example": "MyVariable"
                },
                "value": {
                    "description": "入参示例",
                    "type": "string",
                    "example": "123"
                }
            }
        },
        "webhookrouters.WebHookExampleResponse": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 200
                },
                "data": {
                    "type": "string",
                    "example": "webhook example"
                },
                "message": {
                    "type": "string",
                    "example": "webhook example success"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "OPC-UA Open API",
	Description:      "OPC-UA转http协议",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
