![GitHub Badge Blue](https://user-images.githubusercontent.com/25454503/157903512-a9be0f7b-9255-4f88-9b00-9d50539dd901.svg)

# Pluralith CLI

Pluralith is a tool to visualize your Terraform state and automate infrastructure documentation. The **Pluralith CLI** is a tool written in **`go`** to:

- Integrate Pluralith with Terraform
- Interact with the **[Pluralith UI](https://www.pluralith.com)**
- Run Pluralith in **Pipelines** for automated infrastructure documentation
- Ship other useful little features

![flow-illustration](https://user-images.githubusercontent.com/25454503/157021111-816c9936-3232-455f-9709-c3a65f5f8dfe.svg)

`Pluralith is currently in Private Alpha`

&nbsp;

## ⚙️ Getting Started

We are currently working on getting `pluralith` into **homebrew/core** and **apt**! For homebrew we need 75 stars, 30 forks and 30 watchers. We'd greatly appreciate if you can help us out with that.

Until we manage to get into these package managers you can manually install the **Pluralith CLI** by downloading your OS's binary from the latest release or by installing the **Pluralith UI** available [here](https://www.pluralith.com). The UI will install the CLI and keep it updated automatically.

&nbsp;

## 🛰️ CLI Overview

### → Terraform Commands

- `pluralith plan`: Creates a Terraform execution plan and opens it as a Diagram in the **Pluralith UI**<sup>1</sup>.
- `pluralith apply`: Essentially the same as `pluralith plan` with the intention of actually applying the execution plan<sup>1</sup>.
- `pluralith destroy`: Creates a Terraform execution plan in destroy mode and opens it as a Diagram in the **Pluralith UI**

### → Diagram Export Command

- `pluralith graph`: Creates a Terraform execution plan, draws a graph and exports it as a PDF in a headless way<sup>2</sup>
  - `--title`: The title to be shown in the diagram _(e.g. "Pluralith EKS Cluster")_
  - `--author`: The title to be shown in the diagram _(e.g. "Tony Stark")_
  - `--version`: The title to be shown in the diagram _(e.g. "1.0.5")_
  - `--show-changes`: Enables change highlighting in the output diagram. When enabled, resources that have been added, updated, deleted etc. will be highlighted with special colors
  - `--show-drift`: Enables drift highlighting in the output diagram. When enabled, resources that Terraform has detected drift for will be highlighted with a special badge and color
  - `--skip-plan`: Skip the generation of a new execution plan
    - Only works if Pluralith has run in the current directory before
  - `--out-dir`: The path your exported diagram PDF gets saved to _(e.g. "~/pluralith-infra/eks")_
    - Saved to current directory by default
  - `--file-name`: The path your exported diagram PDF gets saved to
    - The value passed for `--title` is used by default
  - `--generate-md`: Generates markdown for GitHub pull request / commit comment _(used in our [Pluralith GitHub actions](https://github.com/Pluralith/actions))_

&nbsp;

### 📍 Here's an example output for one of our test projects. View the PDF version **[here](https://github.com/Pluralith/pluralith-cli/files/8197192/HighlyAvailableIaaS.pdf)**

![HighlyAvailableIaaS](https://user-images.githubusercontent.com/25454503/157020490-8dadf7a2-ccb6-4323-a5d1-596d264bb06e.png)

### → Strip Command

- `pluralith strip`: Strips and hashes your plan state to make it shareable with us for debugging
  - Takes an existing **Pluralith Plan state** and subjects it to rigorous hashing of values
    - The Pluralith Plan state is located in the file _pluralith.state.stripped_ in your project directory
  - The purpose of this command is to strip the state of all sensitive data while keeping the structure intact, making it shareable
  - This is meant for us to debug edge cases on user state without the security hazard

### → Module Commands

- `pluralith install`: Installs/updates the specific module whose name is passed (e.g. `pluralith install graph-module`)
- `pluralith update`: Essentially the same as `install`. Updates existing modules, if not installed it downloads the latest release
  - If no value is passed, the latest version of the CLI itself will be installed

### → Utility Commands

- `pluralith login`: Authenticate with your API key (necessary for the CLI to work without the UI)<sup>2</sup>
- `pluralith version`: Shows information about the current CLI version as well as additional, installed modules

&nbsp;

<sup>1</sup> The UI then shows a prompt that lets you confirm or deny an `apply` with hotkeys.  
<sup>2</sup> You need to be authenticated with your **API key** via `pluralith login`. Currently only available for closed alpha testers. Interested? Shoot us an email dan@pluralith.com

&nbsp;

## 📦 Modules

The **Pluralith CLI** works with modules under the hood to extend its functionality. Below you can see a table of modules.

| **Module**       | **Installation**                 | **Description**                                 |
| ---------------- | -------------------------------- | ----------------------------------------------- |
| **Graph Module** | `pluralith install graph-module` | Installs the latest version of the graph module |
| **Pluralith UI** | `pluralith install ui`           | Installs the Pluralith UI. **`Coming Soon`**    |

&nbsp;

## 👩‍🚀 Looking to become a tester or talk about the project?

- Sign up for the `alpha` over on our **[Website](https://www.pluralith.com)**
- Join our **[Subreddit](https://www.reddit.com/r/Pluralith/)**
- Check out our **[Roadmap](https://roadmap.pluralith.com)** and upvote features you'd like to see
- Or just shoot us a message on Linkedin:
  - [Dan's Linkedin](https://www.linkedin.com/in/danielputzer/)
  - [Phi's Linkedin](https://www.linkedin.com/in/philipp-weber-a8517b231/)

_Disclaimer: To properly use this CLI you **will need** the **Pluralith UI** and/or an **API key**. [Sign up](https://www.pluralith.com) for the private alpha to get access!_

![Subreddit subscribers](https://img.shields.io/reddit/subreddit-subscribers/pluralith?style=social)
