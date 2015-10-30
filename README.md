# pando-cloud  [![Build Status](https://travis-ci.org/PandoCloud/pando-cloud.svg)](https://travis-ci.org/PandoCloud/pando-cloud)  [![Coverage Status](https://coveralls.io/repos/PandoCloud/pando-cloud/badge.svg?branch=master&service=github)](https://coveralls.io/github/PandoCloud/pando-cloud?branch=master)

Pando Cloud is the cloud part of Pando IoT solution.

[中文文档](docs/zh-cn/README.md)

## What is Pando?

Pando is an open source solution to help you build your IoT application from scratch. It's made of a bunch of tools, protocols and frameworks below:

* Pando Cloud
* [Pando Embeded Framework](https://github.com/PandoCloud/pando-embeded-framework)
* [Pando Protocol](https://github.com/PandoCloud/pando-protocol)
* Mobile SDK for [Android](https://github.com/PandoCloud/pando-android-sdk) and [IOS](https://github.com/PandoCloud/pando-ios-sdk)
* IoT Explorer APP FreeIOT for [Android](https://github.com/PandoCloud/freeiot-android) and [IOS](https://github.com/PandoCloud/freeiot-ios)
* IoT Kit [Tisan](https://github.com/tisan-kit)

## Design Concepts

Pando Cloud aims at Providing a production ready IoT Cloud Solution, NOT just an IoT server demo. Here are some system key points we focus on:

* **IoT**: should use best practice for IoT business scenarios.
* **stability**: should be well tested and bug free.
* **distributive**: can be deployed on cluster as well as on single machine.
* **availability**: mechanisms for fault tolerance, self healing, and more.
* **extensibility**: modular, pluggable, hookable and composable.
* **accessibility**: easy to learn, develop, and deploy.
* **universality**: highly optimized and well designed for general use.
* **efficiency**: lower comsuption on hardware and bandwidth.
* **security**: away of common invasion and cracking.
* **scalability**: capable to manage thousands to billions of IoT devices.

## Architecture Overview

![architecture](docs/img/architecture.jpeg)

The system can be divided into three layers.

### access
Access layer provides different protocol brokers to accept device connections, including but not limited to HTTP, XMPP, MQTT, COAP. HTTP is REQUIRED for device registration, authentication and discovery. We currently suport MQTT for device data exchange.
### logic
Logic layer implements IoT related business like device online, device management and api, etc.

* **registry**: registry serves system level meta data and configurations. 
* **devicemanger**: device manager is a service to processing device data and requests, and keepping device status.
* **apiprovider**: provide restful API for applications.
* **notifier**: notify application when device status changes.

### core 

We relay on a couple of open source solution as core services:

* **mongodb**: we use mongo as device data storage engine and benefits from its simplicity and scalebility. 
* **mysql**: global configuration like product and device info are stored in mysql.
* **redis**: widely used as cache service.
* **nsq**: message queue for asynchronous communication.
* **etcd**: server registration and discovery service.
* **ELK**: the elasticsearch, logstash and kibana stack privodes a efficient way for log collection and analytics.

## Methodology and Technology Reference

* [micro-services](http://martinfowler.com/articles/microservices.html)
* [12-factor apps](http://12factor.net/)
* [GoLang](http://golang.org)
* [MQTT](http://mqtt.org/)
* [Docker](http://www.docker.com/)
* [MySQL](http://www.mysql.com/)
* [MongoDB](https://www.mongodb.org/)
* [redis](http://redis.io/)
* [etcd](https://github.com/coreos/etcd)
* [nsq](http://nsq.io/)
* [ELK](https://www.elastic.co/products)