# Server Manager
A site for server members to start certain docker composes. Useful for running multiple gameservers on a single low end server, where not all servers can be online at once. This allows server members to click a button to start the server, setting it to live for a defined time. The server's time to live can be extended by a defined amount of time, when it has at most a certain amount of time left to live.

# Configurable Details
The directory `backend/config.yaml` allows you to define, in yaml, the properties of servers to control. This is read on runtime, so its recommended to set a volume here.

```yaml
servers:
    satisfactory:
        directory: "~/containers/satisfactory"
        initialTTL: "24h" # when started, how long before automatic shutdown.
        extendedTTL: "6h" # when extended, how much time to add
        maxTimeBeforeExtend: "3h" #the maximum time left on the server before its life can be extended (i.e. you can only extend this servers TTL when it has less than 3 hours left)
        maxExtensions: 2 # Maximum number of times the servers TTL can be extended
    minecraft:
        directory: "~/containers/minecraft"
        initialTTL: "1w" # accepts h, d, w, m for hours, days, weeks, months
        extendedTTL: "3d"
        maxTimeBeforeExtend: "1d"
        maxExtensions: -1 # -1 means infinite extensions allowed
config:
    maxServers: 1 # max number of servers that can be online at once.

defaults: # set any server property default here. The default defaults are as follows:
    initialTTL: "24h"
    extendedTTL: "6h"
    maxTimeBeforeExtend: "6h"
    maxExtensions: 4
```
