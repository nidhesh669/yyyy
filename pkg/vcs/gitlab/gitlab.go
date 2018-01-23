package gitlab

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/while-loop/todo/pkg/issue"
	"github.com/while-loop/todo/pkg/log"
	"github.com/while-loop/todo/pkg/vcs/config"
)

const (
	name = "gitlab"
)

type Service struct {
	router   *mux.Router
	glClient interface{}
	issueCh  <-chan []*issue.Issue
}

func NewService(config *config.GitlabConfig, issueChan <-chan []*issue.Issue) *Service {
	//ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: app.Config.Gitlab.AccessToken})
	//oauthClient := oauth2.NewClient(context.Background(), ts)
	s := &Service{
		issueCh:  issueChan,
		router:   mux.NewRouter(),
		glClient: nil,
	}

	return s
}

func (s *Service) GetRepository(user, project string) error {
	panic("implement me")
}

func (s *Service) Name() string {
	return name
}

func (s *Service) Handler() http.Handler {
	s.router.HandleFunc("/webhook/"+name, s.webhook)
	return s.router
}

func (s *Service) webhook(w http.ResponseWriter, r *http.Request) {
	log.Info(name, r.URL, w, r)
}
