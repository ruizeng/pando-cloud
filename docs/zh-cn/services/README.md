# 服务部署指南

Pando物联网云平台采用微服务架构。除了系统依赖的MySQL，MongoDB等开源服务，平台自身目前有以下几个服务组成：

- **[registry](registry.md)**：维护平台全局配置和信息的服务。
- **[devicemanager](devicemanager.md)**: 设备信息和设备状态维护。
- **[controller](controller.md)**: 和设备进行信息交互的路由服务。
- **[apiprovider](apiprovider.md)**: 为应用提供REST接口。
- **[httpaccess](httpaccess.md)**: 设备API服务，提供设备登陆、注册等逻辑。
- **[mqttaccess](mqttaccess.md)**：MQTT接入服务。