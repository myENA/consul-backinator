job "consul-backinator" {
	datacenters = [
		"dc1"
	]
	type = "service"

	update {
		stagger = "10s"
		max_parallel = 1
	}

	group "consul-backinator" {
		count = 1
		restart {
			# The number of attempts to run the job within the specified interval.
			# Runs the job every hour, adjust depending on your needs.
			attempts = 1
			interval = "1h"

			# The "delay" parameter specifies the duration to wait before restarting
			# a task after it has failed.
			delay = "1h"
			mode = "delay"
		}

		ephemeral_disk {
			size = 300
		}

		task "consul-backinator" {
			driver = "docker"

			# The "config" stanza specifies the driver configuration, which is passed
			# directly to the driver to start the task. The details of configurations
			# are specific to each driver, so please see specific driver
			# documentation for more information.
			config {
				image = "myena/consul-backinator"
				args = [
					"backup", "-file", "s3://name-of-your-bucket/backups"
				]
			}

			env {
				"AWS_ACCESS_KEY_ID"     = "AWS ACCESS KEY GOES HERE"
				"AWS_SECRET_ACCESS_KEY" = "AWS SECRET KEY GOES HERE"
				"AWS_REGION"            = "us-east-1"
				"CONSUL_HTTP_ADDR"      = "${NOMAD_IP_consulbackinator}:8500"
			}

			resources {
				cpu = 500 # 500 MHz
				memory = 256 # 256MB
				network {
					mbits = 10
					port "consulbackinator" {}
				}
			}
		}
	}
}
