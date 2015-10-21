# core services to serve iot devices.

- **mqttaccess**: mqtt access service which accepts device mqtt connections. 
- **deviceapi**: device api service which offers device http apis like authentication, registration, etc.
- **controller**: core logic and route service.
- **onlinemanager**: keep device online information. 
- **apiprovidor**: http apis for applications. 
- **notifier**: service that notify device status changes to applications.
- **registry**: service that keep global configuration and info.
