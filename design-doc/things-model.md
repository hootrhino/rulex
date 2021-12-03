下面描述一个支持  MQTT  协议的设备:
```yml
name: "my-custom-device-profile"
manufacturer: "iot"
model: "MQTT-DEVICE"
description: "Test device profile"
labels:
  - "mqtt"
  - "test"
deviceResources:
  -
    name: randnum
    isHidden: true
    description: "device random number"
    properties:
      valueType: "Float32"
      readWrite: "R"
  -
    name: ping
    isHidden: true
    description: "device awake"
    properties:
      valueType: "String"
      readWrite: "R"
  -
    name: message
    isHidden: false
    description: "device message"
    properties:
      valueType: "String"
      readWrite: "RW"

deviceCommands:
  -
    name: values
    readWrite: "R"
    isHidden: false
    resourceOperations:
        - { deviceResource: "randnum" }
        - { deviceResource: "ping" }
        - { deviceResource: "message" }
```

下面描述一个支持 ModBus 的设备
```yml
name: "Ethernet-Temperature-Sensor"
manufacturer: "Audon Electronics"
model: "Temperature"
labels:
  - "Web"
  - "Modbus TCP"
  - "SNMP"
description: "The NANO_TEMP is a Ethernet Thermometer measuring from -55°C to 125°C with a web interface and Modbus TCP communications."

deviceResources:
  -
    name: "ThermostatL"
    isHidden: true
    description: "Lower alarm threshold of the temperature"
    attributes:
      { primaryTable: "HOLDING_REGISTERS", startingAddress: 3999, rawType: "Int16" }
    properties:
      valueType: "Float32"
      readWrite: "RW"
      scale: "0.1"
  -
    name: "ThermostatH"
    isHidden: true
    description: "Upper alarm threshold of the temperature"
    attributes:
      { primaryTable: "HOLDING_REGISTERS", startingAddress: 4000, rawType: "Int16" }
    properties:
      valueType: "Float32"
      readWrite: "RW"
      scale: "0.1"
  -
    name: "AlarmMode"
    isHidden: true
    description: "1 - OFF (disabled), 2 - Lower, 3 - Higher, 4 - Lower or Higher"
    attributes:
      { primaryTable: "HOLDING_REGISTERS", startingAddress: 4001 }
    properties:
      valueType: "Int16"
      readWrite: "RW"
  -
    name: "Temperature"
    isHidden: false
    description: "Temperature x 10 (np. 10,5 st.C to 105)"
    attributes:
      { primaryTable: "HOLDING_REGISTERS", startingAddress: 4003, rawType: "Int16" }
    properties:
      valueType: "Float32"
      readWrite: "R"
      scale: "0.1"

deviceCommands:
  -
    name: "AlarmThreshold"
    readWrite: "RW"
    isHidden: false
    resourceOperations:
      - { deviceResource: "ThermostatL" }
      - { deviceResource: "ThermostatH" }
  -
    name: "AlarmMode"
    readWrite: "RW"
    isHidden: false
    resourceOperations:
      - { deviceResource: "AlarmMode", mappings: { "1":"OFF","2":"Lower","3":"Higher","4":"Lower or Higher"} }
```