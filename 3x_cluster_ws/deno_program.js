import {connect, headers} from "https://raw.githubusercontent.com/nats-io/nats.ws/main/src/mod.ts"

const nc = await connect({servers: "ws://localhost:8080", user: "acc1", pass: "password"})
console.log(`connected`)


// create 8MB buffer and fill it with a
const buf = new Uint8Array(8 * 1000000);
for (let i = 0; i < buf.length; i++) {
    buf[i] = 42;
}

const h = headers();
h.append("health", '\t\t\\"123456');
nc.publish("foo.xxx", buf, { headers: h });
// console.log(buf.length);

await nc.flush();

self.close();