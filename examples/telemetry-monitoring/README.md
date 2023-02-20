StoneWork: Telemetry & Monitoring
==========================

This example is a modification of the [StoneWork as a Cross-Connect example](https://github.com/PANTHEONtech/StoneWork/tree/main/examples/testing/010-xconnect). It additionally enables the telemetry plugin and adds Prometheus and Grafana containers to the deployment.

Run with:
```shell
docker compose up -d
```

Grafana is accessible at http://localhost:3000.

The default login credentials are:  
```
Username: admin  
Password: admin
```

Optionally, you can run the `test-stonework.sh` script to populate the dashboard with metrics.