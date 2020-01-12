# Alkemi Switch Board Service

---

![API](https://img.shields.io/website?down_color=red&down_message=offline&label=API&up_color=green&up_message=online&url=https%3A%2F%2Fswitchboard.alkemi.tech)
[![Documentation](https://img.shields.io/badge/Docs-yes-blue)](https://alkemi.atlassian.net/wiki/spaces/ARCH/pages/65994893/Switchboard+Service)

![Alkemi Switch Board Service](https://raw.githubusercontent.com/project-alkemi/api-docs/master/assets/img/switchboard.gif)

## Service Overview

The Switchboard Service maintains a running whitelist for every exchange on the Alkemi Network and the assets that are being settled on it. The Switchboard Service is operated by the Alkemi team and is adjusted and modified as exchanges join the network.

- When an exchange joins the network for the first time the switchboard will generate an “Auth Token” to be used by the Exchange Containers when they first initiate a connection with the Data Relay and Price Index services. Without a valid auth token, the container will not be allowed to connect to the services to retrieve or submit data. This auth token will periodically be updated to protect against leakage.

- After connection, the Exchange Container will begin to send price data and completed order data to the Data Relay Service for every market it is configured to participate in. The Data Relay Service will retrieve this data and verify it against a whitelist on the Exchange Container for each asset it is submitting data for. If the auth token isn’t allowed to send that instrument data then it will be ignored and a flag will be reported that the container is sending data it isn’t allowed to.

- Once the data is verified through proper authentication the Data Relay Service passes the necessary information to the BookKeeper Nodes and the Price Index Service.

---

## Latest Deployment URLS

- [Switch Board Service API](https://switchboard-api.alkemi.tech/v1)
- [Switch Board Dashboard](https://switchboard.alkemi.tech)
  - User/Pass: admin/admin
- [Service Status](https://status.alkemi.tech/service/switch-board-service)

## Usage

### To Build

    go build bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

### To Run

    go run bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

### To Test

    go test bitbucket.org/project-alkemi/switch-board-service/test/switch-board-service

### To Clean

    go clean bitbucket.org/project-alkemi/switch-board-service/cmd/switch-board-service

### To Generate Documentation

    swag init -g cmd/switch-board-service/main.go --output docs/swagger/
