package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"creswoodcornerscards/config"
)

type Store struct {
	mu            sync.RWMutex
	cfg           config.Config
	nextUserID    int64
	nextProductID int64
	nextLiveID    int64
	nextOrderID   int64
	nextItemID    int64
	users         map[int64]User
	usersByEmail  map[string]int64
	sessions      map[string]session
	products      map[int64]Product
	productBySlug map[string]int64
	live          LiveSession
	orders        map[int64]Order
}

type session struct {
	UserID    int64
	ExpiresAt time.Time
}

type User struct {
	ID        int64
	Email     string
	Password  string
	Name      string
	IsAdmin   bool
	CreatedAt time.Time
}

type Product struct {
	ID            int64
	Slug          string
	Title         string
	Summary       string
	Description   string
	Sport         string
	Player        string
	Team          string
	Brand         string
	SetName       string
	Year          int
	CardNumber    string
	Grade         string
	Condition     string
	PriceCents    int
	Quantity      int
	Accent        string
	Featured      bool
	LiveExclusive bool
	Status        string
}

type LiveSession struct {
	ID        int64
	Title     string
	Pitch     string
	Callout   string
	StreamURL string
	Platform  string
	IsActive  bool
}

type Order struct {
	ID                int64
	OrderNumber       string
	Email             string
	CustomerName      string
	Source            string
	PaymentStatus     string
	FulfillmentStatus string
	ShippingStatus    string
	TotalCents        int
	Notes             string
	CreatedAt         time.Time
	Items             []OrderItem
}

type OrderItem struct {
	ID           int64
	ProductID    int64
	ProductTitle string
	ProductSlug  string
	Quantity     int
	UnitPrice    int
}

type DashboardStats struct {
	ProductCount int
	OrderCount   int
	AccountCount int
	GrossSales   int
}

type CartLine struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func Open(cfg config.Config) (*Store, error) {
	store := &Store{
		cfg:           cfg,
		users:         map[int64]User{},
		usersByEmail:  map[string]int64{},
		sessions:      map[string]session{},
		products:      map[int64]Product{},
		productBySlug: map[string]int64{},
		orders:        map[int64]Order{},
	}
	if err := store.ensureBootstrapAdmin(context.Background()); err != nil {
		return nil, err
	}
	store.seedProducts()
	store.seedLiveSession()
	return store, nil
}

func (s *Store) Close() error {
	return nil
}

func (s *Store) ensureBootstrapAdmin(ctx context.Context) error {
	if s.cfg.AdminBootstrapEmail == "" {
		return nil
	}
	if user, err := s.FindUserByEmail(ctx, s.cfg.AdminBootstrapEmail); err == nil {
		user.IsAdmin = true
		s.mu.Lock()
		s.users[user.ID] = user
		s.mu.Unlock()
		return nil
	}
	_, err := s.CreateUser(ctx, "Store Owner", s.cfg.AdminBootstrapEmail, s.cfg.AdminBootstrapPassword)
	return err
}

func (s *Store) seedProducts() {
	products := []Product{
		{Slug: "jayden-daniels-prizm-rc-gold-wave", Title: "Jayden Daniels Prizm RC Gold Wave", Summary: "Low-number rookie color built to headline both the shop and live stream.", Description: "A flagship rookie listing with clear condition notes, sharp presentation, and immediate buying confidence.", Sport: "Football", Player: "Jayden Daniels", Team: "Washington Commanders", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "312", Grade: "Raw", Condition: "Near Mint", PriceCents: 32900, Quantity: 1, Accent: "gold", Featured: true, LiveExclusive: true, Status: "active"},
		{Slug: "cj-stroud-downtown-sgc-10", Title: "C.J. Stroud Downtown SGC 10", Summary: "Premium slab inventory to anchor the homepage and featured rails.", Description: "High-end inventory gives the storefront trust and average order value. This slab acts as a premium centerpiece.", Sport: "Football", Player: "C.J. Stroud", Team: "Houston Texans", Brand: "Donruss", SetName: "Downtown", Year: 2023, CardNumber: "DT-7", Grade: "SGC 10", Condition: "Gem Mint", PriceCents: 79900, Quantity: 1, Accent: "crimson", Featured: true, Status: "active"},
		{Slug: "bo-nix-mosaic-genesis-rc", Title: "Bo Nix Mosaic Genesis RC", Summary: "Mid-tier chase card with fast add-to-cart appeal.", Description: "A strong chase rookie that supports quick catalog browsing and live stream call-outs.", Sport: "Football", Player: "Bo Nix", Team: "Denver Broncos", Brand: "Panini Mosaic", SetName: "Mosaic Football", Year: 2024, CardNumber: "287", Grade: "Raw", Condition: "Near Mint-Mint", PriceCents: 18500, Quantity: 2, Accent: "emerald", Featured: true, Status: "active"},
		{Slug: "drake-maye-select-field-level-auto", Title: "Drake Maye Select Field Level Auto", Summary: "Autograph inventory that works for direct checkout and stream mentions.", Description: "A balanced autograph listing that supports both normal storefront conversion and live merchandising.", Sport: "Football", Player: "Drake Maye", Team: "New England Patriots", Brand: "Panini Select", SetName: "Select Football", Year: 2024, CardNumber: "FLA-DM", Grade: "Raw", Condition: "Near Mint", PriceCents: 26900, Quantity: 1, Accent: "blue", Status: "active"},
		{Slug: "malik-nabers-optic-holo-rc", Title: "Malik Nabers Optic Holo RC", Summary: "Popular rookie profile that keeps the grid feeling current.", Description: "Designed as a flexible inventory piece for homepage, catalog, and live modules.", Sport: "Football", Player: "Malik Nabers", Team: "New York Giants", Brand: "Donruss Optic", SetName: "Optic Football", Year: 2024, CardNumber: "214", Grade: "Raw", Condition: "Near Mint", PriceCents: 12900, Quantity: 3, Accent: "violet", Status: "active"},
		{Slug: "rome-odunze-prizm-silver-rc", Title: "Rome Odunze Prizm Silver RC", Summary: "Affordable parallel for a healthier catalog price spread.", Description: "Not every listing should be a grail. This card supports basket building and easier first purchases.", Sport: "Football", Player: "Rome Odunze", Team: "Chicago Bears", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "341", Grade: "Raw", Condition: "Near Mint", PriceCents: 8900, Quantity: 4, Accent: "slate", Status: "active"},
	}
	for _, product := range products {
		_ = s.UpsertProduct(context.Background(), product)
	}
}

