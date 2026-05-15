# Flaky Fix Attempt (chaos test - run5)

Persona: flaky-fixer
Strategy: single-commit

Intermittent test failures tracked here:
- timeout in TestCompile under load
- race in TestAudit with concurrent writes

Proposed: increase timeout constants, serialize writes
