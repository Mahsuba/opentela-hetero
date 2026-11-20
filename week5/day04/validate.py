# validate.py - checks JSON files against schema.json
# Usage: python3 validate.py schema.json file1.json [file2.json ...]
import json
import sys

from jsonschema import Draft202012Validator

schema = json.load(open(sys.argv[1]))
Draft202012Validator.check_schema(schema)  # is the schema itself valid?
validator = Draft202012Validator(schema)

all_ok = True
for path in sys.argv[2:]:
    errors = list(validator.iter_errors(json.load(open(path))))
    if not errors:
        print("PASS", path)
    for e in errors:
        all_ok = False
        where = "/".join(str(p) for p in e.absolute_path) or "(top)"
        print("FAIL", path, "->", where + ":", e.message)

sys.exit(0 if all_ok else 1)
