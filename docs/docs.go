// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/segment": {
            "post": {
                "description": "Create a new segment.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "create a new segment",
                "parameters": [
                    {
                        "type": "string",
                        "description": "admin",
                        "name": "User-role",
                        "in": "header"
                    },
                    {
                        "description": "Segment attributes",
                        "name": "segment_attrs",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/domain.Segment"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "403": {
                        "description": "Forbidden"
                    }
                }
            }
        },
        "/api/segment/{id}": {
            "delete": {
                "description": "Delete a segment.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Segment"
                ],
                "summary": "delete a new segment",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Segment ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "admin",
                        "name": "User-role",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "404": {
                        "description": "Not found"
                    }
                }
            }
        },
        "/api/user/get/{user_id}": {
            "get": {
                "description": "Get list of user` + "`" + `s segments.",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get list",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Segment ID",
                        "name": "user_id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "204": {
                        "description": "No content"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "404": {
                        "description": "Not found"
                    }
                }
            }
        },
        "/api/user/get_history": {
            "get": {
                "description": "Get file with history of user` + "`" + `s segments",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "get file with history",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Segment ID",
                        "name": "user_id",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Time start",
                        "name": "start",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "Time end",
                        "name": "end",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        },
        "/api/user/update": {
            "patch": {
                "description": "Update user` + "`" + `s segments.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "update user",
                "parameters": [
                    {
                        "description": "Update attributes",
                        "name": "segment_attrs",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/domain.UpdateRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad request"
                    },
                    "403": {
                        "description": "Forbidden"
                    },
                    "404": {
                        "description": "Not found"
                    },
                    "500": {
                        "description": "Internal server error"
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Segment": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "percent": {
                    "type": "number"
                },
                "ttl": {
                    "type": "string"
                }
            }
        },
        "domain.UpdateRequest": {
            "type": "object",
            "properties": {
                "segments_add": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Segment"
                    }
                },
                "segments_del": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.Segment"
                    }
                },
                "user_id": {
                    "type": "integer"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "API",
	Description:      "This is an auto-generated API Docs.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
