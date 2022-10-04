package cost

import "time"

type CostMap struct {
	Version  string `json:"version"`
	Currency string `json:"currency"`
	Projects []struct {
		Name     string `json:"name"`
		Metadata struct {
			Path              string    `json:"path"`
			InfracostCommand  string    `json:"infracostCommand"`
			Type              string    `json:"type"`
			Branch            string    `json:"branch"`
			Commit            string    `json:"commit"`
			CommitAuthorName  string    `json:"commitAuthorName"`
			CommitAuthorEmail string    `json:"commitAuthorEmail"`
			CommitTimestamp   time.Time `json:"commitTimestamp"`
			CommitMessage     string    `json:"commitMessage"`
			VcsRepoURL        string    `json:"vcsRepoUrl"`
			VcsSubPath        string    `json:"vcsSubPath"`
		} `json:"metadata"`
		PastBreakdown struct {
			Resources        []interface{} `json:"resources"`
			TotalHourlyCost  string        `json:"totalHourlyCost"`
			TotalMonthlyCost string        `json:"totalMonthlyCost"`
		} `json:"pastBreakdown"`
		Breakdown struct {
			Resources []struct {
				Name string `json:"name"`
				Tags struct {
					Environment string `json:"Environment"`
					Name        string `json:"Name"`
				} `json:"tags"`
				Metadata struct {
				} `json:"metadata"`
				HourlyCost   interface{} `json:"hourlyCost"`
				MonthlyCost  interface{} `json:"monthlyCost"`
				Subresources []struct {
					Name     string `json:"name"`
					Metadata struct {
					} `json:"metadata"`
					HourlyCost     interface{} `json:"hourlyCost"`
					MonthlyCost    interface{} `json:"monthlyCost"`
					CostComponents []struct {
						Name            string      `json:"name"`
						Unit            string      `json:"unit"`
						HourlyQuantity  interface{} `json:"hourlyQuantity"`
						MonthlyQuantity interface{} `json:"monthlyQuantity"`
						Price           string      `json:"price"`
						HourlyCost      interface{} `json:"hourlyCost"`
						MonthlyCost     interface{} `json:"monthlyCost"`
					} `json:"costComponents"`
				} `json:"subresources"`
			} `json:"resources"`
			TotalHourlyCost  string `json:"totalHourlyCost"`
			TotalMonthlyCost string `json:"totalMonthlyCost"`
		} `json:"breakdown"`
		Diff struct {
			Resources []struct {
				Name string `json:"name"`
				Tags struct {
					Environment string `json:"Environment"`
					Name        string `json:"Name"`
				} `json:"tags"`
				Metadata struct {
				} `json:"metadata"`
				HourlyCost   string `json:"hourlyCost"`
				MonthlyCost  string `json:"monthlyCost"`
				Subresources []struct {
					Name     string `json:"name"`
					Metadata struct {
					} `json:"metadata"`
					HourlyCost     string `json:"hourlyCost"`
					MonthlyCost    string `json:"monthlyCost"`
					CostComponents []struct {
						Name            string `json:"name"`
						Unit            string `json:"unit"`
						HourlyQuantity  string `json:"hourlyQuantity"`
						MonthlyQuantity string `json:"monthlyQuantity"`
						Price           string `json:"price"`
						HourlyCost      string `json:"hourlyCost"`
						MonthlyCost     string `json:"monthlyCost"`
					} `json:"costComponents"`
				} `json:"subresources"`
			} `json:"resources"`
			TotalHourlyCost  string `json:"totalHourlyCost"`
			TotalMonthlyCost string `json:"totalMonthlyCost"`
		} `json:"diff"`
		Summary struct {
			TotalDetectedResources    int `json:"totalDetectedResources"`
			TotalSupportedResources   int `json:"totalSupportedResources"`
			TotalUnsupportedResources int `json:"totalUnsupportedResources"`
			TotalUsageBasedResources  int `json:"totalUsageBasedResources"`
			TotalNoPriceResources     int `json:"totalNoPriceResources"`
			UnsupportedResourceCounts struct {
			} `json:"unsupportedResourceCounts"`
			NoPriceResourceCounts struct {
				AwsInternetGateway       int `json:"aws_internet_gateway"`
				AwsRouteTable            int `json:"aws_route_table"`
				AwsRouteTableAssociation int `json:"aws_route_table_association"`
				AwsSubnet                int `json:"aws_subnet"`
				AwsVpc                   int `json:"aws_vpc"`
			} `json:"noPriceResourceCounts"`
		} `json:"summary"`
	} `json:"projects"`
	TotalHourlyCost      string    `json:"totalHourlyCost"`
	TotalMonthlyCost     string    `json:"totalMonthlyCost"`
	PastTotalHourlyCost  string    `json:"pastTotalHourlyCost"`
	PastTotalMonthlyCost string    `json:"pastTotalMonthlyCost"`
	DiffTotalHourlyCost  string    `json:"diffTotalHourlyCost"`
	DiffTotalMonthlyCost string    `json:"diffTotalMonthlyCost"`
	TimeGenerated        time.Time `json:"timeGenerated"`
	Summary              struct {
		TotalDetectedResources    int `json:"totalDetectedResources"`
		TotalSupportedResources   int `json:"totalSupportedResources"`
		TotalUnsupportedResources int `json:"totalUnsupportedResources"`
		TotalUsageBasedResources  int `json:"totalUsageBasedResources"`
		TotalNoPriceResources     int `json:"totalNoPriceResources"`
		UnsupportedResourceCounts struct {
		} `json:"unsupportedResourceCounts"`
		NoPriceResourceCounts struct {
			AwsInternetGateway       int `json:"aws_internet_gateway"`
			AwsRouteTable            int `json:"aws_route_table"`
			AwsRouteTableAssociation int `json:"aws_route_table_association"`
			AwsSubnet                int `json:"aws_subnet"`
			AwsVpc                   int `json:"aws_vpc"`
		} `json:"noPriceResourceCounts"`
	} `json:"summary"`
}
