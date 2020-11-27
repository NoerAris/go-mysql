package router

import (
	"github.com/kataras/iris/context"
	"log"
)

func CheckPublic(ctx context.Context) {

	if checkAccessLevel(ctx, RolePublic) {
		//check,err:=checkAccessLevel(ctx,RolePublic);err==nil&&
		//	check
		//	{
		ctx.Next()
	} else {
		ctx.StatusCode(RESPONSE_CODE_UNAUTHORIZED)
		/*if err!=nil{
			ctx.JSON(GetErrorResponse(err.Error()))
			return
		}*/
		ctx.JSON(GetErrorResponse("Access Forbidden"))
	}
}
func CheckRegistered(ctx context.Context) {
	if checkAccessLevel(ctx, RoleMemberRegistered) {
		//check,err:=checkAccessLevel(ctx,RoleMemberRegistered);err==nil&&check{
		ctx.Next()
	} else {
		ctx.StatusCode(RESPONSE_CODE_UNAUTHORIZED)
		/*if err!=nil{
			ctx.JSON(GetErrorResponse(err.Error()))
			return
		}*/
		ctx.JSON(GetErrorResponse("Access Forbidden"))
	}
}
func CheckApproved(ctx context.Context) {
	if checkAccessLevel(ctx, RoleMemberApproved) {
		//check,err:=checkAccessLevel(ctx,RoleMemberApproved);err==nil&&check{
		ctx.Next()
	} else {
		ctx.StatusCode(RESPONSE_CODE_UNAUTHORIZED)
		/*if err!=nil{
			ctx.JSON(GetErrorResponse(err.Error()))
			return
		}*/
		ctx.JSON(GetErrorResponse("Access Forbidden"))
	}
}
func checkAccessLevel(ctx context.Context, level string) bool {
	accessGranted := false
	claims, _ := GetClaimFromContext(ctx)
	/*if err!=nil{
		return false,err
	}*/
	if rolesEncrypted, ok := claims[JwtRoles].([]interface{}); ok {
		roles := make([]string, 0)
		for _, r := range rolesEncrypted {
			if role, err := DecryptClaimValue(r.(string)); err == nil {
				roles = append(roles, role)
			} else {
				log.Println("ERROR DECRYPTING ROLES|" + err.Error())
			}
		}
		for _, role := range roles {
			if level == role {
				accessGranted = true
				break
			}
		}
		//log.Println("DEBUG|"+ctx.GetHeader(JwtHeaderKey)+"|",roles,level,accessGranted)
	}
	return accessGranted
}
