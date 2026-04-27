## 1. Remove LT/LE methods from expectMethods

- [x] 1.1 Remove `is_lt` method from `expectMethods` map (lines 628-640)
- [x] 1.2 Remove `not_lt` method from `expectMethods` map (lines 642-655)
- [x] 1.3 Remove `is_le` method from `expectMethods` map (lines 656-668)
- [x] 1.4 Remove `not_le` method from `expectMethods` map (lines 670-683)

## 2. Verify

- [x] 2.1 Run existing tests to confirm nothing breaks

## 3. Remove LT/LE tests from mod_testing_test.go

- [x] 3.1 Remove LT/LE pass tests (lines 140-162)
- [x] 3.2 Remove LT/LE fail tests (lines 214-244)
