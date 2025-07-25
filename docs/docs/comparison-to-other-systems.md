---
title: Comparison to Other Systems
---

Often the easiest way to understand a new language is by comparing
it to languages you already know. Here we show how policies from
several existing policy systems can be implemented with the Open
Policy Agent.

## Role-based access control (RBAC)

Role-based access control (RBAC) is pervasive today for authorization.
To use RBAC for authorization, you write down two different kinds of
information.

- Which users have which roles
- Which roles have which permissions

Once you provide RBAC with both those assignments, RBAC tells you
how to make an authorization decision. A user is authorized for
all those permissions assigned to any of the roles she is assigned to.

For example, we might have the following user/role assignments:

| User    | Role          |
| ------- | ------------- |
| `alice` | `engineering` |
| `alice` | `webdev`      |
| `bob`   | `hr`          |

And the following role/permission assignments:

| Role          | Permission | Resource      |
| ------------- | ---------- | ------------- |
| `engineering` | `read`     | `server123`   |
| `webdev`      | `write`    | `server123`   |
| `webdev`      | `read`     | `server123`   |
| `hr`          | `read`     | `database456` |

In this example, RBAC makes the following authorization decisions:

| User    | Operation | Resource      | Decision                                                 |
| ------- | --------- | ------------- | -------------------------------------------------------- |
| `alice` | `read`    | `server123`   | `allow` because `alice` is in `engineering`              |
| `alice` | `write`   | `server123`   | `allow` because `alice` is in `webdev`                   |
| `bob`   | `read`    | `database456` | `allow` because `bob` is in `hr`                         |
| `bob`   | `read`    | `server123`   | `deny` because `bob` is not in `engineering` or `webdev` |

With OPA, you can write the following snippets to implement the
example RBAC policy shown above.

```rego
package rbac.authz

# user-role assignments
user_roles := {
    "alice": ["engineering", "webdev"],
    "bob": ["hr"],
}

# role-permissions assignments
role_permissions := {
    "engineering": [{"action": "read",  "object": "server123"}],
    "webdev":      [{"action": "read",  "object": "server123"},
                    {"action": "write", "object": "server123"}],
    "hr":          [{"action": "read",  "object": "database456"}],
}

# logic that implements RBAC.
default allow := false

allow if {
    # lookup the list of roles for the user
    roles := user_roles[input.user]
    # for each role in that list
    r := roles[_]
    # lookup the permissions list for role r
    permissions := role_permissions[r]
    # for each permission
    p := permissions[_]
    # check if the permission granted to r matches the user's request
    p == {"action": input.action, "object": input.object}
}
```

<RunSnippet files="#input.json" command="data.rbac.authz" />

```json
{
  "user": "bob",
  "action": "read",
  "object": "server123"
}
```

<RunSnippet id="input.json"/>

### RBAC Separation of duty (SOD)

Separation of duty (SOD) refers to the idea that there are certain
combinations of permissions that no one should have at the same time.
For example, no one should be able to both create payments and approve payments.

In RBAC, that means there are some pairs of roles that no one should be
assigned simultaneously. For example, any user assigned both of the roles
in each pair below would violate SOD.

- create-payment and approve-payment
- create-vendor and pay-vendor

OPA's API does not yet let you enforce SOD by rejecting improper role-assignments,
but it does let you express SOD constraints and ask for all SOD violations,
as shown below. (Here we assume the statements below are added to the RBAC
statements above.)

```rego
# Pairs of roles that no user can be assigned to simultaneously
sod_roles := [
	["create-payment", "approve-payment"],
	["create-vendor", "pay-vendor"],
]

# Find all users violating SOD
sod_violation contains user if {
	some user
	# grab one role for a user
	role1 := user_roles[user][_]
	# grab another role for that same user
	role2 := user_roles[user][_]
	# check if those roles are forbidden by SOD
	sod_roles[_] == [role1, role2]
}
```

(For those familiar with SOD, this is the static version since SOD violations
happen whenever a user is assigned two conflicting roles. The dynamic version of SOD allows
a single user to be assigned two conflicting roles but requires that the same user not
utilize those roles on the same transaction, which is out of scope for this document.)

## Attribute-based access control (ABAC)

