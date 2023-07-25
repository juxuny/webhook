package main

import (
	"fmt"
	"github.com/juxuny/webhook/config"
	"github.com/juxuny/webhook/executor"
	"log"
	"net/http"
	"os"
	"strings"
)

type Handler struct {
	logger executor.Logger
}

func NewHandler() http.Handler {
	return &Handler{
		logger: NewDefaultLogger(),
	}
}

func (t *Handler) resp(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func (t *Handler) respFail(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func (t *Handler) respSuccess(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
}

func (t *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		_ = r.Body.Close()
	}()

	// deny all request methods except POST
	if r.Method != http.MethodPost {
		log.Println("method not allow: ", r.Method)
		t.resp(w, http.StatusMethodNotAllowed)
		return
	}
	initConfig()
	// validate request token
	token := r.FormValue("token")
	if token == "" {
		t.resp(w, http.StatusUnauthorized)
		return
	}
	auth, found := config.Get().GetAuthByToken(token)
	if !found {
		log.Println("token not found:", token)
		t.resp(w, http.StatusUnauthorized)
		return
	}
	log.Println("authorized:", auth.Name)

	// check deployment config
	log.Println(r.RequestURI)
	deploymentName := strings.Replace(r.URL.Path, config.Get().GetHookPrefix(), "", 1)
	deploymentName = strings.Trim(deploymentName, "/")
	if strings.Contains(deploymentName, "/") {
		t.respFail(w, "invalid deployment name: "+deploymentName)
		return
	}
	log.Printf("deployment name: %s\n", deploymentName)
	deployment, found := config.Get().GetDeployment(deploymentName)
	if !found {
		t.respFail(w, fmt.Sprintf("not found: %s", deploymentName))
		return
	}
	log.Println("deployment found:", deployment.Name)
	log.Println("working dir: ", deployment.WorkDir)

	// parse parameters
	params := make(map[string]string)
	err := r.ParseForm()
	if err != nil {
		t.logger.Println("parse parameters failed:", err)
		t.respFail(w, err.Error())
		return
	}
	for _, variable := range deployment.Variables {
		value := r.PostFormValue(variable.Name)
		if variable.IsRequired() && value == "" {
			t.respFail(w, "missing parameter: "+variable.Name)
			return
		}
		if err = config.VariableValidate(value, variable.Type); err != nil {
			t.respFail(w, err.Error())
			return
		}
		params[variable.Name] = value
	}

	// build bash executor
	builder := executor.NewBashBuilder()
	builder.SetWorkdir(deployment.WorkDir).SetLogger(t.logger).SetOutput(os.Stdout).SetScripts(deployment.Scripts).AddVariables(params)
	err = builder.Build().Exec()
	if err != nil {
		t.respFail(w, err.Error())
		return
	}

	t.respSuccess(w, "success")
	log.Println("success !")
}
