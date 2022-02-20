# Relationship Based Access Control (ReBAC)

This repository is a sample relationship based access control (ReBAC).  

This uses the [Authzed](https://authzed.com/) implementation of [Zanzibar](https://research.google/pubs/pub48190/) (Google's Global Authorization System).

SpiceDB is Authzed Open source implementation of Zanzibar.

The following instructions uses [Authzed playground](https://play.authzed.com/schema)

## Application Context

The application is a simple CRUD based application allowing users to `CREATE`, `UPDATE`, `DELETE` and `READ` Deals.
It is focused on verifying which user is authorized to do actions on a deal.

## Authorization Rules

The following are the authorization rules:

* Each `user` is categorized into `groups`
* There are a total of 6 groups namely:
    * Front office:
        * `front_office_manager` - Users belonging to this group can `create`, `update`, `delete` and `read` the `core` section of a `deal`.
        * `front_office_member`- Users belonging to this group can only `read` the `core` section of a `deal`.
    * Middle office:
        * `middle_office_manager` - Users belonging to this group can only `create`, `update` and `read` the `supplementary` section of a `deal`.
        * `middle_office_member` - Users belonging to this group can only `read` the `core` and `supplementary` section of a `deal`.
    * Back Office:
        * `back_office_manager` - Users belonging to this group can only `create`, `read` and `update` the `servicing` section of a `deal`.  These users can also update the `servicing` section of the `deal.
        * `back_office_member` - Users belonging to this group can only `read` the `core` and `servicing` section of a `deal`
* Any users can read the `core` section of a deal
* Each `deal` should belong to a `group`.  Users belonging to different `group` should not be able to access `deal`s not in their `group`.
* Each `deal` has 3 main sections, namely:
    * `core` section
    * `supplementary` section
    * `servicing` section


## Schema definition

```
definition user {}



definition group {
    relation front_office_manager: user
    relation front_office_member: user
    relation middle_office_manager: user
    relation middle_office_member: user
    relation back_office_manager: user
    relation back_office_member: user

}

definition deal {
    relation group: group
    
    permission update_core_section =  group->front_office_manager
    permission update_supplementary_section =  group->middle_office_manager
    permission update_servicing_section = group->back_office_manager
    permission read_core_section = group->front_office_member + group->middle_office_member + group->back_office_member + update_core_section + update_supplementary_section + update_servicing_section
    permission read_supplementary_section = group->middle_office_manager + group->middle_office_member
    permission read_servicing_section = group->back_office_manager + group->back_office_member
}

```
## Sample Groups and Users

There are 2 groups:
* `france`
* `singapore`

* james is a `front_office_manager` of the group `singapore`
* mofarrell is a `middle_office_manager` of the group `singapore`
* boban is a `back_office_manager` of the group `singapore`

* loki is a `front_office_manager` of the group `france`
* magneto is a `middle_office_manager` of the group `france`
* logan is a `back_office_manager` of the group `france`

### Using `zed` command to create the group relationships

```
zed --insecure --log-level=trace relationship create group:singapore front_office_manager user:james
zed --insecure --log-level=trace relationship create group:singapore middle_office_manager user:mofarrell
zed --insecure --log-level=trace relationship create group:singapore back_office_manager user:boban

zed --insecure --log-level=trace relationship create group:france front_office_manager user:loki
zed --insecure --log-level=trace relationship create group:france middle_office_manager user:magneto
zed --insecure --log-level=trace relationship create group:france back_office_manager user:logan

```

Sample to delete relationships using `zed`

``` shell
zed --insecure --log-level=trace relationship delete group:singapore front_office_manager user:loki
zed --insecure --log-level=trace relationship delete group:singapore middle_office_manager user:magneto
zed --insecure --log-level=trace relationship delete group:singapore back_office_manager user:logan
```


# Run the spicedb server

## Install PostgresDB.

If you are running on kubernetes, the sample manifest files are located [here](manifest/postgres/postgres.yaml)

* Create a new `postgres` namespace.  

```shell
kubectl create ns postgres
```

* Apply the manifest

```shell
kubectl -n postgres apply -f postgres.yaml
```


Follow docs from [here](https://github.com/authzed/spicedb#installing-spicedb)

```shell
spicedb serve --datastore-engine postgres --datastore-conn-uri="postgres://amazinguser:perfectpassword@localhost:31747/awesomedb?sslmode=disable" --grpc-preshared-key "supersecretthingy"
```

Run spiceDB schema migration with PostgreSQL.

```shell

spicedb migrate head --datastore-engine postgres --datastore-conn-uri="postgres://amazinguser:perfectpassword@localhost:31747/awesomedb?sslmode=disable"
```

## Download zed tool

### Import the schema yamls

```
zed import --insecure --endpoint=localhost:50051 --token=supersecretthingy --relationships=false file:///home/thor/workspace/rebac-samples/files/schema.yaml
```

### Import the relationships

```
zed import --insecure --endpoint=localhost:50051 --token=supersecretthingy --schema=false file:///home/thor/workspace/rebac-samples/files/relationships.yaml
```

### Write schema

``` shell
zed --insecure --log-level=trace schema write hack/schema.zed
```

### Get relationships

#### To get all relationships related to `portfolio:sgx`

```shell
zed relationship  read --insecure --endpoint=localhost:50051 --token=supersecretthingy portfolio:sgx

portfolio:sgx associate_agent user:minime
portfolio:sgx industry industry:financial
portfolio:sgx portfolio_manager user:james
```

#### To get all relationships related to `portfolio`

```shell
zed relationship  read --insecure --endpoint=localhost:50051 --token=supersecretthingy portfolio

portfolio:sgx associate_agent user:minime
portfolio:sgx industry industry:financial
portfolio:sgx portfolio_manager user:james
portfolio:shell industry industry:oil_and_gas
portfolio:shell portfolio_manager user:agentsmith
```

### Check permission

```
zed --insecure --endpoint=localhost:50051 --token=supersecretthingy permission check "portfolio:shell" "update" "user:topdawg"
```
This returns `true`

```
zed --insecure --endpoint=localhost:50051 --token=supersecretthingy permission check "portfolio:sgx" "update" "user:topdawg"
```

This returns `false`

To see the go code client in action to verify for permissions check out the [main.go](main.go)

```shell
go run main.go

2022/01/31 13:35:52 create permission is true for user topdawg on portfolio shell
2022/01/31 13:35:52 create permission is false for user topdawg on portfolio sgx
2022/01/31 13:35:52 create permission is true for user madame_oracle on portfolio sgx
2022/01/31 13:35:52 create permission is false for user minime on portfolio sgx
2022/01/31 13:35:52 read permission is true for user minime on portfolio sgx
2022/01/31 13:35:52 read permission is false for user minime on portfolio shell
2022/01/31 13:35:52 read permission is true for user minime on document findoc
2022/01/31 13:35:52 update permission is false for user minime on document findoc
```