package appcontext

import (
	"context"
)

type contextKey string

const (
	// KeyURLPath represents the url path key in http server context
	KeyURLPath contextKey = "URLPath"

	// KeyHTTPMethodName represents the method name key in http server context
	KeyHTTPMethodName contextKey = "HTTPMethodName"

	// KeySessionID represents the current logged-in SessionID
	KeySessionID contextKey = "SessionID"

	// KeyCurrentAccount represents the CurrentAccountId key in http server context
	KeyCurrentAccount contextKey = "CurrentAccount"

	// KeyOwner represents the OwnerID key in http server context
	KeyOwner contextKey = "Owner"

	// KeyUserID represents the current logged-in UserID
	KeyUserID contextKey = "UserID"

	// KeyLoginToken represents the current logged-in token
	KeyLoginToken contextKey = "LoginToken"

	// KeyWarehouseID represents the current prefered warehouseID of CustomerID
	KeyWarehouseID contextKey = "WarehouseID"

	// KeyVersionCode represents the current version code of request
	KeyVersionCode contextKey = "VersionCode"

	// KeyCurrentClientAccess represents the Current Client access key in http server context
	KeyCurrentClientAccess contextKey = "CurrentClientAccess"

	// KeyClientID represents the Current Client in http server context
	KeyClientID contextKey = "ClientID"

	// KeyIsSales represents the current type of customer
	KeyIsSales contextKey = "IsSales"

	// KeyWarehouseProvider represents the Current Client in http server context
	KeyWarehouseProvider contextKey = "WarehouseProvider"

	// KeyLogString represents the key Log String in server context
	KeyLogString contextKey = "KeyLogString"

	// KeyAllLog represents the key Log String in server context
	KeyAllLog contextKey = "KeyAllLog"
)

// Owner gets the data owner from the context
func Owner(ctx context.Context) *int {
	owner := ctx.Value(KeyOwner)
	if owner != nil {
		v := owner.(int)
		return &v
	}
	return nil
}

// URLPath gets the data url path from the context
func URLPath(ctx context.Context) *string {
	urlPath := ctx.Value(KeyURLPath)
	if urlPath != nil {
		v := urlPath.(string)
		return &v
	}
	return nil
}

// HTTPMethodName gets the data http method from the context
func HTTPMethodName(ctx context.Context) *string {
	httpMethodName := ctx.Value(KeyHTTPMethodName)
	if httpMethodName != nil {
		v := httpMethodName.(string)
		return &v
	}
	return nil
}

// SessionID gets the data session id from the context
func SessionID(ctx context.Context) *string {
	sessionID := ctx.Value(KeySessionID)
	if sessionID != nil {
		v := sessionID.(string)
		return &v
	}
	return nil
}

// CurrentAccount gets current account from the context
func CurrentAccount(ctx context.Context) *int {
	currentAccount := ctx.Value(KeyCurrentAccount)
	if currentAccount != nil {
		v := currentAccount.(int)
		return &v
	}
	return nil
}

// UserID gets current userId logged in from the context
func UserID(ctx context.Context) int {
	userID := ctx.Value(KeyUserID)
	if userID != nil {
		v := userID.(int)
		return v
	}
	return 0
}

// WarehouseID gets current prefered warehouseID of CustomerID
func WarehouseID(ctx context.Context) int {
	warehouseID := ctx.Value(KeyWarehouseID)
	if warehouseID != nil {
		v := warehouseID.(int)
		return v
	}
	return 0
}

// VersionCode gets current version code of request
func VersionCode(ctx context.Context) int {
	versionCode := ctx.Value(KeyVersionCode)
	if versionCode != nil {
		v := versionCode.(int)
		return v
	}
	return 0
}

// CurrentClientAccess gets current client id from the context
func CurrentClientAccess(ctx context.Context) []string {
	currentClientAccess := ctx.Value(KeyCurrentClientAccess)
	// datas := reflect.ValueOf(currentClientAccess)
	// if datas.Kind() != reflect.Slice {
	// 	return nil
	// }
	// if currentClientAccess != nil || datas.Len() > 0 {
	// 	v := currentClientAccess.([]string)
	// 	return v
	// }
	if currentClientAccess != nil {
		v := currentClientAccess.([]string)
		return v
	}
	return nil
}

// ClientID gets current client from the context
func ClientID(ctx context.Context) *int {
	currentClientAccess := ctx.Value(KeyClientID)
	if currentClientAccess != nil {
		v := currentClientAccess.(int)
		return &v
	}
	return nil
}

// IsSales gets current type of customer
func IsSales(ctx context.Context) bool {
	isSales := ctx.Value(KeyIsSales)
	if isSales != nil {
		v := isSales.(bool)
		return v
	}
	return false
}

// WarehouseProvider gets current client from the context
func WarehouseProvider(ctx context.Context) *int {
	warehouseProvider := ctx.Value(KeyWarehouseProvider)
	if warehouseProvider != nil {
		v := warehouseProvider.(int)
		return &v
	}
	return nil
}

// LogString gets log String from context
func LogString(ctx context.Context) *string {
	logString := ctx.Value(KeyLogString)
	if logString != nil {
		v := logString.(*string)
		return v
	}
	return nil
}

// AllLog gets log String from context
func AllLog(ctx context.Context) *string {
	logString := ctx.Value(KeyAllLog)
	if logString != nil {
		v := logString.(string)
		return &v
	}
	return nil
}
