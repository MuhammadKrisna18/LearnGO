package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login": {
            "post": {
                "description": "Authenticates a user and returns a JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "User Login",
                "parameters": [
                    {
                        "description": "Login Credentials",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_modules_auth_domain.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/siakad-pro_internal_modules_auth_domain.LoginResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                        }
                    }
                }
            }
        },
        "/auth/me": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Returns the profile of the user based on the provided JWT token",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Get User Profile",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "allOf": [
                                {
                                    "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                                },
                                {
                                    "type": "object",
                                    "properties": {
                                        "data": {
                                            "$ref": "#/definitions/siakad-pro_internal_modules_auth_domain.UserProfileResponse"
                                        }
                                    }
                                }
                            ]
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/siakad-pro_internal_shared_response.APIResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "siakad-pro_internal_modules_auth_domain.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "siakad-pro_internal_modules_auth_domain.LoginResponse": {
            "type": "object",
            "properties": {
                "role": {
                    "type": "string"
                },
                "token": {
                    "type": "string"
                }
            }
        },
        "siakad-pro_internal_modules_auth_domain.UserProfileResponse": {
            "type": "object",
            "properties": {
                "created_at": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "role": {
                    "type": "string"
                }
            }
        },
        "siakad-pro_internal_shared_response.APIResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "errors": {},
                "message": {
                    "type": "string"
                },
                "success": {
                    "type": "boolean"
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

var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "SIAKAD Pro API",
	Description:      "This is a sample API for a SIAKAD Pro built with Go Fiber.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
