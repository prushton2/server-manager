# Server Manager
A site for server members to start certain docker composes. Useful for running multiple gameservers on a single low end server, where not all servers can be online at once. This allows server members to click a button to start the server, setting it to live for a defined time. The server's time to live can be extended by a defined amount of time, when it has at most a certain amount of time left to live.

# Configurable Details
The file `backend/config.yaml` allows you to define, in yaml, the properties of servers to control. This is read on runtime, so its recommended to set a volume here.

```yaml
servers:
    satisfactory:
        directory: "/home/user/containers/satisfactory"
        initialTTL: "24h" # when started, how long before automatic shutdown.
        extendedTTL: "6h" # when extended, how much time to add
        maxTimeBeforeExtend: "3h" #the maximum time left on the server before its life can be extended (i.e. you can only extend this servers TTL when it has less than 3 hours left)
        maxExtensions: 2 # Maximum number of times the servers TTL can be extended
    minecraft:
        directory: "/home/user/containers/minecraft"
        initialTTL: "1w" # accepts h, d, w, m for hours, days, weeks, months
        extendedTTL: "3d"
        maxTimeBeforeExtend: "1d"
        maxExtensions: -1 # -1 means infinite extensions allowed
config:
    maxServers: 1 # max number of servers that can be online at once.

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
