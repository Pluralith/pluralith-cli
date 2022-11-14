package auxiliary

type Vendor struct {
	Name       string      //`json:Name`
	Constant   string      //`json:Constant`
	Env        []string    //`json:Env`
	Pr         interface{} //`json:"pr,omitempty"`
	Branch     string      //`json:"branch,omitempty"`
	PipelineId string      //`json:"branch,omitempty"`
}

var GeneralEnvVars = []string{
	"CI", // Travis CI, CircleCI, Cirrus CI, Gitlab CI, Appveyor, CodeShip, dsari
	"CI_NAME",
	"CONTINUOUS_INTEGRATION", // Travis CI, Cirrus CI
	"BUILD_NUMBER",           // Jenkins, TeamCity
	"RUN_ID",
	"PLURALITH_CI", // Pluralith CI Docker Image
}

var CIVendors = []Vendor{
	{
		Name:     "GitHub Actions",
		Constant: "GITHUB_ACTIONS",
		Env:      []string{"GITHUB_ACTIONS"},
		Pr: map[string]interface{}{
			"GITHUB_EVENT_NAME": "pull_request",
		},
		Branch:     "GITHUB_HEAD_REF",
		PipelineId: "GITHUB_RUN_ID",
	},
	{
		Name:       "GitLab CI",
		Constant:   "GITLAB",
		Env:        []string{"GITLAB_CI"},
		Pr:         "CI_MERGE_REQUEST_ID",
		Branch:     "CI_MERGE_REQUEST_SOURCE_BRANCH_NAME",
		PipelineId: "CI_PIPELINE_ID",
	},
	{
		Name:     "AppVeyor",
		Constant: "APPVEYOR",
		Env:      []string{"APPVEYOR"},
		Pr:       "APPVEYOR_PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Azure Pipelines",
		Constant: "AZURE_PIPELINES",
		Env:      []string{"SYSTEM_TEAMFOUNDATIONCOLLECTIONURI"},
		Pr:       "SYSTEM_PULLREQUEST_PULLREQUESTID",
	},
	{
		Name:     "Appcircle",
		Constant: "APPCIRCLE",
		Env:      []string{"AC_APPCIRCLE"},
	},
	{
		Name:     "Bamboo",
		Constant: "BAMBOO",
		Env:      []string{"bamboo_planKey"},
	},
	{
		Name:     "Bitbucket Pipelines",
		Constant: "BITBUCKET",
		Env:      []string{"BITBUCKET_COMMIT"},
		Pr:       "BITBUCKET_PR_ID",
	},
	{
		Name:     "Bitrise",
		Constant: "BITRISE",
		Env:      []string{"BITRISE_IO"},
		Pr:       "BITRISE_PULL_REQUEST",
	},
	{
		Name:     "Buddy",
		Constant: "BUDDY",
		Env:      []string{"BUDDY_WORKSPACE_ID"},
		Pr:       "BUDDY_EXECUTION_PULL_REQUEST_ID",
	},
	{
		Name:     "Buildkite",
		Constant: "BUILDKITE",
		Env:      []string{"BUILDKITE"},
		Pr: map[string]interface{}{
			"env": "BUILDKITE_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "CircleCI",
		Constant: "CIRCLE",
		Env:      []string{"CIRCLECI"},
		Pr:       "CIRCLE_PULL_REQUEST",
	},
	{
		Name:     "Cirrus CI",
		Constant: "CIRRUS",
		Env:      []string{"CIRRUS_CI"},
		Pr:       "CIRRUS_PR",
	},
	{
		Name:     "AWS CodeBuild",
		Constant: "CODEBUILD",
		Env:      []string{"CODEBUILD_BUILD_ARN"},
	},
	{
		Name:     "Codefresh",
		Constant: "CODEFRESH",
		Env:      []string{"CF_BUILD_ID"},
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
		Env:      []string{"codeship"},
	},
	{
		Name:     "Drone",
		Constant: "DRONE",
		Env:      []string{"DRONE"},
		Pr: map[string]interface{}{
			"DRONE_BUILD_EVENT": "pull_request",
		},
	},
	{
		Name:     "dsari",
		Constant: "DSARI",
		Env:      []string{"DSARI"},
	},
	{
		Name:     "Expo Application Services",
		Constant: "EAS",
		Env:      []string{"EAS_BUILD"},
	},
	{
		Name:     "GoCD",
		Constant: "GOCD",
		Env:      []string{"GO_PIPELINE_LABEL"},
	},
	{
		Name:     "LayerCI",
		Constant: "LAYERCI",
		Env:      []string{"LAYERCI"},
		Pr:       "LAYERCI_PULL_REQUEST",
	},
	{
		Name:     "Hudson",
		Constant: "HUDSON",
		Env:      []string{"HUDSON_URL"},
	},
	{
		Name:     "Jenkins",
		Constant: "JENKINS",
		Env: []string{
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
		Env:      []string{"MAGNUM"},
	},
	{
		Name:     "Netlify CI",
		Constant: "NETLIFY",
		Env:      []string{"NETLIFY"},
		Pr: map[string]interface{}{
			"env": "PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Nevercode",
		Constant: "NEVERCODE",
		Env:      []string{"NEVERCODE"},
		Pr: map[string]interface{}{
			"env": "NEVERCODE_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Render",
		Constant: "RENDER",
		Env:      []string{"RENDER"},
		Pr: map[string]interface{}{
			"IS_PULL_REQUEST": "true",
		},
	},
	{
		Name:     "Sail CI",
		Constant: "SAIL",
		Env:      []string{"SAILCI"},
		Pr:       "SAIL_PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Semaphore",
		Constant: "SEMAPHORE",
		Env:      []string{"SEMAPHORE"},
		Pr:       "PULL_REQUEST_NUMBER",
	},
	{
		Name:     "Screwdriver",
		Constant: "SCREWDRIVER",
		Env:      []string{"SCREWDRIVER"},
		Pr: map[string]interface{}{
			"env": "SD_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Shippable",
		Constant: "SHIPPABLE",
		Env:      []string{"SHIPPABLE"},
		Pr: map[string]interface{}{
			"IS_PULL_REQUEST": "true",
		},
	},
	{
		Name:     "Solano CI",
		Constant: "SOLANO",
		Env:      []string{"TDDIUM"},
		Pr:       "TDDIUM_PR_ID",
	},
	{
		Name:     "Strider CD",
		Constant: "STRIDER",
		Env:      []string{"STRIDER"},
	},
	{
		Name:     "TaskCluster",
		Constant: "TASKCLUSTER",
		Env: []string{
			"TASK_ID",
			"RUN_ID",
		},
	},
	{
		Name:     "TeamCity",
		Constant: "TEAMCITY",
		Env:      []string{"TEAMCITY_VERSION"},
	},
	{
		Name:     "Travis CI",
		Constant: "TRAVIS",
		Env:      []string{"TRAVIS"},
		Pr: map[string]interface{}{
			"env": "TRAVIS_PULL_REQUEST",
			"ne":  "false",
		},
	},
	{
		Name:     "Vercel",
		Constant: "VERCEL",
		Env:      []string{"NOW_BUILDER"},
	},
	{
		Name:     "Visual Studio App Center",
		Constant: "APPCENTER",
		Env:      []string{"APPCENTER_BUILD_ID"},
	},
}
