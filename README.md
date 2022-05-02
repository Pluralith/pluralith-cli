![GitHub Badge Blue](https://user-images.githubusercontent.com/25454503/157903512-a9be0f7b-9255-4f88-9b00-9d50539dd901.svg)

# Pluralith CLI

Pluralith is a tool to visualise your Terraform state and automate infrastructure documentation.

`Pluralith is currently in Private Alpha`

&nbsp;

## ⚡ Highlights

- Create beautiful `infrastructure diagrams` with one single command: `pluralith plan`
- Highlight current plan `changes` in the diagram
- Detect and visualise infrastructure `drift`
- Visualise the `cost` of your infrastructure in the diagram (via Infracost)
- Do all of the above in complete automation `in CI pipelines`

&nbsp;

![flow-illustration](https://user-images.githubusercontent.com/25454503/157021111-816c9936-3232-455f-9709-c3a65f5f8dfe.svg)

The **Pluralith CLI** is a tool written in **`Go`** to:

- Integrate Pluralith with Terraform
- Interact with the **[Pluralith UI](https://www.pluralith.com)**
- Run Pluralith in **Pipelines** for automated infrastructure documentation
- Ship other useful little features

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

All three of the above commands share the same flags:

- `pluralith plan | apply | destroy`:
  - `--var`: Specify a variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')
  - `--var-file`: Specify a path to a var file to pass to Terraform. Can be specified multiple times.
  - `--cost-usage-file`: Specify a path to an infracost usage file to be used for the cost breakdown.
  - `--no-costs`: Turn off internal `infracost breakdown` on current run

&nbsp;

### → Diagram Export Command

- `pluralith graph`: Creates a Terraform execution plan, draws a graph and exports it as a PDF in a headless way<sup>2</sup>
  - `--title`: The title to be shown in the diagram *(e.g. "Pluralith EKS Cluster")*
  - `--author`: The author to be shown in the diagram *(e.g. "Tony Stark")*
  - `--version`: The version to be shown in the diagram *(e.g. "1.0.5")*
  - `--show-changes`: Enables change highlighting in the output diagram. When enabled, resources that have been added, updated, deleted etc. will be highlighted with special colors
  - `--show-drift`: Enables drift highlighting in the output diagram. When enabled, resources that Terraform has detected drift for will be highlighted with a special badge and color
  - `--skip-plan`: Skip the generation of a new execution plan
    - Only works if Pluralith has run in the current directory before
  - `--out-dir`: The path your exported diagram PDF gets saved to *(e.g. "~/pluralith-infra/eks")*
    - Saved to current directory by default
  - `--file-name`: The path your exported diagram PDF gets saved to
    - The value passed for `--title` is used by default
  - `--generate-md`: Generates markdown for GitHub pull request / commit comment *(used in our [Pluralith GitHub actions](https://github.com/Pluralith/actions))*
  - `--var`: Specify a variable to pass to Terraform. Can be specified multiple times. (Format: --var='NAME=VALUE')
  - `--var-file`: Specify a path to a var file to pass to Terraform. Can be specified multiple times.
  - `--cost-usage-file`: Specify a path to an infracost usage file to be used for the cost breakdown.
  - `--no-costs`: Turn off internal `infracost breakdown` on current run

&nbsp;

### 📍 Here's an example output for one of our test projects. View the PDF version **[here](https://github.com/Pluralith/pluralith-cli/files/8197192/HighlyAvailableIaaS.pdf)**

![HighlyAvailableIaaS](https://user-images.githubusercontent.com/25454503/157020490-8dadf7a2-ccb6-4323-a5d1-596d264bb06e.png)

&nbsp;

### → Init Command

This command initializes Pluralith in the current directory and creates a `pluralith.yml` config file for you to configure Pluralith to your liking.

- `pluralith init`: Initializes Pluralith in the current directory
  - `--empty`: Creates a config file with all the options commented out for you to edit to your liking
  - `--api-key`: If not authenticated through `pluralith login` yet, you can pass your API directly to `pluralith init` with this flag *(for CI contexts)*
  - `--project-id`: Links runs in the current directory to a specific project and posts them to the Pluralith dashboard to share with collaborators *(for CI contexts)*

&nbsp;

### → Strip Command

- `pluralith strip`: Strips and hashes your plan state to make it shareable with us for debugging
  - Takes an existing **Pluralith Plan state** and subjects it to rigorous hashing of values
    - The Pluralith Plan state is located in the file *pluralith.state.json* in your project directory
  - The purpose of this command is to strip the state of all sensitive data while keeping the structure intact, making it shareable
  - This is meant for us to debug edge cases on user state without the security hazard

&nbsp;

### → Module Commands

- `pluralith install`: Installs/updates the specific module whose name is passed (e.g. `pluralith install graph-module`)
- `pluralith update`: Essentially the same as `install`. Updates existing modules, if not installed it downloads the latest release
  - If no value is passed, the latest version of the CLI itself will be installed

&nbsp;

### → Utility Commands

- `pluralith login`: Authenticate with your API key (necessary for the CLI to work without the UI)<sup>2</sup>
- `pluralith version`: Shows information about the current CLI version as well as additional, installed modules

&nbsp;

<sup>1</sup> The UI then shows a prompt that lets you confirm or deny an `apply` with hotkeys.  
<sup>2</sup> You need to be authenticated with your **API key** via `pluralith login`. Currently only available for closed alpha testers. Interested? Shoot us an email dan@pluralith.com

&nbsp;

## 💰 Infracost Integration

**`Coming Soon`**

Pluralith shows you cost information directly in your infrastructure diagram. We automatically detect any Infracost installation and run `infracost breakdown` under the hood on every plan. We then match the cost data with the resources in the diagram and show you how much you'll pay. You can select from an array of cost modes:

- `Total Mode`: Shows the total cost for each individual resource
- `Diff Mode`: Shows the difference in costs for each individual resource in your current plan
- `Spike Mode`: Highlights only resources whose costs have increased in your current plan

The `plan`, `apply`, `destroy` and `graph` integrate Infracost:

- `--cost-usage-file`: Lets you pass a usage file in YML format as generated by infracost
- `--no-costs`: Lets you skip the infracost step on the current run

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

*Disclaimer: To properly use this CLI you **will need** the **Pluralith UI** and/or an **API key**. [Sign up](https://www.pluralith.com) for the private alpha to get access!*

![Subreddit subscribers](https://img.shields.io/reddit/subreddit-subscribers/pluralith?style=social)
