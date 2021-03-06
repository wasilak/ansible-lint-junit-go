[ansible-lint](https://github.com/willthames/ansible-lint) to JUnit converter [![Build Status](https://travis-ci.org/wasilak/ansible-lint-junit-go.svg?branch=master)](https://travis-ci.org/wasilak/ansible-lint-junit-go) [![Maintainability](https://api.codeclimate.com/v1/badges/c009223d118fa63bc2c0/maintainability)](https://codeclimate.com/github/wasilak/ansible-lint-junit-go/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/c009223d118fa63bc2c0/test_coverage)](https://codeclimate.com/github/wasilak/ansible-lint-junit-go/test_coverage) [![Total alerts](https://img.shields.io/lgtm/alerts/g/wasilak/ansible-lint-junit-go.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/wasilak/ansible-lint-junit-go/alerts/) [![Language grade: Go](https://img.shields.io/lgtm/grade/go/g/wasilak/ansible-lint-junit-go.svg?logo=lgtm&logoWidth=18)](https://lgtm.com/projects/g/wasilak/ansible-lint-junit-go/context:go)
---

### Installation
- download precompiled binary from [releases](https://github.com/wasilak/ansible-lint-junit-go/releases) page
- make it executable with: `chmod +x ansible-lint-junit`

### Usage:
1. you can pipe output of `ansible-lint -p`:
    ```shell
    ansible-lint playbook.yml -p | ansible-lint-junit -output ansible-lint.xml
    ```
3. or run `ansible-lint` on your playbook(s) with parameter `-p` (it is required) and redirect output to file
    ```shell
    ansible-lint -p your_fancy_playbook.yml > ansible-lint.txt
    ```
    and run `ansible-lint-junit` and pass generated file to it
    ```shell
    ansible-lint-junit ansible-lint.txt -output ansible-lint.xml
    ```

### Output
* if there are any lint errors, full JUnit XML will be created
* if there are no errors, empty JUnit XML will be created, this is for i.e. [Bamboo](https://www.atlassian.com/software/bamboo) JUnit parser plugin compatibility.
It will break build if XML is missing or incorrect, and there is really no way of generating XML with *"PASSED"* tests in case of linter.
