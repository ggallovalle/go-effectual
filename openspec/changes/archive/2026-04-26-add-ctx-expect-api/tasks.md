## 1. Expect object infrastructure

- [x] 1.1 Add expect object struct type and shared methods map in mod_testing.go
- [x] 1.2 Implement expect method on testCtx that creates and returns expect table

## 2. Assertion methods

- [x] 2.1 Implement is_nil assertion method
- [x] 2.2 Implement not_nil assertion method
- [x] 2.3 Implement is_true assertion method
- [x] 2.4 Implement is_false assertion method
- [x] 2.5 Implement equals assertion method
- [x] 2.6 Implement not_equals assertion method
- [x] 2.7 Implement is_lt assertion method
- [x] 2.8 Implement not_lt assertion method
- [x] 2.9 Implement is_le assertion method
- [x] 2.10 Implement not_le assertion method

## 3. Failure message helper

- [x] 3.1 Implement shared failure message builder that formats source location + msg + expected vs actual

## 4. Migrate semver_test.lua

- [x] 4.1 Replace nil checks (assert(v ~= nil)) with expect(v):not_nil()
- [x] 4.2 Replace equality checks (assert(v.major == 1)) with expect(v.major):equals(1)
- [x] 4.3 Replace boolean checks (assert(not ok)) with expect(ok):is_false()
- [x] 4.4 Replace comparison checks (assert(v1 < v2)) with expect(v1 < v2):is_true()
- [x] 4.5 Replace negated comparisons (assert(not (v2 < v1))) with expect(v2 < v1):is_false()
- [x] 4.6 Replace tostring equality with expect(tostring(v)):equals("2.3.4")

## 5. Verify

- [x] 5.1 Run semver tests and confirm all pass with ctx.expect assertions
