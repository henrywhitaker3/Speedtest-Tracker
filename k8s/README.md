# Using Kubernetes

## Considerations

Using persistent storage is highly recommended to preserve database information in case of pod failure.

## Installation Methods

### kubectl

Raw spec files for manual alteration and deployment have been supplied if that is your preferred deployment method.

```bash
kubectl apply -f specs/
```

These spec files assume you have the [nfs-provisioner](https://github.com/kubernetes-incubator/external-storage/tree/master/nfs) and `type: LoadBalancer` available on your Kubernetes implementation.

### helm

A helm chart has been included to assist with deployment

```bash
helm install speedtesttracker ./Speedtest-Tracker \
    --values Speedtest-Tracker/values.yaml \
    --set eula.OOKLA_EULA_GDPR=true
```

This will stand up a minimal install for testing and validation on you system. It is advised that you use persistent and not ephemeral storage but no assumptions have been made about your setup.

## Known Issues

- Liveness probes are disabled currently, will look into a working configuration in the future.
- On some Kubernetes installations ndots issues can cause DNS to not resolve properly, you may have to manually configure [ndots](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/) to meet your use case.
- Hugepages are required for FPM Pools to work properly, the needed configuration is included.
