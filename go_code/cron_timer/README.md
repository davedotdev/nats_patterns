# README

This NATS micro service is a really tiny cron job service. Does it scale? No! Is it prod ready? Absolutely not. Will that stop you? Probably not.

1. Start the service as a binary

OR

2. `seq 3 | xargs -P3 -I {} go run .` to start up three

Then use `nats micro stats timerService` to show the load balancing etc and use `nats micro info TimerService` to get info.

### To use the service

Set the `duration` and `respond` fields, then pass the JSON into the request/reply against the service.

You can use normalised time like `30s`, `1h` etc.

`jo duration=3s respond=cron-alert | nats req timer.set`