With attribute-based access control, you make policy decisions using the
attributes of the users, objects, and actions involved in the request.
It has three main components:

- Attributes for users
- Attributes for objects
- Logic dictating which attribute combinations are authorized

For example, we might know the following attributes for our users

- alice
  - joined the company 15 years ago
  - is a trader
- bob
  - joined the company 5 years ago
  - is an analyst

We would also have attributes for the objects, in this case stock ticker symbols.

- MSFT
  - is sold on NASDAQ
  - sells at $59.20 per share
- AMZN
  - is sold on NASDAQ
  - sells at $813.64 per share

An example ABAC policy in natural language might be:

- Traders may purchase NASDAQ stocks for under $2M
- Traders with 10+ years experience may purchase NASDAQ stocks for under $5M

OPA supports ABAC policies as shown below.

```rego
package abac

# User attributes
user_attributes := {
	"alice": {"tenure": 15, "title": "trader"},
	"bob": {"tenure": 5, "title": "analyst"},
}

# Stock attributes
ticker_attributes := {
	"MSFT": {"exchange": "NASDAQ", "price": 59.20},
	"AMZN": {"exchange": "NASDAQ", "price": 813.64},
}

default allow := false

# all traders may buy NASDAQ under $2M
allow if {
	# lookup the user's attributes
	user := user_attributes[input.user]
	# check that the user is a trader
	user.title == "trader"
	# check that the stock being purchased is sold on the NASDAQ
	ticker_attributes[input.ticker].exchange == "NASDAQ"
	# check that the purchase amount is under $2M
	input.amount <= 2000000
}

# traders with 10+ years experience may buy NASDAQ under $5M
allow if {
	# lookup the user's attributes
	user := user_attributes[input.user]
	# check that the user is a trader
	user.title == "trader"
	# check that the stock being purchased is sold on the NASDAQ
	ticker_attributes[input.ticker].exchange == "NASDAQ"
	# check that the user has at least 10 years of experience
	user.tenure > 10
	# check that the purchase amount is under $5M
	input.amount <= 5000000
}
```

<RunSnippet files="#input.abac.json" command="data.abac" />

```json
{
  "user": "alice",
  "ticker": "MSFT",
  "action": "buy",
  "amount": 1000000
}
```

<RunSnippet id="input.abac.json"/>

In OPA, there's nothing special about users and objects. You can attach
attributes to anything. And the attributes can themselves be structured JSON objects
and have attributes on attributes on attributes, etc. Because OPA was designed to work
with arbitrarily nested JSON data, it supports incredibly rich ABAC policies.

## Amazon Web Services IAM

Amazon Web Services (AWS) lets you create policies that can be attached to users, roles, groups,
and selected resources. You write `allow` and `deny` statements to enforce which users/roles can/can't
execute which API calls on which resources under certain conditions.
By default all API access requests are implicitly denied (i.e., not allowed). Policy statements
can explicitly allow or deny API requests. If a request is both allowed and denied, it is always denied.
Let's assume that the following [customer managed policy](https://docs.aws.amazon.com/IAM/latest/UserGuide/access_policies_managed-vs-inline.html#customer-managed-policies) is defined in AWS:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "FirstStatement",
      "Effect": "Allow",
      "Action": ["iam:ChangePassword"],
      "Resource": "*"
    },
    {
      "Sid": "SecondStatement",
      "Effect": "Allow",
      "Action": "s3:ListAllMyBuckets",
      "Resource": "*"
    },
    {
      "Sid": "ThirdStatement",
      "Effect": "Allow",
      "Action": [
        "s3:List*",
        "s3:Get*"
      ],
      "Resource": [
        "arn:aws:s3:::confidential-data",
        "arn:aws:s3:::confidential-data/*"
      ]
    }
  ]
}
```

And the above policy is attached to principal alice in AWS using
[attach-user-policy](https://docs.aws.amazon.com/cli/latest/reference/iam/attach-user-policy.html) API.
In OPA, you write each of the AWS `allow` statements as a separate statement, and you
expect the input to have `principal`, `action`, and `resource` fields.

```rego
package aws

default allow := false

# FirstStatement
allow if {
	principals_match
	input.action == "iam:ChangePassword"
}

# SecondStatement
allow if {
	principals_match
	input.action == "s3:ListAllMyBuckets"
}

