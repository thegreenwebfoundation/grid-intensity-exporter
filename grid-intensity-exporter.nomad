
job "grid-intensity-exporter" {

  # The "datacenters" parameter specifies the list of datacenters which should
  # be considered when placing this task. This must be provided.
  datacenters = ["dc1"]

  # system jobs run on every client node in a cluster, which is what we
  # want, but this bug here means that running this job as a system on
  # means that network ports aren't allocated properly
  # https://github.com/hashicorp/nomad/issues/8934
  # until it's resolved, we need to run this as a service instead
  type = "service"

  group "grid-intensity-exporter" {

    network {
      # for testing, use host network so integration test can connect
      mode "host"

      # for testing, we can get away with having a fixed port
      # but in production we'd let nomad allocated a port instead
      port "exporter" {
        static = 8000
        to = 8000
      }
    }

    task "grid-intensity-exporter" {

      driver = "docker"
      
      config {
        image = "grid-intensity-exporter:integration-test"
        ports = ["exporter"]
      }
    }
  }
}
