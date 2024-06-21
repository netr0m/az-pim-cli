# Contributing to az-pim-cli

Welcome! We are glad that you want to contribute to our project! 💖

There are many ways to contribute:

* Suggest [features](https://github.com/netr0m/az-pim-cli/issues/new?labels=enhancement&template=feature_request.yaml)
* Report [bugs](https://github.com/netr0m/az-pim-cli/issues/new?labels=bug&template=bug_report.yaml)
* Develop [features](https://github.com/netr0m/az-pim-cli/issues?q=label%3Aenhancement)
* Improve the documentation


## Table of Contents

- [I Have a Question](#i-have-a-question)
- [I Want To Contribute](#i-want-to-contribute)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Enhancements](#suggesting-enhancements)
- [Your First Code Contribution](#your-first-code-contribution)
- [Styleguides](#styleguides)
- [Commit Messages](#commit-messages)


## I Have a Question

> If you want to ask a question, we assume that you have read the available [README](https://github.com/netr0m/az-pim-cli/blob/main/README.md).

Before you ask a question, it is best to search for existing [Issues](https://github.com/netr0m/az-pim-cli/issues) that might help you. In case you have found a suitable issue and still need clarification, you can write your question in this issue. It is also advisable to search the internet for answers first.

If you then still feel the need to ask a question and need clarification, we recommend the following:
- Create a post in the [Discussions](https://github.com/netr0m/az-pim-cli/discussions)
- Open an [Issue](https://github.com/netr0m/az-pim-cli/issues/new).
- Provide as much context as you can about what you're running into.
- Provide project and platform versions (OS, golang version, etc), depending on what seems relevant.


## I Want To Contribute

### Reporting Bugs

#### Before Submitting a Bug Report

A good bug report shouldn't leave others needing to chase you up for more information. Therefore, we ask you to investigate carefully, collect information and describe the issue in detail in your report. Please complete the following steps in advance to help us fix any potential bug as fast as possible.

- Make sure that you are using the latest version.
- Determine if your bug is really a bug and not an error on your side e.g. using incompatible environment components/versions (Make sure that you have read the [README](https://github.com/netr0m/az-pim-cli/blob/main/README.md). If you are looking for support, you might want to check [this section](#i-have-a-question)).
- To see if other users have experienced (and potentially already solved) the same issue you are having, check if there is not already a bug report existing for your bug or error in the [issue tracker](https://github.com/netr0m/az-pim-cli/issues?q=label%3Abug).

<!-- omit in toc -->
#### How Do I Submit a Good Bug Report?

> You must never report security related issues, vulnerabilities or bugs including sensitive information to the issue tracker, or elsewhere in public. Instead sensitive bugs must be sent by email to <>.

We use GitHub issues to track bugs and errors. If you run into an issue with the project:
- Open an [Issue with the 'Bug Report' template](https://github.com/netr0m/az-pim-cli/issues/new?labels=bug&template=bug_report.yaml)

Once it's filed, a maintainer will try to investigate the issue. Please make sure that you respond within reasonable time if there are any follow-up questions.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for az-pim-cli, **including completely new features and minor improvements to existing functionality**. Following these guidelines will help maintainers and the community to understand your suggestion and find related suggestions.

<!-- omit in toc -->
#### Before Submitting an Enhancement

- Make sure that you are using the latest version.
- Read the [README](https://github.com/netr0m/az-pim-cli/blob/main/README.md) carefully and find out if the functionality is already covered.
- Perform a [search](https://github.com/netr0m/az-pim-cli/issues) to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one.

<!-- omit in toc -->
#### How Do I Submit a Good Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://github.com/netr0m/az-pim-cli/issues/new?labels=enhancement&template=feature_request.yaml).

- Use a **clear and descriptive title** for the issue to identify the suggestion.
- Make sure to fill out all of the required fields in the issue (following the template for feature requests).

### Your First Code Contribution

In general, we follow the ["fork-and-pull" Git workflow](https://github.com/susam/gitpr).

Here's a quick guide:

1. Create your own fork of the repository
2. Clone the forked project to your machine
3. To keep track of the original repository, add another remote named 'upstream'
```shell
git remote add upstream git@github.com:netr0m/az-pim-cli.git
```
4. Create a branch locally with a succinct but descriptive name and prefixed with change type.
```shell
git checkout -b feature/my-new-feature
```
5. Make the changes in the newly created branch.
6. Stage the changed files
```shell
git add path/to/filename
```
7. Commit your changes; but make sure to follow our [Styleguide for commit messages](#commit-messages)
```shell
git commit -m "<type>[optional scope]: <description>"
```
8. Before you create the pull request, be sure to rebase onto the upstream source. This ensures your code is running on the latest available code.
```shell
git fetch upstream
git rebase upstream/main
```
9. Push to your fork.
```shell
git push origin feature/my-new-feature
```
10. Submit a pull request to the original repository (via the GitHub web interface)
  - Make sure that you have an issue associated to the change you wish to introduce, and reference the issue in the pull request

That's it... thank you for your contribution!

#### Automated checks

The GitHub repository includes a few workflows that will run [on new pull requests](./.github/workflows/on-pull-request.yaml). Make sure to resolve any issues if any of the tests are failing.

#### Code review process

The maintainer(s) will look try to review pull requests on a regular basis. After feedback has been given, we expect responses within reasonable time. If there is no response from the contributor, the pull request may be closed.

#### pre-commit

When contributing to this project, pre-commits are necessary, as they run certain checks, sanitisers, and formatters.

The project provides a `.pre-commit-config.yaml` file that is used to setup git _pre-commit hooks_.

On commit locally, code is automatically formatted and checked for security vulnerabilities using pre-commit git hooks.

##### Installation

To initialize pre-commit in your local repository, run

```shell
pre-commit install
```

This tells pre-commit to run for this repository on every commit.

##### Usage

Pre-commit will run on every commit, but can also be run manually on all files:

```shell
pre-commit run --all-files
```


## Styleguides
### Commit Messages

This repository defines precise rules for how the git commit messages may be formatted.

We enforce the [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/) for commit messages, both through `pre-commit` hooks and through GitHub Actions workflows.

The commit messages should adhere to the following structure:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

where *type* must be one of `feat, fix, refactor, chore, test, ci, build, lint, docs`

For instance,

*Changes to the documentation*

```
docs: correct spelling of CHANGELOG
```
*A breaking change introducing a new feature*
```
feat!: list eligible roles
```

*Summary should:*
- Be written in imperative, present tense, e.g. write `add` instead of `added` or `adds`.
- Don't capitalize the first letter.
- Don't write dot (.) at the end.

<details>

```
git commit -m "fix(parser): fix broken config parser"

git commit -m "feat(parser): add config file parser"

git commit -m "feat(new-parser): change to config-parser-lib
BREAKING CHANGE: new library does not support foo-construct"
```

</details>

This results in more readable messages that are easy to follow when looking through the project history, and faciliates for automated changelog generation.
