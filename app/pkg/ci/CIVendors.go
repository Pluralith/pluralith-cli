package ci

type CIVendorStruct struct {
	Name     string      //`json:Name`
	Constant string      //`json:Constant`
	Env      interface{} //`json:Env`
	Pr       interface{} //`json:"pr,omitempty"`
}

var CIVendors = []CIVendorStruct{
	{
		Name:     "AppVeyor",
		Constant: "APPVEYOR",
		Env:      "APPVEYOR",
		Pr:       "APPVEYOR_PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Azure Pipelines",
		Constant: "AZURE_PIPELINES",
		Env:      "SYSTEM_TEAMFOUNDATIONCOLLECTIONURI",
		Pr:       "SYSTEM_PULLREQUEST_PULLREQUESTID",
	},
	{
		Name:     "Appcircle",
		Constant: "APPCIRCLE",
		Env:      "AC_APPCIRCLE",
	},
	{
		Name:     "Bamboo",
		Constant: "BAMBOO",
		Env:      "bamboo_planKey",
	},
	{
		Name:     "Bitbucket Pipelines",
		Constant: "BITBUCKET",
		Env:      "BITBUCKET_COMMIT",
		Pr:       "BITBUCKET_PR_ID",
	},
	{
		Name:     "Bitrise",
		Constant: "BITRISE",
		Env:      "BITRISE_IO",
		Pr:       "BITRISE_PULL_REQUEST",
	},
	{
		Name:     "Buddy",
		Constant: "BUDDY",
		Env:      "BUDDY_WORKSPACE_ID",
		Pr:       "BUDDY_EXECUTION_PULL_REQUEST_ID",
	},
	{
		Name:     "Buildkite",
		Constant: "BUILDKITE",
		Env:      "BUILDKITE",
		Pr: map[string]interface{}{
			"env": "BUILDKITE_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "CircleCI",
		Constant: "CIRCLE",
		Env:      "CIRCLECI",
		Pr:       "CIRCLE_PULL_REQUEST",
	},
	{
		Name:     "Cirrus CI",
		Constant: "CIRRUS",
		Env:      "CIRRUS_CI",
		Pr:       "CIRRUS_PR",
	},
	{
		Name:     "AWS CodeBuild",
		Constant: "CODEBUILD",
		Env:      "CODEBUILD_BUILD_ARN",
	},
	{
		Name:     "Codefresh",
		Constant: "CODEFRESH",
		Env:      "CF_BUILD_ID",
		Pr: map[string]interface{}{
			"any": []interface{}{
				"CF_PULL_REQUEST_NUMBER",
				"CF_PULL_REQUEST_ID",
			},
		},
	},
	{
		Name:     "Codeship",
		Constant: "CODESHIP",
		Env: map[string]interface{}{
			"CI_NAME": "codeship",
		},
	},
	{
		Name:     "Drone",
		Constant: "DRONE",
		Env:      "DRONE",
		Pr: map[string]interface{}{
			"DRONE_BUILD_EVENT": "pull_request",
		},
	},
	{
		Name:     "dsari",
		Constant: "DSARI",
		Env:      "DSARI",
	},
	{
		Name:     "Expo Application Services",
		Constant: "EAS",
		Env:      "EAS_BUILD",
	},
	{
		Name:     "GitHub Actions",
		Constant: "GITHUB_ACTIONS",
		Env:      "GITHUB_ACTIONS",
		Pr: map[string]interface{}{
			"GITHUB_EVENT_NAME": "pull_request",
		},
	},
	{
		Name:     "GitLab CI",
		Constant: "GITLAB",
		Env:      "GITLAB_CI",
		Pr:       "CI_MERGE_REQUEST_ID",
	},
	{
		Name:     "GoCD",
		Constant: "GOCD",
		Env:      "GO_PIPELINE_LABEL",
	},
	{
		Name:     "LayerCI",
		Constant: "LAYERCI",
		Env:      "LAYERCI",
		Pr:       "LAYERCI_PULL_REQUEST",
	},
	{
		Name:     "Hudson",
		Constant: "HUDSON",
		Env:      "HUDSON_URL",
	},
	{
		Name:     "Jenkins",
		Constant: "JENKINS",
		Env: []interface{}{
			"JENKINS_URL",
			"BUILD_ID",
		},
		Pr: map[string]interface{}{
			"any": []interface{}{
				"ghprbPullId",
				"CHANGE_ID",
			},
		},
	},
	{
		Name:     "Magnum CI",
		Constant: "MAGNUM",
		Env:      "MAGNUM",
	},
	{
		Name:     "Netlify CI",
		Constant: "NETLIFY",
		Env:      "NETLIFY",
		Pr: map[string]interface{}{
			"env": "PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Nevercode",
		Constant: "NEVERCODE",
		Env:      "NEVERCODE",
		Pr: map[string]interface{}{
			"env": "NEVERCODE_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Render",
		Constant: "RENDER",
		Env:      "RENDER",
		Pr: map[string]interface{}{
			"IS_PULL_REQUEST": "true",
		},
	},
	{
		Name:     "Sail CI",
		Constant: "SAIL",
		Env:      "SAILCI",
		Pr:       "SAIL_PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Semaphore",
		Constant: "SEMAPHORE",
		Env:      "SEMAPHORE",
		Pr:       "PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Screwdriver",
		Constant: "SCREWDRIVER",
		Env:      "SCREWDRIVER",
		Pr: map[string]interface{}{
			"env": "SD_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Shippable",
		Constant: "SHIPPABLE",
		Env:      "SHIPPABLE",
		Pr: map[string]interface{}{
			"IS_PULL_REQUEST": "true",
		},
	},
	{
		Name:     "Solano CI",
		Constant: "SOLANO",
		Env:      "TDDIUM",
		Pr:       "TDDIUM_PR_ID",
	},
	{
		Name:     "Strider CD",
		Constant: "STRIDER",
		Env:      "STRIDER",
	},
	{
		Name:     "TaskCluster",
		Constant: "TASKCLUSTER",
		Env: []interface{}{
			"TASK_ID",
			"RUN_ID",
		},
	},
	{
		Name:     "TeamCity",
		Constant: "TEAMCITY",
		Env:      "TEAMCITY_VERSION",
	},
	{
		Name:     "Travis CI",
		Constant: "TRAVIS",
		Env:      "TRAVIS",
		Pr: map[string]interface{}{
			"env": "TRAVIS_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Vercel",
		Constant: "VERCEL",
		Env:      "NOW_BUILDER",
	},
	{
		Name:     "Visual Studio App Center",
		Constant: "APPCENTER",
		Env:      "APPCENTER_BUILD_ID",
	},
}
