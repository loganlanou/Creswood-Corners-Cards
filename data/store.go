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
	PlayerName    string
	Team          string
	TeamLogoURL   string
	TeamLogoAlt   string
	Brand         string
	SetName       string
	Year          int
	CardNumber    string
	Grade         string
	Condition     string
	CardType      string
	PriceCents    int
	Quantity      int
	Accent        string
	ImageURL      string
	Featured      bool
	LivePriority  bool
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

type LiveEvent struct {
	Name          string
	DateTime      string
	Platform      string
	Description   string
	FeaturedCards string
	Status        string
	LinkURL       string
}

type SalesChannel struct {
	Name     string
	Type     string
	Schedule string
	Notes    string
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
	SubtotalCents     int
	DiscountCents     int
	TaxCents          int
	TotalCents        int
	CouponCode        string
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
		{Slug: "jayden-daniels-prizm-rc-gold-wave", Title: "Jayden Daniels Prizm RC Gold Wave Raw", Summary: "Low-number rookie color built to headline both the shop and live stream.", Description: "A flagship raw rookie listing with clear condition notes, sharp presentation, and immediate buying confidence.", Sport: "Football", PlayerName: "Jayden Daniels", Team: "Washington Commanders", TeamLogoAlt: "Washington Commanders logo", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "312", Grade: "Raw", Condition: "Near Mint", CardType: "Rookie Parallel", PriceCents: 32900, Quantity: 1, Accent: "gold", Featured: true, LivePriority: true, LiveExclusive: true, Status: "active"},
		{Slug: "cj-stroud-downtown-raw", Title: "C.J. Stroud Downtown Raw", Summary: "Short-print insert built to anchor Featured Drops without implying grading.", Description: "A raw Downtown insert placeholder with condition-focused copy, strong photography support, and a clear path to checkout.", Sport: "Football", PlayerName: "C.J. Stroud", Team: "Houston Texans", TeamLogoAlt: "Houston Texans logo", Brand: "Donruss", SetName: "Downtown", Year: 2023, CardNumber: "DT-7", Grade: "Raw", Condition: "Near Mint", CardType: "Insert", PriceCents: 79900, Quantity: 1, Accent: "crimson", Featured: true, Status: "active"},
		{Slug: "bo-nix-mosaic-genesis-rc", Title: "Bo Nix Mosaic Genesis RC Raw", Summary: "Mid-tier chase card with fast add-to-cart appeal.", Description: "A strong raw chase rookie that supports quick catalog browsing and live stream call-outs.", Sport: "Football", PlayerName: "Bo Nix", Team: "Denver Broncos", TeamLogoAlt: "Denver Broncos logo", Brand: "Panini Mosaic", SetName: "Mosaic Football", Year: 2024, CardNumber: "287", Grade: "Raw", Condition: "Near Mint", CardType: "Rookie Parallel", PriceCents: 18500, Quantity: 2, Accent: "emerald", Featured: true, Status: "active"},
		{Slug: "drake-maye-select-field-level-auto", Title: "Drake Maye Select Field Level Auto Raw", Summary: "Autograph inventory that works for direct checkout and stream mentions.", Description: "A balanced raw autograph listing that supports both normal storefront conversion and live merchandising.", Sport: "Football", PlayerName: "Drake Maye", Team: "New England Patriots", TeamLogoAlt: "New England Patriots logo", Brand: "Panini Select", SetName: "Select Football", Year: 2024, CardNumber: "FLA-DM", Grade: "Raw", Condition: "Near Mint", CardType: "Auto", PriceCents: 26900, Quantity: 1, Accent: "blue", LivePriority: true, Status: "active"},
		{Slug: "malik-nabers-optic-holo-rc", Title: "Malik Nabers Optic Holo RC Raw", Summary: "Popular rookie profile that keeps the grid feeling current.", Description: "Designed as a flexible raw rookie inventory piece for homepage, catalog, and live modules.", Sport: "Football", PlayerName: "Malik Nabers", Team: "New York Giants", TeamLogoAlt: "New York Giants logo", Brand: "Donruss Optic", SetName: "Optic Football", Year: 2024, CardNumber: "214", Grade: "Raw", Condition: "Near Mint", CardType: "Rookie Parallel", PriceCents: 12900, Quantity: 3, Accent: "violet", Status: "active"},
		{Slug: "rome-odunze-prizm-silver-rc", Title: "Rome Odunze Prizm Silver RC Raw", Summary: "Affordable parallel for a healthier catalog price spread.", Description: "Not every listing should be a grail. This raw card supports basket building and easier first purchases.", Sport: "Football", PlayerName: "Rome Odunze", Team: "Chicago Bears", TeamLogoAlt: "Chicago Bears logo", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "341", Grade: "Raw", Condition: "Near Mint", CardType: "Rookie Parallel", PriceCents: 8900, Quantity: 4, Accent: "slate", Status: "active"},
	}
	for _, product := range products {
		_ = s.UpsertProduct(context.Background(), product)
	}
}

