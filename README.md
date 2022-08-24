## [eb-worker-scheduler](https://github.com/shivanshs9/eb-worker-scheduler)

Provides an open-source implementation to [AWS in-built crond daemon](https://docs.aws.amazon.com/elasticbeanstalk/latest/dg/using-features-managing-env-tiers.html#worker-periodictasks) in its EB worker platform.
Allows executing periodic tasks via SQS and crontab format.

### Use Cases

1. In Kubernetes: Can run as a independent pod (from deployment/statefulset) mounted with the `cron.yaml` file, detailing all the periodic tasks in the required format.

### How to use?

As a go program:

```bash
go run main.go -q ${QUEUE_URL}
```

As a docker container:

```bash
docker run eb-crond -q ${QUEUE_URL}
```

### Sample cron.yaml file

```yaml
version: 1
cron:
  # Note that all the `schedule` specified here are specified in Server Time
  - name: "Job name"
    url: "/path/to/POST/API"
    schedule: "*/1 * * * *" # every minute
```

### Reference

```bash
Start the SQS scheduler, that triggers jobs via SQS based on given crontab

Usage:
  crond [flags]

Flags:
  -h, --help              help for crond
  -p, --path string       Provide the path to cron.yaml file. (default "cron.yaml")
  -q, --queueUrl string   Provide the queue URL (required)
  -r, --retry int         Number of re-attempts for every failed message push (default 3)
  -v, --verbose           verbose logging
```
