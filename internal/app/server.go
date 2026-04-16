package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"creswoodcornerscards"
	"creswoodcornerscards/internal/auth"
	"creswoodcornerscards/internal/config"
	"creswoodcornerscards/internal/data"
	"creswoodcornerscards/internal/web"
)

type Server struct {
	cfg      config.Config
	store    *data.Store
	auth     *auth.Manager
	renderer *web.Renderer
	mux      *http.ServeMux
}

type ViewData struct {
	Title            string
	AppName          string
	Path             string
	CurrentUser      *data.User
	Products         []data.Product
	Product          data.Product
	Featured         []data.Product
	Live             data.LiveSession
	Orders           []data.Order
	Users            []data.User
	Stats            data.DashboardStats
	Teams            []string
	Order            data.Order
	Flash            string
	ClerkEnabled     bool
	StripeEnabled    bool
	BootstrapEmail   string
	ReferenceWebsite string
}

func New(cfg config.Config, store *data.Store) (*Server, error) {
	templateFS, err := fs.Sub(creswoodcornerscards.Assets, "web/templates")
	if err != nil {
		return nil, err
	}
	renderer, err := web.NewRenderer(templateFS)
	if err != nil {
		return nil, err
	}

	srv := &Server{
		cfg:      cfg,
		store:    store,
		auth:     auth.NewManager(cfg, store),
		renderer: renderer,
		mux:      http.NewServeMux(),
	}
	srv.routes()
	return srv, nil
}

func (s *Server) routes() {
	staticFS, err := fs.Sub(creswoodcornerscards.Assets, "web/static")
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(staticFS))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	s.mux.HandleFunc("/", s.home)
	s.mux.HandleFunc("/shop", s.shop)
	s.mux.HandleFunc("/shop/", s.product)
	s.mux.HandleFunc("/live", s.live)
	s.mux.HandleFunc("/cart", s.cartPage)
	s.mux.HandleFunc("/checkout", s.checkout)
	s.mux.HandleFunc("/checkout/success", s.checkoutSuccess)
	s.mux.HandleFunc("/sign-in", s.signIn)
	s.mux.HandleFunc("/sign-up", s.signUp)
	s.mux.HandleFunc("/sign-out", s.signOut)
	s.mux.HandleFunc("/account", s.account)
	s.mux.HandleFunc("/admin", s.admin)
	s.mux.HandleFunc("/admin/product/save", s.saveProduct)
	s.mux.HandleFunc("/admin/live/save", s.saveLive)
	s.mux.HandleFunc("/admin/order/update", s.updateOrder)
}

func (s *Server) Handler() http.Handler {
	return loggingMiddleware(s.mux)
}

func (s *Server) baseData(r *http.Request, title string) ViewData {
	user, _ := s.auth.CurrentUser(r)
	return ViewData{
		Title:            title,
		AppName:          s.cfg.AppName,
		Path:             r.URL.Path,
		CurrentUser:      user,
		ClerkEnabled:     s.cfg.ClerkSecretKey != "" && s.cfg.ClerkPublishableKey != "",
		StripeEnabled:    s.cfg.StripeSecretKey != "" && s.cfg.StripePublishableKey != "",
		BootstrapEmail:   s.cfg.AdminBootstrapEmail,
		ReferenceWebsite: "https://www.logans3dcreations.com/",
	}
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	ctx := r.Context()
	featured, _ := s.store.FeaturedProducts(ctx)
	products, _ := s.store.Products(ctx)
	live, _ := s.store.ActiveLiveSession(ctx)

	data := s.baseData(r, "Home")
	data.Featured = featured
	data.Products = products
	data.Live = live
	data.Teams = dataPkg.UniqueTeams(products)
	s.renderer.HTML(w, "home", data)
}

func (s *Server) shop(w http.ResponseWriter, r *http.Request) {
	products, _ := s.store.Products(r.Context())
	data := s.baseData(r, "Shop")
	data.Products = products
	data.Teams = dataPkg.UniqueTeams(products)
	s.renderer.HTML(w, "shop", data)
}

func (s *Server) product(w http.ResponseWriter, r *http.Request) {
	slug := strings.TrimPrefix(r.URL.Path, "/shop/")
	if slug == "" || slug == "/" {
		http.Redirect(w, r, "/shop", http.StatusSeeOther)
		return
	}
	product, err := s.store.ProductBySlug(r.Context(), slug)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data := s.baseData(r, product.Title)
	data.Product = product
	s.renderer.HTML(w, "product", data)
}

func (s *Server) live(w http.ResponseWriter, r *http.Request) {
	products, _ := s.store.Products(r.Context())
	live, _ := s.store.ActiveLiveSession(r.Context())
	var featured []data.Product
	for _, product := range products {
		if product.LiveExclusive {
			featured = append(featured, product)
		}
	}
	data := s.baseData(r, "Live Selling")
	data.Live = live
	data.Products = featured
	s.renderer.HTML(w, "live", data)
}

func (s *Server) cartPage(w http.ResponseWriter, r *http.Request) {
	data := s.baseData(r, "Cart")
	s.renderer.HTML(w, "cart", data)
}

type checkoutRequest struct {
	Email        string          `json:"email"`
	CustomerName string          `json:"customer_name"`
	Items        []data.CartLine `json:"items"`
}

func (s *Server) checkout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req checkoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid checkout payload", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(req.Email) == "" || len(req.Items) == 0 {
		http.Error(w, "email and items are required", http.StatusBadRequest)
		return
	}
	order, err := s.store.CreateOrder(r.Context(), req.Email, req.CustomerName, req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	respondJSON(w, map[string]any{
		"redirect": fmt.Sprintf("/checkout/success?id=%d", order.ID),
	})
}

