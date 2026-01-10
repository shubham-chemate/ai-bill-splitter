## Problem Statement

Given a dinner/shopping bill among friends and a set of rules, split the bill among the friends. Calculate each friend's share for each item and their total share.

### Input

The bill is expected in the following format:

```csv
item, price, quantity, tax, total
it1, 100, 2, 0, 200
it2, 50, 2, 10, 110
it3, 500, 5, 0, 2500
it4, 10, 3, 10, 40
```

### Rules

- There are 5 friends: `fr1`, `fr2`, `fr3`, `fr4`, `fr5`.
- An item is shared by everyone if it is not specifically mentioned in the rules.
- Item `it1` is shared by `fr1` and `fr2`.
- Item `it2` is shared by `fr2` and `fr3`.

### Output

The output format is as follows:

**fr1**
```text
it1 100
it3 500
it4 8
total 608
```

**fr2**
```text
it1 100
it2 105
it3 500
it4 8
total 713
```

**fr3**
```text
it2 105
it3 500
it4 8
total 613
```

**fr4**
```text
it3 500
it4 8
total 508
```

**fr5**
```text
it3 500
it4 8
total 508
```