# ThirdStatement
#  Use helpers to handle implicit OR in the AWS policy.
#  Below all of the 'principals_match', 'actions_match' and 'resources_match' must be true.
allow if {
	principals_match
	actions_match
	resources_match
}

# principals_match is true if input.principal matches
principals_match if {
	input.principal == "alice"
}

# actions_match is true if input.action matches one in the list
actions_match if {
	# iterate over the actions in the list
	actions := ["s3:List.*", "s3:Get.*"]
	action := actions[_]
	# check if input.action matches an action
	regex.globs_match(input.action, action)
}

# resources_match is true if input.resource matches one in the list
resources_match if {
	# iterate over the resources in the list
	resources := ["arn:aws:s3:::confidential-data", "arn:aws:s3:::confidential-data/.*"]
	resource := resources[_]
	# check if input.resource matches a resource
	regex.globs_match(input.resource, resource)
}
```

<RunSnippet files="#input.aws.json" command="data.aws" />

```json
{
  "principal": "alice",
  "action": "ec2:StartInstance",
  "resource": "arn:aws:ec2:::instance/i78999879"
}
```

<RunSnippet id="input.aws.json"/>

## XACML

eXtensible Access Control Markup Language (XACML) was designed to express security policies: allow/deny decisions using attributes of users, resources, actions, and the environment.
The following policy says that users from the organization Curtiss or Packard who are US or Great Britain nationals and who work on DetailedDesign or Simulation are permitted access to documents about NavigationSystems.

```xml
<Policy xmlns="urn:oasis:names:tc:xacml:3.0:core:schema:wd-17" PolicyId="urn:curtiss:ba:taa"
  Version="1.1"
  RuleCombiningAlgId="urn:oasis:names:tc:xacml:3.0:rule-combining-algorithm:deny-unless-permit">
  <Description>Policy for Business Authorization category TAA-1.1</Description>
  <Target />
  <Rule RuleId="Rule for NavigationSystems" Effect="Permit">
    <Target>
      <AnyOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">NavigationSystem</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:3.0:attribute-category:resource"
              AttributeId="urn:curtiss:names:tc:xacml:1.0:resource:Topics"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
      </AnyOf>
      <AnyOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Packard</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/OrganizationID"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Curtiss</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/OrganizationID"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
      </AnyOf>
      <AnyOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">GB</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/Nationality"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">US</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/Nationality"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
      </AnyOf>
      <AnyOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">DetailedDesign</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/Work-Effort"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
        <AllOf>
          <Match MatchId="urn:oasis:names:tc:xacml:1.0:function:string-equal">
            <AttributeValue DataType="http://www.w3.org/2001/XMLSchema#string">Simulation</AttributeValue>
            <AttributeDesignator
              Category="urn:oasis:names:tc:xacml:1.0:subject-category:access-subject"
              AttributeId="http://schemas.tscp.org/2012-03/claims/Work-Effort"
              DataType="http://www.w3.org/2001/XMLSchema#string" MustBePresent="true" />
          </Match>
        </AllOf>
      </AnyOf>
    </Target>
  </Rule>
</Policy>
```

The same statement is shown below in OPA. Here the inputs are assumed to be
roughly the same as for XACML: attributes of users, actions, and resources.

```rego
package xacml

# METADATA
# title: urn:curtiss:ba:taa:taa-1.1
# description: Policy for Business Authorization category TAA-1.1
default permit := false
permit if {
    # Check that resource has a "NavigationSystem" entry
    input.resource["NavigationSystem"]

    # Check that organization is one of the options
    input.user.organization in ["Packard", "Curtiss"]

    # Check that nationality is one of the options
    input.user.nationality in ["GB", "US"]

    # Check that work_effort is one of the options
    input.user.work_effort in ["DetailedDesign", "Simulation"]
}
```

<RunSnippet files="#input.xacml.json" command="data.xacml" />

```json
{
  "user": {
    "name": "alice",
    "organization": "Packard",
    "nationality": "GB",
    "work_effort": "DetailedDesign"
  },
  "resource": {
    "NavigationSystem": true
  },
  "action": {
    "name": "read"
  }
}
```

<RunSnippet id="input.xacml.json"/>
