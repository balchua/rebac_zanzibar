# Relationship Based Access Control (ReBAC)

This repository is a sample relationship based access control (ReBAC).  

This uses the [Authzed](https://authzed.com/) implementation of [Zanzibar](https://research.google/pubs/pub48190/) (Google's Global Authorization System).

SpiceDB is Authzed Open source implementation of Zanzibar.

The following instructions uses [Authzed playground](https://play.authzed.com/schema)

## Sample Authorization

In this particular example, the main authorization is to make sure only authorized users can do the following:

* `portfolio_managers` can `read`, `create` and `update` Portoflios for a specific industry.
* Industry `relationship_manager` can `update`, `create` and `read` portfolios belonging to its industry.
* `senior_agents` can **only** `read` and `update` portfolios, but not `create`.
* `associate_agents` can **only** `read` portfolios
* Documents belonging to a portfolio follows the same rule as permissions defined in the portfolio.


## Schema definition

```
definition user {}

definition industry {
    relation relationship_manager: user
}
/**
 * portfolio resource.
 */
definition portfolio {
    relation portfolio_manager: user
    relation senior_agent: user
    relation associate_agent: user

    relation industry: industry

    permission update = portfolio_manager + industry->relationship_manager + senior_agent
    permission create = portfolio_manager + industry->relationship_manager
    permission read = portfolio_manager + associate_agent + industry->relationship_manager
}

/**
 * document is a sub resource of portfolio
 */

definition document {
    relation parent: portfolio

    permission create = parent->create
    permission update = parent->update
    permission read = parent->update + parent->read + parent->create

}

```
## Sample relationships

### Industry relationships
* User `topdawg` is the _Oil and Gas_ Industry relationship manager
* User `madame_oracle` is the _Financial_industry_ relationship manager

Let us create the relationship

```
// industries
industry:oil_and_gas#relationship_manager@user:topdawg#...
industry:financial#relationship_manager@user:madame_oracle#...
```

Let us link the `portfolio`s to the `industry`

```
portfolio:shell#industry@industry:oil_and_gas#...
portfolio:sgx#industry@industry:financial#...
```

Retrieve all the expected permissions, by specifying these:

```
portfolio:sgx#create: 
portfolio:sgx#read:
portfolio:sgx#update:
portfolio:shell#create:
portfolio:shell#read:
portfolio:shell#update:

```

Then "Regenerate", it will result to these `Portfolio` permissions

```
portfolio:sgx#create:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#read:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#update:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:shell#create:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#read:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#update:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'

```

#### Asserting permissions

Let us check if `topdawg` have `update` permission to the Portfolio `shell`
On the tab Assertions.

```
assertTrue:
- "portfolio:shell#update@user:topdawg#..."
```

Let us check if `madame_oracle` has `create` permission to the `sgx` Portfolio

```
assertTrue:
- "portfolio:shell#update@user:topdawg#..."
- "portfolio:sgx#create@user:madame_oracle#..."
```
Both of these should yield to True.

How about not having permissions.  Let us try that by making sure that `topdawg` do not have permission on the `sgx` portfolio and vice versa for `madame_oracle`

```
assertFalse:
- "portfolio:sgx#create@user:topdawg#..."
- "portfolio:sgx#update@user:topdawg#..."
- "portfolio:sgx#read@user:topdawg#..."
- "portfolio:shell#create@user:madame_oracle#..."
- "portfolio:shell#update@user:madame_oracle#..."
- "portfolio:shell#read@user:madame_oracle#..."
```

### Portfolio documents

Let us try to check if the documents linked to the portfolios are accessible to `topdawg` and `madame_oracle`.

Create these additional relationships

```
// document relationship with portfolio
document:findoc#parent@portfolio:sgx#...
document:envdoc#parent@portfolio:shell#...
```

Go to the expected permissions and add the following then generate.

```
document:findoc#create: 
document:findoc#update: 
document:findoc#read: 
document:envdoc#create: 
document:envdoc#update: 
document:envdoc#read: 
```

You should now see `topdawg` and `madame_oracle` to have permissions to the documents, due to its transitive relationship permissions defined in the schema above.

```
document:envdoc#create:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:envdoc#read:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:envdoc#update:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:findoc#create:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
document:findoc#read:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
document:findoc#update:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#create:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#read:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#update:
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:shell#create:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#read:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#update:
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
```

## Add portfolio users

Let us add the following:
* User `agentsmith` is the `portfolio_manager` for `shell` portfolio.
* User `james` is the `portfolio_manager` for the `sgx` portfolio.

```
portfolio:shell#portfolio_manager@user:agentsmith#...
portfolio:sgx#portfolio_manager@user:james#...
```

Followed by Regenerate in the expected permissions tab.

This will be the resulting relationships.

```
document:envdoc#create:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:envdoc#read:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:envdoc#update:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
document:findoc#create:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
document:findoc#read:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
document:findoc#update:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#create:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#read:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:sgx#update:
- '[user:james] is <portfolio:sgx#portfolio_manager>'
- '[user:madame_oracle] is <industry:financial#relationship_manager>'
portfolio:shell#create:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#read:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
portfolio:shell#update:
- '[user:agentsmith] is <portfolio:shell#portfolio_manager>'
- '[user:topdawg] is <industry:oil_and_gas#relationship_manager>'
```

### Add associate agents to portfolios

Let us assign User `minime` as `associate_agent` of `sgx` portfolio and verify that he has no permission to create portfolio

In the Assertions tab add this to the `assertFalse`

```
- "portfolio:sgx#create@user:minime#..."
```

Click Run

# Run the spicedb server

Follow docs from [here](https://github.com/authzed/spicedb#installing-spicedb)

```shell
spicedb serve --grpc-preshared-key "supersecretthingy"
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