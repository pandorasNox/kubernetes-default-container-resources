# resource patch strategies

default compute resources
```yaml
resources:
    limits:
        memory: 256Mi
        cpu: 20m
        # ephemeral-storage: 1024Mi
    requests:
        memory: 128Mi
        cpu: 10m
        # ephemeral-storage: 1024Mi
```

## base cases

__________________________________________________

request:
```yaml
resources: {}
```

what to patch:
```yaml
resources:
    limits:
        memory: 256Mi
        cpu: 20m
    requests:
        memory: 128Mi
        cpu: 10m
```

__________________________________________________

request:
```yaml
resources:
    limits:
        memory: 256Mi
```

what to patch:
```yaml
resources:
    limits:
        cpu: 20m
```

__________________________________________________

request:
```yaml
resources:
    limits:
        cpu: 20m
```

what to patch:
```yaml
resources:
    limits:
        memory: 256Mi
```

__________________________________________________

request:
```yaml
resources:
    limits:
        memory: 256Mi
        cpu: 20m
```

what to patch:
```yaml
resources: {}
```

__________________________________________________

request:
```yaml
resources:
    requests:
        memory: 128Mi
```

what to patch:
```yaml
resources:
    limits:
        memory: resources.requests.memory
        cpu: 20m
```

__________________________________________________

request:
```yaml
resources:
    requests:
        cpu: 10m
```

what to patch:
```yaml
resources:
    limits:
        memory: 256Mi
        cpu: resources.requests.cpu
```

__________________________________________________

request:
```yaml
resources:
    requests:
        memory: 128Mi
        cpu: 1    limits:
        memory: 128Mi
        cpu: 20m
        # ephemeral-storage: 1024Mi
```

what to patch:    limits:
        memory: 128Mi
        cpu: 20m
        # ephemeral-storage: 1024Mi
```yaml
resources:
    limits:
        memory: resources.requests.memory
        cpu: resources.requests.cpu
```

__________________________________________________

rc 
rm 
lc
lm

```yaml
resources:
    limits:
        memory: 128Mi
        cpu: 20m
        # ephemeral-storage: 1024Mi
    requests:
        memory: 256Mi
        cpu: 10m
        # ephemeral-storage: 1024Mi
```

______________________________________________

/*

mem
	limit || requests
		nothing set		=> default
		both isSet		=> nil
		limit isSet		=> nil
		request isSet	=> nil
cou
	limit || requests
		nothing set		=> default
		both isSet		=> nil
		limit isSet		=> nil
		request isSet	=> nil

*/

