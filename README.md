# ChirpBird

Welcome to Haraj take home challenge!

This is my forked repo for my own solution. 

It requires `docker-compose` command available within your machine.

### How to run
- Take a look at [.env](.env) file on this repo root directory, change it to your host docker internal IP which can you get by `ifconfig` command on Linux. [Reason why this is needed](documentation/failsafe.md#demonstrate-the-outage)
- run `docker-compose build`
- run `docker-compose up`
- Access the static client web page on [http://localhost:4000](http://localhost:4000)

I've hosted the app itself on [https://harunalfat.site](https://harunalfat.site)

All of the rest of documentation of the project can be seen on [this page](documentation/README.md)