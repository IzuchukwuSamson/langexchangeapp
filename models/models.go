package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FirstName   string             `json:"firstname" bson:"firstname"`
	LastName    string             `json:"lastname" bson:"lastname"`
	Email       string             `json:"email" bson:"email"`
	Password    string             `json:"password" bson:"password"`
	PhoneNumber string             `json:"phonenumber" bson:"phonenumber"`
	Role        string             `json:"role" bson:"role"`
	IsActive    int64              `json:"is_active" bson:"is_active"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type UserDTO struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	FirstName   string    `json:"firstname" bson:"firstname"`
	LastName    string    `json:"lastname" bson:"lastname"`
	Email       string    `json:"email" bson:"email"`
	PhoneNumber string    `json:"phonenumber" bson:"phonenumber"`
	Role        string    `json:"role" bson:"role"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type EmailVerification struct {
	ID     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID primitive.ObjectID `json:"userid" bson:"_userid,omitempty"`
	Email  string             `json:"email" bson:"email"`
	Code   string             `json:"code" bson:"code"`
}

type PasswordReset struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID    primitive.ObjectID `json:"userid" bson:"_userid,omitempty"`
	Code      string             `json:"code" bson:"code"`
	ExpiresAt time.Time          `json:"expires_at" bson:"expires_at"`
}

// type lexibuddy struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Network        string             `bson:"network" json:"network"`
// 	BuyerDiscount  float64            `bson:"network" json:"buyer_discount"`
// 	UserDiscount   float64            `bson:"user_discount" json:"user_discount"`
// 	AgentDiscount  float64            `bson:"agent_discount" json:"agent_discount"`
// 	VendorDiscount float64            `bson:"vendor_discount" json:"vendor_discount"`
// 	lexibuddyType    string             `bson:"lexibuddy_type" json:"lexibuddy_type"`
// }

// type lexibuddyPinPrice struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Network        string             `bson:"name" json:"network"`
// 	UserDiscount   float64            `bson:"user_discount" json:"user_discount"`
// 	AgentDiscount  float64            `bson:"agent_discount" json:"agent_discount"`
// 	VendorDiscount float64            `bson:"vendor_discount" json:"vendor_discount"`
// }

// type TopUpPrice struct {
// 	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	BuyingPrice  int                `bson:"buying_price" json:"buying_price"`
// 	SellingPrice int                `bson:"provider_id" json:"provider_id"`
// 	Agent        string             `bson:"agent" json:"agent"`
// 	Vendor       float64            `bson:"vendor" json:"vendor"`
// 	Description  string             `bson:"description" json:"description"`
// 	DatePosted   time.Time          `bson:"date_posted" json:"date_posted"`
// }

// type ApiConfigs struct {
// 	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name  string             `bson:"name" json:"name"`
// 	Value string             `bson:"value" json:"value"`
// }

// type ApiLinks struct {
// 	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name  string             `bson:"name" json:"name"`
// 	Value string             `bson:"value" json:"value"`
// 	Type  string             `bson:"type" json:"type"`
// }

// type Cable struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	CableID        primitive.ObjectID `bson:"_cableid,omitempty" json:"cableid,omitempty"`
// 	Provider       string             `bson:"provider" json:"provider"`
// 	ProviderStatus string             `bson:"provider_status" json:"provider_status"`
// }

// type CablePlans struct {
// 	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name          string             `bson:"name" json:"name"`
// 	Price         string             `bson:"provider" json:"provider"`
// 	UserPrice     string             `bson:"user_price" json:"user_price"`
// 	CablePrice    string             `bson:"cable_price" json:"cable_price"`
// 	AgentPrice    string             `bson:"agent_price" json:"agent_price"`
// 	VendorPrice   string             `bson:"vendor_price" json:"vendor_price"`
// 	PlanID        string             `bson:"plan_id" json:"plan_id"`
// 	Type          string             `bson:"type" json:"type"`
// 	CableProvider Cable              `bson:"cable_provider" json:"cable_provider"`
// 	Duration      string             `bson:"duration" json:"duration"`
// }

// type Contact struct {
// 	// for support messages
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	User       User               `bson:"user_id" json:"user_id"`
// 	Name       string             `bson:"name" json:"name"`
// 	Contact    string             `bson:"contact" json:"contact"`
// 	Subject    string             `bson:"subject" json:"subject"`
// 	Message    string             `bson:"message" json:"message"`
// 	DatePosted time.Time          `bson:"date_posted" json:"date_posted"`
// }

// type DataPins struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name        string             `bson:"name" json:"name"`
// 	Price       string             `bson:"provider" json:"provider"`
// 	UserPrice   string             `bson:"user_price" json:"user_price"`
// 	AgentPrice  string             `bson:"agent_price" json:"agent_price"`
// 	PlanID      string             `bson:"plan_id" json:"plan_id"`
// 	Type        string             `bson:"type" json:"type"`
// 	DataNetwork string             `bson:"data_network" json:"data_network"`
// 	Duration    string             `bson:"duration" json:"duration"`
// }

