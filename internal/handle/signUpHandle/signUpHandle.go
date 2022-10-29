package signUpHandle

import (
	"adi-sign-up-be/internal/controller/SignUps"
	serviceErr "adi-sign-up-be/internal/dto/err"
	"adi-sign-up-be/internal/dto/signUp"
	"adi-sign-up-be/internal/middleware"
	"adi-sign-up-be/internal/model/Mysql"
	"adi-sign-up-be/pkg/utils/captcha"
	"adi-sign-up-be/pkg/utils/check"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

//TODO:
//1. 覆盖报名
//2. 数据库ID设计
//3. 信息为空不存储

func HandleAddSignUp(c *gin.Context) {
	var req signUp.AddSignUpRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		middleware.Fail(c, serviceErr.RequestErr)
		return
	}
	if len(req.CaptchaId) == 0 || len(req.CaptchaValue) == 0 || !captcha.VerifyCaptcha(req.CaptchaId, req.CaptchaValue) {
		middleware.FailWithCode(c, 40201, "验证码错误")
		return
	}
	// TODO 访问限制
	//检查参数
	if len(req.MemberArr) != 3 && len(req.MemberArr) != 2 {
		middleware.FailWithCode(c, 40202, "队伍人数必须是3人或2人")
		return
	}
	if len(req.TeamName) > 90 {
		middleware.FailWithCode(c, 40203, "队伍名称过长")
		return
	}
	if len(req.School) == 0 && !req.IsHDU {
		middleware.FailWithCode(c, 40218, "请填写参赛学校")
		return
	}
	if len(req.School) > 90 {
		middleware.FailWithCode(c, 40203, "参赛学校过长")
		return
	}
	if len(req.TeamName) == 0 {
		middleware.FailWithCode(c, 40204, "队伍名称不能为空")
		return
	}
	if len(req.MemberArr) == 3 {
		if req.MemberArr[0].IDNumber == req.MemberArr[1].IDNumber || req.MemberArr[0].IDNumber == req.MemberArr[2].IDNumber {

			middleware.FailWithCode(c, 40216, fmt.Sprint("队员身份信息不能相同"))
			return
		}
	} else if len(req.MemberArr) == 2 {
		if req.MemberArr[0].IDNumber == req.MemberArr[1].IDNumber {
			middleware.FailWithCode(c, 40216, fmt.Sprint("队员身份信息不能相同"))
			return
		}
	}
	tempFlag := 0
	if req.MemberArr[2].IDNumber == "" && req.MemberArr[2].QQ == "" && req.MemberArr[2].Phone == "" && req.MemberArr[2].Name == "" && req.MemberArr[2].HDUID == "" {
		tempFlag = 1
	}
	for i, v := range req.MemberArr[:len(req.MemberArr)-tempFlag] {
		if len(v.Phone) == 0 {
			middleware.FailWithCode(c, 40206, fmt.Sprint(i+1, "号队员电话为空"))
			return
		}
		if !check.CheckMobile(v.Phone) {
			middleware.FailWithCode(c, 40205, fmt.Sprint(i+1, "号队员电话格式错误"))
			return
		}
		if len(v.QQ) > 30 {
			middleware.FailWithCode(c, 40207, fmt.Sprint(i+1, "号队员QQ号过长"))
			return
		}
		if len(v.QQ) == 0 {
			middleware.FailWithCode(c, 40208, fmt.Sprint(i+1, "号队员QQ号为空"))
			return
		}
		if len(v.Name) > 30 {
			middleware.FailWithCode(c, 40209, fmt.Sprint(i+1, "号队员名字过长"))
			return
		}
		if len(v.Name) == 0 {
			middleware.FailWithCode(c, 40210, fmt.Sprint(i+1, "号队员名字为空"))
			return
		}
		if check.CheckIdCard(v.IDNumber) == false {
			middleware.FailWithCode(c, 40211, fmt.Sprint(i+1, "号队员身份证号格式不正确"))
			return
		}
		var tempFlag bool
		err, tempFlag = SignUps.CheckIdNumberExist(v.IDNumber)
		if tempFlag {
			middleware.FailWithCode(c, 40217, fmt.Sprint(i+1, "号队员身份证号已报名"))
			return
		}
		if req.IsHDU {
			if len(v.HDUID) != 8 {
				middleware.FailWithCode(c, 40212, fmt.Sprint(i+1, "号队员学号格式不正确"))
				return
			}
			//} else { //默认非杭电队伍
			//	if len(v.BankCardNumber) > 30 {
			//		middleware.FailWithCode(c, 40213, fmt.Sprint(i+1, "号队员银行卡号过长"))
			//		return
			//	}
			//	if len(v.BankCardNumber) == 0 {
			//		middleware.FailWithCode(c, 40214, fmt.Sprint(i+1, "号队员银行卡号为空"))
			//		return
			//	}
			//	if len(v.BankName) > 50 {
			//		middleware.FailWithCode(c, 40215, fmt.Sprint(i+1, "号队员开户行名字过长"))
			//		return
			//	}
			//	if len(v.BankName) == 0 {
			//		middleware.FailWithCode(c, 40216, fmt.Sprint(i+1, "号队员开户行名字为空"))
			//		return
			//	}

		}
	}
	memberIDArr := make([]string, 3)
	for i, v := range req.MemberArr {
		if req.IsHDU {
			req.School = "杭州电子科技大学"
			err, memberIDArr[i] = SignUps.AddMember(&Mysql.Member{
				Phone:    v.Phone,
				QQ:       v.QQ,
				Name:     v.Name,
				IDNumber: v.IDNumber,
				HDUID:    v.HDUID,
				Role:     "队员",
			})
			if err != nil {
				middleware.Fail(c, serviceErr.InternalErr)
				return
			}
		} else {
			err, memberIDArr[i] = SignUps.AddMember(&Mysql.Member{
				Phone:          v.Phone,
				QQ:             v.QQ,
				Name:           v.Name,
				IDNumber:       v.IDNumber,
				BankCardNumber: "", //暂时不搜集银行卡信息
				BankName:       "",
				Role:           "队员",
			})
		}
	}
	var tempLen int64
	if err, tempLen = SignUps.AddSignUp(&Mysql.SignUp{
		TeamName:  req.TeamName,
		Teacher:   req.Teacher,
		IsHDU:     req.IsHDU,
		School:    req.School,
		Member1ID: memberIDArr[0],
		Member2ID: memberIDArr[1],
		Member3ID: memberIDArr[2],
	}); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	if tempLen < 100 {
		middleware.Success(c, "20220"+strconv.FormatInt(tempLen, 10))
		return
	} else {
		middleware.Success(c, "2022"+strconv.FormatInt(tempLen, 10))
		return
	}
}

