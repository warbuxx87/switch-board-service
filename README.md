# Alkemi | Switchboard Service

![Version](https://img.shields.io/badge/version-0.1.0-blue.svg?cacheSeconds=2592000)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
![API](https://img.shields.io/website?down_color=red&down_message=offline&label=API&up_color=green&up_message=online&url=https%3A%2F%2Fswitchboard.alkemi.tech)
[![Documentation](https://img.shields.io/badge/documentation-yes-brightgreen.svg)](https://alkemi.atlassian.net/wiki/spaces/ARCH/pages/12648485/Switch+Board+Service)
[![Twitter: AlkemiOfficial](https://img.shields.io/twitter/follow/AlkemiOfficial.svg?style=social)](https://twitter.com/AlkemiOfficial)

![Alkemi Switchboard Service](https://raw.githubusercontent.com/project-alkemi/api-docs/master/assets/img/switchboard.gif)

---

## Service Overview

The Switchboard Service maintains a running whitelist for every exchange on the Alkemi Network and the assets that are being settled on it. The Switchboard Service is operated by the Alkemi team and is adjusted and modified as exchanges join the network.

- When an exchange joins the network for the first time the switchboard will generate an “Auth Token” to be used by the Exchange Containers when they first initiate a connection with the Data Relay and Price Index services. Without a valid auth token, the container will not be allowed to connect to the services to retrieve or submit data. This auth token will periodically be updated to protect against leakage.

- After connection, the Exchange Container will begin to send price data and completed order data to the Data Relay Service for every market it is configured to participate in. The Data Relay Service will retrieve this data and verify it against a whitelist on the Exchange Container for each asset it is submitting data for. If the auth token isn’t allowed to send that instrument data then it will be ignored and a flag will be reported that the container is sending data it isn’t allowed to.

- Once the data is verified through proper authentication the Data Relay Service passes the necessary information to the BookKeeper Nodes and the Price Index Service.

---

## Latest Deployment URLS

- [Switchboard Service API](https://switchboard-api.alkemi.tech/v1)
- [Switchboard Dashboard](https://switchboard.alkemi.tech)
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

## Author

👤 **Alkemi**

- Website: <https://alkemi.network>
- Twitter: [@AlkemiOfficial](https://twitter.com/AlkemiOfficial)
- Github: [@project-alkemi](https://github.com/project-alkemi)

## Issues

Feel free to check [issues page](https://alkemi.atlassian.net/jira/software/projects/CD/boards/4).

If you come across an issue with the Switchboard Service, do a search in the [Issues](https://alkemi.atlassian.net/jira/software/projects/CD/boards/4) tab of this repo to make sure it hasn't been reported before. Follow these steps to help us prevent duplicate issues and unnecessary notifications going to the many people watching this repo:

- If the issue you found has been reported and is still open, and the details match your issue, give a "thumbs up" to the relevant posts in the issue thread to signal that you have the same issue. No further action is required on your part.
- If the issue you found has been reported and is still open, but the issue is missing some details, you can add a comment to the issue thread describing the additional details.
- If the issue you found has been reported but has been closed, you can comment on the closed issue thread and ask to have the issue reopened because you are still experiencing the issue. Alternatively, you can open a new issue, reference the closed issue by number or link, and state that you are still experiencing the issue. Provide any additional details in your post so we can better understand the issue and how to fix it.

## 🤝 Contributing

Contributions, issues and feature requests are welcome! Please be sure to abide by our [Code Of Conduct](CODE_OF_CONDUCT.md)

## 📝 License

Copyright © 2020 [Alkemi Network](https://github.com/project-alkemi).

This project is [MIT](https://opensource.org/licenses/MIT) licensed.

MIT License

Copyright (c) 2020 Alkemi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
