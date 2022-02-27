# Relationship Based Access Control (ReBAC)

This repository is a sample relationship based access control (ReBAC).  

This uses the [Authzed](https://authzed.com/) implementation of [Zanzibar](https://research.google/pubs/pub48190/) (Google's Global Authorization System).

SpiceDB is Authzed Open source implementation of Zanzibar.

The following instructions uses [Authzed playground](https://play.authzed.com/schema)

## Pre-requisites

* MicroK8s
* PostgreSQL
* SpiceDB
* zed CLI

Follow the [getting started section](docs/getting_started.md#getting_started.md)

## Application Context

The application is a simple CRUD based application allowing users to `CREATE`, `UPDATE`, `DELETE` and `READ` Deals.
It is focused on verifying which user is authorized to do actions on a deal.

## Authorization Rules

Basic relationships:

* There is a top level group called `organization`
* Users must belong to an `organization`.
* Users must have a `thirdparty_role`
* Each `thirdparty_role` must have an associated `operational_role`.

Below are the protected resources:
* `deal`

Organizations are:
* `singapore`
* `france`

Thirdparty roles:
* `agent`
* `auditor`
* `loan_officer`

Below are the operational roles:
* `agents`
    * `front_office_manager`
    * `front_office_originator`
* `auditor`
    * `middle_office_manager`
    * `middle_office_validator`
* `loan_officer`
    * `back_office_manager`
    * `back_office_validator`

The following are the authorization rules:

### On deal creation:
* Only users with `agent` with `operational_role`=`front_office_originator` for a particlar `organization` can `create_deal`

Example:
User: james belonging to a organization `singapore` is an `agent` having operational role `front_office_originator` 

### Read a deal

* User who created the deal can read the deal.
* User with `front_office_manager` role can read the deal.
* Deal must be in the same organization as the subject

### Deal State Transitions

Initial state `created`
From `created` to `reviewed`
From `reviewed` to `validated`
From `validated` to `processed`


Under the state `created`, `agents` with operational role `front_office_manager` can set the state to `reviewed`
When the deal is in `reviewed` state, only `auditors` with operational_role `middle_office_manager` can set the state to `reviewed`
When the deal is in `reviewed` state, any `auditors` can read the deal

When the deal is in `validated` state, only `loan_officer` with operational_role `back_office_manager` and `back_office_validator` can view
When the deal is in `processed` state, all users can view.

### Deal fields access rights

`field1` - writable
* When the deal is in `created` state and in `agent` thirdparty role with operational role `front_office_manager` or `front_office_originator`

`field1` - readable
* When the deal is in `created` state and in `agent` thirdparty role with operational role `front_office_manager` or `front_office_originator`
* When the deal is in `reviewed` state and in `auditor` thirdparty role with operational role `middle_office_manager` or `middle_office_validator`
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

`field2` is writable
* When the deal is in `created` state and in `agent` thirdparty role with operational role `front_office_manager` or `front_office_originator`

`field2` - readable
* When the deal is in `created` state and in `agent` thirdparty role with operational role `front_office_manager` or `front_office_originator`
* When the deal is in `reviewed` state and in `auditor` thirdparty role with operational role `middle_office_manager` or `middle_office_validator`
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

`field3` - writable
* When the deal is in `reviewed` state and in `auditor` thirdparty role with operational role `middle_office_manager` 

`field3` - readable
* When the deal is in `reviewed` state and in `auditor` thirdparty role with operational role `middle_office_manager` or `middle_office_validator`
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

`field4` - writable
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator` 

`field4` - readable
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

`field5` - writable
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

`field5` - readable
* When the deal is in `validated` state and in `loan_officer` thirdparty role with operational role `back_office_manager` or `back_office_validator`

## Schema definition

```


```

Loading the schema using the `zed` tool.

``` shell
zed --insecure --log-level=trace schema write hack/schema.zed
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

## Building the application

``` shell
go build .
```

## Running the application

``` shell
./demo-spicedb serve
```

The application will use the port `8181`.

Supported endpoints:

| METHOD | URI | COMMENTS |
|--------|-----|----------|
| `GET` | `/api/v1/deal/:id` | Gets a deal |
| `POST` | `/api/v1/deal` | Create a new deal |
