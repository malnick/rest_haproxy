# REST HaProxy
This golang project serves ```server``` lines in haproxy.cfg as JSON REST endpoints on :3000.

I use it to expose running services which Puppet statically provisions. In the future I won't need this interface because we'll use Marathon to query dynamically provisioned services. But for now it's low hanging fruit to get IP and port assignements via this method, on the fly.
