---
# Noop Reminder - Prompt skill reminding agents to always signal completion via safe-outputs.
# Without this reminder, missing safe-output calls are the github/gh-aw#1 cause of workflow failures.
---

**Important**: If no action is needed after completing your analysis, you **MUST** call the `noop` safe-output tool with a brief explanation. Failing to call any safe-output tool is the most common cause of safe-output workflow failures.

```json
{"noop": {"message": "No action needed: [brief explanation of what was analyzed and why]"}}
```