func (s *Store) seedLiveSession() {
	s.live = LiveSession{
		ID:        1,
		Title:     "Friday Night Football Heat Check",
		Pitch:     "Live singles, rookie color, and quick-hit offers built to move inventory fast.",
		Callout:   "Running live now and mirrored from the storefront banner.",
		StreamURL: "https://www.whatnot.com/",
		Platform:  "Whatnot",
		IsActive:  true,
	}
	s.nextLiveID = 1
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.usersByEmail[normalizeEmail(email)]
	if !ok {
		return User{}, sql.ErrNoRows
	}
	return s.users[id], nil
}

func (s *Store) FindUserByID(ctx context.Context, id int64) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return User{}, sql.ErrNoRows
	}
	return user, nil
}

func (s *Store) CreateUser(ctx context.Context, name, email, password string) (User, error) {
	email = normalizeEmail(email)
	if email == "" {
		return User{}, errors.New("email is required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.usersByEmail[email]; exists {
		return User{}, errors.New("email already registered")
	}
	s.nextUserID++
	_, allowAdmin := s.cfg.AllowedAdminEmails[email]
	user := User{
		ID:        s.nextUserID,
		Email:     email,
		Password:  string(hash),
		Name:      strings.TrimSpace(name),
		IsAdmin:   allowAdmin || email == s.cfg.AdminBootstrapEmail,
		CreatedAt: time.Now(),
	}
	if user.Name == "" {
		user.Name = "Registered account"
	}
	s.users[user.ID] = user
	s.usersByEmail[email] = user.ID
	return user, nil
}

func (s *Store) AuthenticateUser(ctx context.Context, email, password string) (User, error) {
	user, err := s.FindUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return User{}, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *Store) CreateSession(ctx context.Context, userID int64) (string, error) {
	token := uuid.NewString()
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.users[userID]; !ok {
		return "", sql.ErrNoRows
	}
	s.sessions[token] = session{UserID: userID, ExpiresAt: time.Now().Add(30 * 24 * time.Hour)}
	return token, nil
}

func (s *Store) FindUserBySession(ctx context.Context, token string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[token]
	if !ok || time.Now().After(session.ExpiresAt) {
		return User{}, sql.ErrNoRows
	}
	user, ok := s.users[session.UserID]
	if !ok {
		return User{}, sql.ErrNoRows
	}
	return user, nil
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
	return nil
}

func (s *Store) FeaturedProducts(ctx context.Context) ([]Product, error) {
	products, err := s.Products(ctx)
	if err != nil {
		return nil, err
	}
	var featured []Product
	for _, product := range products {
		if product.Featured {
			featured = append(featured, product)
		}
	}
	if len(featured) > 4 {
		featured = featured[:4]
	}
	return featured, nil
}

func (s *Store) Products(ctx context.Context) ([]Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	products := make([]Product, 0, len(s.products))
	for _, product := range s.products {
		if product.Status != "archived" {
			products = append(products, product)
		}
	}
	sort.Slice(products, func(i, j int) bool {
		if products[i].Featured != products[j].Featured {
			return products[i].Featured
		}
		return products[i].ID < products[j].ID
	})
	return products, nil
}

func (s *Store) ProductBySlug(ctx context.Context, slug string) (Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.productBySlug[slug]
	if !ok {
		return Product{}, sql.ErrNoRows
	}
	return s.products[id], nil
}

func (s *Store) ProductByID(ctx context.Context, id int64) (Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	product, ok := s.products[id]
	if !ok {
		return Product{}, sql.ErrNoRows
	}
	return product, nil
}

func (s *Store) ActiveLiveSession(ctx context.Context) (LiveSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.live, nil
}

func (s *Store) SaveLiveSession(ctx context.Context, live LiveSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if live.ID == 0 {
		live.ID = s.live.ID
	}
	if live.ID == 0 {
		s.nextLiveID++
		live.ID = s.nextLiveID
	}
	s.live = live
	return nil
}

func (s *Store) UpsertProduct(ctx context.Context, product Product) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if product.Slug == "" || product.Title == "" {
		return errors.New("product slug and title are required")
	}
	if product.Sport == "" {
		product.Sport = "Football"
	}
	if product.Status == "" {
		product.Status = "active"
	}
	if product.Accent == "" {
		product.Accent = "gold"
	}
	if product.ID == 0 {
		s.nextProductID++
		product.ID = s.nextProductID
	} else if existing, ok := s.products[product.ID]; ok && existing.Slug != product.Slug {
		delete(s.productBySlug, existing.Slug)
	}
	s.products[product.ID] = product
	s.productBySlug[product.Slug] = product.ID
	return nil
}

