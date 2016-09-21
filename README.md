# Ensure

Ensure Docker Service run what you expect it to run.
Forget about remembering what to add or what to remove to a service, just send
all the stuff you expect to be in that service and `ensure` will do the
heavylifting for you

If you want to create a `test` service:

```bash
ensure --name test --image registry/test --replicas 30 --publish 80:9090
```

But if you want to update everything:

```bash
ensure --name test --image registry/test:v2 --replicas 20 --publish 80:8080
--env DAY=today
```

Everything you send to `ensure` is what the service is going to be. No more
`--env-add`, `--env-rm`, `--publish-add`, `--publish-rm`, etc.
