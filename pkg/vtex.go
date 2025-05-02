package pkg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"web-server-pnb/config"
)

type GetTokenVtex struct {
	AuthStatus string `json:"authStatus"`
	Token      string `json:"token"`
	Expires    int    `json:"expires"`
}

type LojasContaFranquia []struct {
	Account      string `json:"account"`
	Active       bool   `json:"active"`
	AffiliateID  string `json:"affiliateId"`
	Cep          string `json:"cep"`
	DNS          string `json:"dns"`
	IDLoja       int    `json:"idLoja"`
	ID           string `json:"id"`
	AccountID    string `json:"accountId"`
	AccountName  string `json:"accountName"`
	DataEntityID string `json:"dataEntityId"`
}

type SacolasAbandonadas []int

type Receipt struct {
	Date    time.Time `json:"date"`
	OrderID string    `json:"orderId"`
	Receipt string    `json:"receipt"`
}

type InvoicedItem struct {
	Sku            int         `json:"sku"`
	Quantity       int         `json:"quantity"`
	Price          int         `json:"price"`
	Description    interface{} `json:"description"`
	UnitMultiplier float64     `json:"unitMultiplier"`
}
type Invoiced struct {
	Items           []InvoicedItem
	Courier         interface{} `json:"courier"`
	InvoiceNumber   string      `json:"invoiceNumber"`
	InvoiceValue    int         `json:"invoiceValue"`
	InvoiceURL      string      `json:"invoiceUrl"`
	IssuanceDate    time.Time   `json:"issuanceDate"`
	TrackingNumber  string      `json:"trackingNumber"`
	InvoiceKey      string      `json:"invoiceKey"`
	TrackingURL     interface{} `json:"trackingUrl"`
	EmbeddedInvoice string      `json:"embeddedInvoice"`
	NumberOrder     string      `json:"numberOrder"`
	Type            string      `json:"type"`
	CourierStatus   interface{} `json:"courierStatus"`
	Cfop            interface{} `json:"cfop"`
	Restitutions    struct {
	} `json:"restitutions"`
	Volumes          interface{} `json:"volumes"`
	EnableInferItems interface{} `json:"EnableInferItems"`
}

