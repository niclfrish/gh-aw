# Agent Performance - 2026-04-19
Run: §24621102435 | Q:76↑2 E:75↑2

## Ecosystem Overview (Apr 18-19)
- Overall success rate: ~85% (↑2% from Apr 18)
- AWF bumped to v0.25.25 ✅; MCP Gateway v0.2.25 ✅
- 196 workflows total (stable)
- Codex 401 auth still P0; Copilot/Claude stable

## Top Performers
1. **[aw] Failure Investigator** (Q:91 E:88) - Outstanding RCA on Codex 401 (#27127); identified gpt-5.3-codex endpoint issue
2. **Agentic Maintenance** (Q:88 E:100) - 2/2 success; caught lock file drift post-AWF bump (#27140)
3. **CLI Version Checker** (Q:87 E:100) - Copilot 1.0.32 + Claude 2.1.114 upgrade issue (#27143)
4. **Issue Monster** (Q:85 E:95) - 3/3 runs today, consistent
5. **Daily CLI Performance** (Q:84 E:85) - Caught BenchmarkFindIncludesInContent +51.4% (#26995) + BenchmarkValidation +24% (#26993)
6. **Copilot Optimization Agent** (Q:82 E:90) - Data-driven: branch proliferation (#27131), reviewer fan-out (#27130)

## Recovery
- **Agent Persona Explorer** ✅ RECOVERED (was 100% fail Apr 18 with 1.68M wasted tokens; now successful after AWF v0.25.25)
- Monitor 3+ more runs before removing from watch list

## Watch List
- **AI Moderator** (Q:45 E:45) - Codex 401 (#27127, #27122); P0 unresolved
- **Daily Observability Report** (Q:35 E:25) - Same Codex 401; zero output
- **GitHub Remote MCP Auth Test** (Q:50 E:0) - New failure today (#24620886472)
- **Smoke Claude** (Q:55 E:40) - Issue group #27030; failures since Apr 14
- **Smoke Copilot** (Q:60 E:55) - Issue group #27028

## P0 Active Issues
- **Codex 401 auth** (#27127, OPEN): OPENAI_API_KEY / gpt-5.3-codex access; needs admin rotation

## Issues Created This Run
- None new (all patterns already tracked; report in discussion)

## Key Findings
- AWF v0.25.25 resolved Agent Persona Explorer failures (Node.js binary)
- Copilot CLI 1.0.32 upgrade critical (11 versions behind)
- 20+ auto-created failure issues closed in 24h — over-creation pattern
- Performance regressions detected and assigned to Copilot for fix

Last updated: 2026-04-19T04:41Z by agent-performance-manager
