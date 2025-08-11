package user

type service struct {
	repo Repository
}

func (s *service) GetUserByID(id uint) (*User, error) {

	user, err := s.repo.GetById(uint(id))
	if err != nil {
		return &User{}, err
	}
	return &user, nil
}

func (s *service) CreateUser(user *User) error {

	return s.repo.CreateUser(user)
}

func (s *service) GetUserByEmailAndProvider(email string, provider SocialProvider) (*User, error) {
	user, err := s.repo.GetByEmailAndProvider(email, provider)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