func (s *Store) CreateOrder(ctx context.Context, email, customerName string, cart []CartLine) (Order, error) {
	if strings.TrimSpace(email) == "" {
		return Order{}, errors.New("email is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var items []OrderItem
	total := 0
	for _, line := range cart {
		product, ok := s.products[line.ProductID]
		if !ok {
			return Order{}, errors.New("product not found")
		}
		if line.Quantity < 1 {
			line.Quantity = 1
		}
		s.nextItemID++
		items = append(items, OrderItem{
			ID:           s.nextItemID,
			ProductID:    product.ID,
			ProductTitle: product.Title,
			ProductSlug:  product.Slug,
			Quantity:     line.Quantity,
			UnitPrice:    product.PriceCents,
		})
		total += product.PriceCents * line.Quantity
	}
	if len(items) == 0 {
		return Order{}, errors.New("cart is empty")
	}
	s.nextOrderID++
	order := Order{
		ID:                s.nextOrderID,
		OrderNumber:       fmt.Sprintf("CCC-%d", time.Now().UnixNano()/1_000_000),
		Email:             normalizeEmail(email),
		CustomerName:      strings.TrimSpace(customerName),
		Source:            "web_demo",
		PaymentStatus:     "paid",
		FulfillmentStatus: "ready_to_pack",
		ShippingStatus:    "pending",
		TotalCents:        total,
		CreatedAt:         time.Now(),
		Items:             items,
	}
	s.orders[order.ID] = order
	return order, nil
}

func (s *Store) OrderByID(ctx context.Context, id int64) (Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	order, ok := s.orders[id]
	if !ok {
		return Order{}, sql.ErrNoRows
	}
	return order, nil
}

func (s *Store) Orders(ctx context.Context) ([]Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]Order, 0, len(s.orders))
	for _, order := range s.orders {
		orders = append(orders, order)
	}
	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.After(orders[j].CreatedAt)
	})
	return orders, nil
}

func (s *Store) UpdateOrder(ctx context.Context, orderID int64, payment, fulfillment, shipping, notes string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	order, ok := s.orders[orderID]
	if !ok {
		return sql.ErrNoRows
	}
	order.PaymentStatus = payment
	order.FulfillmentStatus = fulfillment
	order.ShippingStatus = shipping
	order.Notes = notes
	s.orders[orderID] = order
	return nil
}

func (s *Store) Users(ctx context.Context) ([]User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})
	return users, nil
}

func (s *Store) DashboardStats(ctx context.Context) (DashboardStats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stats := DashboardStats{
		ProductCount: len(s.products),
		OrderCount:   len(s.orders),
		AccountCount: len(s.users),
	}
	for _, order := range s.orders {
		stats.GrossSales += order.TotalCents
	}
	return stats, nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func UniqueTeams(products []Product) []string {
	seen := map[string]struct{}{}
	var teams []string
	for _, product := range products {
		if product.Team == "" {
			continue
		}
		if _, ok := seen[product.Team]; ok {
			continue
		}
		seen[product.Team] = struct{}{}
		teams = append(teams, product.Team)
	}
	sort.Strings(teams)
	return teams
}
