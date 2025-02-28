# README

Projects in here are for running against NATS.io. 

__100msConsumer__ is a setup to test a 100mS period NATS consumers. This was an experiment for a game. I start off by trying out a consumer with a timeout, then iterate to 100mS checks.

__Cron_Timer__ is a simple NATS micro service that behaves as a cron timer. Access the service over a `req` and feed it a JSON structure with `duration` and `respond` fields.

