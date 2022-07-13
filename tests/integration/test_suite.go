package integration

import "fmt"

const (
	localHost = "0.0.0.0"
	testPort  = 8080

	pathProducts          = "/api/v1/products"
	pathSubscriptionPlans = "/api/v1/subscription-plans"

	voucherCode10Off       = "10-off"
	voucherCode20Off       = "20-off"
	voucherCodeNonExistent = "i_dont_exist"

	userIDWithoutExistingPlan1 = "qo6_0keqGKqDA9EB5obql"
	userIDWithoutExistingPlan2 = "MFE0MTZqGKqDA9EB5obql"

	productIDBasic30Days   = "Ft2GgLRgN3FbMveklTy-W"
	productIDBasic1Year    = "mApy9b9Fqpt_WjghgUkSY"
	productIDPremium30Days = "7Yn_IvvYsfkeo7-ysixd7"
	productIDPremium1Year  = "UMp3k41eV5mY_iOkiElGm"
)

var (
	baseUrlProducts          = fmt.Sprintf("%s%s", getBaseUrl(), pathProducts)
	baseUrlSubscriptionPlans = fmt.Sprintf("%s%s", getBaseUrl(), pathSubscriptionPlans)
)

func getBaseUrl() string {
	return fmt.Sprintf("http://%s:%d", localHost, testPort)
}
