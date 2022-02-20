# Installing spicedb server

## Install PostgresDB with Kubernetes

**The steps below are purely for development purposes only**

Enable `hostpath-storage` addon in MicroK8s

If you are running on kubernetes, the sample manifest files are located [here](manifest/postgres/postgres.yaml)

* Create a new `postgres` namespace.  

```shell
kubectl create ns postgres
```

* Apply the manifest

```shell
kubectl -n postgres apply -f postgres.yaml
```

## Download spiceDB

Follow docs from [here](https://github.com/authzed/spicedb#installing-spicedb)

## Start spiceDB

### Run spiceDB schema migration

**Do not forget to change the port of your PostgreSQL, in this example the port used is `31747`**

Schema migration with PostgreSQL.

```shell

spicedb migrate head --datastore-engine postgres --datastore-conn-uri="postgres://amazinguser:perfectpassword@localhost:31747/awesomedb?sslmode=disable"
```

### Start spiceDB

```shell
spicedb serve --datastore-engine postgres --datastore-conn-uri="postgres://amazinguser:perfectpassword@localhost:31747/awesomedb?sslmode=disable" --grpc-preshared-key "supersecretthingy"
```

## Download zed tool

Go to [zed github](https://github.com/authzed/zed) and download `spiceDB` binary from the [releases page](https://github.com/authzed/zed/releases)

## Write authorization schema into spiceDB

``` shell
zed --insecure --log-level=trace schema write hack/schema.zed
```