func (s *Store) seedLiveSession() {
	s.live = LiveSession{
		ID:       1,
		Title:    "Rookie Rush Night",
		Pitch:    "Live singles, rookie color, autos, and quick-hit offers built for early-stage selling.",
		Callout:  "Next live link will appear here when a stream platform is configured.",
		Platform: "Platform TBD",
		IsActive: false,
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

func (s *Store) LiveEvents(ctx context.Context) ([]LiveEvent, error) {
	return []LiveEvent{
		{
			Name:          "Rookie Rush Night",
			DateTime:      "Schedule coming soon",
			Platform:      "Whatnot / TikTok TBD",
			Description:   "Fast-paced football rookie singles, color, and affordable first-buyer lots.",
			FeaturedCards: "Jayden Daniels, Bo Nix, Rome Odunze",
			Status:        "Upcoming",
		},
		{
			Name:          "Sunday Singles Drop",
			DateTime:      "Sundays, time TBD",
			Platform:      "Website + live stream",
			Description:   "Fresh raw-card listings and simple cart checkout support for direct buyers.",
			FeaturedCards: "Featured Drops and team stacks",
			Status:        "Upcoming",
		},
		{
			Name:          "Downtown Hunt",
			DateTime:      "Pop-up event",
			Platform:      "Live link coming soon",
			Description:   "Short-print inserts, chase cards, and spotlight teams when inventory is ready.",
			FeaturedCards: "C.J. Stroud Downtown Raw",
			Status:        "Upcoming",
		},
		{
			Name:          "Auto Chase Live",
			DateTime:      "Date TBD",
			Platform:      "Live link coming soon",
			Description:   "Autographs and live-priority cards grouped for easy claims and checkout.",
			FeaturedCards: "Drake Maye Select Auto Raw",
			Status:        "Upcoming",
		},
	}, nil
}

func (s *Store) SalesChannels(ctx context.Context) ([]SalesChannel, error) {
	return []SalesChannel{
		{Name: "Website Shop", Type: "Online", Schedule: "Open 24/7", Notes: "Direct storefront for football singles, Featured Drops, cart testing, and customer accounts."},
		{Name: "Live Stream Drops", Type: "Live Stream", Schedule: "Schedule coming soon", Notes: "Live link and platform can be added once the selling channel is ready."},
		{Name: "Local Card Events", Type: "Local Event", Schedule: "Dates coming soon", Notes: "In-person sales channel placeholder for future event announcements."},
	}, nil
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
	if product.PlayerName == "" {
		product.PlayerName = product.Player
	}
	if product.Player == "" {
		product.Player = product.PlayerName
	}
	if product.TeamLogoAlt == "" && product.Team != "" {
		product.TeamLogoAlt = product.Team + " logo"
	}
	if product.LivePriority {
		product.LiveExclusive = true
	}
	if product.LiveExclusive {
		product.LivePriority = true
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

func (s *Store) CreateOrder(ctx context.Context, email, customerName string, cart []CartLine, couponCode string) (Order, error) {
	if strings.TrimSpace(email) == "" {
		return Order{}, errors.New("email is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	var items []OrderItem
	subtotal := 0
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
		subtotal += product.PriceCents * line.Quantity
	}
	if len(items) == 0 {
		return Order{}, errors.New("cart is empty")
	}
	couponCode = strings.ToUpper(strings.TrimSpace(couponCode))
	discount := int(float64(subtotal) * demoDiscountRate(couponCode))
	tax := 0
	total := subtotal - discount + tax
	s.nextOrderID++
	order := Order{
		ID:                s.nextOrderID,
		OrderNumber:       fmt.Sprintf("CCC-%d", time.Now().UnixNano()/1_000_000),
		Email:             normalizeEmail(email),
		CustomerName:      strings.TrimSpace(customerName),
		Source:            "web_demo",
		PaymentStatus:     "demo_checkout",
		FulfillmentStatus: "not_started",
		ShippingStatus:    "pending",
		SubtotalCents:     subtotal,
		DiscountCents:     discount,
		TaxCents:          tax,
		TotalCents:        total,
		CouponCode:        couponCode,
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

func demoDiscountRate(code string) float64 {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case "WELCOME10":
		return 0.10
	case "LIVE5":
		return 0.05
	default:
		return 0
	}
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
