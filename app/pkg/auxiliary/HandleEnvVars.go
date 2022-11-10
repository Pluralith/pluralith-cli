package auxiliary

import "os"

func (S *State) GetEnvVars() error {
	// functionName := "GetEnvVars"

	if S.APIKey == "" {
		S.APIKey, _ = os.LookupEnv("PLURALITH_API_KEY")
	}
	if S.PluralithConfig.OrgId == "" {
		S.PluralithConfig.OrgId, _ = os.LookupEnv("PLURALITH_ORG_ID")
	}
	if S.PluralithConfig.ProjectId == "" {
		S.PluralithConfig.ProjectId, _ = os.LookupEnv("PLURALITH_PROJECT_ID")
	}
	if S.PluralithConfig.ProjectName == "" {
		S.PluralithConfig.ProjectName, _ = os.LookupEnv("PLURALITH_PROJECT_NAME")
	}

	return nil
}
