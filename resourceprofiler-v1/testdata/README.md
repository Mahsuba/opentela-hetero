# Test data

Go skips any folder named testdata when it builds the program.
Only the tests read these files.

## kaggle-t4x2.csv
- Real nvidia-smi output from a Kaggle notebook ("GPU T4 x2"),
  captured 2026-11-20T20:37:34Z (same lines as week5/day05/fixtures/).
- Day 7 added the "#" header lines (capture time, uname -m = x86_64)
  so the profile reports when and where it was captured.
- uuids were cut to GPU- + 8 hex before saving. No full uuid here.
- Used by TestKaggleFixture in fixture_test.go.

All other test inputs are written inside the _test.go files and are
made-up test text, not real hardware readings.
