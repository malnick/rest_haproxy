# REST HaProxy
REST HaProxy exposes the server inputs in haproxy.cfg as JSON parameters so you can query any haproxy and know what services are available, on what ports.

Defaults expose on :3000 @ /services

We use this for exposing haproxy services to query internal /info endpoints that we can use to aggregate known running versions of microservices within our infrastrcture. 

## Design
Given a haproxy configuration at path ```/etc/haproxy/haproxy.cfg```:

```ini
global
  daemon
  group  haproxy
  log  127.0.0.1 local0
  log  127.0.0.1 local1 notice
  maxconn  4096

defaults
  log global
  mode  http

frontend http-in
  bind *:80
  acl service path_beg -t /service
  acl other_service path_beg -t /other_service
  default backend my_service
  use_backend service if my_service
  use_backend other_service if other_service

backend service
  balance leastconn
  server service-01 10.0.1.10:31501 check port 32501
  server service-02 10.0.1.10:21502 check port 32502
  server service-03 10.0.2.10:31500 check port 32501
  server service-04 10.0.2.11:31500 check port 32502

backend other_service
  balance leastconn
  server service-01 10.0.5.10:31501 check 
  server service-02 10.0.5.10:21502 check
  server service-03 10.0.5.10:31503 check 
  server service-04 10.0.5.11:31500 check
```

Will result in the following JSON endpoint available at: ```localhost:3000/services```

```json
{
  "Service": {
    "other_service": [
      "10.0.5.10:31501",
      "10.0.5.10:21502",
      "10.0.5.10:31503",
      "10.0.5.11:31500"
    ],
    "service": [
      "10.0.1.10:31501",
      "10.0.1.10:21502",
      "10.0.2.10:31500",
      "10.0.2.11:31500"
    ]
  }
}
```


