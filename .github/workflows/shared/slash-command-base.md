---
# Slash Command Base — Standard safe-output and UX config for /slash command workflows.
# Provides: add-comment (max configurable), noop.
#
# Usage:
#   imports:
#     - uses: shared/slash-command-base.md
#       with:
#         reaction: eyes          # optional, default: eyes
#         status-comment: true    # optional, default: true
#         add-comment-max: 1      # optional, default: 1

import-schema:
  reaction:
    type: string
    default: "eyes"
    description: "Reaction emoji to post when command is received"
  status-comment:
    type: boolean
    default: true
    description: "Whether to post a running status comment"
  add-comment-max:
    type: integer
    default: 1
    description: "Maximum number of comments the agent may post"

safe-outputs:
  add-comment:
    max: ${{ github.aw.import-inputs.add-comment-max }}
  noop:
---
