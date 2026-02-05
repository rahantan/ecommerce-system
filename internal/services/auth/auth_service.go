package authservices

// type AuthServiceImpl struct {
// 	userservices.UserService
// 	userrepositories.UserRepositories
// 	*gorm.
// }

// func NewAuthService(user userservices.UserService) AuthServices {
// 	return &AuthServiceImpl{
// 		UserService: user,
// 	}
// }
// func (auth *AuthServiceImpl) Login(req *request.ReqLogin) (*response.ResUser, error) {
// 	password, err := auth.UserService.GetUserPasswordByEmail(req.Email)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
// 		return nil, exceptions.ErrCustomInvalidCredential
// 	}

// 	return auth.UserService.GetUserByEmail(req.Email)
// }
// func (auth *AuthServiceImpl) Register(req *request.ReqCreateUser) (*response.ResUser, error) {

// 	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}
// 	req.Password = string(hash)
// 	result, err := auth.UserService.Create(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return result, nil
// }
