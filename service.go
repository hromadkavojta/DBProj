package DBproj

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"

	"github.com/labstack/echo"
	"net/http"
	"regexp"
)

type Service struct {
	access IAccesser
}

func NewService(access IAccesser) *Service {
	return &Service{
		access: access,
	}
}

func validateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return re.MatchString(email)
}

func (s *Service) SignUp(ctx echo.Context) error {

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	//Validate email
	if !validateEmail(usr.Email) {
		fmt.Println("email not valid")
		return err
	}

	s.access.CreateUser(usr)
	return nil
}

func (s *Service) SignIn(ctx echo.Context) error {

	//return login token

	var usr User

	err := ctx.Bind(&usr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	query := s.access.GetUser(usr)
	if query == (User{}) {
		return ctx.JSON(http.StatusBadRequest, "user does not exists")
	}

	//create token

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = query.Email
	claims["id"] = query.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (s *Service) SetTask(ctx echo.Context) error {

	var tsk Task

	err := ctx.Bind(&tsk)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	tsk.UserId = ctx.Get("id").(int)

	s.access.CreateTask(tsk)

	return nil
}

func (s *Service) GetTaskByUserId(ctx echo.Context) error {

	query := s.access.GetTask(ctx.Get("id").(int))

	//var results[] Task
	//for _,v := range query{
	//	results = append(results, Task{
	//		Name:v.Name,
	//		ExecuteAt:v.ExecuteAt,
	//	})
	//}

	return ctx.JSON(http.StatusOK, query)
}