// type DataPlans struct {
// 	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Name        string             `bson:"name" json:"name"`
// 	Price       string             `bson:"provider" json:"provider"`
// 	UserPrice   string             `bson:"user_price" json:"user_price"`
// 	AgentPrice  string             `bson:"agent_price" json:"agent_price"`
// 	VendorPrice string             `bson:"vendor_price" json:"vendor_price"`
// 	PlanID      string             `bson:"plan_id" json:"plan_id"`
// 	Type        string             `bson:"type" json:"type"`
// 	DataNetwork string             `bson:"data_network" json:"data_network"`
// 	Duration    string             `bson:"duration" json:"duration"`
// }

// type DataTokens struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	Sid      primitive.ObjectID `bson:"_sid,omitempty" json:"sid,omitempty"`
// 	TRef     string             `bson:"tRef" json:"tRef"`
// 	Business string             `bson:"business" json:"business"`
// 	Network  string             `bson:"network" json:"network"`
// 	Datasize string             `bson:"datasize" json:"datasize"`
// 	Quantity string             `bson:"quantity" json:"quantity"`
// 	Serial   string             `bson:"serial" json:"serial"`
// 	Tokens   string             `bson:"tokens" json:"tokens"`
// 	Date     time.Time          `bson:"date" json:"date"`
// }

// type Electricity struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	ElectricityId  int                `bson:"electricity_id" json:"electricity_id"`
// 	Provider       string             `bson:"provider" json:"provider"`
// 	Abbreviation   string             `bson:"abbreviation" json:"abbreviation"`
// 	ProviderStatus string             `bson:"providerStatus" json:"providerStatus"`
// }

// type Exam struct {
// 	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	ExamId         string             `bson:"exam_id" json:"exam_id"`
// 	Provider       string             `bson:"provider" json:"provider"`
// 	Price          int                `bson:"price" json:"price"`
// 	BuyingPrice    int                `bson:"buying_price" json:"buying_price"`
// 	ProviderStatus string             `bson:"providerStatus" json:"providerStatus"`
// }

// type Network struct {
// 	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	NetworkId        string             `bson:"network_id" json:"network_id"`
// 	SmeId            string             `bson:"sme_id" json:"sme_id"`
// 	GiftingId        string             `bson:"gifting_id" json:"gifting_id"`
// 	CorporateId      int                `bson:"corporate_id" json:"corporate_id"`
// 	VtuId            string             `bson:"vtu_id" json:"vtu_id"`
// 	ShareSell        int                `bson:"sharesell_id" json:"sharesell_id"`
// 	Network          int                `bson:"network" json:"network"`
// 	NetworkStatus    string             `bson:"networkStatus" json:"networkStatus"`
// 	VtuStatus        string             `bson:"vtuStatus" json:"vtuStatus"`
// 	SharesellStatus  string             `bson:"sharesellStatus" json:"sharesellStatus"`
// 	lexibuddypinStatus string             `bson:"lexibuddypinStatus" json:"lexibuddypinStatus"`
// 	SmeStatus        string             `bson:"smeStatus" json:"smeStatus"`
// 	GiftingStatus    string             `bson:"giftingStatus" json:"giftingStatus"`
// 	CorporateStatus  string             `bson:"corporateStatus" json:"corporateStatus"`
// 	DatapinStatus    string             `bson:"datapinStatus" json:"datapinStatus"`
// }

// type Notifications struct {
// 	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	MsgFor     User               `bson:"msgfor" json:"msgfor"`
// 	Subject    string             `bson:"subject" json:"subject"`
// 	Message    string             `bson:"message" json:"message"`
// 	Status     int                `bson:"status" json:"status"`
// 	DatePosted time.Time          `bson:"date_posted" json:"date_posted"`
// }