func HandleGetAllSignUp(c *gin.Context) {
	var err error
	var res signUp.GetSignUpResponse
	var tempSignUpStruct []Mysql.SignUp
	if err, tempSignUpStruct = SignUps.GetAllSignUp(); err != nil {
		middleware.Fail(c, serviceErr.InternalErr)
		return
	}
	for _, v := range tempSignUpStruct {
		var tempMember1, tempMember2, tempMember3 Mysql.Member
		if err, tempMember1 = SignUps.GetMemberByID(v.Member1ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err, tempMember2 = SignUps.GetMemberByID(v.Member2ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		if err, tempMember3 = SignUps.GetMemberByID(v.Member3ID); err != nil {
			middleware.Fail(c, serviceErr.InternalErr)
			return
		}
		res.SignUpFormArr = append(res.SignUpFormArr, signUp.SignUpForm{
			TeamName: v.TeamName,
			IsHDU:    v.IsHDU,
			Teacher:  v.Teacher,
			School:   v.School,
			MemberArr: []signUp.Member{
				{
					Phone:    tempMember1.Phone,
					QQ:       tempMember1.QQ,
					Name:     tempMember1.Name,
					IDNumber: tempMember1.IDNumber,
					HDUID:    tempMember1.HDUID,
				}, {
					Phone:    tempMember2.Phone,
					QQ:       tempMember2.QQ,
					Name:     tempMember2.Name,
					IDNumber: tempMember2.IDNumber,
					HDUID:    tempMember2.HDUID,
				}, {
					Phone:    tempMember3.Phone,
					QQ:       tempMember3.QQ,
					Name:     tempMember3.Name,
					IDNumber: tempMember3.IDNumber,
					HDUID:    tempMember3.HDUID,
				},
			},
		})
	}
	res.Number = len(res.SignUpFormArr)
	middleware.Success(c, res)
	return
}
