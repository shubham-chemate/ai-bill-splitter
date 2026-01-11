## Problem Statement

Given a dinner/shopping bill among friends and a set of rules, split the bill among the friends. Calculate each friend's share for each item and their total share.

### Input

The bill is expected in the following format:

```
| Item | Price | Quantity | Tax  | Total |
|------|-------|----------|------|-------|
| it1  | 100   | 2        | 0    | 200   |
| it2  | 50    | 2        | 10   | 110   |
| it3  | 500   | 5        | 0    | 2500  |
| it4  | 10    | 3        | 10   | 40    |
```

you will also be given rules  

- There are 5 friends: `fr1`, `fr2`, `fr3`, `fr4`, `fr5`.
- An item is shared by everyone if it is not specifically mentioned in the rules.
- Item `it1` is shared by `fr1` and `fr2`.
- Item `it2` is shared by `fr2` and `fr3`.

### Output

The output format is as follows:
```
| Person | it1    | it2    | it3    | it4   | Total  |
|--------|--------|--------|--------|-------|--------|
| fr1    | 100.00 | -      | 500.00 | 8.00  | 608.00 |
| fr2    | 100.00 | 105.00 | 500.00 | 8.00  | 713.00 |
| fr3    | -      | 105.00 | 500.00 | 8.00  | 613.00 |
| fr4    | -      | -      | 500.00 | 8.00  | 508.00 |
| fr5    | -      | -      | 500.00 | 8.00  | 508.00 |
```