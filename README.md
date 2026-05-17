<!--

Hello fellow agent! Welcome to GitHub Agentic Workflows = Actions + Agent + Safety. 
Here are some pointers to get you started in using this tool.

- Create a new workflow: https://raw.githubusercontent.com/github/gh-aw/main/create.md
- Install: https://raw.githubusercontent.com/github/gh-aw/main/install.md
- Package a workflow repository: https://raw.githubusercontent.com/github/gh-aw/main/package.md
- Reference: https://raw.githubusercontent.com/github/gh-aw/main/.github/aw/github-agentic-workflows.md

-->

# GitHub Agentic Workflows

Write agentic workflows in natural language markdown, and run them in GitHub Actions.

> [!NOTE]
> **Releases 0.68.4 through 0.71.3 are being retired** due to a bug that impacts billing. If you are running one of these versions, please upgrade to the latest release as soon as possible.

## Contents

- [Quick Start](#quick-start)
- [Overview](#overview)
- [Guardrails](#guardrails)
- [Documentation](#documentation)
- [FAQ](#faq)
- [Contributing](#contributing)
- [Community Contributions](#-community-contributions)
- [Share Feedback](#share-feedback)
- [Peli's Agent Factory](#pelis-agent-factory)
- [Related Projects](#related-projects)

## Quick Start

Ready to get your first agentic workflow running? Follow our step-by-step [Quick Start Guide](https://github.github.com/gh-aw/setup/quick-start/) to install the extension, add a sample workflow, and see it in action.

## Overview

Learn about the concepts behind agentic workflows, explore available workflow types, and understand how AI can automate your repository tasks. See [How It Works](https://github.github.com/gh-aw/introduction/how-they-work/).

## Guardrails

Guardrails, safety and security are foundational to GitHub Agentic Workflows. Workflows run with read-only permissions by default, with write operations only allowed through sanitized `safe-outputs`. The system implements multiple layers of protection including sandboxed execution, input sanitization, network isolation, supply chain security (SHA-pinned dependencies), tool allow-listing, and compile-time validation. Access can be gated to team members only, with human approval gates for critical operations, ensuring AI agents operate safely within controlled boundaries. See the [Security Architecture](https://github.github.com/gh-aw/introduction/architecture/) for comprehensive details on threat modeling, implementation guidelines, and best practices.

Using agentic workflows in your repository requires careful attention to security considerations and careful human supervision, and even then things can still go wrong. Use it with caution, and at your own risk.

## Documentation

For complete documentation, examples, and guides, see the [Documentation](https://github.github.com/gh-aw/). If you are an agent, download the [llms.txt](https://github.github.com/gh-aw/llms.txt).

If you are running a version between 0.68.4 and 0.71.3, upgrading is strongly recommended due to a bug that impacts billing.

## FAQ

For answers to common questions—billing, engine selection, workflow lock files, safe outputs, and more—see the [full FAQ](https://github.github.com/gh-aw/reference/faq/) in the documentation.

## Contributing

For development setup and contribution guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

### Custom Go linters

To build and test repository custom linters:

- `go test ./pkg/linters/<linter-name>/...`
- `go build ./cmd/linters`
- `make golint-custom`

`make golint-custom` builds `cmd/linters` and runs the custom analyzers against `./cmd/...` and `./pkg/...`.

## 🌍 Community Contributions

<details>
<summary>Thank you to the community members whose issue reports were resolved in this project! This list is updated automatically and reflects all attributed contributions.</summary>

- @abillingsley: [#23736](https://github.com/github/gh-aw/issues/23736) _(report)_
- @adamhenson: [#25345](https://github.com/github/gh-aw/issues/25345) _(report)_, [#24282](https://github.com/github/gh-aw/issues/24282) _(report)_
- @ahmadabdalla: [#27473](https://github.com/github/gh-aw/issues/27473) _(report)_
- @ajfeldman6: [#23924](https://github.com/github/gh-aw/issues/23924) _(report)_
- @AlexDeMichieli: [#26645](https://github.com/github/gh-aw/issues/26645) _(report)_
- @alexsiilvaa: [#20781](https://github.com/github/gh-aw/issues/20781) _(report)_, [#20664](https://github.com/github/gh-aw/issues/20664) _(report)_
- @alondahari: [#21207](https://github.com/github/gh-aw/issues/21207) _(report)_
- @anthonymastreanvae: [#30897](https://github.com/github/gh-aw/issues/30897) _(report)_, [#30841](https://github.com/github/gh-aw/issues/30841) _(report)_
- @apenab: [#25626](https://github.com/github/gh-aw/issues/25626) _(report)_
- @app/github-actions: [#31288](https://github.com/github/gh-aw/issues/31288) _(report)_, [#30740](https://github.com/github/gh-aw/issues/30740) _(report)_, [#29561](https://github.com/github/gh-aw/issues/29561) _(report)_, [#29343](https://github.com/github/gh-aw/issues/29343) _(report)_, [#26257](https://github.com/github/gh-aw/issues/26257) _(report)_, [#26256](https://github.com/github/gh-aw/issues/26256) _(report)_, [#26255](https://github.com/github/gh-aw/issues/26255) _(report)_, [#26254](https://github.com/github/gh-aw/issues/26254) _(report)_, [#26253](https://github.com/github/gh-aw/issues/26253) _(report)_
- @arezero: [#20515](https://github.com/github/gh-aw/issues/20515) _(report)_, [#20514](https://github.com/github/gh-aw/issues/20514) _(report)_, [#20513](https://github.com/github/gh-aw/issues/20513) _(report)_, [#20512](https://github.com/github/gh-aw/issues/20512) _(report)_, [#20511](https://github.com/github/gh-aw/issues/20511) _(report)_, [#20510](https://github.com/github/gh-aw/issues/20510) _(report)_
- @arthurfvives: [#30356](https://github.com/github/gh-aw/issues/30356) _(report)_, [#30088](https://github.com/github/gh-aw/issues/30088) _(report)_, [#26223](https://github.com/github/gh-aw/issues/26223) _(report)_, [#25993](https://github.com/github/gh-aw/issues/25993) _(report)_, [#25294](https://github.com/github/gh-aw/issues/25294) _(report)_
- @askpaisa: [#29240](https://github.com/github/gh-aw/issues/29240) _(report)_
- @b2pacific: [#28720](https://github.com/github/gh-aw/issues/28720) _(report)_
- @bartul: [#29499](https://github.com/github/gh-aw/issues/29499) _(report)_
- @bbonafed: [#29174](https://github.com/github/gh-aw/issues/29174) _(report)_, [#29173](https://github.com/github/gh-aw/issues/29173) _(report)_, [#29172](https://github.com/github/gh-aw/issues/29172) _(report)_, [#29171](https://github.com/github/gh-aw/issues/29171) _(report)_, [#27670](https://github.com/github/gh-aw/issues/27670) _(report)_, [#27472](https://github.com/github/gh-aw/issues/27472) _(report)_, [#26719](https://github.com/github/gh-aw/issues/26719) _(report)_, [#26045](https://github.com/github/gh-aw/issues/26045) _(report)_, [#26043](https://github.com/github/gh-aw/issues/26043) _(report)_, [#25646](https://github.com/github/gh-aw/issues/25646) _(report)_, [#25224](https://github.com/github/gh-aw/issues/25224) _(report)_, [#24949](https://github.com/github/gh-aw/issues/24949) _(report)_, [#24918](https://github.com/github/gh-aw/issues/24918) _(report)_, [#24896](https://github.com/github/gh-aw/issues/24896) _(report)_, [#24323](https://github.com/github/gh-aw/issues/24323) _(report)_, [#23900](https://github.com/github/gh-aw/issues/23900) _(report)_, [#23724](https://github.com/github/gh-aw/issues/23724) _(report)_, [#23566](https://github.com/github/gh-aw/issues/23566) _(report)_, [#22564](https://github.com/github/gh-aw/issues/22564) _(report)_, [#21990](https://github.com/github/gh-aw/issues/21990) _(report)_, [#20801](https://github.com/github/gh-aw/issues/20801) _(report)_, [#20378](https://github.com/github/gh-aw/issues/20378) _(report)_
- @benvillalobos: [#25717](https://github.com/github/gh-aw/issues/25717) _(report)_, [#20885](https://github.com/github/gh-aw/issues/20885) _(report)_
- @bmerkle: [#31689](https://github.com/github/gh-aw/issues/31689) _(report)_, [#26621](https://github.com/github/gh-aw/issues/26621) _(report)_, [#20646](https://github.com/github/gh-aw/issues/20646) _(report)_
- @bryanchen-d: [#30866](https://github.com/github/gh-aw/issues/30866) _(report)_, [#30704](https://github.com/github/gh-aw/issues/30704) _(report)_, [#30695](https://github.com/github/gh-aw/issues/30695) _(report)_, [#30472](https://github.com/github/gh-aw/issues/30472) _(report)_, [#28774](https://github.com/github/gh-aw/issues/28774) _(report)_, [#26696](https://github.com/github/gh-aw/issues/26696) _(report)_, [#26487](https://github.com/github/gh-aw/issues/26487) _(report)_, [#25719](https://github.com/github/gh-aw/issues/25719) _(report)_, [#23265](https://github.com/github/gh-aw/issues/23265) _(report)_
- @bryanknox: [#25351](https://github.com/github/gh-aw/issues/25351) _(report)_
- @Calidus: [#26923](https://github.com/github/gh-aw/issues/26923) _(report)_
- @camposbrunocampos: [#23726](https://github.com/github/gh-aw/issues/23726) _(report)_, [#22897](https://github.com/github/gh-aw/issues/22897) _(report)_
- @carlincherry: [#22017](https://github.com/github/gh-aw/issues/22017) _(report)_
- @chepa92: [#20322](https://github.com/github/gh-aw/issues/20322) _(report)_
- @chrisfregly: [#25349](https://github.com/github/gh-aw/issues/25349) _(report)_, [#23963](https://github.com/github/gh-aw/issues/23963) _(report)_
- @chrizbo: [#31399](https://github.com/github/gh-aw/issues/31399) _(report)_, [#28158](https://github.com/github/gh-aw/issues/28158) _(report)_, [#22510](https://github.com/github/gh-aw/issues/22510) _(report)_, [#21863](https://github.com/github/gh-aw/issues/21863) _(report)_
- @CiscoRob: [#20416](https://github.com/github/gh-aw/issues/20416) _(report)_
- @clementbolin: [#28888](https://github.com/github/gh-aw/issues/28888) _(report)_
- @Corb3nik: [#21306](https://github.com/github/gh-aw/issues/21306) _(report)_
- @corygehr: [#27638](https://github.com/github/gh-aw/issues/27638) _(report)_, [#26539](https://github.com/github/gh-aw/issues/26539) _(report)_, [#26270](https://github.com/github/gh-aw/issues/26270) _(report)_, [#26268](https://github.com/github/gh-aw/issues/26268) _(report)_, [#25680](https://github.com/github/gh-aw/issues/25680) _(report)_, [#24355](https://github.com/github/gh-aw/issues/24355) _(report)_, [#23944](https://github.com/github/gh-aw/issues/23944) _(report)_, [#23753](https://github.com/github/gh-aw/issues/23753) _(report)_
- @corymhall: [#19839](https://github.com/github/gh-aw/issues/19839) _(report)_
- @dagecko: [#24743](https://github.com/github/gh-aw/issues/24743) _(report)_
- @Daidanny008: [#27402](https://github.com/github/gh-aw/issues/27402) _(report)_
- @Dan-Co: [#22707](https://github.com/github/gh-aw/issues/22707) _(report)_
- @danielmeppiel: [#29076](https://github.com/github/gh-aw/issues/29076) _(report)_, [#28678](https://github.com/github/gh-aw/issues/28678) _(report)_, [#20663](https://github.com/github/gh-aw/issues/20663) _(report)_, [#20380](https://github.com/github/gh-aw/issues/20380) _(report)_, [#19810](https://github.com/github/gh-aw/issues/19810) _(report)_
- @danquirk: [#30403](https://github.com/github/gh-aw/issues/30403) _(report)_
- @dbudym-cs: [#22913](https://github.com/github/gh-aw/issues/22913) _(report)_
- @devantler: [#25768](https://github.com/github/gh-aw/issues/25768) _(report)_, [#25767](https://github.com/github/gh-aw/issues/25767) _(report)_
- @deyaaeldeen: [#28966](https://github.com/github/gh-aw/issues/28966) _(report)_, [#26486](https://github.com/github/gh-aw/issues/26486) _(report)_, [#25573](https://github.com/github/gh-aw/issues/25573) _(report)_, [#25359](https://github.com/github/gh-aw/issues/25359) _(report)_, [#23198](https://github.com/github/gh-aw/issues/23198) _(report)_, [#23024](https://github.com/github/gh-aw/issues/23024) _(report)_, [#23020](https://github.com/github/gh-aw/issues/23020) _(report)_, [#22957](https://github.com/github/gh-aw/issues/22957) _(report)_, [#19773](https://github.com/github/gh-aw/issues/19773) _(report)_, [#19770](https://github.com/github/gh-aw/issues/19770) _(report)_
- @dholmes: [#29228](https://github.com/github/gh-aw/issues/29228) _(report)_, [#23578](https://github.com/github/gh-aw/issues/23578) _(report)_
- @DimaBir: [#20483](https://github.com/github/gh-aw/issues/20483) _(report)_
- @dkurepa: [#25511](https://github.com/github/gh-aw/issues/25511) _(report)_
- @DogeAmazed: [#22703](https://github.com/github/gh-aw/issues/22703) _(report)_
- @doughgle: [#23655](https://github.com/github/gh-aw/issues/23655) _(report)_
- @drehelis: [#25304](https://github.com/github/gh-aw/issues/25304) _(report)_
- @dsyme: [#23936](https://github.com/github/gh-aw/issues/23936) _(report)_, [#22340](https://github.com/github/gh-aw/issues/22340) _(report)_, [#20953](https://github.com/github/gh-aw/issues/20953) _(report)_, [#20952](https://github.com/github/gh-aw/issues/20952) _(report)_, [#20950](https://github.com/github/gh-aw/issues/20950) _(report)_, [#20787](https://github.com/github/gh-aw/issues/20787) _(report)_, [#20578](https://github.com/github/gh-aw/issues/20578) _(report)_, [#20420](https://github.com/github/gh-aw/issues/20420) _(report)_, [#20243](https://github.com/github/gh-aw/issues/20243) _(report)_, [#20241](https://github.com/github/gh-aw/issues/20241) _(report)_, [#20108](https://github.com/github/gh-aw/issues/20108) _(report)_, [#20103](https://github.com/github/gh-aw/issues/20103) _(report)_, [#19976](https://github.com/github/gh-aw/issues/19976) _(report)_, [#19708](https://github.com/github/gh-aw/issues/19708) _(report)_, [#19468](https://github.com/github/gh-aw/issues/19468) _(report)_, [#19465](https://github.com/github/gh-aw/issues/19465) _(report)_
- @duncankmckinnon: [#25944](https://github.com/github/gh-aw/issues/25944) _(report)_
- @eaftan: [#23257](https://github.com/github/gh-aw/issues/23257) _(report)_, [#20457](https://github.com/github/gh-aw/issues/20457) _(report)_
- @edburns: [#26920](https://github.com/github/gh-aw/issues/26920) _(report)_
- @edgeq: [#28315](https://github.com/github/gh-aw/issues/28315) _(report)_, [#28308](https://github.com/github/gh-aw/issues/28308) _(report)_
- @ericchansen: [#20222](https://github.com/github/gh-aw/issues/20222) _(report)_
- @ericstj: [#30260](https://github.com/github/gh-aw/issues/30260) _(report)_, [#23766](https://github.com/github/gh-aw/issues/23766) _(report)_
- @Esomoire-consultancy-Company: [#20207](https://github.com/github/gh-aw/issues/20207) _(report)_
- @ferryhinardi: [#24128](https://github.com/github/gh-aw/issues/24128) _(report)_
- @flatiron32: [#22469](https://github.com/github/gh-aw/issues/22469) _(report)_
- @fr4nc1sc0-r4m0n: [#20657](https://github.com/github/gh-aw/issues/20657) _(report)_
- @G1Vh: [#20308](https://github.com/github/gh-aw/issues/20308) _(report)_
- @glitch-ux: [#24403](https://github.com/github/gh-aw/issues/24403) _(report)_
- @grahame-white: [#23643](https://github.com/github/gh-aw/issues/23643) _(report)_, [#23093](https://github.com/github/gh-aw/issues/23093) _(report)_, [#23092](https://github.com/github/gh-aw/issues/23092) _(report)_, [#23088](https://github.com/github/gh-aw/issues/23088) _(report)_, [#23083](https://github.com/github/gh-aw/issues/23083) _(report)_, [#20868](https://github.com/github/gh-aw/issues/20868) _(report)_, [#20719](https://github.com/github/gh-aw/issues/20719) _(report)_, [#20629](https://github.com/github/gh-aw/issues/20629) _(report)_, [#20299](https://github.com/github/gh-aw/issues/20299) _(report)_
- @h3y6e: [#27794](https://github.com/github/gh-aw/issues/27794) _(report)_
- @haavamoa: [#30191](https://github.com/github/gh-aw/issues/30191) _(report)_
- @heiskr: [#20394](https://github.com/github/gh-aw/issues/20394) _(report)_
- @hermanho: [#32197](https://github.com/github/gh-aw/issues/32197) _(report)_
- @holwerda: [#21243](https://github.com/github/gh-aw/issues/21243) _(report)_
- @hrishikeshathalye: [#19547](https://github.com/github/gh-aw/issues/19547) _(report)_
- @IEvangelist: [#32536](https://github.com/github/gh-aw/issues/32536) _(report)_, [#32354](https://github.com/github/gh-aw/issues/32354) _(report)_, [#30848](https://github.com/github/gh-aw/issues/30848) _(report)_, [#26908](https://github.com/github/gh-aw/issues/26908) _(report)_, [#25467](https://github.com/github/gh-aw/issues/25467) _(report)_
- @Infinnerty: [#21957](https://github.com/github/gh-aw/issues/21957) _(report)_
- @insop: [#21686](https://github.com/github/gh-aw/issues/21686) _(report)_
- @j-srodka: [#25199](https://github.com/github/gh-aw/issues/25199) _(report)_, [#23485](https://github.com/github/gh-aw/issues/23485) _(report)_, [#23484](https://github.com/github/gh-aw/issues/23484) _(report)_, [#23483](https://github.com/github/gh-aw/issues/23483) _(report)_, [#23482](https://github.com/github/gh-aw/issues/23482) _(report)_, [#23461](https://github.com/github/gh-aw/issues/23461) _(report)_
- @jamesadevine: [#28957](https://github.com/github/gh-aw/issues/28957) _(report)_, [#26407](https://github.com/github/gh-aw/issues/26407) _(report)_, [#26406](https://github.com/github/gh-aw/issues/26406) _(report)_
- @JamesNK: [#29310](https://github.com/github/gh-aw/issues/29310) _(report)_, [#28867](https://github.com/github/gh-aw/issues/28867) _(report)_, [#28704](https://github.com/github/gh-aw/issues/28704) _(report)_
- @JanKrivanek: [#25656](https://github.com/github/gh-aw/issues/25656) _(report)_, [#25439](https://github.com/github/gh-aw/issues/25439) _(report)_, [#20187](https://github.com/github/gh-aw/issues/20187) _(report)_
- @jaroslawgajewski: [#31678](https://github.com/github/gh-aw/issues/31678) _(report)_, [#25593](https://github.com/github/gh-aw/issues/25593) _(report)_, [#24373](https://github.com/github/gh-aw/issues/24373) _(report)_, [#24372](https://github.com/github/gh-aw/issues/24372) _(report)_, [#24371](https://github.com/github/gh-aw/issues/24371) _(report)_, [#24259](https://github.com/github/gh-aw/issues/24259) _(report)_, [#24036](https://github.com/github/gh-aw/issues/24036) _(report)_, [#23779](https://github.com/github/gh-aw/issues/23779) _(report)_, [#23558](https://github.com/github/gh-aw/issues/23558) _(report)_, [#22647](https://github.com/github/gh-aw/issues/22647) _(report)_, [#21816](https://github.com/github/gh-aw/issues/21816) _(report)_, [#20813](https://github.com/github/gh-aw/issues/20813) _(report)_, [#20811](https://github.com/github/gh-aw/issues/20811) _(report)_, [#19732](https://github.com/github/gh-aw/issues/19732) _(report)_
- @JasonYeMSFT: [#27424](https://github.com/github/gh-aw/issues/27424) _(report)_
- @jbaruch: [#30832](https://github.com/github/gh-aw/issues/30832) _(report)_
- @jeffhandley: [#30232](https://github.com/github/gh-aw/issues/30232) _(report)_, [#30204](https://github.com/github/gh-aw/issues/30204) _(report)_, [#26799](https://github.com/github/gh-aw/issues/26799) _(report)_, [#26788](https://github.com/github/gh-aw/issues/26788) _(report)_, [#24384](https://github.com/github/gh-aw/issues/24384) _(report)_
- @jfomhover: [#25420](https://github.com/github/gh-aw/issues/25420) _(report)_
- @johnpreed: [#25687](https://github.com/github/gh-aw/issues/25687) _(report)_, [#23777](https://github.com/github/gh-aw/issues/23777) _(report)_, [#23212](https://github.com/github/gh-aw/issues/23212) _(report)_, [#21334](https://github.com/github/gh-aw/issues/21334) _(report)_
- @johnwilliams-12: [#21205](https://github.com/github/gh-aw/issues/21205) _(report)_, [#21074](https://github.com/github/gh-aw/issues/21074) _(report)_, [#21071](https://github.com/github/gh-aw/issues/21071) _(report)_, [#21062](https://github.com/github/gh-aw/issues/21062) _(report)_, [#20821](https://github.com/github/gh-aw/issues/20821) _(report)_, [#20779](https://github.com/github/gh-aw/issues/20779) _(report)_, [#20697](https://github.com/github/gh-aw/issues/20697) _(report)_, [#20694](https://github.com/github/gh-aw/issues/20694) _(report)_, [#20658](https://github.com/github/gh-aw/issues/20658) _(report)_, [#20567](https://github.com/github/gh-aw/issues/20567) _(report)_, [#20508](https://github.com/github/gh-aw/issues/20508) _(report)_
- @jonathanpeppers: [#30662](https://github.com/github/gh-aw/issues/30662) _(report)_
- @jsoref: [#27230](https://github.com/github/gh-aw/issues/27230) _(report)_
- @jtracey93: [#26176](https://github.com/github/gh-aw/issues/26176) _(report)_
- @kaovilai: [#32596](https://github.com/github/gh-aw/issues/32596) _(report)_, [#32587](https://github.com/github/gh-aw/issues/32587) _(report)_, [#32482](https://github.com/github/gh-aw/issues/32482) _(report)_, [#32467](https://github.com/github/gh-aw/issues/32467) _(report)_
- @kbreit-insight: [#24930](https://github.com/github/gh-aw/issues/24930) _(report)_, [#23940](https://github.com/github/gh-aw/issues/23940) _(report)_, [#23725](https://github.com/github/gh-aw/issues/23725) _(report)_, [#22430](https://github.com/github/gh-aw/issues/22430) _(report)_, [#21978](https://github.com/github/gh-aw/issues/21978) _(report)_
- @kkruel8100: [#30867](https://github.com/github/gh-aw/issues/30867) _(report)_
- @kthompson: [#25550](https://github.com/github/gh-aw/issues/25550) _(report)_
- @labudis: [#30846](https://github.com/github/gh-aw/issues/30846) _(report)_
- @look: [#23258](https://github.com/github/gh-aw/issues/23258) _(report)_
- @lpcox: [#30634](https://github.com/github/gh-aw/issues/30634) _(report)_, [#29353](https://github.com/github/gh-aw/issues/29353) _(report)_, [#29191](https://github.com/github/gh-aw/issues/29191) _(report)_, [#22281](https://github.com/github/gh-aw/issues/22281) _(report)_
- @lukeed: [#20552](https://github.com/github/gh-aw/issues/20552) _(report)_
- @lupinthe14th: [#26542](https://github.com/github/gh-aw/issues/26542) _(report)_, [#26441](https://github.com/github/gh-aw/issues/26441) _(report)_
- @mark-hingston: [#20335](https://github.com/github/gh-aw/issues/20335) _(report)_
- @mason-tim: [#31489](https://github.com/github/gh-aw/issues/31489) _(report)_, [#30336](https://github.com/github/gh-aw/issues/30336) _(report)_, [#29301](https://github.com/github/gh-aw/issues/29301) _(report)_, [#21562](https://github.com/github/gh-aw/issues/21562) _(report)_, [#19765](https://github.com/github/gh-aw/issues/19765) _(report)_
- @MatthewLabasan-NBCU: [#26289](https://github.com/github/gh-aw/issues/26289) _(report)_, [#19500](https://github.com/github/gh-aw/issues/19500) _(report)_
- @MattSkala: [#24567](https://github.com/github/gh-aw/issues/24567) _(report)_, [#21203](https://github.com/github/gh-aw/issues/21203) _(report)_
- @MauroDruwel: [#30178](https://github.com/github/gh-aw/issues/30178) _(report)_, [#30169](https://github.com/github/gh-aw/issues/30169) _(report)_, [#29379](https://github.com/github/gh-aw/issues/29379) _(report)_, [#29378](https://github.com/github/gh-aw/issues/29378) _(report)_
- @mcantrell: [#20592](https://github.com/github/gh-aw/issues/20592) _(report)_
- @mdashrraf: [#28657](https://github.com/github/gh-aw/issues/28657) _(report)_
- @MH0386: [#20997](https://github.com/github/gh-aw/issues/20997) _(report)_
- @mhavelock: [#22110](https://github.com/github/gh-aw/issues/22110) _(report)_
- @michen00: [#31869](https://github.com/github/gh-aw/issues/31869) _(report)_
- @microsasa: [#27715](https://github.com/github/gh-aw/issues/27715) _(report)_, [#21103](https://github.com/github/gh-aw/issues/21103) _(report)_, [#21098](https://github.com/github/gh-aw/issues/21098) _(report)_, [#20851](https://github.com/github/gh-aw/issues/20851) _(report)_, [#20833](https://github.com/github/gh-aw/issues/20833) _(report)_, [#20586](https://github.com/github/gh-aw/issues/20586) _(report)_
- @mlinksva: [#22533](https://github.com/github/gh-aw/issues/22533) _(report)_
- @mnkiefer: [#22409](https://github.com/github/gh-aw/issues/22409) _(report)_, [#19836](https://github.com/github/gh-aw/issues/19836) _(report)_
- @molson504x: [#21834](https://github.com/github/gh-aw/issues/21834) _(report)_, [#21615](https://github.com/github/gh-aw/issues/21615) _(report)_
- @Mossaka: [#21644](https://github.com/github/gh-aw/issues/21644) _(report)_, [#21630](https://github.com/github/gh-aw/issues/21630) _(report)_
- @mrjf: [#32271](https://github.com/github/gh-aw/issues/32271) _(report)_, [#32069](https://github.com/github/gh-aw/issues/32069) _(report)_, [#31600](https://github.com/github/gh-aw/issues/31600) _(report)_, [#29152](https://github.com/github/gh-aw/issues/29152) _(report)_, [#28955](https://github.com/github/gh-aw/issues/28955) _(report)_, [#28471](https://github.com/github/gh-aw/issues/28471) _(report)_, [#28197](https://github.com/github/gh-aw/issues/28197) _(report)_
- @mvdbos: [#20411](https://github.com/github/gh-aw/issues/20411) _(report)_, [#20249](https://github.com/github/gh-aw/issues/20249) _(report)_
- @neta-vega: [#26447](https://github.com/github/gh-aw/issues/26447) _(report)_, [#25895](https://github.com/github/gh-aw/issues/25895) _(report)_
- @NicoAvanzDev: [#21542](https://github.com/github/gh-aw/issues/21542) _(report)_, [#20540](https://github.com/github/gh-aw/issues/20540) _(report)_, [#20528](https://github.com/github/gh-aw/issues/20528) _(report)_
- @NicolasRannou: [#31701](https://github.com/github/gh-aw/issues/31701) _(report)_
- @NikolajBjorner: [#28812](https://github.com/github/gh-aw/issues/28812) _(report)_
- @norrietaylor: [#32312](https://github.com/github/gh-aw/issues/32312) _(report)_, [#32310](https://github.com/github/gh-aw/issues/32310) _(report)_, [#30733](https://github.com/github/gh-aw/issues/30733) _(report)_, [#30392](https://github.com/github/gh-aw/issues/30392) _(report)_
- @octatone: [#31918](https://github.com/github/gh-aw/issues/31918) _(report)_
- @petercort: [#28281](https://github.com/github/gh-aw/issues/28281) _(report)_
- @pethers: [#28470](https://github.com/github/gh-aw/issues/28470) _(report)_
- @pgaskin: [#26156](https://github.com/github/gh-aw/issues/26156) _(report)_
- @pholleran: [#25313](https://github.com/github/gh-aw/issues/25313) _(report)_, [#23572](https://github.com/github/gh-aw/issues/23572) _(report)_, [#21313](https://github.com/github/gh-aw/issues/21313) _(report)_
- @PureWeen: [#28767](https://github.com/github/gh-aw/issues/28767) _(report)_, [#27655](https://github.com/github/gh-aw/issues/27655) _(report)_, [#23769](https://github.com/github/gh-aw/issues/23769) _(report)_, [#23567](https://github.com/github/gh-aw/issues/23567) _(report)_
- @rabo-unumed: [#31660](https://github.com/github/gh-aw/issues/31660) _(report)_, [#31578](https://github.com/github/gh-aw/issues/31578) _(report)_, [#31513](https://github.com/github/gh-aw/issues/31513) _(report)_, [#20679](https://github.com/github/gh-aw/issues/20679) _(report)_
- @rhardouin: [#30840](https://github.com/github/gh-aw/issues/30840) _(report)_, [#30838](https://github.com/github/gh-aw/issues/30838) _(report)_
- @romainh-betclic: [#28143](https://github.com/github/gh-aw/issues/28143) _(report)_
- @rspurgeon: [#26475](https://github.com/github/gh-aw/issues/26475) _(report)_, [#19451](https://github.com/github/gh-aw/issues/19451) _(report)_
- @Rubyj: [#31542](https://github.com/github/gh-aw/issues/31542) _(report)_, [#21432](https://github.com/github/gh-aw/issues/21432) _(report)_, [#20283](https://github.com/github/gh-aw/issues/20283) _(report)_
- @ruokun-niu: [#24961](https://github.com/github/gh-aw/issues/24961) _(report)_
- @ryckmansm: [#31501](https://github.com/github/gh-aw/issues/31501) _(report)_
- @salekseev: [#25137](https://github.com/github/gh-aw/issues/25137) _(report)_, [#25122](https://github.com/github/gh-aw/issues/25122) _(report)_, [#24135](https://github.com/github/gh-aw/issues/24135) _(report)_
- @samuelkahessay: [#24756](https://github.com/github/gh-aw/issues/24756) _(report)_, [#24755](https://github.com/github/gh-aw/issues/24755) _(report)_, [#24754](https://github.com/github/gh-aw/issues/24754) _(report)_, [#22380](https://github.com/github/gh-aw/issues/22380) _(report)_, [#22364](https://github.com/github/gh-aw/issues/22364) _(report)_, [#22161](https://github.com/github/gh-aw/issues/22161) _(report)_, [#22138](https://github.com/github/gh-aw/issues/22138) _(report)_, [#21975](https://github.com/github/gh-aw/issues/21975) _(report)_, [#21955](https://github.com/github/gh-aw/issues/21955) _(report)_, [#21784](https://github.com/github/gh-aw/issues/21784) _(report)_, [#21501](https://github.com/github/gh-aw/issues/21501) _(report)_, [#21304](https://github.com/github/gh-aw/issues/21304) _(report)_, [#20035](https://github.com/github/gh-aw/issues/20035) _(report)_, [#20031](https://github.com/github/gh-aw/issues/20031) _(report)_, [#20030](https://github.com/github/gh-aw/issues/20030) _(report)_, [#19605](https://github.com/github/gh-aw/issues/19605) _(report)_, [#19476](https://github.com/github/gh-aw/issues/19476) _(report)_, [#19475](https://github.com/github/gh-aw/issues/19475) _(report)_, [#19474](https://github.com/github/gh-aw/issues/19474) _(report)_, [#19473](https://github.com/github/gh-aw/issues/19473) _(report)_
- @sbodapati-gfm: [#29417](https://github.com/github/gh-aw/issues/29417) _(report)_
- @seangibeault: [#26910](https://github.com/github/gh-aw/issues/26910) _(report)_, [#24905](https://github.com/github/gh-aw/issues/24905) _(report)_
- @sg650: [#32044](https://github.com/github/gh-aw/issues/32044) _(report)_, [#31617](https://github.com/github/gh-aw/issues/31617) _(report)_, [#31616](https://github.com/github/gh-aw/issues/31616) _(report)_, [#29009](https://github.com/github/gh-aw/issues/29009) _(report)_, [#28612](https://github.com/github/gh-aw/issues/28612) _(report)_
- @shiran-gutsy: [#27641](https://github.com/github/gh-aw/issues/27641) _(report)_
- @srgibbs99: [#22939](https://github.com/github/gh-aw/issues/22939) _(report)_, [#19640](https://github.com/github/gh-aw/issues/19640) _(report)_, [#19622](https://github.com/github/gh-aw/issues/19622) _(report)_
- @stacktick: [#21361](https://github.com/github/gh-aw/issues/21361) _(report)_
- @stefankrzyz: [#27260](https://github.com/github/gh-aw/issues/27260) _(report)_
- @straub: [#24569](https://github.com/github/gh-aw/issues/24569) _(report)_, [#19631](https://github.com/github/gh-aw/issues/19631) _(report)_
- @strawgate: [#24422](https://github.com/github/gh-aw/issues/24422) _(report)_, [#24199](https://github.com/github/gh-aw/issues/24199) _(report)_, [#23935](https://github.com/github/gh-aw/issues/23935) _(report)_, [#23768](https://github.com/github/gh-aw/issues/23768) _(report)_, [#21157](https://github.com/github/gh-aw/issues/21157) _(report)_, [#21144](https://github.com/github/gh-aw/issues/21144) _(report)_, [#21135](https://github.com/github/gh-aw/issues/21135) _(report)_, [#21028](https://github.com/github/gh-aw/issues/21028) _(report)_, [#20910](https://github.com/github/gh-aw/issues/20910) _(report)_, [#20259](https://github.com/github/gh-aw/issues/20259) _(report)_, [#20168](https://github.com/github/gh-aw/issues/20168) _(report)_, [#20125](https://github.com/github/gh-aw/issues/20125) _(report)_, [#20033](https://github.com/github/gh-aw/issues/20033) _(report)_, [#19982](https://github.com/github/gh-aw/issues/19982) _(report)_, [#19972](https://github.com/github/gh-aw/issues/19972) _(report)_
- @susmahad: [#26276](https://github.com/github/gh-aw/issues/26276) _(report)_, [#25866](https://github.com/github/gh-aw/issues/25866) _(report)_, [#25710](https://github.com/github/gh-aw/issues/25710) _(report)_
- @szabta89: [#29064](https://github.com/github/gh-aw/issues/29064) _(report)_, [#29063](https://github.com/github/gh-aw/issues/29063) _(report)_, [#24037](https://github.com/github/gh-aw/issues/24037) _(report)_
- @tadelesh: [#26001](https://github.com/github/gh-aw/issues/26001) _(report)_
- @theletterf: [#30964](https://github.com/github/gh-aw/issues/30964) _(report)_, [#30365](https://github.com/github/gh-aw/issues/30365) _(report)_, [#30327](https://github.com/github/gh-aw/issues/30327) _(report)_, [#28898](https://github.com/github/gh-aw/issues/28898) _(report)_, [#28895](https://github.com/github/gh-aw/issues/28895) _(report)_, [#28691](https://github.com/github/gh-aw/issues/28691) _(report)_, [#28672](https://github.com/github/gh-aw/issues/28672) _(report)_, [#28221](https://github.com/github/gh-aw/issues/28221) _(report)_, [#27566](https://github.com/github/gh-aw/issues/27566) _(report)_, [#25494](https://github.com/github/gh-aw/issues/25494) _(report)_
- @thi-feonir: [#21426](https://github.com/github/gh-aw/issues/21426) _(report)_
- @tinytelly: [#27282](https://github.com/github/gh-aw/issues/27282) _(report)_
- @tomasmed: [#20157](https://github.com/github/gh-aw/issues/20157) _(report)_
- @tore-unumed: [#31909](https://github.com/github/gh-aw/issues/31909) _(report)_, [#31650](https://github.com/github/gh-aw/issues/31650) _(report)_, [#30550](https://github.com/github/gh-aw/issues/30550) _(report)_, [#30324](https://github.com/github/gh-aw/issues/30324) _(report)_, [#29312](https://github.com/github/gh-aw/issues/29312) _(report)_, [#28019](https://github.com/github/gh-aw/issues/28019) _(report)_, [#20780](https://github.com/github/gh-aw/issues/20780) _(report)_, [#19703](https://github.com/github/gh-aw/issues/19703) _(report)_
- @trask: [#31612](https://github.com/github/gh-aw/issues/31612) _(report)_, [#31241](https://github.com/github/gh-aw/issues/31241) _(report)_, [#31098](https://github.com/github/gh-aw/issues/31098) _(report)_, [#31097](https://github.com/github/gh-aw/issues/31097) _(report)_
- @tsm-harmoney: [#31695](https://github.com/github/gh-aw/issues/31695) _(report)_, [#27880](https://github.com/github/gh-aw/issues/27880) _(report)_
- @tspascoal: [#20597](https://github.com/github/gh-aw/issues/20597) _(report)_
- @UncleBats: [#20359](https://github.com/github/gh-aw/issues/20359) _(report)_
- @verkyyi: [#27407](https://github.com/github/gh-aw/issues/27407) _(report)_, [#27259](https://github.com/github/gh-aw/issues/27259) _(report)_
- @veverkap: [#22362](https://github.com/github/gh-aw/issues/22362) _(report)_, [#21260](https://github.com/github/gh-aw/issues/21260) _(report)_, [#21257](https://github.com/github/gh-aw/issues/21257) _(report)_
- @virenpepper: [#23765](https://github.com/github/gh-aw/issues/23765) _(report)_
- @wtgodbe: [#26057](https://github.com/github/gh-aw/issues/26057) _(report)_, [#25130](https://github.com/github/gh-aw/issues/25130) _(report)_, [#24921](https://github.com/github/gh-aw/issues/24921) _(report)_
- @yaananth: [#24125](https://github.com/github/gh-aw/issues/24125) _(report)_
- @yskopets: [#32022](https://github.com/github/gh-aw/issues/32022) _(report)_, [#31831](https://github.com/github/gh-aw/issues/31831) _(report)_, [#31086](https://github.com/github/gh-aw/issues/31086) _(report)_, [#31073](https://github.com/github/gh-aw/issues/31073) _(report)_, [#30872](https://github.com/github/gh-aw/issues/30872) _(report)_, [#30705](https://github.com/github/gh-aw/issues/30705) _(report)_, [#27935](https://github.com/github/gh-aw/issues/27935) _(report)_, [#27898](https://github.com/github/gh-aw/issues/27898) _(report)_, [#27881](https://github.com/github/gh-aw/issues/27881) _(report)_, [#27773](https://github.com/github/gh-aw/issues/27773) _(report)_, [#27757](https://github.com/github/gh-aw/issues/27757) _(report)_, [#26922](https://github.com/github/gh-aw/issues/26922) _(report)_, [#26569](https://github.com/github/gh-aw/issues/26569) _(report)_, [#26468](https://github.com/github/gh-aw/issues/26468) _(report)_, [#26358](https://github.com/github/gh-aw/issues/26358) _(report)_, [#26346](https://github.com/github/gh-aw/issues/26346) _(report)_, [#26345](https://github.com/github/gh-aw/issues/26345) _(report)_, [#26280](https://github.com/github/gh-aw/issues/26280) _(report)_, [#26279](https://github.com/github/gh-aw/issues/26279) _(report)_, [#26120](https://github.com/github/gh-aw/issues/26120) _(report)_, [#26101](https://github.com/github/gh-aw/issues/26101) _(report)_, [#26085](https://github.com/github/gh-aw/issues/26085) _(report)_, [#26080](https://github.com/github/gh-aw/issues/26080) _(report)_, [#26067](https://github.com/github/gh-aw/issues/26067) _(report)_, [#25959](https://github.com/github/gh-aw/issues/25959) _(report)_, [#25946](https://github.com/github/gh-aw/issues/25946) _(report)_, [#25833](https://github.com/github/gh-aw/issues/25833) _(report)_, [#25363](https://github.com/github/gh-aw/issues/25363) _(report)_, [#25362](https://github.com/github/gh-aw/issues/25362) _(report)_, [#25125](https://github.com/github/gh-aw/issues/25125) _(report)_, [#24897](https://github.com/github/gh-aw/issues/24897) _(report)_, [#24573](https://github.com/github/gh-aw/issues/24573) _(report)_, [#23914](https://github.com/github/gh-aw/issues/23914) _(report)_
- @zkoppert: [#27741](https://github.com/github/gh-aw/issues/27741) _(report)_

</details>
## Share Feedback

We welcome your feedback on GitHub Agentic Workflows! 

- [Community Feedback Discussions](https://github.com/orgs/community/discussions/186451)
- [GitHub Next](https://githubnext.com/)

## Peli's Agent Factory

See the [Peli's Agent Factory](https://github.github.com/gh-aw/blog/2026-01-12-welcome-to-pelis-agent-factory/) for a guided tour through many uses of agentic workflows.

## Related Projects

GitHub Agentic Workflows is supported by companion projects that provide additional security and integration capabilities:

- **[Agent Workflow Firewall (AWF)](https://github.com/github/gh-aw-firewall)** - Network egress control for AI agents, providing domain-based access controls and activity logging for secure workflow execution
- **[MCP Gateway](https://github.com/github/gh-aw-mcpg)** - Routes Model Context Protocol (MCP) server calls through a unified HTTP gateway for centralized access management
- **[gh-aw-actions](https://github.com/github/gh-aw-actions)** - Shared library of custom GitHub Actions used by compiled workflows, providing functionality such as MCP server file management