func (s *Server) checkoutSuccess(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	order, _ := s.store.OrderByID(r.Context(), id)
	data := s.baseData(r, "Order Received")
	data.Order = order
	s.renderer.HTML(w, "success", data)
}

func (s *Server) signIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.renderer.HTML(w, "sign_in", s.baseData(r, "Sign In"))
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	user, err := s.store.AuthenticateUser(r.Context(), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		data := s.baseData(r, "Sign In")
		data.Flash = "Invalid email or password."
		s.renderer.HTML(w, "sign_in", data)
		return
	}
	if err := s.auth.SignIn(r.Context(), w, user.ID); err != nil {
		http.Error(w, "unable to sign in", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (s *Server) signUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.renderer.HTML(w, "sign_up", s.baseData(r, "Sign Up"))
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	user, err := s.store.CreateUser(r.Context(), r.FormValue("name"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		data := s.baseData(r, "Sign Up")
		data.Flash = "Could not create account. That email may already be registered."
		s.renderer.HTML(w, "sign_up", data)
		return
	}
	if err := s.auth.SignIn(r.Context(), w, user.ID); err != nil {
		http.Error(w, "unable to start session", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/account", http.StatusSeeOther)
}

func (s *Server) signOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	_ = s.auth.SignOut(r.Context(), w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) account(w http.ResponseWriter, r *http.Request) {
	user, err := s.requireUser(w, r)
	if err != nil {
		return
	}
	data := s.baseData(r, "Account")
	data.CurrentUser = user
	s.renderer.HTML(w, "account", data)
}

func (s *Server) admin(w http.ResponseWriter, r *http.Request) {
	user, err := s.requireAdmin(w, r)
	if err != nil {
		return
	}

	ctx := r.Context()
	products, _ := s.store.Products(ctx)
	users, _ := s.store.Users(ctx)
	orders, _ := s.store.Orders(ctx)
	live, _ := s.store.ActiveLiveSession(ctx)
	stats, _ := s.store.DashboardStats(ctx)

	data := s.baseData(r, "Admin")
	data.CurrentUser = user
	data.Products = products
	data.Users = users
	data.Orders = orders
	data.Live = live
	data.Stats = stats
	s.renderer.HTML(w, "admin", data)
}

func (s *Server) saveProduct(w http.ResponseWriter, r *http.Request) {
	_, err := s.requireAdmin(w, r)
	if err != nil {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	year, _ := strconv.Atoi(r.FormValue("year"))
	price, _ := strconv.Atoi(r.FormValue("price_cents"))
	quantity, _ := strconv.Atoi(r.FormValue("quantity"))

	err = s.store.UpsertProduct(r.Context(), data.Product{
		ID:            id,
		Slug:          strings.TrimSpace(r.FormValue("slug")),
		Title:         strings.TrimSpace(r.FormValue("title")),
		Summary:       strings.TrimSpace(r.FormValue("summary")),
		Description:   strings.TrimSpace(r.FormValue("description")),
		Sport:         strings.TrimSpace(r.FormValue("sport")),
		Player:        strings.TrimSpace(r.FormValue("player")),
		Team:          strings.TrimSpace(r.FormValue("team")),
		Brand:         strings.TrimSpace(r.FormValue("brand")),
		SetName:       strings.TrimSpace(r.FormValue("set_name")),
		Year:          year,
		CardNumber:    strings.TrimSpace(r.FormValue("card_number")),
		Grade:         strings.TrimSpace(r.FormValue("grade")),
		Condition:     strings.TrimSpace(r.FormValue("condition")),
		PriceCents:    price,
		Quantity:      quantity,
		Accent:        strings.TrimSpace(r.FormValue("accent")),
		Featured:      r.FormValue("featured") == "on",
		LiveExclusive: r.FormValue("live_exclusive") == "on",
		Status:        strings.TrimSpace(r.FormValue("status")),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) saveLive(w http.ResponseWriter, r *http.Request) {
	_, err := s.requireAdmin(w, r)
	if err != nil {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	err = s.store.SaveLiveSession(r.Context(), data.LiveSession{
		Title:     r.FormValue("title"),
		Pitch:     r.FormValue("pitch"),
		Callout:   r.FormValue("callout"),
		StreamURL: r.FormValue("stream_url"),
		Platform:  r.FormValue("platform"),
		IsActive:  r.FormValue("is_active") == "on",
	})
	if err != nil {
		http.Error(w, "could not save live session", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) updateOrder(w http.ResponseWriter, r *http.Request) {
	_, err := s.requireAdmin(w, r)
	if err != nil {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	orderID, _ := strconv.ParseInt(r.FormValue("order_id"), 10, 64)
	err = s.store.UpdateOrder(r.Context(), orderID, r.FormValue("payment_status"), r.FormValue("fulfillment_status"), r.FormValue("shipping_status"), r.FormValue("notes"))
	if err != nil {
		http.Error(w, "could not update order", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (s *Server) requireUser(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	user, err := s.auth.CurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return nil, errors.New("authentication required")
	}
	return user, nil
}

func (s *Server) requireAdmin(w http.ResponseWriter, r *http.Request) (*data.User, error) {
	user, err := s.requireUser(w, r)
	if err != nil {
		return nil, err
	}
	if !user.IsAdmin {
		http.Error(w, "admin access required", http.StatusForbidden)
		return nil, errors.New("admin required")
	}
	return user, nil
}

func respondJSON(w http.ResponseWriter, payload any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

var dataPkg = struct {
	UniqueTeams func([]data.Product) []string
}{
	UniqueTeams: data.UniqueTeams,
}
