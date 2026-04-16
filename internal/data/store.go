package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	_ "modernc.org/sqlite"

	"creswoodcornerscards/internal/config"
)

type Store struct {
	DB  *sql.DB
	cfg config.Config
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

func Open(cfg config.Config) (*Store, error) {
	db, err := sql.Open("sqlite", cfg.DatabasePath)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)
	store := &Store{DB: db, cfg: cfg}
	if err := store.init(context.Background()); err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Store) init(ctx context.Context) error {
	schema := []string{
		`create table if not exists users (
			id integer primary key autoincrement,
			email text not null unique,
			password_hash text not null,
			name text not null,
			is_admin integer not null default 0,
			created_at datetime not null default current_timestamp
		);`,
		`create table if not exists sessions (
			token text primary key,
			user_id integer not null,
			created_at datetime not null default current_timestamp,
			expires_at datetime not null,
			foreign key(user_id) references users(id) on delete cascade
		);`,
		`create table if not exists products (
			id integer primary key autoincrement,
			slug text not null unique,
			title text not null,
			summary text not null,
			description text not null,
			sport text not null,
			player text not null default '',
			team text not null default '',
			brand text not null default '',
			set_name text not null default '',
			year integer not null default 0,
			card_number text not null default '',
			grade text not null default '',
			card_condition text not null default '',
			price_cents integer not null,
			quantity integer not null default 1,
			accent text not null default 'gold',
			featured integer not null default 0,
			live_exclusive integer not null default 0,
			status text not null default 'active',
			created_at datetime not null default current_timestamp,
			updated_at datetime not null default current_timestamp
		);`,
		`create table if not exists live_sessions (
			id integer primary key autoincrement,
			title text not null,
			pitch text not null,
			callout text not null,
			stream_url text not null,
			platform text not null,
			is_active integer not null default 0,
			created_at datetime not null default current_timestamp,
			updated_at datetime not null default current_timestamp
		);`,
		`create table if not exists orders (
			id integer primary key autoincrement,
			order_number text not null unique,
			email text not null,
			customer_name text not null default '',
			source text not null,
			payment_status text not null,
			fulfillment_status text not null,
			shipping_status text not null,
			total_cents integer not null,
			notes text not null default '',
			created_at datetime not null default current_timestamp
		);`,
		`create table if not exists order_items (
			id integer primary key autoincrement,
			order_id integer not null,
			product_id integer not null,
			product_title text not null,
			product_slug text not null,
			quantity integer not null,
			unit_price integer not null,
			foreign key(order_id) references orders(id) on delete cascade
		);`,
	}

	for _, stmt := range schema {
		if _, err := s.DB.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	if err := s.ensureBootstrapAdmin(ctx); err != nil {
		return err
	}
	if err := s.seedProducts(ctx); err != nil {
		return err
	}
	if err := s.seedLiveSession(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Store) ensureBootstrapAdmin(ctx context.Context) error {
	if s.cfg.AdminBootstrapEmail == "" {
		return nil
	}

	_, err := s.FindUserByEmail(ctx, s.cfg.AdminBootstrapEmail)
	if err == nil {
		_, err = s.DB.ExecContext(ctx, `update users set is_admin = 1 where email = ?`, s.cfg.AdminBootstrapEmail)
		return err
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.cfg.AdminBootstrapPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.ExecContext(ctx,
		`insert into users (email, password_hash, name, is_admin) values (?, ?, ?, 1)`,
		s.cfg.AdminBootstrapEmail,
		string(hash),
		"Store Owner",
	)
	return err
}

func (s *Store) seedProducts(ctx context.Context) error {
	var count int
	if err := s.DB.QueryRowContext(ctx, `select count(*) from products`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	products := []Product{
		{Slug: "jayden-daniels-prizm-rc-gold-wave", Title: "Jayden Daniels Prizm RC Gold Wave", Summary: "Low-number rookie color built to headline both the shop and live stream.", Description: "A flagship rookie listing with clear condition notes, sharp presentation, and immediate buying confidence.", Sport: "Football", Player: "Jayden Daniels", Team: "Washington Commanders", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "312", Grade: "Raw", Condition: "Near Mint", PriceCents: 32900, Quantity: 1, Accent: "gold", Featured: true, LiveExclusive: true, Status: "active"},
		{Slug: "cj-stroud-downtown-sgc-10", Title: "C.J. Stroud Downtown SGC 10", Summary: "Premium slab inventory to anchor the homepage and featured rails.", Description: "High-end inventory gives the storefront trust and average order value. This slab acts as a premium centerpiece.", Sport: "Football", Player: "C.J. Stroud", Team: "Houston Texans", Brand: "Donruss", SetName: "Downtown", Year: 2023, CardNumber: "DT-7", Grade: "SGC 10", Condition: "Gem Mint", PriceCents: 79900, Quantity: 1, Accent: "crimson", Featured: true, Status: "active"},
		{Slug: "bo-nix-mosaic-genesis-rc", Title: "Bo Nix Mosaic Genesis RC", Summary: "Mid-tier chase card with fast add-to-cart appeal.", Description: "A strong chase rookie that supports quick catalog browsing and live stream call-outs.", Sport: "Football", Player: "Bo Nix", Team: "Denver Broncos", Brand: "Panini Mosaic", SetName: "Mosaic Football", Year: 2024, CardNumber: "287", Grade: "Raw", Condition: "Near Mint-Mint", PriceCents: 18500, Quantity: 2, Accent: "emerald", Featured: true, Status: "active"},
		{Slug: "drake-maye-select-field-level-auto", Title: "Drake Maye Select Field Level Auto", Summary: "Autograph inventory that works for direct checkout and stream mentions.", Description: "A balanced autograph listing that supports both normal storefront conversion and live merchandising.", Sport: "Football", Player: "Drake Maye", Team: "New England Patriots", Brand: "Panini Select", SetName: "Select Football", Year: 2024, CardNumber: "FLA-DM", Grade: "Raw", Condition: "Near Mint", PriceCents: 26900, Quantity: 1, Accent: "blue", Status: "active"},
		{Slug: "malik-nabers-optic-holo-rc", Title: "Malik Nabers Optic Holo RC", Summary: "Popular rookie profile that keeps the grid feeling current.", Description: "Designed as a flexible inventory piece for homepage, catalog, and live modules.", Sport: "Football", Player: "Malik Nabers", Team: "New York Giants", Brand: "Donruss Optic", SetName: "Optic Football", Year: 2024, CardNumber: "214", Grade: "Raw", Condition: "Near Mint", PriceCents: 12900, Quantity: 3, Accent: "violet", Status: "active"},
		{Slug: "rome-odunze-prizm-silver-rc", Title: "Rome Odunze Prizm Silver RC", Summary: "Affordable parallel for a healthier catalog price spread.", Description: "Not every listing should be a grail. This card supports basket building and easier first purchases.", Sport: "Football", Player: "Rome Odunze", Team: "Chicago Bears", Brand: "Panini Prizm", SetName: "Prizm Football", Year: 2024, CardNumber: "341", Grade: "Raw", Condition: "Near Mint", PriceCents: 8900, Quantity: 4, Accent: "slate", Status: "active"},
	}

	for _, p := range products {
		_, err := s.DB.ExecContext(ctx, `insert into products
			(slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.Slug, p.Title, p.Summary, p.Description, p.Sport, p.Player, p.Team, p.Brand, p.SetName, p.Year, p.CardNumber, p.Grade, p.Condition, p.PriceCents, p.Quantity, p.Accent, boolInt(p.Featured), boolInt(p.LiveExclusive), p.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) seedLiveSession(ctx context.Context) error {
	var count int
	if err := s.DB.QueryRowContext(ctx, `select count(*) from live_sessions`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	_, err := s.DB.ExecContext(ctx,
		`insert into live_sessions (title, pitch, callout, stream_url, platform, is_active) values (?, ?, ?, ?, ?, 1)`,
		"Friday Night Football Heat Check",
		"Live singles, rookie color, and quick-hit offers built to move inventory fast.",
		"Running live now and mirrored from the storefront banner.",
		"https://www.whatnot.com/",
		"Whatnot",
	)
	return err
}

func (s *Store) FindUserByEmail(ctx context.Context, email string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `select id, email, password_hash, name, is_admin, created_at from users where lower(email) = ?`, normalizeEmail(email))
	return scanUser(row)
}

func (s *Store) FindUserByID(ctx context.Context, id int64) (User, error) {
	row := s.DB.QueryRowContext(ctx, `select id, email, password_hash, name, is_admin, created_at from users where id = ?`, id)
	return scanUser(row)
}

func scanUser(scanner interface{ Scan(dest ...any) error }) (User, error) {
	var u User
	var isAdmin int
	err := scanner.Scan(&u.ID, &u.Email, &u.Password, &u.Name, &isAdmin, &u.CreatedAt)
	u.IsAdmin = isAdmin == 1
	return u, err
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

	isAdmin := 0
	if _, ok := s.cfg.AllowedAdminEmails[email]; ok || email == s.cfg.AdminBootstrapEmail {
		isAdmin = 1
	}

	result, err := s.DB.ExecContext(ctx,
		`insert into users (email, password_hash, name, is_admin) values (?, ?, ?, ?)`,
		email, string(hash), strings.TrimSpace(name), isAdmin,
	)
	if err != nil {
		return User{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, err
	}

	return s.FindUserByID(ctx, id)
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
	_, err := s.DB.ExecContext(ctx,
		`insert into sessions (token, user_id, expires_at) values (?, ?, ?)`,
		token, userID, time.Now().Add(30*24*time.Hour),
	)
	return token, err
}

func (s *Store) FindUserBySession(ctx context.Context, token string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `select users.id, users.email, users.password_hash, users.name, users.is_admin, users.created_at
		from sessions join users on users.id = sessions.user_id
		where sessions.token = ? and sessions.expires_at > ?`, token, time.Now())
	return scanUser(row)
}

func (s *Store) DeleteSession(ctx context.Context, token string) error {
	_, err := s.DB.ExecContext(ctx, `delete from sessions where token = ?`, token)
	return err
}

func (s *Store) FeaturedProducts(ctx context.Context) ([]Product, error) {
	rows, err := s.DB.QueryContext(ctx, `select id, slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status from products where status = 'active' and featured = 1 order by id asc limit 4`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (s *Store) Products(ctx context.Context) ([]Product, error) {
	rows, err := s.DB.QueryContext(ctx, `select id, slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status from products where status != 'archived' order by featured desc, id asc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanProducts(rows)
}

func (s *Store) ProductBySlug(ctx context.Context, slug string) (Product, error) {
	row := s.DB.QueryRowContext(ctx, `select id, slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status from products where slug = ?`, slug)
	return scanProduct(row)
}

func (s *Store) ProductByID(ctx context.Context, id int64) (Product, error) {
	row := s.DB.QueryRowContext(ctx, `select id, slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status from products where id = ?`, id)
	return scanProduct(row)
}

func scanProducts(rows *sql.Rows) ([]Product, error) {
	var items []Product
	for rows.Next() {
		item, err := scanProduct(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func scanProduct(scanner interface{ Scan(dest ...any) error }) (Product, error) {
	var p Product
	var featured, live int
	err := scanner.Scan(&p.ID, &p.Slug, &p.Title, &p.Summary, &p.Description, &p.Sport, &p.Player, &p.Team, &p.Brand, &p.SetName, &p.Year, &p.CardNumber, &p.Grade, &p.Condition, &p.PriceCents, &p.Quantity, &p.Accent, &featured, &live, &p.Status)
	p.Featured = featured == 1
	p.LiveExclusive = live == 1
	return p, err
}

func (s *Store) ActiveLiveSession(ctx context.Context) (LiveSession, error) {
	row := s.DB.QueryRowContext(ctx, `select id, title, pitch, callout, stream_url, platform, is_active from live_sessions order by id desc limit 1`)
	var live LiveSession
	var active int
	err := row.Scan(&live.ID, &live.Title, &live.Pitch, &live.Callout, &live.StreamURL, &live.Platform, &active)
	live.IsActive = active == 1
	return live, err
}

func (s *Store) SaveLiveSession(ctx context.Context, live LiveSession) error {
	var count int
	if err := s.DB.QueryRowContext(ctx, `select count(*) from live_sessions`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		_, err := s.DB.ExecContext(ctx, `insert into live_sessions (title, pitch, callout, stream_url, platform, is_active) values (?, ?, ?, ?, ?, ?)`,
			live.Title, live.Pitch, live.Callout, live.StreamURL, live.Platform, boolInt(live.IsActive))
		return err
	}
	_, err := s.DB.ExecContext(ctx, `update live_sessions set title = ?, pitch = ?, callout = ?, stream_url = ?, platform = ?, is_active = ?, updated_at = current_timestamp where id = (select id from live_sessions order by id desc limit 1)`,
		live.Title, live.Pitch, live.Callout, live.StreamURL, live.Platform, boolInt(live.IsActive))
	return err
}

func (s *Store) UpsertProduct(ctx context.Context, p Product) error {
	if p.ID == 0 {
		_, err := s.DB.ExecContext(ctx, `insert into products (slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status)
			values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			p.Slug, p.Title, p.Summary, p.Description, p.Sport, p.Player, p.Team, p.Brand, p.SetName, p.Year, p.CardNumber, p.Grade, p.Condition, p.PriceCents, p.Quantity, p.Accent, boolInt(p.Featured), boolInt(p.LiveExclusive), p.Status)
		return err
	}

	_, err := s.DB.ExecContext(ctx, `update products set slug=?, title=?, summary=?, description=?, sport=?, player=?, team=?, brand=?, set_name=?, year=?, card_number=?, grade=?, card_condition=?, price_cents=?, quantity=?, accent=?, featured=?, live_exclusive=?, status=?, updated_at=current_timestamp where id=?`,
		p.Slug, p.Title, p.Summary, p.Description, p.Sport, p.Player, p.Team, p.Brand, p.SetName, p.Year, p.CardNumber, p.Grade, p.Condition, p.PriceCents, p.Quantity, p.Accent, boolInt(p.Featured), boolInt(p.LiveExclusive), p.Status, p.ID)
	return err
}

func (s *Store) CreateOrder(ctx context.Context, email, customerName string, cart []CartLine) (Order, error) {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return Order{}, err
	}
	defer tx.Rollback()

	total := 0
	lines := make([]OrderItem, 0, len(cart))
	for _, line := range cart {
		product, err := productByIDTx(ctx, tx, line.ProductID)
		if err != nil {
			return Order{}, err
		}
		if line.Quantity < 1 {
			line.Quantity = 1
		}
		total += product.PriceCents * line.Quantity
		lines = append(lines, OrderItem{
			ProductID:    product.ID,
			ProductTitle: product.Title,
			ProductSlug:  product.Slug,
			Quantity:     line.Quantity,
			UnitPrice:    product.PriceCents,
		})
	}

	result, err := tx.ExecContext(ctx, `insert into orders (order_number, email, customer_name, source, payment_status, fulfillment_status, shipping_status, total_cents) values (?, ?, ?, 'web_demo', 'paid', 'ready_to_pack', 'pending', ?)`,
		fmt.Sprintf("CCC-%d", time.Now().UnixNano()/1_000_000), normalizeEmail(email), customerName, total)
	if err != nil {
		return Order{}, err
	}
	orderID, err := result.LastInsertId()
	if err != nil {
		return Order{}, err
	}

	for _, line := range lines {
		if _, err := tx.ExecContext(ctx, `insert into order_items (order_id, product_id, product_title, product_slug, quantity, unit_price) values (?, ?, ?, ?, ?, ?)`,
			orderID, line.ProductID, line.ProductTitle, line.ProductSlug, line.Quantity, line.UnitPrice); err != nil {
			return Order{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Order{}, err
	}
	return s.OrderByID(ctx, orderID)
}

func productByIDTx(ctx context.Context, tx *sql.Tx, id int64) (Product, error) {
	row := tx.QueryRowContext(ctx, `select id, slug, title, summary, description, sport, player, team, brand, set_name, year, card_number, grade, card_condition, price_cents, quantity, accent, featured, live_exclusive, status from products where id = ?`, id)
	return scanProduct(row)
}

type CartLine struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func (s *Store) OrderByID(ctx context.Context, id int64) (Order, error) {
	row := s.DB.QueryRowContext(ctx, `select id, order_number, email, customer_name, source, payment_status, fulfillment_status, shipping_status, total_cents, notes, created_at from orders where id = ?`, id)
	var order Order
	if err := row.Scan(&order.ID, &order.OrderNumber, &order.Email, &order.CustomerName, &order.Source, &order.PaymentStatus, &order.FulfillmentStatus, &order.ShippingStatus, &order.TotalCents, &order.Notes, &order.CreatedAt); err != nil {
		return Order{}, err
	}
	items, err := s.orderItems(ctx, id)
	if err != nil {
		return Order{}, err
	}
	order.Items = items
	return order, nil
}

func (s *Store) Orders(ctx context.Context) ([]Order, error) {
	rows, err := s.DB.QueryContext(ctx, `select id, order_number, email, customer_name, source, payment_status, fulfillment_status, shipping_status, total_cents, notes, created_at from orders order by created_at desc`)
	if err != nil {
		return nil, err
	}

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.OrderNumber, &order.Email, &order.CustomerName, &order.Source, &order.PaymentStatus, &order.FulfillmentStatus, &order.ShippingStatus, &order.TotalCents, &order.Notes, &order.CreatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	rows.Close()

	for i := range orders {
		items, err := s.orderItems(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Items = items
	}

	return orders, nil
}

func (s *Store) orderItems(ctx context.Context, orderID int64) ([]OrderItem, error) {
	rows, err := s.DB.QueryContext(ctx, `select id, product_id, product_title, product_slug, quantity, unit_price from order_items where order_id = ?`, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		if err := rows.Scan(&item.ID, &item.ProductID, &item.ProductTitle, &item.ProductSlug, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateOrder(ctx context.Context, orderID int64, payment, fulfillment, shipping, notes string) error {
	_, err := s.DB.ExecContext(ctx, `update orders set payment_status=?, fulfillment_status=?, shipping_status=?, notes=? where id=?`,
		payment, fulfillment, shipping, notes, orderID)
	return err
}

func (s *Store) Users(ctx context.Context) ([]User, error) {
	rows, err := s.DB.QueryContext(ctx, `select id, email, password_hash, name, is_admin, created_at from users order by created_at desc`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (s *Store) DashboardStats(ctx context.Context) (DashboardStats, error) {
	var stats DashboardStats
	if err := s.DB.QueryRowContext(ctx, `select count(*) from products`).Scan(&stats.ProductCount); err != nil {
		return stats, err
	}
	if err := s.DB.QueryRowContext(ctx, `select count(*) from orders`).Scan(&stats.OrderCount); err != nil {
		return stats, err
	}
	if err := s.DB.QueryRowContext(ctx, `select count(*) from users`).Scan(&stats.AccountCount); err != nil {
		return stats, err
	}
	if err := s.DB.QueryRowContext(ctx, `select coalesce(sum(total_cents), 0) from orders`).Scan(&stats.GrossSales); err != nil {
		return stats, err
	}
	return stats, nil
}

func normalizeEmail(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func UniqueTeams(products []Product) []string {
	seen := map[string]struct{}{}
	var teams []string
	for _, p := range products {
		if p.Team == "" {
			continue
		}
		if _, ok := seen[p.Team]; ok {
			continue
		}
		seen[p.Team] = struct{}{}
		teams = append(teams, p.Team)
	}
	sort.Strings(teams)
	return teams
}
