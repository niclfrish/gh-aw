---
"gh-aw": patch
---

Fixed AWF resolution in the PR Sous Chef detection job by making the installed `awf` binary available under a secure sudo path on Linux and validating the install-time compatibility symlink.
