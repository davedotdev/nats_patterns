# README

These configurations are for decentralised auth with a stateless backbone and stateful edge.

### Explanation
This topology is a set of layers.

```mermaid
flowchart
    direction TB

    subgraph Backbone
        direction LR
        BB-US-EAST
        BB-US-WEST
        BB-UK
        
        BB-US-WEST --- BB-US-EAST
        BB-US-EAST --- BB-UK
        BB-US-WEST --- BB-UK

    end

    subgraph App_Pricing
        direction LR                
        JS_US_EAST
        JS_US_WEST
        JS_EU_LON
        
        JS_US_WEST -.Cluster Route.- JS_US_EAST
        JS_US_EAST -.Cluster Route.- JS_EU_LON
        JS_US_WEST -.Cluster Route.- JS_EU_LON
        
        
    end

    Backbone -.gateway.- App_Pricing
```

[Y] One outcome is for the JS enabled clusters to communicate with backbone producers and consumers over gateway connections.

[Y] A goal is for a stateful JS cluster to be provisioned without any configuration changes to the backbone.

[Y] A security team can manage credentials and therefore access to all NATS clusters, whether stateless or stateful.
