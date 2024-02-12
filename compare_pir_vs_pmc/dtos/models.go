package dtos

import "time"

type Collector struct {
	CollectorID     string          `json:"collector_id"`
	PaymentMethods  []PaymentMethod `json:"custom_payment_methods"`
	Exclusions      []interface{}   `json:"exclusions"`
	Groups          interface{}     `json:"groups"`
	AmountAllowed   []interface{}   `json:"amount_allowed"`
	OwnPromosByUser interface{}     `json:"own_promos_by_user"`
}

type OriginalDataKVS struct {
	IsCompressed bool   `json:"IsInStorage"`
	IsInStorage  bool   `json:"IsInStorage"`
	Data         string `json:"Data"`
	LastUpdated  string `json:"LastUpdated"`
}

type PaymentMethod struct {
	Key              Key
	ID               string            `json:"id"`
	BasicInfo        BasicInfo         `json:"basic_info"`
	PayerCosts       []Installment     `json:"payer_costs"`
	Issuer           Issuer            `json:"issuer"`
	Promos           []Promo           `json:"promos"`
	Misc             Misc              `json:"misc"`
	PayerRegulations *PayerRegulations `json:"payer_regulations,omitempty"`
	Rules            *Rule             `json:"rules,omitempty"`
	Assets           []Asset           `json:"assets,omitempty"`
}

type Key struct {
	Version           string
	SiteID            string `json:"site_id"`
	IssuerID          string `json:"issuer_id"`
	PaymentMethodID   string `json:"payment_method_id"`
	Channel           string `json:"channel"`
	Marketplace       string `json:"marketplace"`
	ProcessingMode    string `json:"processing_mode"`
	MerchantAccountID string `json:"merchant_account_id"`
}

type BasicInfo struct {
	Name           string `json:"name"`
	SiteID         string `json:"site_id"`
	PaymentTypeID  string `json:"payment_type_id"`
	Marketplace    string `json:"marketplace"`
	ProcessingMode string `json:"processing_mode"`
	Status         string `json:"status"`
	Ordering       int    `json:"ordering"`
}

type Installment struct {
	ID                       int64                       `json:"id"`
	Installments             int                         `json:"installments"`
	Rate                     float64                     `json:"rate"`
	MinAmount                float64                     `json:"min_amount"`
	MaxAmount                float64                     `json:"max_amount"`
	InstallmentFullTea       *float64                    `json:"installment_full_tea"`
	InstallmentFullCft       *float64                    `json:"installment_full_cft"`
	InstallmentReducedTea    *float64                    `json:"installment_reduced_tea"`
	InstallmentReducedCft    *float64                    `json:"installment_reduced_cft"`
	ReimbursementRate        *float64                    `json:"reimbursement_rate"`
	BaseInstallmentRate      *float64                    `json:"base_installment_rate"`
	DefaultInstallmentRate   *float64                    `json:"default_installment_rate"`
	DiscountRate             *float64                    `json:"discount_rate"`
	InstallmentRateCollector []string                    `json:"installment_rate_collector"`
	RealInstallments         int                         `json:"real_installments"`
	PaymentMethodOptionID    string                      `json:"payment_method_option_id,omitempty"`
	Labels                   []string                    `json:"labels"`
	ConsumerCredits          *ConsumerCreditsInstallment `json:"consumer_credits,omitempty"`
	InstallmentAmount        *float64                    `json:"installment_amount,omitempty"` // attribute for paypal
	TotalAmount              *float64                    `json:"total_amount,omitempty"`       // attribute for paypal
	SpecialPlan              string                      `json:"special_plan,omitempty"`
	ModoMango                *ModoMango                  `json:"modo_mango,omitempty"`
}

type ConsumerCreditsInstallment struct {
	Conditions *ConsumerCreditsInstallmentCondition `json:"conditions,omitempty"`
}

type ModoMango struct {
	Name                string `json:"name,omitempty"`
	InstallmentQuantity int    `json:"installment_quantity,omitempty"`
	Label               string `json:"label,omitempty"`
}

type ConsumerCreditsInstallmentCondition struct {
	Cat           *string `json:"cat,omitempty"`
	Cft           *string `json:"cft,omitempty"`
	Iva           *string `json:"iva,omitempty"`
	AdditionalIof *string `json:"additional_iof,omitempty"`
	Ceta          *string `json:"ceta,omitempty"`
	Cetm          *string `json:"cetm,omitempty"`
	Iof           *string `json:"iof,omitempty"`
	IofRate       *string `json:"iof_rate,omitempty"`
	Td            *string `json:"td,omitempty"`
	Tea           *string `json:"tea,omitempty"`
	Tem           *string `json:"tem,omitempty"`
	Tna           *string `json:"tna,omitempty"`
	Cftea         *string `json:"cftea,omitempty"`
}

type Issuer struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Default bool     `json:"default"`
	Tags    []string `json:"tags"`
}

type Promo struct {
	ID               string            `json:"id"`
	Name             string            `json:"name"` // TODO joel: esto no de donde saldria. La v1 no lo devuelves
	Legals           *string           `json:"legals"`
	DiscountRate     float64           `json:"discount_rate"` // TODO joel: entiendo que esto lo tengo que hablar con clau
	StartDate        *time.Time        `json:"start_date"`
	ExpirationDate   *time.Time        `json:"expiration_date"`
	ActiveWeekdays   []time.Weekday    `json:"active_weekdays"`
	BinPattern       string            `json:"bin_pattern"`
	PayerCosts       []Installment     `json:"payer_costs"` // TODO joel: entiendo que esto lo tengo que hablar con clau
	MerchantAccounts []MerchantAccount `json:"merchant_accounts"`
	Type             string            `json:"type"` // indica si en la v1 es un agreement o una promo propiamente dicho
	Status           string            `json:"status"`
	CategoryID       *string           `json:"category_id"`
}

