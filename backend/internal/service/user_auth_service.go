package service

type CodeStore interface {
	Save(phone, code string) error
	Verify(phone, code string) bool
}

type UserAuthService struct{ store CodeStore }

func NewUserAuthService(store CodeStore) *UserAuthService { return &UserAuthService{store: store} }

func (s *UserAuthService) SendSMS(phone, code string) error {
	return s.store.Save(phone, code)
}

func (s *UserAuthService) VerifySMS(phone, code string) bool {
	return s.store.Verify(phone, code)
}

func (s *UserAuthService) WechatLogin(openID string) string {
	return openID
}
