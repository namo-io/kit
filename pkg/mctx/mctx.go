package mctx

import "context"

func getStringFromContext(ctx context.Context, key Key) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}

	return v.(string)
}

func GetAuthorization(ctx context.Context) string {
	return getStringFromContext(ctx, AuthorizationKey)
}

func WithAuthorization(ctx context.Context, authorization string) context.Context {
	return context.WithValue(ctx, AuthorizationKey, authorization)
}

func GetRequestId(ctx context.Context) string {
	return getStringFromContext(ctx, RequestIdKey)
}

func WithRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, RequestIdKey, requestId)
}

func GetAppName(ctx context.Context) string {
	return getStringFromContext(ctx, AppNameKey)
}

func WithAppName(ctx context.Context, appName string) context.Context {
	return context.WithValue(ctx, AppNameKey, appName)
}

func GetAppId(ctx context.Context) string {
	return getStringFromContext(ctx, AppIdKey)
}

func WithAppId(ctx context.Context, appId string) context.Context {
	return context.WithValue(ctx, AppIdKey, appId)
}

func GetAppVersion(ctx context.Context) string {
	return getStringFromContext(ctx, AppVersionKey)
}

func WithAppVersion(ctx context.Context, appVersion string) context.Context {
	return context.WithValue(ctx, AppVersionKey, appVersion)
}
