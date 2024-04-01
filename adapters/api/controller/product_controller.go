package controller

import (
	"encoding/json"
	"errors"
	"golang-api-hexagonal/adapters/api/dto"
	"golang-api-hexagonal/adapters/api/middleware"
	"golang-api-hexagonal/adapters/api/router"
	"golang-api-hexagonal/adapters/opa"
	"golang-api-hexagonal/core/domain"
	"golang-api-hexagonal/core/ports"
	"net/http"

	"github.com/go-chi/chi/v5"
	serverMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// ProductController controller for product API
type ProductController struct {
	log           *zap.SugaredLogger
	counterMetric prometheus.Counter
	validate      *validator.Validate
	service       ports.IProductService
	jwtVerify     *middleware.JWTVerify
	policyService *opa.PolicyService
}

// NewProductController Create a new http product controller API
func NewProductController(httpRouter *router.HTTPRouter, log *zap.SugaredLogger, validator *validator.Validate, prometheusRegistry *middleware.CustomMetricRegistry,
	service ports.IProductService, jwtVerify *middleware.JWTVerify, policyService *opa.PolicyService) {
	controller := &ProductController{
		log:           log,
		validate:      validator,
		service:       service,
		jwtVerify:     jwtVerify,
		policyService: policyService,
		counterMetric: promauto.With(prometheusRegistry).NewCounter(prometheus.CounterOpts{
			Name: "products_reqs_total",
			Help: "The total number of request for products endpoints",
		}),
	}

	httpRouter.Router.Group(func(r chi.Router) {
		r.Use(controller.jwtVerify.JWTVerifyHandler())
		r.Post("/v1/product", controller.createProduct)
		r.Get("/v1/product/{id}", controller.getProduct)
	})
}

// createProduct create the product
func (pc *ProductController) createProduct(writer http.ResponseWriter, request *http.Request) {
	pc.counterMetric.Inc()
	traceID := request.Context().Value(serverMiddleware.RequestIDKey).(string)
	claims := request.Context().Value(domain.ClaimsKey).(domain.AuthClaims)
	pc.log.With("traceId", traceID).Infof("User %v is creating a product.", claims.Username)

	allowed := pc.policyService.EvaluateApiPolicy(request.Context(), claims, "createProduct", "")
	if !allowed {
		pc.log.With("traceId", traceID).Errorf("Forbidden access role")
		dto.RenderErrorResponse(request.Context(), writer, http.StatusForbidden, errors.New("forbidden access"))
		return
	}

	productRequest := &domain.Product{}
	err := json.NewDecoder(request.Body).Decode(productRequest)
	if err != nil {
		pc.log.With("traceId", traceID).Errorf("Error to parsing the product payload body. Maformed: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	_ = pc.validate.RegisterValidation("not_blank", validators.NotBlank)
	err = pc.validate.Struct(productRequest)
	if err != nil {
		pc.log.With("traceId", traceID).Errorf("Product validation error: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	response, err := pc.service.CreateProduct(request.Context(), productRequest, claims.Username, traceID)
	if err != nil {
		dto.RenderErrorResponse(request.Context(), writer, 0, err)
		return
	}
	dto.RenderResponse(request.Context(), writer, http.StatusCreated, response)
}

// getProduct get the product by id
func (pc *ProductController) getProduct(writer http.ResponseWriter, request *http.Request) {
	pc.counterMetric.Inc()
	traceID := request.Context().Value(serverMiddleware.RequestIDKey).(string)
	claims := request.Context().Value(domain.ClaimsKey).(domain.AuthClaims)
	id := chi.URLParam(request, "id")
	pc.log.With("traceId", traceID).Infof("User %v is searching a product.", claims.Username)

	allowed := pc.policyService.EvaluateApiPolicy(request.Context(), claims, "viewProduct", "")
	if !allowed {
		pc.log.With("traceId", traceID).Errorf("Forbidden access role")
		dto.RenderErrorResponse(request.Context(), writer, http.StatusForbidden, errors.New("forbidden access"))
		return
	}

	response, err := pc.service.GetProduct(request.Context(), id, traceID)
	if err != nil {
		dto.RenderErrorResponse(request.Context(), writer, 0, err)
		return
	}
	dto.RenderResponse(request.Context(), writer, http.StatusOK, response)
}