type MerchantAccount struct {
	ID                    string  `json:"id"`
	BranchID              *string `json:"branch_id"`
	PaymentMethodOptionID string  `json:"payment_method_option_id"`
}

type Misc struct {
	Thumbnail             string                 `json:"thumbnail"`
	SecureThumbnail       string                 `json:"secure_thumbnail"`
	TotalFinancialCost    *float64               `json:"total_financial_cost"`
	MinAccreditationDays  int                    `json:"min_accreditation_days"`
	MaxAccreditationDays  int                    `json:"max_accreditation_days"`
	Labels                []string               `json:"labels"`
	Bins                  []int                  `json:"bins"`
	Owner                 string                 `json:"owner"`
	PmIssuerRelation      PMIssuerRelation       `json:"pm_issuer_relation"`
	AdditionalInfoNeeded  []string               `json:"additional_info_needed"`
	FinancialInstitutions []FinancialInstitution `json:"financial_institutions"`
	Settings              []Settings             `json:"settings"`
	FinancingDeal         FinancingDeal          `json:"financing_deals"`
	PayerCreditLine       PayerCreditLine        `json:"payer_credit_line"`
	DifferentialPricingID *int64                 `json:"differential_pricing_id"`
	Cccs                  []CCC                  `json:"cccs"`
	AccountMoneyCredit    *AccountMoneyCredit    `json:"account_money_credit"`
	Credit                *Credit                `json:"credit"`
}

type PMIssuerRelation struct {
	ID                string     `json:"id"`
	DeferredCapture   string     `json:"deferred_capture"`
	AccreditationTime *int       `json:"accreditation_time"`
	MerchantAccountID *string    `json:"merchant_account_id"`
	DateCreated       *time.Time `json:"date_created"`
}

type FinancialInstitution struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Settings struct {
	SecurityCode SecurityCode `json:"security_code"`
	CardNumber   CardNumber   `json:"card_number"`
	Bin          Bin          `json:"bin"`
	ID           int          `json:"id"`
}

type SecurityCode struct {
	Mode         string `json:"mode"`
	CardLocation string `json:"card_location"`
	Length       int    `json:"length"`
}

type CardNumber struct {
	Length     int    `json:"length"`
	Validation string `json:"validation"`
}

type Bin struct {
	ID                  int       `json:"id"`
	Bin                 int       `json:"bin"`
	SiteID              string    `json:"site_id"`
	PaymentMethodID     string    `json:"payment_method_id"`
	IssuerID            int       `json:"issuer_id"`
	Marketplace         string    `json:"marketplace"`
	Owner               string    `json:"owner"`
	MerchantAccountID   string    `json:"merchant_account_id"`
	FinanceDealID       int       `json:"finance_deal_id"`
	DiffPricingID       int       `json:"diff_pricing_id"`
	Status              string    `json:"status"`
	DateCreated         time.Time `json:"date_created"`
	DateLastUpdated     time.Time `json:"date_last_updated"`
	Pattern             *string   `json:"pattern"`
	InstallmentsPattern *string   `json:"installments_pattern"`
	ExclusionPattern    *string   `json:"exclusion_pattern"`
}

type FinancingDeal struct {
	Legals         *string    `json:"legals"`
	Installments   []int      `json:"installments"`
	ExpirationDate *time.Time `json:"expiration_date"`
	StartDate      *time.Time `json:"start_date"`
	Status         string     `json:"status"`
}

type PayerCreditLine struct {
	AvailableBalance *float64 `json:"available_balance"`
	NextDueDate      *string  `json:"next_due_date"`
}

type CCC struct {
	Code   string `json:"code"`
	Source string `json:"source"`
	Type   string `json:"type"`
}

type AccountMoneyCredit struct {
	BlockReason []string `json:"block_reason"`
	Amount      float64  `json:"amount"`
}

type Credit struct {
	AvailableLimit interface{} `json:"available_limit"`
	TotalLimit     int         `json:"total_limit"`
	CardStatus     string      `json:"card_status"`
	Account        Account     `json:"account"`
}

type Account struct {
	StatusDetail string `json:"status_detail"`
	Status       string `json:"status"`
}

type PayerRegulations struct {
	Name        string
	IsCompliant bool
	Regulations []Regulation `json:"regulations"`
}

type Regulation struct {
	Name             string    `json:"name"`
	Status           string    `json:"status"`
	Level            string    `json:"level"`
	EvaluationResult string    `json:"evaluation_result"`
	LastUpdated      time.Time `json:"last_updated"`
}

type Rule struct {
	Name map[string]Category `json:"partitions"`
}

type Category struct {
	Type []CategoryType `json:"categories"`
}

type CategoryType struct {
	Name string  `json:"name"`
	Mcc  []int64 `json:"mcc"`
}

type Asset struct {
	ID               int       `json:"id"`
	Type             string    `json:"type"`
	Code             string    `json:"code"`
	Status           string    `json:"status"`
	Partitions       []string  `json:"partitions,omitempty"`
	Synchronize      string    `json:"synchronize,omitempty"`
	MinAllowedAmount string    `json:"min_allowed_amount,omitempty"`
	MaxAllowedAmount string    `json:"max_allowed_amount,omitempty"`
	Segments         []Segment `json:"segments,omitempty"`
	AssetType        string    `json:"asset_type"`
}

type Segment struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Value string `json:"value"`
}