// type SiteSettings struct {
// 	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	SiteName              string             `bson:"sitename" json:"sitename"`
// 	SiteUrl               string             `bson:"siteurl" json:"siteurl"`
// 	Agentupgrade          string             `bson:"agentupgrade" json:"agentupgrade"`
// 	Vendorupgrade         string             `bson:"vendorupgrade" json:"vendorupgrade"`
// 	Apidocumentation      string             `bson:"apidocumentation" json:"apidocumentation"`
// 	Phone                 string             `bson:"phone" json:"phone"`
// 	Email                 string             `bson:"email" json:"email"`
// 	Whatsapp              string             `bson:"whatsapp" json:"whatsapp"`
// 	Whatsappgroup         string             `bson:"whatsappgroup" json:"whatsappgroup"`
// 	Facebook              string             `bson:"facebook" json:"facebook"`
// 	Twitter               string             `bson:"twitter" json:"twitter"`
// 	Instagram             string             `bson:"instagram" json:"instagram"`
// 	Telegram              string             `bson:"telegram" json:"telegram"`
// 	Referalupgradebonus   float64            `bson:"referalupgradebonus" json:"referalupgradebonus"`
// 	Referallexibuddybonus   float64            `bson:"referallexibuddybonus" json:"referallexibuddybonus"`
// 	Referaldatabonus      float64            `bson:"referaldatabonus" json:"referaldatabonus"`
// 	Referalwalletbonus    float64            `bson:"referalwalletbonus" json:"referalwalletbonus"`
// 	Referalcablebonus     float64            `bson:"referalcablebonus" json:"referalcablebonus"`
// 	Referalexambonus      float64            `bson:"referalexambonus" json:"referalexambonus"`
// 	Referalmeterbonus     float64            `bson:"referalmeterbonus" json:"referalmeterbonus"`
// 	Wallettowalletcharges float64            `bson:"wallettowalletcharges" json:"wallettowalletcharges"`
// 	NotificationStatus    string             `bson:"notificationStatus" json:"notificationStatus"`
// 	Accountname           string             `bson:"accountname" json:"accountname"`
// 	Accountno             string             `bson:"accountno" json:"accountno"`
// 	BankName              string             `bson:"bank_name" json:"bank_name"`
// 	Electricitycharges    string             `bson:"electricitycharges" json:"electricitycharges"`
// 	Mirtimemin            string             `bson:"lexibuddymin" json:"lexibuddymin"`
// 	lexibuddymax            string             `bson:"lexibuddymax" json:"lexibuddymax"`
// }

// type Customers struct {
// 	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	ApiKey                string             `bson:"api_key" json:"api_key"` //unique
// 	FirstName             string             `bson:"firstname" json:"firstname"`
// 	LastName              string             `bson:"lastname" json:"lastname"`
// 	Email                 string             `bson:"email" json:"email"` //unique
// 	Password              string             `bson:"password" json:"password"`
// 	PhoneNumber           string             `bson:"phone_number" json:"phone_number"` //unique
// 	Country               string             `bson:"country" json:"country"`
// 	State                 string             `bson:"state" json:"state"`
// 	Pin                   int                `bson:"pin" json:"pin"`
// 	PinStatus             int                `bson:"pin_status" json:"pin_status"`
// 	CustomerType          int                `bson:"customer_type" json:"customer_type"` // 1 for vendor 2 = ageent, 3 = user
// 	CustomerWallet        float64            `bson:"customer_wallet" json:"customer_wallet"`
// 	CustomerReferalWallet float64            `bson:"customer_Referalwallet" json:"customer_Referalwallet"`
// 	BankNumber            string             `bson:"bank_number" json:"bank_number"`
// 	RolexBank             string             `bson:"rolex_bank" json:"rolex_bank"`
// 	FidelityBank          string             `bson:"fidelity_bank" json:"fidelity_bank"`
// 	SterlingBank          string             `bson:"sterling_bank" json:"sterling_bank"`
// 	BankName              string             `bson:"bank_name" json:"bank_name"`
// 	RegistrationStatus    int                `bson:"registration_status" json:"registration_status"`
// 	VerificationCode      int                `bson:"verification_code" json:"verification_code"`
// 	RegistrationDate      time.Time          `bson:"registration_date" json:"registration_date"`
// 	LastActive            time.Time          `bson:"last_active" json:"last_active"`
// 	Referals              string             `bson:"referals" json:"referals"`
// }

// type Transactions struct {
// 	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	UserID               primitive.ObjectID `bson:"_user_id,omitempty" json:"user_id,omitempty"`
// 	TransactionReference string             `bson:"transaction_reference" json:"transaction_reference"`
// 	ServiceName          string             `bson:"service_name" json:"service_name"`
// 	ServiceDescription   string             `bson:"service_desc" json:"service_desc"`
// 	Amount               string             `bson:"amount" json:"amount"`
// 	Status               int                `bson:"status" json:"status"`
// 	OldBalance           string             `bson:"old_balance" json:"old_balance"`
// 	NewBalance           string             `bson:"new_balance" json:"new_balance"`
// 	Profit               float64            `bson:"profit" json:"profit"`
// 	Date                 time.Time          `bson:"date" json:"date"`
// 	// TransactionReference should be unique
// }

// type UserLogin struct {
// 	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	User  User               `bson:"user" json:"user"`
// 	Token string             `bson:"token" json:"token"`
// }

// type UserVisits struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
// 	User      User               `bson:"user" json:"user"`
// 	Country   string             `bson:"country" json:"country"`
// 	State     string             `bson:"state" json:"state"`
// 	VisitTime time.Time          `bson:"visit_time" json:"visit_time"`
// }
