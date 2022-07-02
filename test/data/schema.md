# Golang

```go
package test

type Mapping struct {
	Type  string
	Value string
}
type Spec []struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	DataType struct {
		Type    string  `json:"type"`
		Mapping Mapping `json:"mapping"`
	} `json:"dataType"`
}
type Define struct {
	Type  string `json:"type"`
	Specs []Spec `json:"specs"`
}
type Property struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Mode     string `json:"mode"`
	Define   Define `json:"define"`
	Required bool   `json:"required"`
}
type Param struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Define Define `json:"define"`
}
type Event struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Type     string  `json:"type"`
	Params   []Param `json:"params"`
	Required bool    `json:"required"`
}
type InDefine struct {
	Type    string  `json:"type"`
	Mapping Mapping `json:"mapping"`
}
type Input struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Define InDefine `json:"define"`
}
type OutDefine struct {
	Type    string  `json:"type"`
	Mapping Mapping `json:"mapping"`
}
type Output struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Define OutDefine `json:"define"`
}
type Action struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Input    []Input  `json:"input"`
	Output   []Output `json:"output"`
	Required bool     `json:"required"`
}
type Profile struct {
	ProductID  string `json:"ProductId"`
	CategoryID string `json:"CategoryId"`
}
type Schema struct {
	Version        string     `json:"version"`        // 版本
	Profile        Profile    `json:"profile"`        // 物的元描述
	Configurations []Property `json:"configurations"` // 配置
	Properties     []Property `json:"properties"`     // 属性
	Events         []Event    `json:"events"`         // 事件
	Actions        []Action   `json:"actions"`        // 动作
}

```

## JSON

```json
{
    "version": "1.0",
    "properties": [
        {
            "id": "switchers",
            "name": "继电器开关状态",
            "desc": "控制器的开关",
            "mode": "rw",
            "define": {
                "type": "struct",
                "specs": [
                    {
                        "id": "sw1",
                        "name": "sw1",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw2",
                        "name": "sw2",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw3",
                        "name": "sw3",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw4",
                        "name": "sw4",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw5",
                        "name": "sw5",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw6",
                        "name": "sw6",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw7",
                        "name": "sw7",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    },
                    {
                        "id": "sw8",
                        "name": "sw8",
                        "dataType": {
                            "type": "bool",
                            "mapping": {
                                "0": "关",
                                "1": "开"
                            }
                        }
                    }
                ]
            },
            "required": false
        }
    ],
    "events": [
        {
            "id": "status",
            "name": "状态上报",
            "desc": "",
            "type": "info",
            "params": [
                {
                    "id": "sw1",
                    "name": "sw1",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw2",
                    "name": "sw2",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw3",
                    "name": "sw3",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw4",
                    "name": "sw4",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw5",
                    "name": "sw5",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw6",
                    "name": "sw6",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw7",
                    "name": "sw7",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw8",
                    "name": "sw8",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                }
            ],
            "required": false
        }
    ],
    "actions": [
        {
            "id": "control",
            "name": "控制",
            "desc": "",
            "input": [
                {
                    "id": "sw1",
                    "name": "sw1",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw2",
                    "name": "sw2",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw3",
                    "name": "sw3",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw4",
                    "name": "sw4",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw5",
                    "name": "sw5",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw6",
                    "name": "sw6",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw7",
                    "name": "sw7",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                },
                {
                    "id": "sw8",
                    "name": "sw8",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                }
            ],
            "output": [
                {
                    "id": "result",
                    "name": "result",
                    "define": {
                        "type": "bool",
                        "mapping": {
                            "0": "关",
                            "1": "开"
                        }
                    }
                }
            ],
            "required": false
        }
    ],
    "profile": {
        "ProductId": "Y0ST19XLP1",
        "CategoryId": "1"
    }
}
```

## 证书

为每个创建的产品分配唯一标识 ProductID，用户可以自定义 Devicename 标识设备，用产品标识 + 设备标识 + 设备证书/密钥来验证设备的合法性。