type Delivered struct {
	IsDelivered   bool            `json:"isDelivered"`
	DeliveredDate string          `json:"deliveredDate"`
	Events        DeliveredEvents `json:"events"`
}
type DeliveredEvents struct {
	City        string `json:"city"`
	State       string `json:"state"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type SaldoPicku struct {
	Data struct {
		InventoryProducts struct {
			Products []struct {
				ID                string `json:"id"`
				Sku               string `json:"sku"`
				Name              string `json:"name"`
				WarehouseID       string `json:"warehouseId"`
				AvailableQuantity int    `json:"availableQuantity"`
				Unlimited         bool   `json:"unlimited"`
			} `json:"products"`
			Paging struct {
				Pages   int `json:"pages"`
				PerPage int `json:"perPage"`
				Total   int `json:"total"`
			} `json:"paging"`
		} `json:"inventoryProducts"`
	} `json:"data"`
}

type ResponseSendPush struct {
	Enviados    []int `json:"enviados"`
	NaoEnviados []int `json:"naoEnviados"`
}

type SaldoEventoPickup struct {
	Data struct {
		ProductHistory struct {
			SkuName        string `json:"skuName"`
			Sku            string `json:"sku"`
			CatalogProduct struct {
				Name string `json:"name"`
			} `json:"catalogProduct"`
			Quantity          int `json:"quantity"`
			ReservedQuantity  int `json:"reservedQuantity"`
			AvailableQuantity int `json:"availableQuantity"`
			ChangelogHistory  []struct {
				User           string    `json:"user"`
				QuantityBefore int       `json:"quantityBefore"`
				QuantityAfter  int       `json:"quantityAfter"`
				Date           time.Time `json:"date"`
			} `json:"changelogHistory"`
		} `json:"productHistory"`
	} `json:"data"`
}

type CliqueRetire []struct {
	DeliveredDate                       time.Time
	UrlRastreio                         string
	CourierName                         string
	FinishedCourier                     bool
	AddressType                         any       `json:"address_type"`
	ArquivoTermo                        any       `json:"arquivo_termo"`
	ChapaColaborador                    any       `json:"chapaColaborador"`
	CodStatus                           int       `json:"cod_status"`
	CodTeste                            any       `json:"cod_teste"`
	CustomerAddressNumber               any       `json:"customer_address_number"`
	CustomerBs64Biometry                string    `json:"customer_bs64_biometry"`
	CustomerBs64Signature               any       `json:"customer_bs64_signature"`
	CustomerCity                        any       `json:"customer_city"`
	CustomerCpf                         string    `json:"customer_cpf"`
	CustomerCpfPickup                   string    `json:"customer_cpf_pickup"`
	CustomerEmail                       string    `json:"customer_email"`
	CustomerName                        string    `json:"customer_name"`
	CustomerNameReceived                any       `json:"customer_name_received"`
	CustomerNeighborhood                any       `json:"customer_neighborhood"`
	CustomerPhone                       string    `json:"customer_phone"`
	CustomerPostalcode                  any       `json:"customer_postalcode"`
	CustomerState                       any       `json:"customer_state"`
	CustomerStreet                      any       `json:"customer_street"`
	DataEventoRetirada                  any       `json:"data_evento_retirada"`
	DataHoraEventoRetirada              any       `json:"data_hora_evento_retirada"`
	EntregaTermo                        any       `json:"entrega_termo"`
	EstimateShippingInStore             time.Time `json:"estimate_shipping_in_store"`
	IDOrder                             string    `json:"id_order"`
	IDStorePickup                       string    `json:"id_store_pickup"`
	InvoiceKey                          string    `json:"invoice_key"`
	IsShippingCompany                   any       `json:"is_shipping_company"`
	IsShippingEmployee                  any       `json:"is_shipping_employee"`
	IsShippingFromStore                 any       `json:"is_shipping_from_store"`
	IsThirdPickup                       any       `json:"is_third_pickup"`
	ItemUniforme                        any       `json:"item_uniforme"`
	LockInChapa                         string    `json:"lock_in_chapa"`
	LockInDate                          time.Time `json:"lock_in_date"`
	LockNumber                          string    `json:"lock_number"`
	LockOutChapa                        string    `json:"lock_out_chapa"`
	LockOutDate                         time.Time `json:"lock_out_date"`
	LockRegion                          string    `json:"lock_region"`
	NfeOrder                            string    `json:"nfe_order"`
	NotificationEstimateShippingInStore any       `json:"notification_estimate_shipping_in_store"`
	NotificationOutLimitDate            any       `json:"notification_out_limit_date"`
	OrderHasTechno                      any       `json:"order_has_techno"`
	OrderTechitemIsland                 any       `json:"order_techitem_island"`
	PaymentSystemName                   any       `json:"paymentSystemName"`
	PedidoTratado                       any       `json:"pedido_Tratado"`
	QrCode                              any       `json:"qrCode"`
	Search                              any       `json:"search"`
	SellerChapa                         any       `json:"seller_chapa"`
	SendNps                             any       `json:"send_nps"`
	SendWorkchat                        any       `json:"send_workchat"`
	Status                              string    `json:"status"`
	StoreAddress                        string    `json:"store_address"`
	ThirdPickupCpf                      any       `json:"third_pickup_cpf"`
	TrackingCode                        string    `json:"tracking_code"`
	ID                                  string    `json:"id"`
	AccountID                           string    `json:"accountId"`
	AccountName                         string    `json:"accountName"`
	DataEntityID                        string    `json:"dataEntityId"`
	CreatedBy                           string    `json:"createdBy"`
	CreatedIn                           time.Time `json:"createdIn"`
	UpdatedBy                           string    `json:"updatedBy"`
	UpdatedIn                           time.Time `json:"updatedIn"`
	LastInteractionBy                   string    `json:"lastInteractionBy"`
	LastInteractionIn                   time.Time `json:"lastInteractionIn"`
	Followers                           []any     `json:"followers"`
	Tags                                []any     `json:"tags"`
	AutoFilter                          any       `json:"auto_filter"`
}

type SimulationVtex struct {
	DataProcessa  string
	Cep           string
	PickupId      int
	PickupName    string
	MetodoEntrega string
	Items         []struct {
		ID                    string    `json:"id"`
		RequestIndex          int       `json:"requestIndex"`
		Quantity              int       `json:"quantity"`
		Seller                string    `json:"seller"`
		SellerChain           []string  `json:"sellerChain"`
		Tax                   int       `json:"tax"`
		PriceValidUntil       time.Time `json:"priceValidUntil"`
		Price                 int       `json:"price"`
		ListPrice             int       `json:"listPrice"`
		RewardValue           int       `json:"rewardValue"`
		SellingPrice          int       `json:"sellingPrice"`
		Offerings             []any     `json:"offerings"`
		PriceTags             []any     `json:"priceTags"`
		MeasurementUnit       string    `json:"measurementUnit"`
		UnitMultiplier        float64   `json:"unitMultiplier"`
		ParentItemIndex       any       `json:"parentItemIndex"`
		ParentAssemblyBinding any       `json:"parentAssemblyBinding"`
		Availability          string    `json:"availability"`
		CatalogProvider       string    `json:"catalogProvider"`
		PriceDefinition       struct {
			CalculatedSellingPrice int `json:"calculatedSellingPrice"`
			Total                  int `json:"total"`
			SellingPrices          []struct {
				Value    int `json:"value"`
				Quantity int `json:"quantity"`
			} `json:"sellingPrices"`
		} `json:"priceDefinition"`
	} `json:"items"`
	RatesAndBenefitsData struct {
		RateAndBenefitsIdentifiers []any `json:"rateAndBenefitsIdentifiers"`
		Teaser                     []any `json:"teaser"`
	} `json:"ratesAndBenefitsData"`
	PaymentData struct {
		InstallmentOptions []struct {
			PaymentSystem    string `json:"paymentSystem"`
			Bin              any    `json:"bin"`
			PaymentName      string `json:"paymentName"`
			PaymentGroupName string `json:"paymentGroupName"`
			Value            int    `json:"value"`
			Installments     []struct {
				Count                      int  `json:"count"`
				HasInterestRate            bool `json:"hasInterestRate"`
				InterestRate               int  `json:"interestRate"`
				Value                      int  `json:"value"`
				Total                      int  `json:"total"`
				SellerMerchantInstallments []struct {
					ID              string `json:"id"`
					Count           int    `json:"count"`
					HasInterestRate bool   `json:"hasInterestRate"`
					InterestRate    int    `json:"interestRate"`
					Value           int    `json:"value"`
					Total           int    `json:"total"`
				} `json:"sellerMerchantInstallments"`
			} `json:"installments"`
		} `json:"installmentOptions"`
		PaymentSystems []struct {
			ID                     int       `json:"id"`
			Name                   string    `json:"name"`
			GroupName              string    `json:"groupName"`
			Validator              any       `json:"validator"`
			StringID               string    `json:"stringId"`
			Template               string    `json:"template"`
			RequiresDocument       bool      `json:"requiresDocument"`
			DisplayDocument        bool      `json:"displayDocument"`
			IsCustom               bool      `json:"isCustom"`
			Description            any       `json:"description"`
			RequiresAuthentication bool      `json:"requiresAuthentication"`
			DueDate                time.Time `json:"dueDate"`
			AvailablePayments      any       `json:"availablePayments"`
		} `json:"paymentSystems"`
		Payments              []any `json:"payments"`
		GiftCards             []any `json:"giftCards"`
		GiftCardMessages      []any `json:"giftCardMessages"`
		AvailableAccounts     []any `json:"availableAccounts"`
		AvailableTokens       []any `json:"availableTokens"`
		AvailableAssociations struct {
		} `json:"availableAssociations"`
	} `json:"paymentData"`
	SelectableGifts []any  `json:"selectableGifts"`
	MarketingData   any    `json:"marketingData"`
	PostalCode      string `json:"postalCode"`
	Country         string `json:"country"`
	LogisticsInfo   []struct {
		ItemIndex               int      `json:"itemIndex"`
		AddressID               any      `json:"addressId"`
		SelectedSLA             any      `json:"selectedSla"`
		SelectedDeliveryChannel any      `json:"selectedDeliveryChannel"`
		Quantity                int      `json:"quantity"`
		ShipsTo                 []string `json:"shipsTo"`
		Slas                    []struct {
			ID              string `json:"id"`
			DeliveryChannel string `json:"deliveryChannel"`
			Name            string `json:"name"`
			DeliveryIds     []struct {
				CourierID      string `json:"courierId"`
				WarehouseID    string `json:"warehouseId"`
				DockID         string `json:"dockId"`
				CourierName    string `json:"courierName"`
				Quantity       int    `json:"quantity"`
				KitItemDetails []any  `json:"kitItemDetails"`
			} `json:"deliveryIds"`
			ShippingEstimate         string `json:"shippingEstimate"`
			ShippingEstimateDate     any    `json:"shippingEstimateDate"`
			LockTTL                  any    `json:"lockTTL"`
			AvailableDeliveryWindows []any  `json:"availableDeliveryWindows"`
			DeliveryWindow           any    `json:"deliveryWindow"`
			Price                    int    `json:"price"`
			ListPrice                int    `json:"listPrice"`
			Tax                      int    `json:"tax"`
			PickupStoreInfo          struct {
				IsPickupStore  bool `json:"isPickupStore"`
				FriendlyName   any  `json:"friendlyName"`
				Address        any  `json:"address"`
				AdditionalInfo any  `json:"additionalInfo"`
				DockID         any  `json:"dockId"`
			} `json:"pickupStoreInfo"`
			PickupPointID  string  `json:"pickupPointId"`
			PickupDistance float64 `json:"pickupDistance"`
			PolygonName    string  `json:"polygonName"`
			TransitTime    string  `json:"transitTime"`
		} `json:"slas"`
		DeliveryChannels []struct {
			ID string `json:"id"`
		} `json:"deliveryChannels"`
	} `json:"logisticsInfo"`
	Messages           []any `json:"messages"`
	PurchaseConditions struct {
		ItemPurchaseConditions []struct {
			ID          string   `json:"id"`
			Seller      string   `json:"seller"`
			SellerChain []string `json:"sellerChain"`
			Slas        []struct {
				ID              string `json:"id"`
				DeliveryChannel string `json:"deliveryChannel"`
				Name            string `json:"name"`
				DeliveryIds     []struct {
					CourierID      string `json:"courierId"`
					WarehouseID    string `json:"warehouseId"`
					DockID         string `json:"dockId"`
					CourierName    string `json:"courierName"`
					Quantity       int    `json:"quantity"`
					KitItemDetails []any  `json:"kitItemDetails"`
				} `json:"deliveryIds"`
				ShippingEstimate         string `json:"shippingEstimate"`
				ShippingEstimateDate     any    `json:"shippingEstimateDate"`
				LockTTL                  any    `json:"lockTTL"`
				AvailableDeliveryWindows []any  `json:"availableDeliveryWindows"`
				DeliveryWindow           any    `json:"deliveryWindow"`
				Price                    int    `json:"price"`
				ListPrice                int    `json:"listPrice"`
				Tax                      int    `json:"tax"`
				PickupStoreInfo          struct {
					IsPickupStore  bool `json:"isPickupStore"`
					FriendlyName   any  `json:"friendlyName"`
					Address        any  `json:"address"`
					AdditionalInfo any  `json:"additionalInfo"`
					DockID         any  `json:"dockId"`
				} `json:"pickupStoreInfo"`
				PickupPointID  any     `json:"pickupPointId"`
				PickupDistance float64 `json:"pickupDistance"`
				PolygonName    string  `json:"polygonName"`
				TransitTime    string  `json:"transitTime"`
			} `json:"slas"`
			Price     int `json:"price"`
			ListPrice int `json:"listPrice"`
		} `json:"itemPurchaseConditions"`
	} `json:"purchaseConditions"`
	PickupPoints []struct {
		FriendlyName string `json:"friendlyName"`
		Address      struct {
			AddressType    string    `json:"addressType"`
			ReceiverName   any       `json:"receiverName"`
			AddressID      string    `json:"addressId"`
			IsDisposable   bool      `json:"isDisposable"`
			PostalCode     string    `json:"postalCode"`
			City           string    `json:"city"`
			State          string    `json:"state"`
			Country        string    `json:"country"`
			Street         string    `json:"street"`
			Number         string    `json:"number"`
			Neighborhood   string    `json:"neighborhood"`
			Complement     string    `json:"complement"`
			Reference      any       `json:"reference"`
			GeoCoordinates []float64 `json:"geoCoordinates"`
		} `json:"address"`
		AdditionalInfo string `json:"additionalInfo"`
		ID             string `json:"id"`
		BusinessHours  []struct {
			DayOfWeek   int    `json:"DayOfWeek"`
			OpeningTime string `json:"OpeningTime"`
			ClosingTime string `json:"ClosingTime"`
		} `json:"businessHours"`
	} `json:"pickupPoints"`
	SubscriptionData any `json:"subscriptionData"`
	Totals           []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"totals"`
	ItemMetadata            any  `json:"itemMetadata"`
	AllowMultipleDeliveries bool `json:"allowMultipleDeliveries"`
}
type AllPickupPoints []struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Instructions     string `json:"instructions"`
	FormattedAddress string `json:"formatted_address"`
	Address          struct {
		PostalCode string `json:"postalCode"`
		Country    struct {
			Acronym string `json:"acronym"`
			Name    string `json:"name"`
		} `json:"country"`
		City         string `json:"city"`
		State        string `json:"state"`
		Neighborhood string `json:"neighborhood"`
		Street       string `json:"street"`
		Number       string `json:"number"`
		Complement   string `json:"complement"`
		Reference    string `json:"reference"`
		Location     struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	} `json:"address"`
	IsActive bool    `json:"isActive"`
	Distance float64 `json:"distance"`
	Seller   string  `json:"seller"`
	//Sort          []int64 `json:"_sort"`
	BusinessHours []struct {
		DayOfWeek   int    `json:"dayOfWeek"`
		OpeningTime string `json:"openingTime"`
		ClosingTime string `json:"closingTime"`
	} `json:"businessHours"`
	TagsLabel          []any  `json:"tagsLabel"`
	PickupHolidays     []any  `json:"pickupHolidays"`
	IsThirdPartyPickup bool   `json:"isThirdPartyPickup"`
	AccountOwnerName   string `json:"accountOwnerName"`
	AccountOwnerID     string `json:"accountOwnerId"`
	ParentAccountName  any    `json:"parentAccountName"`
	OriginalID         any    `json:"originalId"`
}
type OrderDetails struct {
	OrderID                     string      `json:"orderId"`
	Sequence                    string      `json:"sequence"`
	MarketplaceOrderID          string      `json:"marketplaceOrderId"`
	MarketplaceServicesEndpoint string      `json:"marketplaceServicesEndpoint"`
	SellerOrderID               string      `json:"sellerOrderId"`
	Origin                      string      `json:"origin"`
	AffiliateID                 string      `json:"affiliateId"`
	SalesChannel                string      `json:"salesChannel"`
	MerchantName                interface{} `json:"merchantName"`
	Status                      string      `json:"status"`
	WorkflowIsInError           bool        `json:"workflowIsInError"`
	StatusDescription           string      `json:"statusDescription"`
	Value                       int         `json:"value"`
	CreationDate                time.Time   `json:"creationDate"`
	LastChange                  time.Time   `json:"lastChange"`
	OrderGroup                  interface{} `json:"orderGroup"`
	Totals                      []struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"totals"`
	Items []struct {
		UniqueID       string `json:"uniqueId"`
		ID             string `json:"id"`
		ProductID      string `json:"productId"`
		Ean            string `json:"ean"`
		LockID         string `json:"lockId"`
		ItemAttachment struct {
			Content struct {
			} `json:"content"`
			Name interface{} `json:"name"`
		} `json:"itemAttachment"`
		Attachments []interface{} `json:"attachments"`
		Quantity    int           `json:"quantity"`
		Seller      string        `json:"seller"`
		Name        string        `json:"name"`
		RefID       string        `json:"refId"`
		Price       int           `json:"price"`
		ListPrice   int           `json:"listPrice"`
		ManualPrice interface{}   `json:"manualPrice"`
		PriceTags   []struct {
			Name         string      `json:"name"`
			Value        int         `json:"value"`
			IsPercentual bool        `json:"isPercentual"`
			Identifier   interface{} `json:"identifier"`
			RawValue     float64     `json:"rawValue"`
			Rate         interface{} `json:"rate"`
			JurisCode    interface{} `json:"jurisCode"`
			JurisType    interface{} `json:"jurisType"`
			JurisName    interface{} `json:"jurisName"`
		} `json:"priceTags"`
		ImageURL            string        `json:"imageUrl"`
		DetailURL           interface{}   `json:"detailUrl"`
		Components          []interface{} `json:"components"`
		BundleItems         []interface{} `json:"bundleItems"`
		Params              []interface{} `json:"params"`
		Offerings           []interface{} `json:"offerings"`
		AttachmentOfferings []interface{} `json:"attachmentOfferings"`
		SellerSku           string        `json:"sellerSku"`
		PriceValidUntil     interface{}   `json:"priceValidUntil"`
		Commission          int           `json:"commission"`
		Tax                 int           `json:"tax"`
		PreSaleDate         interface{}   `json:"preSaleDate"`
		AdditionalInfo      struct {
			BrandName     string `json:"brandName"`
			BrandID       string `json:"brandId"`
			CategoriesIds string `json:"categoriesIds"`
			Categories    []struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"categories"`
			ProductClusterID      string `json:"productClusterId"`
			CommercialConditionID string `json:"commercialConditionId"`
			Dimension             struct {
				Cubicweight float64 `json:"cubicweight"`
				Height      float64 `json:"height"`
				Length      float64 `json:"length"`
				Weight      float64 `json:"weight"`
				Width       float64 `json:"width"`
			} `json:"dimension"`
			OfferingInfo   interface{} `json:"offeringInfo"`
			OfferingType   interface{} `json:"offeringType"`
			OfferingTypeID interface{} `json:"offeringTypeId"`
		} `json:"additionalInfo"`
		MeasurementUnit   string      `json:"measurementUnit"`
		UnitMultiplier    float64     `json:"unitMultiplier"`
		SellingPrice      int         `json:"sellingPrice"`
		IsGift            bool        `json:"isGift"`
		ShippingPrice     interface{} `json:"shippingPrice"`
		RewardValue       int         `json:"rewardValue"`
		FreightCommission int         `json:"freightCommission"`
		PriceDefinition   struct {
			SellingPrices []struct {
				Value    int `json:"value"`
				Quantity int `json:"quantity"`
			} `json:"sellingPrices"`
			CalculatedSellingPrice int         `json:"calculatedSellingPrice"`
			Total                  int         `json:"total"`
			Reason                 interface{} `json:"reason"`
		} `json:"priceDefinition"`
		TaxCode               interface{}   `json:"taxCode"`
		ParentItemIndex       interface{}   `json:"parentItemIndex"`
		ParentAssemblyBinding interface{}   `json:"parentAssemblyBinding"`
		CallCenterOperator    interface{}   `json:"callCenterOperator"`
		SerialNumbers         interface{}   `json:"serialNumbers"`
		Assemblies            []interface{} `json:"assemblies"`
		CostPrice             int           `json:"costPrice"`
	} `json:"items"`
	MarketplaceItems  []interface{} `json:"marketplaceItems"`
	ClientProfileData struct {
		ID                 string      `json:"id"`
		Email              string      `json:"email"`
		FirstName          string      `json:"firstName"`
		LastName           string      `json:"lastName"`
		DocumentType       string      `json:"documentType"`
		Document           string      `json:"document"`
		Phone              string      `json:"phone"`
		CorporateName      interface{} `json:"corporateName"`
		TradeName          interface{} `json:"tradeName"`
		CorporateDocument  interface{} `json:"corporateDocument"`
		StateInscription   interface{} `json:"stateInscription"`
		CorporatePhone     interface{} `json:"corporatePhone"`
		IsCorporate        bool        `json:"isCorporate"`
		UserProfileID      interface{} `json:"userProfileId"`
		UserProfileVersion interface{} `json:"userProfileVersion"`
		CustomerClass      interface{} `json:"customerClass"`
	} `json:"clientProfileData"`
	GiftRegistryData     interface{} `json:"giftRegistryData"`
	MarketingData        interface{} `json:"marketingData"`
	RatesAndBenefitsData struct {
		ID                         string        `json:"id"`
		RateAndBenefitsIdentifiers []interface{} `json:"rateAndBenefitsIdentifiers"`
	} `json:"ratesAndBenefitsData"`
	ShippingData struct {
		ID      string `json:"id"`
		Address struct {
			AddressType    string        `json:"addressType"`
			ReceiverName   string        `json:"receiverName"`
			AddressID      string        `json:"addressId"`
			VersionID      interface{}   `json:"versionId"`
			EntityID       interface{}   `json:"entityId"`
			PostalCode     string        `json:"postalCode"`
			City           string        `json:"city"`
			State          string        `json:"state"`
			Country        string        `json:"country"`
			Street         string        `json:"street"`
			Number         string        `json:"number"`
			Neighborhood   string        `json:"neighborhood"`
			Complement     string        `json:"complement"`
			Reference      string        `json:"reference"`
			GeoCoordinates []interface{} `json:"geoCoordinates"`
		} `json:"address"`
		LogisticsInfo []struct {
			ItemIndex            int         `json:"itemIndex"`
			SelectedSLA          string      `json:"selectedSla"`
			LockTTL              string      `json:"lockTTL"`
			Price                int         `json:"price"`
			ListPrice            int         `json:"listPrice"`
			SellingPrice         int         `json:"sellingPrice"`
			DeliveryWindow       interface{} `json:"deliveryWindow"`
			DeliveryCompany      string      `json:"deliveryCompany"`
			ShippingEstimate     string      `json:"shippingEstimate"`
			ShippingEstimateDate time.Time   `json:"shippingEstimateDate"`
			Slas                 interface{} `json:"slas"`
			ShipsTo              interface{} `json:"shipsTo"`
			DeliveryIds          []struct {
				CourierID          string        `json:"courierId"`
				CourierName        string        `json:"courierName"`
				DockID             string        `json:"dockId"`
				Quantity           int           `json:"quantity"`
				WarehouseID        string        `json:"warehouseId"`
				AccountCarrierName string        `json:"accountCarrierName"`
				KitItemDetails     []interface{} `json:"kitItemDetails"`
			} `json:"deliveryIds"`
			DeliveryChannels interface{} `json:"deliveryChannels"`
			DeliveryChannel  string      `json:"deliveryChannel"`
			PickupStoreInfo  struct {
				AdditionalInfo interface{} `json:"additionalInfo"`
				Address        interface{} `json:"address"`
				DockID         interface{} `json:"dockId"`
				FriendlyName   interface{} `json:"friendlyName"`
				IsPickupStore  bool        `json:"isPickupStore"`
			} `json:"pickupStoreInfo"`
			AddressID     string      `json:"addressId"`
			VersionID     interface{} `json:"versionId"`
			EntityID      interface{} `json:"entityId"`
			PolygonName   string      `json:"polygonName"`
			PickupPointID interface{} `json:"pickupPointId"`
			TransitTime   string      `json:"transitTime"`
		} `json:"logisticsInfo"`
		TrackingHints     []interface{} `json:"trackingHints"`
		SelectedAddresses []struct {
			AddressID      string        `json:"addressId"`
			VersionID      interface{}   `json:"versionId"`
			EntityID       interface{}   `json:"entityId"`
			AddressType    string        `json:"addressType"`
			ReceiverName   string        `json:"receiverName"`
			Street         string        `json:"street"`
			Number         string        `json:"number"`
			Complement     string        `json:"complement"`
			Neighborhood   string        `json:"neighborhood"`
			PostalCode     string        `json:"postalCode"`
			City           string        `json:"city"`
			State          string        `json:"state"`
			Country        string        `json:"country"`
			Reference      string        `json:"reference"`
			GeoCoordinates []interface{} `json:"geoCoordinates"`
		} `json:"selectedAddresses"`
	} `json:"shippingData"`
	PaymentData struct {
		GiftCards    []interface{} `json:"giftCards"`
		Transactions []struct {
			IsActive      bool        `json:"isActive"`
			TransactionID interface{} `json:"transactionId"`
			MerchantName  interface{} `json:"merchantName"`
			Payments      []struct {
				ID                 interface{} `json:"id"`
				PaymentSystem      string      `json:"paymentSystem"`
				PaymentSystemName  string      `json:"paymentSystemName"`
				Value              int         `json:"value"`
				Installments       int         `json:"installments"`
				ReferenceValue     int         `json:"referenceValue"`
				CardHolder         interface{} `json:"cardHolder"`
				CardNumber         interface{} `json:"cardNumber"`
				FirstDigits        interface{} `json:"firstDigits"`
				LastDigits         interface{} `json:"lastDigits"`
				Cvv2               interface{} `json:"cvv2"`
				ExpireMonth        interface{} `json:"expireMonth"`
				ExpireYear         interface{} `json:"expireYear"`
				URL                interface{} `json:"url"`
				GiftCardID         interface{} `json:"giftCardId"`
				GiftCardName       interface{} `json:"giftCardName"`
				GiftCardCaption    interface{} `json:"giftCardCaption"`
				RedemptionCode     interface{} `json:"redemptionCode"`
				Group              interface{} `json:"group"`
				Tid                interface{} `json:"tid"`
				DueDate            interface{} `json:"dueDate"`
				ConnectorResponses struct {
				} `json:"connectorResponses"`
				GiftCardProvider                               interface{} `json:"giftCardProvider"`
				GiftCardAsDiscount                             interface{} `json:"giftCardAsDiscount"`
				KoinURL                                        interface{} `json:"koinUrl"`
				AccountID                                      interface{} `json:"accountId"`
				ParentAccountID                                interface{} `json:"parentAccountId"`
				BankIssuedInvoiceIdentificationNumber          interface{} `json:"bankIssuedInvoiceIdentificationNumber"`
				BankIssuedInvoiceIdentificationNumberFormatted interface{} `json:"bankIssuedInvoiceIdentificationNumberFormatted"`
				BankIssuedInvoiceBarCodeNumber                 interface{} `json:"bankIssuedInvoiceBarCodeNumber"`
				BankIssuedInvoiceBarCodeType                   interface{} `json:"bankIssuedInvoiceBarCodeType"`
				BillingAddress                                 interface{} `json:"billingAddress"`
				PaymentOrigin                                  interface{} `json:"paymentOrigin"`
			} `json:"payments"`
		} `json:"transactions"`
	} `json:"paymentData"`
	PackageAttachment struct {
		Packages []struct {
			Items []struct {
				ItemIndex      int         `json:"itemIndex"`
				Quantity       int         `json:"quantity"`
				Price          int         `json:"price"`
				Description    interface{} `json:"description"`
				UnitMultiplier float64     `json:"unitMultiplier"`
			} `json:"items"`
			Courier         string      `json:"courier"`
			InvoiceNumber   string      `json:"invoiceNumber"`
			Description     interface{} `json:"description"`
			InvoiceValue    int         `json:"invoiceValue"`
			InvoiceURL      interface{} `json:"invoiceUrl"`
			IssuanceDate    time.Time   `json:"issuanceDate"`
			TrackingNumber  string      `json:"trackingNumber"`
			InvoiceKey      string      `json:"invoiceKey"`
			TrackingURL     string      `json:"trackingUrl"`
			EmbeddedInvoice string      `json:"embeddedInvoice"`
			Type            string      `json:"type"`
			CourierStatus   struct {
				Status        string    `json:"status"`
				Finished      bool      `json:"finished"`
				DeliveredDate time.Time `json:"deliveredDate"`
				Data          []struct {
					LastChange  time.Time `json:"lastChange"`
					City        string    `json:"city"`
					State       string    `json:"state"`
					Description string    `json:"description"`
					CreateDate  time.Time `json:"createDate"`
				} `json:"data"`
			} `json:"courierStatus"`
			Cfop         interface{} `json:"cfop"`
			Restitutions struct {
			} `json:"restitutions"`
			Volumes          int         `json:"volumes"`
			EnableInferItems interface{} `json:"EnableInferItems"`
		} `json:"packages"`
	} `json:"packageAttachment"`
	Sellers []struct {
		ID                  string      `json:"id"`
		Name                string      `json:"name"`
		Logo                interface{} `json:"logo"`
		FulfillmentEndpoint interface{} `json:"fulfillmentEndpoint"`
	} `json:"sellers"`
	CallCenterOperatorData interface{} `json:"callCenterOperatorData"`
	FollowUpEmail          string      `json:"followUpEmail"`
	LastMessage            interface{} `json:"lastMessage"`
	Hostname               string      `json:"hostname"`
	InvoiceData            struct {
		Address         interface{} `json:"address"`
		UserPaymentInfo interface{} `json:"userPaymentInfo"`
	} `json:"invoiceData"`
	ChangesAttachment interface{} `json:"changesAttachment"`
	OpenTextField     struct {
		Value string `json:"value"`
	} `json:"openTextField"`
	RoundingError           int         `json:"roundingError"`
	OrderFormID             interface{} `json:"orderFormId"`
	CommercialConditionData interface{} `json:"commercialConditionData"`
	IsCompleted             bool        `json:"isCompleted"`
	CustomData              struct {
		CustomApps []struct {
			Fields struct {
				MarketplacePaymentMethod string `json:"marketplacePaymentMethod"`
				MarketplaceFreightPrice  string `json:"marketplaceFreightPrice"`
			} `json:"fields"`
			ID    string `json:"id"`
			Major int    `json:"major"`
		} `json:"customApps"`
	} `json:"customData"`
	StorePreferencesData struct {
		CountryCode        string `json:"countryCode"`
		CurrencyCode       string `json:"currencyCode"`
		CurrencyFormatInfo struct {
			CurrencyDecimalDigits    int    `json:"CurrencyDecimalDigits"`
			CurrencyDecimalSeparator string `json:"CurrencyDecimalSeparator"`
			CurrencyGroupSeparator   string `json:"CurrencyGroupSeparator"`
			CurrencyGroupSize        int    `json:"CurrencyGroupSize"`
			StartsWithCurrencySymbol bool   `json:"StartsWithCurrencySymbol"`
		} `json:"currencyFormatInfo"`
		CurrencyLocale int    `json:"currencyLocale"`
		CurrencySymbol string `json:"currencySymbol"`
		TimeZone       string `json:"timeZone"`
	} `json:"storePreferencesData"`
	AllowCancellation bool `json:"allowCancellation"`
	AllowEdition      bool `json:"allowEdition"`
	IsCheckedIn       bool `json:"isCheckedIn"`
	Marketplace       struct {
		BaseURL     string `json:"baseURL"`
		IsCertified bool   `json:"isCertified"`
		Name        string `json:"name"`
	} `json:"marketplace"`
	AuthorizedDate time.Time   `json:"authorizedDate"`
	InvoicedDate   time.Time   `json:"invoicedDate"`
	CancelReason   interface{} `json:"cancelReason"`
	ItemMetadata   struct {
		Items []struct {
			ID              string        `json:"Id"`
			Seller          string        `json:"Seller"`
			Name            string        `json:"Name"`
			SkuName         string        `json:"SkuName"`
			ProductID       string        `json:"ProductId"`
			RefID           string        `json:"RefId"`
			Ean             string        `json:"Ean"`
			ImageURL        string        `json:"ImageUrl"`
			DetailURL       string        `json:"DetailUrl"`
			AssemblyOptions []interface{} `json:"AssemblyOptions"`
		} `json:"Items"`
	} `json:"itemMetadata"`
	SubscriptionData       interface{} `json:"subscriptionData"`
	TaxData                interface{} `json:"taxData"`
	CheckedInPickupPointID interface{} `json:"checkedInPickupPointId"`
	CancellationData       interface{} `json:"cancellationData"`
	ClientPreferencesData  interface{} `json:"clientPreferencesData"`
}

