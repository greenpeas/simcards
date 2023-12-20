# SIM-card importer

### Host
`192.168.158.81`

### Config
Path `/configs/simcard-importer/app.yml`
app.yml example:
```
enableService: true
# Database credentials
database:
  psql:
    url: "postgres://wwwspo:_my_password_@192.168.158.68:5432/wwwspo?application_name=SIM-card importer (GoLang)"
  redis:
    addr: "spo_redis_1:6379"
    password: ""
    db:       0
```
### Logging
Logs path `/var/lib/docker/volumes/go-sim-cards-import_logs/_data`

### Control
Portainer `http://192.168.158.81:9005/#!/2/docker/stacks/go-sim-cards-import` [click here](http://192.168.158.81:9005/#!/2/docker/stacks/go-sim-cards-import)


### Redis

```redis-cli```

```
HGETALL ICCIDs:list
HGETALL ICCIDs:lost
FLUSHDB
KEYS *

HGET ICCIDs:list 89701015388560972019

HSET ICCIDs:lost 89701015388560972019 '{"status":"blocked","imei":"867459044597532"}'
```
