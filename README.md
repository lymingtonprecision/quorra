quorra
======

A utility to make dealing with CoreOS hosts running in our vSphere
cluster standardized and simple.

## Overview

`quorra` is reliant on **two** configuration files:

* `quorra.conf` the configuration of `quorra` itself
* `cloud-config.tmpl` the template to use for CoreOS `cloud-config` files

`quorra.conf` is a simple [TOML format] configuration file. We've tried
to keep the number of configuration parameters to a bare minimum. See
the [sample configuration file](./quorra.conf.sample) for details of the
configuration sections and parameters.

Some, but not all, configuration parameters can be specified as
environment variables (taking precedence over the configuration file.)
See the `ENV` Parameters section below for details.

[TOML format]: https://github.com/toml-lang/toml

`cloud-config.tmpl` should be a "standard" CoreOS `cloud-config` file
that you can _optionally_ enhance by inserting Go `text/template`
directives.

Most notably the attributes `.PrivateIP` and `.PublicIP` are available
when processing `cloud-config.tmpl` so you'll almost certainly want, at
a minimum, sections like:

    coreos:
      etcd:
        # you may need to use {{.PublicIP}} depending on your network
        # configuration
        addr: {{.PrivateIP}}:4001
        peer-addr: {{.PrivateIP}}:7001
      units:
        - name: 00-eth0.network
          runtime: true
          content: |
            [Match]
            Name=eth0

            [Network]
            Address={{.PrivateIP}}/24
            Gateway=192.168.1.254
            DNS=192.168.1.1
        - name: 00-eth1.network
          runtime: true
          content: |
            [Match]
            Name=eth1

            [Network]
            Address={{.PublcIP}}/24
            Gateway=192.168.1.254
            DNS=192.168.1.1

The full list of available attributes is:

* `.Hostname`
* `.PrivateIP`
* `.PublicIP`

## `ENV` Parameters

The following environment variables can be specified, overriding the
parameter of the same name in the configuration file:

* `QUORRA_USERNAME`
* `QUORRA_PASSWORD`
* `QUORRA_API_URL`
* `QUORRA_DATACENTER`
* `QUORRA_DATASTORE` (overrides the default datastore only)
* `QUORRA_HOST`
* `QUORRA_OVA`

All other configuration parameters must be specified via the
configuration file.