func ConsultaAllPickup() (resp AllPickupPoints, err error) {
	cfg := &config.Config{}
	err = config.New(cfg)
	if err != nil {
		fmt.Println("Arquivo '.env' não encontrado")
	}
	url := "https://" + config.Env.Account + ".vtexcommercestable.com.br/api/logistics/pvt/configuration/pickuppoints"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", config.Env.KeyVtex)
	req.Header.Add("X-VTEX-API-AppToken", config.Env.TokenVtex)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}

func ConsultaTokenVtex() (resp GetTokenVtex, err error) {
	cfg := &config.Config{}
	err = config.New(cfg)
	if err != nil {
		fmt.Println("Arquivo '.env' não encontrado")
	}
	url := "https://api.myvtex.com.br/api/vtexid/apptoken/login?an=" + cfg.Account
	method := "POST"

	payload := strings.NewReader(`{
		"appkey": "` + cfg.KeyVtex + `",
		"apptoken": "` + cfg.TokenVtex + `"
	  }`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}

func FaturaProdutosPedido(account string, key string, token string, OrderId string, Invoice string, shipping int, items []InvoicedItem) (resp Receipt, err error) {
	var dadosFatura Invoiced
	url := "https://" + account + ".vtexcommercestable.com.br/api/oms/pvt/orders/" + OrderId + "/invoice"
	method := "POST"
	client := &http.Client{}

	var value = shipping
	for _, produto := range items {
		value = value + produto.Price*produto.Quantity
	}

	dadosFatura.Items = items
	dadosFatura.InvoiceNumber = OrderId
	dadosFatura.InvoiceValue = value
	dadosFatura.IssuanceDate = time.Now()
	dadosFatura.TrackingNumber = ""
	dadosFatura.InvoiceKey = ""
	dadosFatura.NumberOrder = Invoice
	dadosFatura.Type = "Output"

	out, err := json.Marshal(dadosFatura)
	payload := strings.NewReader(string(out))
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}

func MarcarPedidoEntregue(account string, key string, token string, OrderId string, Invoice string) (resp Receipt, err error) {
	var dadosEntrega Delivered
	var dadosEntregaEventos DeliveredEvents
	url := "https://" + account + ".vtexcommercestable.com.br/api/oms/pvt/orders/" + OrderId + "/invoice/" + Invoice + "/tracking"
	method := "PUT"
	client := &http.Client{}
	current_time := time.Now()
	dadosEntrega.DeliveredDate = current_time.Format("2006-01-02 15:04:05")
	dadosEntrega.IsDelivered = true
	dadosEntregaEventos.City = "São Paulo"
	dadosEntregaEventos.Date = current_time.Format("2006-01-02 15:04:05")
	dadosEntregaEventos.State = "SP"
	dadosEntrega.Events = dadosEntregaEventos

	out, err := json.Marshal(dadosEntrega)
	payload := strings.NewReader(string(out))
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}
	return
}

