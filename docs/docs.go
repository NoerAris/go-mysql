// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "support@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/user/deleteUser/{id}": {
            "delete": {
                "description": "Delete a user by id was input",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Delete User"
                ],
                "summary": "Delete a user",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Delete user by id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ResponseUser"
                        }
                    }
                }
            }
        },
        "/user/findAll": {
            "get": {
                "description": "Get all user",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get All User"
                ],
                "summary": "Get all user",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ResponseUser"
                        }
                    }
                }
            }
        },
        "/user/findOne": {
            "get": {
                "description": "Tes url with query param",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Test Param"
                ],
                "summary": "Tes url with query param",
                "parameters": [
                    {
                        "type": "string",
                        "description": "First Name",
                        "name": "firstName",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Last Name",
                        "name": "lastName",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/user/findOne/{id}": {
            "get": {
                "description": "Get user by id was input",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Get User"
                ],
                "summary": "Get a user by id",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Get user by id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ResponseUser"
                        }
                    }
                }
            }
        },
        "/user/saveUser": {
            "post": {
                "description": "Create a new user with the input paylod",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Save User"
                ],
                "summary": "Create a new user",
                "parameters": [
                    {
                        "description": "Create user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ResponseUser"
                        }
                    }
                }
            }
        },
        "/user/updateUser/{id}": {
            "put": {
                "description": "Update a user with the id and input paylod",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Update User"
                ],
                "summary": "Update a user",
                "parameters": [
                    {
                        "description": "Update user",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/user.User"
                        }
                    },
                    {
                        "type": "integer",
                        "description": "Update user",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/user.ResponseUser"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "user.ResponseUser": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/user.User"
                    }
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "user.User": {
            "type": "object",
            "properties": {
                "age": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:9093",
	BasePath:    "/",
	Schemes:     []string{},
	Title:       "CRUD MySql API",
	Description: "This is a service for CRUD",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
