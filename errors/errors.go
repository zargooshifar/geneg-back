package errors

//TODO: personilize error messages!
const (
	ACCESS_DENIED              = "شما دسترسی لازم را ندارید!"
	ERROR_WRONG_PASSWORD       = "کلمه عبور اشتباه است!"
	WRONG_INPUT                = "ورودی اشتباه!"
	WRONG_VERIFICATION_ID      = "خطا در تایید هویت!"
	WRONG_PIN                  = "کد وارد شده صحیح نمی باشد!"
	EXPIRE_PIN                 = "کد تایید منقضی شده است."
	EMPTY_VERIFICATION         = "کد تایید هویت وارد نشده است."
	EMPTY_FIRST_NAME           = "نام را وارد کنید."
	EMPTY_LAST_NAME            = "نام خانوادگی را وارد کنید."
	USER_EXIST                 = "کاربر با این شماره وجود دارد."
	USER_NOT_EXIST             = "کاربری با این شماره یافت نشد."
	VERIFICATION_NOT_EXIST     = "کد تایید هویت یافت نشد."
	VERIFICATION_NOT_CONFIRMED = "هویت تایید نشده است."
	DB_ERROR_SAVING            = "خطا در ذخبره سازی اطلاعات"
	REFRESH_TOKEN_NOT_EXIST    = "توکن یافت نشد!"
	CANT_HANDLE_REFRESH_TOKEN  = "توکن معتبر نیست."
	REFRESH_TOKEN_UNVALID      = "توکن معتبر نیست."
	REFRESH_TOKEN_NOT_A_TOKEN  = "توکن معتبر نیست."
	REFRESH_TOKEN_EXPIRED      = "توکن باطل شده است!"
	EXPIRE_FOOD                = "این آیتم اکسپابر شده است!"
	NO_RESERVE                 = "رزروی یافت نشد!"
	NOT_SINGLE_FACE            = "NOT_SINGLE_FACE"
	NOT_TRAINED                = "NOT_TRAINED"
	CANT_REGISTER_OUTSIDE_CORP = "ثبت نام نکرده اید. ثبت نام در سایت فقط از طریق اینترنت مجموعه امکان پذیر است."
)
