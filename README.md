## eb-worker-scheduler

Provides an open-source implementation to [AWS in-built sqsd daemon](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/using-features-managing-env-tiers.html#worker-daemon) in its EB worker platform.

### Use Cases

1. In Kubernetes: Can be used as a sidecar container to migrate existing beanstalk workers to K8S env.

### How to use?

As a go program:

```bash
go run main.go -q ${QUEUE_URL} --host http://localhost:80
```

As a docker container:

```bash
docker run sqsd -q ${QUEUE_URL} --host http://localhost:80
```

### Reference

```bash
Start the SQS worker, polling and forwarding the messages via HTTP requests

Usage:
  sqsd  [flags]

Flags:
  -h, --help              help for sqsd
  -a, --host string       Provide the Host on which API is listening (default "http://localhost:80")
  -p, --httpPath string   Provide the HTTP Path of the API to hit with POST request of the job (default "/")
  -m, --maxJobs int       Provide the limit of messages to receive (max. 10) (default 10)
  -q, --queueUrl string   Provide the queue URL (required)
  -v, --verbose           verbose logging
```
