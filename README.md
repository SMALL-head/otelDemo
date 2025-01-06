a simple demo for otel distributed tracing  
using grafana tempo to display  

using docker command to start up container, and then start up two service
```bash
docker compose up
```

```bash
# run in terminal 1
go run ./svc1/...
```

```bash
# run in another terminal
go run ./svc2/... 
```