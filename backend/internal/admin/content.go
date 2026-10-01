package admin

import (
	"astroolog/backend/internal/database"
	"astroolog/backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type relation struct{ Name, Table, FK, Field string }
type entity struct {
	New       func() any
	List      func() any
	Fields    []string
	Relations []relation
	Create    bool
}

func fields(s string) []string { return strings.Fields(s) }

var entities = map[string]entity{
	"pages":        {func() any { return &models.Page{} }, func() any { return &[]models.Page{} }, fields("translations"), []relation{{"Translations", "page_translations", "page_id", "PageID"}}, false},
	"sections":     {func() any { return &models.PageSection{} }, func() any { return &[]models.PageSection{} }, fields("is_active sort_order settings translations"), []relation{{"Translations", "page_section_translations", "page_section_id", "PageSectionID"}}, false},
	"steps":        {func() any { return &models.ProcessStep{} }, func() any { return &[]models.ProcessStep{} }, fields("page_section_id sort_order icon_key translations"), []relation{{"Translations", "process_step_translations", "process_step_id", "ProcessStepID"}}, true},
	"services":     {func() any { return &models.Service{} }, func() any { return &[]models.Service{} }, fields("slug title_ru title_lv title_en short_description_ru short_description_lv short_description_en description_ru description_lv description_en duration_ru duration_lv duration_en image alt_ru alt_lv alt_en is_active sort_order"), nil, true},
	"options":      {func() any { return &models.BookingOption{} }, func() any { return &[]models.BookingOption{} }, fields("code service_id title_ru title_lv title_en format_ru format_lv format_en price image is_active is_addon sort_order"), nil, true},
	"directions":   {func() any { return &models.Direction{} }, func() any { return &[]models.Direction{} }, fields("slug image service_id is_active sort_order translations options"), []relation{{"Translations", "direction_translations", "direction_id", "DirectionID"}, {"Options", "direction_options", "direction_id", "DirectionID"}}, true},
	"reviews":      {func() any { return &models.Review{} }, func() any { return &[]models.Review{} }, fields("is_active is_featured sort_order home_sort_order author_display_name publication_consent translations"), []relation{{"Translations", "review_translations", "review_id", "ReviewID"}}, true},
	"certificates": {func() any { return &models.Certificate{} }, func() any { return &[]models.Certificate{} }, fields("issuer issued_at year is_active sort_order translations assets"), []relation{{"Translations", "certificate_translations", "certificate_id", "CertificateID"}, {"Assets", "certificate_assets", "certificate_id", "CertificateID"}}, true},
	"settings":     {func() any { return &models.SiteSettings{} }, func() any { return &[]models.SiteSettings{} }, fields("expert_name phone email whatsapp_phone telegram_url professional_title_ru professional_title_lv professional_title_en"), nil, false},
}

func query(db *gorm.DB, s entity) *gorm.DB {
	for _, r := range s.Relations {
		db = db.Preload(r.Name)
		if r.Name == "Assets" {
			db = db.Preload("Assets.Media")
		}
	}
	return db
}
func fail(c *gin.Context, e error) { c.JSON(400, gin.H{"error": e.Error()}) }
func pageArgs(c *gin.Context) (int, int) {
	p, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if p < 1 {
		p = 1
	}
	n, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))
	if n < 1 || n > 100 {
		n = 30
	}
	return p, n
}
func (h Handler) List(c *gin.Context) {
	s, ok := entities[c.Param("kind")]
	if !ok {
		c.Status(404)
		return
	}
	p, n := pageArgs(c)
	var count int64
	q := h.DB.Model(s.New())
	if c.Param("kind") == "options" && c.Query("service_id") != "" {
		q = q.Where("service_id=?", c.Query("service_id"))
	}
	if e := q.Count(&count).Error; e != nil {
		c.Status(500)
		return
	}
	if c.Param("kind") == "services" {
		q = q.Select("services.*, (SELECT MIN(o.price) FROM booking_options o WHERE o.service_id=services.id AND o.is_active=true AND o.is_addon=false) AS price")
	}
	out := s.List()
	order := "id"
	if c.Param("kind") == "options" {
		order = "sort_order, code"
	}
	if e := query(q, s).Order(order).Limit(n).Offset((p - 1) * n).Find(out).Error; e != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{"items": out, "total": count, "fields": s.Fields, "can_create": s.Create})
}
func (h Handler) Get(c *gin.Context) {
	spec, ok := entities[c.Param("kind")]
	if !ok {
		c.Status(404)
		return
	}
	record := spec.New()
	if e := query(h.DB, spec).First(record, primaryWhere(c.Param("kind")), c.Param("id")).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			c.Status(404)
		} else {
			c.Status(500)
		}
		return
	}
	c.JSON(200, record)
}
func (h Handler) Schema(c *gin.Context) { c.JSON(200, database.CMSDefinitions().Sections) }
func plain(v any) bool {
	switch x := v.(type) {
	case string:
		return len(x) <= 30000 && !strings.Contains(x, "\x00") && !strings.Contains(x, "<") && !strings.Contains(x, ">")
	case map[string]any:
		for _, v := range x {
			if !plain(v) {
				return false
			}
		}
	case []any:
		for _, v := range x {
			if !plain(v) {
				return false
			}
		}
	}
	return true
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

func safeImage(path string) bool {
	return path == "" || ((strings.HasPrefix(path, "/images/") || strings.HasPrefix(path, "/media/")) && !strings.Contains(path, "..") && !strings.ContainsAny(path, "?#\\\"<>"))
}
func languageList(v any) error {
	b, _ := json.Marshal(v)
	var list []map[string]any
	if json.Unmarshal(b, &list) != nil {
		return errors.New("Некорректные переводы")
	}
	seen := map[string]bool{}
	for _, t := range list {
		l, _ := t["language"].(string)
		if (l != "ru" && l != "lv" && l != "en") || seen[l] {
			return errors.New("Языки должны быть RU, LV, EN без повторений")
		}
		seen[l] = true
	}
	if !seen["ru"] {
		return errors.New("Добавьте основной перевод RU")
	}
	return nil
}
func (h Handler) validate(tx *gorm.DB, kind string, v any, old map[string]any) error {
	b, _ := json.Marshal(v)
	var m map[string]any
	json.Unmarshal(b, &m)
	if !plain(m) {
		return errors.New("Используйте обычный текст без HTML, не более 30 000 символов в поле")
	}
	if t, ok := m["translations"]; ok {
		if e := languageList(t); e != nil {
			return e
		}
	}
	if map[string]bool{"pages": true, "steps": true, "directions": true, "certificates": true}[kind] {
		for _, raw := range m["translations"].([]any) {
			tr := raw.(map[string]any)
			if tr["language"] == "ru" {
				title, _ := tr["title"].(string)
				if strings.TrimSpace(title) == "" {
					return errors.New("Заполните название RU")
				}
				if kind == "pages" {
					for _, key := range []string{"meta_title", "meta_description"} {
						value, _ := tr[key].(string)
						if strings.TrimSpace(value) == "" {
							return errors.New("Заполните SEO-заголовок и описание RU")
						}
					}
				}
			}
		}
	}
	// Lock referenced uploads against concurrent deletion.
	if img, ok := m["image"].(string); ok && strings.HasPrefix(img, "/media/") {
		var media models.MediaAsset
		if tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&media, "storage_path=?", img).Error != nil {
			return errors.New("Изображение удалено; выберите другое")
		}
	}
	if img, ok := m["image"].(string); ok && !safeImage(img) {
		return errors.New("Выберите изображение из медиатеки")
	}
	for _, key := range []string{"slug", "code"} {
		// Page routes are read-only and may contain slashes (including the home route).
		if key == "slug" && kind != "services" && kind != "directions" {
			continue
		}
		if value, ok := m[key].(string); ok {
			if !slugPattern.MatchString(value) || len(value) > 80 {
				return errors.New("Код: латинские буквы, цифры и дефис")
			}
			if old[key] != nil && old[key] != value {
				return errors.New("Существующий адрес или код менять нельзя")
			}
		}
	}
	switch item := v.(type) {
	case *models.Service:
		if strings.TrimSpace(item.TitleRU) == "" {
			return errors.New("Заполните название RU")
		}
		if item.Slug == "growth-point" {
			return errors.New("Этот адрес зарезервирован")
		}
	case *models.BookingOption:
		if item.Price != nil && (*item.Price < 0 || *item.Price > 100000000) {
			return errors.New("Некорректная цена")
		}
		if item.TitleRU == "" {
			return errors.New("Заполните название RU")
		}
		var n int64
		tx.Model(&models.Service{}).Where("id=?", item.ServiceID).Count(&n)
		if n != 1 {
			return errors.New("Выберите услугу")
		}
		if old["code"] != nil && (old["service_id"] != float64(item.ServiceID) || old["is_addon"] != item.IsAddon) {
			return errors.New("Нельзя менять тип и группу существующего варианта; создайте новый")
		}
	case *models.PageSection:
		var def *database.SectionDefinition
		defs := database.CMSDefinitions()
		for i := range defs.Sections {
			if defs.Sections[i].SectionKey == item.SectionKey {
				def = &defs.Sections[i]
			}
		}
		if def == nil {
			return errors.New("Неизвестная секция")
		}
		if def.ReadOnly {
			return errors.New("Технические сведения обновляются вместе с кодом сайта")
		}
		if def.Critical && !item.IsActive {
			return errors.New("Обязательную секцию нельзя скрыть")
		}
		allowed := map[string]bool{}
		for _, f := range def.Fields {
			allowed[f.Key] = true
		}
		for _, tr := range item.Translations {
			for key, value := range tr.Fields {
				if !allowed[key] {
					return errors.New("Неизвестное поле секции")
				}
				if _, ok := value.(string); !ok {
					return errors.New("Поля секции должны содержать текст")
				}
			}
		}
		for key, value := range item.Settings {
			if _, ok := def.Settings[key]; !ok {
				return errors.New("Эта настройка не поддерживается")
			}
			s, ok := value.(string)
			if !ok {
				return errors.New("Некорректная настройка")
			}
			if key == "image" {
				if strings.HasPrefix(s, "/media/") {
					var media models.MediaAsset
					if tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&media, "storage_path=?", s).Error != nil {
						return errors.New("Изображение не найдено")
					}
				}
				if !safeImage(s) {
					return errors.New("Неверный путь изображения")
				}
			} else if !map[string]bool{"booking": true, "question": true, "services": true, "reviews": true, "about": true}[s] {
				return errors.New("Выберите разрешённую цель кнопки")
			}
		}
	case *models.ProcessStep:
		if !map[string]bool{"arrow": true, "eye": true, "lightbulb": true, "target": true}[item.IconKey] {
			return errors.New("Выберите иконку из списка")
		}
		var n int64
		tx.Model(&models.PageSection{}).Where("id=? AND section_key='home.process'", item.PageSectionID).Count(&n)
		if n != 1 {
			return errors.New("Шаг относится к секции «Как мы работаем»")
		}
	case *models.Direction:
		if item.ContentKey == "" {
			item.ContentKey = item.Slug
		}
		var n int64
		tx.Model(&models.Service{}).Where("id=?", item.ServiceID).Count(&n)
		if n != 1 {
			return errors.New("Выберите связанную услугу")
		}
		seen := map[string]bool{}
		for _, o := range item.Options {
			if seen[o.Code] {
				return errors.New("Варианты повторяются")
			}
			seen[o.Code] = true
			tx.Model(&models.BookingOption{}).Where("code=? AND service_id=? AND is_addon=false", o.Code, item.ServiceID).Count(&n)
			if n != 1 {
				return errors.New("Вариант должен быть основным и принадлежать выбранной услуге")
			}
		}
	case *models.Review:
		if item.IsActive && !item.PublicationConsent {
			return errors.New("Для публикации подтвердите согласие автора")
		}
		for _, tr := range item.Translations {
			if tr.Language == "ru" && strings.TrimSpace(tr.FullText) == "" {
				return errors.New("Заполните текст отзыва RU")
			}
		}
	case *models.Certificate:
		seen := map[string]bool{}
		for _, a := range item.Assets {
			if seen[a.Language] || !map[string]bool{"ru": true, "lv": true, "en": true}[a.Language] {
				return errors.New("Некорректный язык изображения")
			}
			seen[a.Language] = true
			var media models.MediaAsset
			if tx.Clauses(clause.Locking{Strength: "SHARE"}).First(&media, a.MediaID).Error != nil {
				return errors.New("Изображение не найдено")
			}
		}
	case *models.SiteSettings:
		address, e := mail.ParseAddress(item.Email)
		if e != nil || address.Address != item.Email {
			return errors.New("Проверьте email")
		}
		if !regexp.MustCompile(`^[1-9][0-9]{6,14}$`).MatchString(item.WhatsappPhone) {
			return errors.New("WhatsApp: международный номер, только цифры")
		}
		if item.TelegramURL != "" {
			u, e := url.Parse(item.TelegramURL)
			if e != nil || u.Scheme != "https" || u.Host != "t.me" || u.RawQuery != "" || !regexp.MustCompile(`^/[a-zA-Z][a-zA-Z0-9_]{3,31}/?$`).MatchString(u.Path) {
				return errors.New("Telegram: https://t.me/username")
			}
		}
	}
	return nil
}
func (h Handler) Save(c *gin.Context) {
	kind := c.Param("kind")
	spec, ok := entities[kind]
	if !ok {
		c.Status(404)
		return
	}
	create := c.Request.Method == "POST"
	if create && !spec.Create {
		c.Status(405)
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 256*1024)
	var input map[string]any
	if c.ShouldBindJSON(&input) != nil {
		fail(c, errors.New("Проверьте поля"))
		return
	}
	clean := map[string]any{}
	for _, f := range spec.Fields {
		if v, ok := input[f]; ok {
			clean[f] = v
		}
	}
	record := spec.New()
	action := "update"
	if create {
		action = "create"
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		old := map[string]any{}
		if !create {
			if e := query(tx.Clauses(clause.Locking{Strength: "UPDATE"}), spec).First(record, primaryWhere(kind), c.Param("id")).Error; e != nil {
				return e
			}
			b, _ := json.Marshal(record)
			json.Unmarshal(b, &old)
			if input["updated_at"] == nil || input["updated_at"] != old["updated_at"] {
				return errors.New("Запись уже изменена. Обновите список перед сохранением")
			}
		}
		b, _ := json.Marshal(clean)
		if e := json.Unmarshal(b, record); e != nil {
			return errors.New("Неверный формат полей")
		}
		if e := h.validate(tx, kind, record, old); e != nil {
			return e
		}
		if e := tx.Omit(clause.Associations).Save(record).Error; e != nil {
			return errors.New("Не удалось сохранить: проверьте уникальность кода и связанные записи")
		}
		value := reflect.ValueOf(record).Elem()
		id := uint(0)
		if f := value.FieldByName("ID"); f.IsValid() {
			id = uint(f.Uint())
		}
		for _, rel := range spec.Relations {
			if e := tx.Table(rel.Table).Where(rel.FK+"=?", id).Delete(map[string]any{}).Error; e != nil {
				return e
			}
			slice := value.FieldByName(rel.Name)
			for i := 0; i < slice.Len(); i++ {
				child := slice.Index(i)
				if f := child.FieldByName("ID"); f.IsValid() {
					f.SetUint(0)
				}
				child.FieldByName(rel.Field).SetUint(uint64(id))
				if e := tx.Omit(clause.Associations).Create(child.Addr().Interface()).Error; e != nil {
					return e
				}
			}
		}
		if active, ok := clean["is_active"].(bool); ok && old["is_active"] != nil && old["is_active"] != active {
			if active {
				action = "restore"
			} else {
				action = "archive"
			}
		}
		entityID := fmt.Sprint(id)
		if kind == "options" {
			entityID = record.(*models.BookingOption).Code
		}
		return audit(tx, uid(c), action, kind, entityID)
	})
	if e != nil {
		fail(c, e)
		return
	}
	c.JSON(200, record)
}
func primaryWhere(kind string) string {
	if kind == "options" {
		return "code=?"
	}
	return "id=?"
}
func (h Handler) Public(c *gin.Context) {
	var pages []models.Page
	var sections []models.PageSection
	var steps []models.ProcessStep
	var dirs []models.Direction
	var reviews []models.Review
	var certificates []models.Certificate
	var settings models.SiteSettings
	queries := []*gorm.DB{h.DB.Preload("Translations").Where("is_active=true").Find(&pages), h.DB.Preload("Translations").Order("sort_order,id").Find(&sections), h.DB.Preload("Translations").Order("sort_order,id").Find(&steps), h.DB.Preload("Translations").Preload("Options").Order("sort_order,id").Where("is_active=true").Find(&dirs), h.DB.Preload("Translations").Order("sort_order,id").Where("is_active=true AND publication_consent=true").Find(&reviews), h.DB.Preload("Translations").Preload("Assets.Media").Order("sort_order,id").Where("is_active=true").Find(&certificates), h.DB.First(&settings, 1)}
	for _, q := range queries {
		if q.Error != nil {
			c.Status(503)
			return
		}
	}
	var assets []models.MediaAsset
	if e := h.DB.Select("storage_path, variants").Find(&assets).Error; e != nil {
		c.Status(503)
		return
	}
	variants := map[string]models.Document{}
	for _, asset := range assets {
		if len(asset.Variants) > 0 {
			variants[asset.StoragePath] = asset.Variants
		}
	}
	c.Header("Cache-Control", "no-cache")
	c.JSON(200, gin.H{"pages": pages, "sections": sections, "steps": steps, "directions": dirs, "reviews": reviews, "certificates": certificates, "settings": settings, "media_variants": variants})
}
func (h Handler) Requests(c *gin.Context) {
	p, n := pageArgs(c)
	kind := c.DefaultQuery("type", "bookings")
	table := "bookings"
	status := "workflow_status"
	if kind == "questions" {
		table = "contact_requests"
		status = "status"
	} else if kind != "bookings" {
		c.Status(400)
		return
	}
	q := h.DB.Table(table)
	if value := c.Query("status"); value != "" {
		q = q.Where(status+"=?", value)
	}
	var count int64
	if e := q.Count(&count).Error; e != nil {
		c.Status(500)
		return
	}
	var rows []map[string]any
	if e := q.Order("created_at DESC,id DESC").Limit(n).Offset((p - 1) * n).Find(&rows).Error; e != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{"items": rows, "total": count})
}
func (h Handler) RequestUpdate(c *gin.Context) {
	table := "bookings"
	column := "workflow_status"
	if c.Param("kind") == "questions" {
		table = "contact_requests"
		column = "status"
	} else if c.Param("kind") != "bookings" {
		c.Status(404)
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 32768)
	var input struct {
		Status string `json:"status"`
		Note   string `json:"internal_note"`
	}
	if c.ShouldBindJSON(&input) != nil || !map[string]bool{"new": true, "in_progress": true, "completed": true, "archived": true}[input.Status] || len(input.Note) > 10000 {
		c.Status(400)
		return
	}
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		r := tx.Table(table).Where("id=?", c.Param("id")).Updates(map[string]any{column: input.Status, "internal_note": input.Note, "updated_at": time.Now()})
		if r.Error != nil {
			return r.Error
		}
		if r.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return audit(tx, uid(c), "update", table, c.Param("id"))
	})
	if e != nil {
		c.Status(400)
		return
	}
	c.Status(204)
}
func (h Handler) Audit(c *gin.Context) {
	p, n := pageArgs(c)
	var rows []models.AdminAuditLog
	var total int64
	if e := h.DB.Model(&models.AdminAuditLog{}).Count(&total).Error; e != nil {
		c.Status(500)
		return
	}
	if e := h.DB.Order("id DESC").Limit(n).Offset((p - 1) * n).Find(&rows).Error; e != nil {
		c.Status(500)
		return
	}
	c.JSON(200, gin.H{"items": rows, "total": total})
}
func (h Handler) Dashboard(c *gin.Context) {
	counts := map[string]int64{}
	for _, x := range []struct{ key, table, where string }{{"bookings", "bookings", "workflow_status='new'"}, {"questions", "contact_requests", "status='new'"}, {"services", "services", "is_active=true"}, {"reviews", "reviews", "is_active=true"}, {"certificates", "certificates", "is_active=true"}} {
		var n int64
		if e := h.DB.Table(x.table).Where(x.where).Count(&n).Error; e != nil {
			c.Status(500)
			return
		}
		counts[x.key] = n
	}
	var recent []models.Booking
	if e := h.DB.Order("id DESC").Limit(5).Find(&recent).Error; e != nil {
		c.Status(500)
		return
	}
	var logs []models.AdminAuditLog
	if e := h.DB.Order("id DESC").Limit(8).Find(&logs).Error; e != nil {
		c.Status(500)
		return
	}
	var publication models.Publication
	if e := h.DB.Order("id DESC").Limit(1).Find(&publication).Error; e != nil {
		c.Status(500)
		return
	}
	if publication.Status == "running" && time.Since(publication.CreatedAt) > 15*time.Minute {
		publication.Status = "failed"
		publication.Error = "Процесс публикации был прерван. Повторите публикацию."
	}
	c.JSON(200, gin.H{"counts": counts, "requests": recent, "audit": logs, "publication": publication})
}
