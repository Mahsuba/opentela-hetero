# Test data (Day 6)

Go skips any folder named testdata when it builds the program.
Only the tests read these files.

## kaggle-t4x2.csv
- Exact copy of week5/day05/fixtures/kaggle-t4x2.csv.
- Real nvidia-smi output from a Kaggle notebook ("GPU T4 x2"), captured 2026-11-20T20:37:34Z.
- uuids were already cut to GPU- + 8 hex before saving. No full uuid here.
- Used by TestKaggleFixture in parse_test.go.

All other test inputs are written inside the _test.go files and are made-up
test text, not real hardware readings.
