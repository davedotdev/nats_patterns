# README

This is a fun one!

A leaf node acting as a bridge, between two three-node clusters

This means that each cluster is leaf-noded into as a remote and the leaf node has transitive properties.

Imports and exports make your brain hurt because of the directionality. Diagram below illustrated the boundings.

![fig1](./fig1.png)

You'll need the NATS server binary and the latest NATS CLI. Create some contexts, one for the green, one for blue and one for both red accounts.

### Requirements

__core.cloud.>__
1. Publishes made on green on `core.cloud.>` can be subscribed to in blue
Note, any other service can publish to `core.cloud.>` to the red accounts or blue.

2. Publishes made on red or blue to `core.cloud.>` are ignored by green

Lockdown connection credentials for further security.


__onsite.123-123.>__
1. Publishes made on blue to `onsite.123-123.>` can be subscribed to by red (UPSTREAM & DOWNSTREAM) and green
2. Publishes made on green and red[UPSTREAM] are ignored by blue

Lockdown connection credentials for further security.