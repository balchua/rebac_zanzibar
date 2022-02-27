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

# Application Context

This is a show case of how spiceDB can be used as an authorization engine.

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

## Creating a Deal:
* Only users with `agent` with `operational_role`=`front_office_originator` for a particlar `organization` can `create_deal`

Example:
User: james belonging to a organization `singapore` is an `agent` having operational role `front_office_originator` 

## Viewing a deal

* User who created the deal can read the deal.
* User with `front_office_manager` role can read the deal.
* Deal must be in the same organization as the subject
* When the deal is in `processed` state, any `agent`, `auditor` or `loan_officer` can view it, as long as the `deal` is part of the organization `singapore`

## Deal State Transitions

Initial state `created`

From `created` --> `reviewed`

From `reviewed` --> `validated`

From `validated` --> `processed`


Under the state `created`, `agents` with operational role `front_office_manager` can set the state to `reviewed`

When the deal is in `reviewed` state, only `auditors` with operational_role `middle_office_manager` can set the state to `reviewed`.

When the deal is in `reviewed` state, any `auditors` can read the deal

When the deal is in `validated` state, only `loan_officer` with operational_role `back_office_manager` and `back_office_validator` can view.

When the deal is in `processed` state, all users can view.

## Deal fields access rights

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
  definition user {}


  /* Every thirdparty role belongs to an organization */

  definition organization {
  	relation member: user
  }


  definition thirdparty_role {
  	relation front_office_manager: user
  	relation front_office_originator: user
  	relation middle_office_manager: user
  	relation middle_office_validator: user
  	relation back_office_manager: user
  	relation back_office_validator: user
  	relation org: organization

  	// check if the subject is a member of the org and a front_office_originator
  	permission create_deal = org->member & front_office_originator
  }


  // resource id format: <deal_id>_<state>

  definition deal {
  	relation thirdparty: thirdparty_role
  	relation org: organization
  	permission can_role_review = thirdparty->front_office_manager
  	permission can_role_validate = thirdparty->middle_office_manager
  	permission can_role_process = thirdparty->back_office_manager
  	permission can_role_view = org->member & thirdparty->back_office_manager + thirdparty->back_office_validator + thirdparty->middle_office_manager + thirdparty->middle_office_validator + thirdparty->front_office_manager + thirdparty->front_office_originator
  }


  // resource id format: <deal_id>_<state>_<field_name>

  definition deal_field {
  	relation reader: thirdparty_role#front_office_manager | thirdparty_role#front_office_originator | thirdparty_role#middle_office_manager | thirdparty_role#middle_office_validator | thirdparty_role#back_office_manager | thirdparty_role#back_office_validator
  	relation writer: thirdparty_role#front_office_manager | thirdparty_role#front_office_originator | thirdparty_role#middle_office_manager | thirdparty_role#middle_office_validator | thirdparty_role#back_office_manager | thirdparty_role#back_office_validator
  	relation org: organization
  	permission write = org->member & writer 
  	permission read = org->member & reader + write
  }

```

Loading the schema using the `zed` tool.

``` shell
zed --insecure --log-level=trace schema write hack/schema.zed
```

## Setting up the roles and organization relationships

``` shell
./setup.sh singapore
```

This will create the relationships between `thirdparty_role` and `organization`, belonging to `singapore`

## Sample relationships between organizations and Users

There are 2 groups:
* `france`
* `singapore`

`singapore` :

* `agent`
  * `james` is a `front_office_manager` for the organization `singapore`
  * `john` is a `front_office_validator` for the organization `singapore`
* `auditor`
  * `mofarrell` is a `middle_office_manager` for the group `singapore`
  * `luke` is a `middle_office_validator` for the groupd `singapore`
* `loan_officer`
  * `boban` is a `back_office_manager` for the group `singapore`
  * `topdawg` is a `back_office_validator` for the group `singapore`


Commands:

``` shell
./add-user.sh james singapore agent front_office_manager
./add-user.sh john singapore agent front_office_originator
./add-user.sh mofarrell singapore auditor middle_office_manager
./add-user.sh luke singapore auditor middle_office_validator
./add-user.sh boban singapore loan_officer back_office_manager
./add-user.sh topdawg singapore loan_officer back_office_validator

```

## Check if user can create deal

Let us check if `james` can create a deal

``` shell
./check_can_create_deal.sh james agent
false
```
The output is `false`, although `james` belongs to the organization `singapore`, but `james` do not have the `front_office_originator` operational role.

Lets us try to see if `james` can create a deal

``` shell
./check_can_create_deal.sh john agent
true
```

As expected, since `john` has the `front_office_originator` operational role, he can now be given the `create_deal` permission.

Given this permission, your application can now allow the user `john` to create a deal.

## Create a first deal

Let us create the first `deal` for the organization `singapore`

deal id: `1`
for organization: `singapore`
for the thirdparty_role: `agent`


``` shell
./create_deal.sh 1 singapore agent
```

Let's check how the deal relationship was created.

``` shell
zed --insecure relationship read deal
deal:1_created org organization:singapore
deal:1_created thirdparty thirdparty_role:agent
```

As you can see we can use the deal resource id to be `<dealId>_<dealState>`

## Check if a user can review a deal

Lets check if `john` can `review` the deal id 1, on status `created`

``` shell
./check_permission.sh john 1 created can_role_review
false
```
Why `false`, because `john` has the operational role `front_office_originator`

Lets check if `james` can `review` the deal id 1, status `created`

``` shell
./check_permission.sh james 1 created can_role_review
true
```

Why `true` because `james` has the operational role `front_office_manager`

## Transitioning a deal state

### Created to Reviewed

From `created` state to `reviewed` state

Lets pretend that `james` has `reviewed` the deal.

``` shell
./deal_transition.sh 1 reviewed singapore
```

Let us check if `james` and `john` can still view the deal id `1` status `reviewed`

``` shell
/check_permission.sh james 1 reviewed can_role_review
false

./check_permission.sh john 1 reviewed can_role_review
false
```

Lets check if the auditors can `view` the deal id `1` status `reviewed`

``` shell

./check_permission.sh luke 1 reviewed can_role_view
true

./check_permission.sh mofarrell 1 reviewed can_role_view
```

### Reviewed to Validated

Lets check if the auditors can `validate` the deal id `1` status `reviewed`

``` shell
./check_permission.sh mofarrell 1 reviewed can_role_validate
true

./check_permission.sh luke 1 reviewed can_role_validate
false
```

Lets pretend that user `mofarrell` has `validated` the deal id `1`

``` shell
./deal_transition.sh 1 validated singapore
```

### Validated to Processed

Lets check if `mofarrell` and `luke` can still view the deal id `1`

``` shell
./check_permission.sh mofarrell 1 validated can_role_view
false

./check_permission.sh luke 1 validated can_role_view
false
```

Now lets check if the loan_officers can view the `validated` deal `1`

``` shell
./check_permission.sh boban 1 validated can_role_view
true
./check_permission.sh topdawg 1 validated can_role_view
true
```

Let topdawg validate the deal.

``` shell
./deal_transition.sh 1 processed singapore
```

At this point the rule states that anyone can view the deal.

``` shell
./check_permission.sh boban 1 processed can_role_view
true

./check_permission.sh mofarrell 1 processed can_role_view
true

./check_permission.sh james 1 processed can_role_view
true

./check_permission.sh john 1 processed can_role_view
true

./check_permission.sh topdawg 1 processed can_role_view
true

./check_permission.sh luke 1 processed can_role_view
true

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
