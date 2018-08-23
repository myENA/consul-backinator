job "consul-backinator" {
  datacenters = [
    "dc1"
  ]
  type = "batch"
  periodic {
    cron             = "*/15 * * * * *"
    prohibit_overlap = true
    time_zone = "America/Chicago"
  }
  group "consul-backinator" {
    task "consul-backinator" {
      driver = "docker"
      config {
        image = "myena/consul-backinator"
        entrypoint = [ "/bin/sh", "-c" ]
        args = [
          "consul-backinator backup -file s3://name-of-your-bucket/backup-$(date +%m%d%Y.%s).bak"
        ]
      }
      env {
        # The "AWS_ACCESS_KEY_ID" and "AWS_SECRET_ACCESS_KEY" environment
        # settings are not needed if you are running this inside an ec2
        # instance with have an associated IAM profile.
        "AWS_ACCESS_KEY_ID"     = "AWS ACCESS KEY GOES HERE"
        "AWS_SECRET_ACCESS_KEY" = "AWS SECRET KEY GOES HERE"
        # The "AWS_REGION" and "CONSUL_HTTP_ADDR" environment settings are
        # always required unless specified in the command section above.
        "AWS_REGION"            = "us-east-1"
        "CONSUL_HTTP_ADDR"      = "${attr.driver.docker.bridge_ip}:8500"
      }
    }
  }
}
