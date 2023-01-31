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
# Learn more about it at https://docs.pluralith.com/docs/more/config

org_id: $PLR_ORG_ID
project_id: $PLR_PROJECT_ID
project_name: $PLR_PROJECT_NAME
pluralith_api_endpoint: $PLR_API_ENDPOINT

# config:
#   title: null
#   version: null
#   sync_to_backend: false
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
# Learn more about it at https://docs.pluralith.com/docs/more/config

# org_id: null
# project_id: null
# pluralith_api_endpoint: null

# config:
#   title: null
#   version: null
#   sync_to_backend: false
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
`
