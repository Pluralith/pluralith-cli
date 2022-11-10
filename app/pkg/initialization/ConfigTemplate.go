package initialization

var ConfigTemplate = `
#  _
# |_)|    _ _ |._|_|_
# |  ||_|| (_||| | | |
#
# Welcome to Pluralith!
# https://www.pluralith.com
#
# This is your Pluralith config file
# Learn more about it at https://docs.pluralith.com/config

org_id: $PLR_ORG_ID
project_id: $PLR_PROJECT_ID
project_name: $PLR_PROJECT_NAME
# config:
#   sensitive_attrs:
#     - "attribute_name"
#     - "attribute_name"
#   vars:
#     - "NAME=VALUE"
#     - "NAME=VALUE"
#   var_files:
#     - "./var_file.tfvars"
#     - "./var_file.tfvars"
#   cost_usage_file: "./usage_file.yml"

# export:
#   title: ""
#   author: ""
#   version: ""
`

var EmtpyConfig = `
#  _
# |_)|    _ _ |._|_|_
# |  ||_|| (_||| | | |
#
# Welcome to Pluralith!
# https://www.pluralith.com
#
# This is your Pluralith config file
# Learn more about it at https://docs.pluralith.com/config

# org_id: null
# project_id: null
# config:
#   sensitive_attrs:
#     - "attribute_name"
#     - "attribute_name"
#   vars:
#     - "NAME=VALUE"
#     - "NAME=VALUE"
#   var_files:
#     - "./var_file.tfvars"
#     - "./var_file.tfvars"
#   cost_usage_file: "./usage_file.yml"

# export:
#   title: ""
#   author: ""
#   version: ""
`
