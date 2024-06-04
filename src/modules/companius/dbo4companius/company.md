# Company model

A company represents a group of 1 or more people.

All modifications of company records should go through
[companies facade](../facade2companies/README.md).

There are 2 kinds of companies:

- private
- public

## Private companies

There are 2 types of private companies:

- personal
- family

They do not have a title.

### Personal company

A personal company always belongs only to a single user.

Any sub-records are accessible just to the owning user unless explicitly shared.

### Family company
A family company can belong to a few users - a family.

Any sub-records are accessible just to all family members unless explicitly shared or explicitly protected by access policy.
