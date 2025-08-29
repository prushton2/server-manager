# Server Manager
A site for server members to start certain docker composes. Useful for running multiple gameservers on a single low end server, where not all servers can be online at once. This allows server members to click a button to start the server, setting it to live for a defined time. The server's time to live can be extended by a defined amount of time, when it has at most a certain amount of time left to live.

# Configurable Details
The file `backend/config.yaml` allows you to define, in yaml, the properties of servers to control. This is read on runtime, so its recommended to set a volume here.

```yaml
servers:
  satisfactory:
      # Directory of the container (this is the directory mounted in the server-manager docker container)
      directory: "/home/user/containers/satisfactory"
      
      # When started, how long before automatic shutdown. 
      #  All time strings accept h, d, w, m for hours, days, weeks, months
      initialTTL: "24h"
      
      # When extended, how much time to add
      extendedTTL: "6h"
      
      # The maximum time left on the server before its life can be extended (i.e. you can only extend this servers TTL when it has less than 3 hours left)
      maxTimeBeforeExtend: "3h" 
      
      # Maximum number of times the servers TTL can be extended
      #  -1 allows for infinite extensions
      maxExtensions: 2 
      
      # Status of the server
      #   enabled: working properly
      #   maintenance: backend will not try to start/stop server, and a message will be shown to users
      #   disabled: same as maintenance, but will say "server is disabled" instead of "server is temporarily down for maintenance"
      #   hidden: backend will not try to start/stop server, and server will be hidden from users
      status: "enabled"

  minecraft:
      directory: "/home/user/containers/minecraft"
      initialTTL: "1w"
      extendedTTL: "3d"
      maxTimeBeforeExtend: "1d"
      maxExtensions: -1
      status: "maintenance"

config:
  # max number of servers that can be online at once.
  maxServers: 1 

users:
  root: # Username for logging purposes 
    password: "testPassword"
    canStart: true # interactions they can perform on the servers they can see
    canExtend: true
    canStop: true
    allowedServers: # servers they can interact with
      - "astroneer"
      - "satisfactory"
  user:
    password: "supersecretpassword" # The user's unique password. You generate these, and they are stored without a hash :(
    canStart: true
    canExtend: true
    canStop: false  
    allowedServers:
      - "satisfactory"
```

# Docker Compose

```yaml
services:
  backend:
    build: ./backend
    ports:
      - 3000:3000
    volumes:
      - ./config.yaml:/app/config.yaml # Config file
      - ./state.json:/app/state.json   # State file to store server state between restarts
      - ./.env:/.env   # State file to store server state between restarts
      - /var/run/docker.sock:/var/run/docker.sock # Mount the docker socket
      - /home/server/containers/satisfactory:/home/user/containers/satisfactory # Mount any containers you want to be controlled to the directory specified in the compose
      - /home/server/containers/minecraft:/home/user/containers/minecraft
  frontend:
    build:
      context: frontend/
      args:
        - VITE_BACKEND_URL=https://servermanager.example.com/backend
    ports:
      - 80:80

```

# .env

This site supports push notifications to let players know when the time on the server is low. You can  acquire your own keys [here](). If you dont want to use this, simply put null in the vapid keys. 

```env
PUBLIC_VAPID_KEY=null
PRIVATE_VAPID_KEY=null
```