func SimuladorCarrinho(resultChan chan<- SimulationVtex, wg *sync.WaitGroup, skuid string, postalCode string, pickupid int, pickupName string, salleschannel string) {
	url := "https://" + config.Env.Account + ".vtexcommercestable.com.br/api/checkout/pub/orderForms/simulation?RnbBehavior=0&sc=" + salleschannel
	method := "POST"
	payload := strings.NewReader(`{
    "isCart": "false",
    "clientProfileData": null,
    "country": "BRA",
    "postalCode": "` + postalCode + `",
    "items": [
        {
            "id": "` + skuid + `",
            "seller": "1",
            "quantity": 1
        }
    ]
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	for {
		res, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			break
		}
		defer res.Body.Close()
		if res.StatusCode == http.StatusTooManyRequests {
			// Atraso de 1 segundo antes de tentar novamente
			time.Sleep(time.Second)
			break
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err)
			break
		}
		//SimulationVtex resp
		var resp SimulationVtex
		err = json.Unmarshal(body, &resp)
		if err != nil {
			fmt.Println("Erro ao converter Json " + err.Error() + " - ")
			break
		}
		if len(resp.Items) == 0 {
			log.Fatalln(resp)
		}
		resp.PickupId = pickupid
		resp.Cep = postalCode
		resp.PickupName = pickupName
		current_time := time.Now()
		resp.DataProcessa = current_time.Format("2006-01-02 15:04:05")
		resultChan <- resp
		defer wg.Done()
		break
	}
}

func GetSaldoSkuPickup(account string, token string, sku string) (resp SaldoEventoPickup, err error) {
	url := "https://" + account + ".myvtex.com/_v/private/graphql/v1?workspace=master&maxAge=long&appsEtag=remove&domain=admin&locale=pt-BR"
	//url = "https://webhook.site/06a8f5ea-1d65-411f-8e0a-f143168ef100"
	method := "POST"

	//payload := strings.NewReader(`{"query":"{inventoryProducts(page: 1, perPage: 10,filter: {skus:` + sku + `}, sorting: {order: ASC}) {products {id,sku      name,warehouseId,availableQuantity,unlimited}paging{pages,perPage,total}  }}","variables":{}}`)
	//payload := strings.NewReader(`{"query":"{productHistory(sku:"` + sku + `",warehouseId:"1_1",accountName:"` + account + `"){skuName, quantity, reservedQuantity, availableQuantity, changelogHistory{user,quantityBefore,quantityAfter,date}}","variables":{}}`)
	payload := strings.NewReader(`{"query":"query{\r\n  productHistory(sku:\"` + sku + `\",warehouseId:\"1_1\",accountName:\"` + account + `\"){\r\n    skuName,\r\n    sku,\r\n    catalogProduct{name\r\n }\r\n    quantity\r\n    reservedQuantity\r\n    availableQuantity\r\n    changelogHistory{user\r\n    quantityBefore\r\n      quantityAfter\r\n      date\r\n      \r\n    }\r\n  }\r\n}","variables":{}}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cookie", "VtexIdclientAutCookie="+token+";")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}

func SendPush(idSacola int) (resp ResponseSendPush, err error) {
	url := "https://psd.pernambucanas.com.br:8443/api/appvarejo/v2/sacola-de-descontos/sacolas-recuperadas/notification"
	method := "POST"

	//payload := strings.NewReader(`{"query":"{inventoryProducts(page: 1, perPage: 10,filter: {skus:` + sku + `}, sorting: {order: ASC}) {products {id,sku      name,warehouseId,availableQuantity,unlimited}paging{pages,perPage,total}  }}","variables":{}}`)
	//payload := strings.NewReader(`{"query":"{productHistory(sku:"` + sku + `",warehouseId:"1_1",accountName:"` + account + `"){skuName, quantity, reservedQuantity, availableQuantity, changelogHistory{user,quantityBefore,quantityAfter,date}}","variables":{}}`)
	payload := strings.NewReader(`{"title": "Psiu!! Não esqueceu nada, não?","body": "Lembra dos produtos que você escolheu na loja ontem? Montamos um carrinho especial pra você. Clique na notificação para conferir e finalizar sua compra! 🛒","idSacolas": [` + strconv.Itoa(idSacola) + `]
	}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}

func GetOrder(account string, key string, token string, OrderId string) (resp OrderDetails, err error) {
	url := "https://" + account + ".vtexcommercestable.com.br/api/oms/pvt/orders/" + OrderId
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}

func GetCliqueRetire(account string, key string, token string, campo string, termo string) (resp CliqueRetire, err error) {
	url := "https://" + account + ".vtexcommercestable.com.br/api/dataentities/CR/search?_fields=_all&" + termo + "=" + campo
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}

func GetSacolaAbandonada() (resp SacolasAbandonadas, err error) {
	var dataAtual = time.Now()
	var dataAnterior = dataAtual.AddDate(0, 0, -1)

	url := "https://psd.pernambucanas.com.br:8443/api/appvarejo/v2/sacola-de-descontos/sacolas-recuperadas/busca-sacolas-abandonadas?dataCriacaoInicio=" + dataAnterior.Format("2006-01-02") + "&dataCriacaoFim=" + dataAnterior.Format("2006-01-02")
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	req.Header.Add("Accept", "application/json")

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json" + err.Error())
		return
	}

	return
}

func GetLojasContaFranquia(account string, key string, token string, min int, max int) (resp LojasContaFranquia, err error) {
	url := "https://" + account + ".vtexcommercestable.com.br/api/dataentities/CF/search?_where=active=true"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)
	req.Header.Add("REST-Range", "resources="+strconv.FormatInt(int64(min), 10)+"-"+strconv.FormatInt(int64(max), 10))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = json.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("Erro ao converter Json do LojaContaFranquia" + err.Error())
		return
	}

	return
}

func RegistraEntregue(account string, key string, token string, OrderId string, Invoice string) {
	url := "https://" + account + ".vtexcommercestable.com.br/api/oms/pvt/orders/" + OrderId + "/invoice/" + Invoice + "/tracking"
	method := "PUT"

	payload := strings.NewReader(`{
     "isDelivered": true,
     "deliveredDate": "` + time.Now().GoString() + `"
	}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}

func RegistraManuseio(account string, key string, token string, OrderId string) {
	url := "https://" + account + ".vtexcommercestable.com.br/api/oms/pvt/orders/" + OrderId + "/start-handling"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-VTEX-API-AppKey", key)
	req.Header.Add("X-VTEX-API-AppToken", token)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}
