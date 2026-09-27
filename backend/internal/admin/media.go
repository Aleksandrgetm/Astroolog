package admin

import (
	"astroolog/backend/internal/models"
	"bytes"
	"errors"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Storage interface {
	Put(string, []byte) error
	Delete(string) error
	Path(string) (string, error)
}
type LocalStorage struct{ Root string }

func (s LocalStorage) Path(name string) (string, error) {
	if filepath.Base(name) != name || name == "." || name == "" {
		return "", errors.New("invalid storage name")
	}
	return filepath.Join(s.Root, name), nil
}
func (s LocalStorage) Put(name string, b []byte) error {
	p, e := s.Path(name)
	if e != nil {
		return e
	}
	if e = os.MkdirAll(s.Root, 0750); e != nil {
		return e
	}
	return os.WriteFile(p, b, 0640)
}
func (s LocalStorage) Delete(name string) error {
	p, e := s.Path(name)
	if e != nil {
		return e
	}
	return os.Remove(p)
}

var imageWorker sync.Mutex

func decodeImage(b []byte) (image.Image, error) {
	var cfg image.Config
	var e error
	kind := http.DetectContentType(b)
	switch {
	case kind == "image/jpeg" || kind == "image/png":
		cfg, _, e = image.DecodeConfig(bytes.NewReader(b))
	case len(b) > 12 && string(b[:4]) == "RIFF" && string(b[8:12]) == "WEBP":
		kind = "image/webp"
		cfg, e = webp.DecodeConfig(bytes.NewReader(b))
	case len(b) > 12 && string(b[4:8]) == "ftyp" && (string(b[8:12]) == "avif" || string(b[8:12]) == "avis"):
		kind = "image/avif"
		cfg, e = avif.DecodeConfig(bytes.NewReader(b))
	default:
		return nil, errors.New("Разрешены только JPEG, PNG, WebP и AVIF")
	}
	if e != nil || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > 12000 || cfg.Height > 12000 || int64(cfg.Width)*int64(cfg.Height) > 20000000 {
		return nil, errors.New("Изображение повреждено или превышает 20 мегапикселей")
	}
	switch kind {
	case "image/webp":
		return webp.Decode(bytes.NewReader(b))
	case "image/avif":
		return avif.Decode(bytes.NewReader(b), avif.Options{AutoRotate: true})
	}
	return imaging.Decode(bytes.NewReader(b), imaging.AutoOrientation(true))
}
func resize(src image.Image, width int) image.Image {
	if width >= src.Bounds().Dx() {
		return src
	}
	return imaging.Resize(src, width, 0, imaging.Lanczos)
}
func (h Handler) Upload(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 11<<20)
	file, e := c.FormFile("file")
	if e != nil || file.Size > 10<<20 {
		fail(c, errors.New("Максимальный размер файла — 10 MB"))
		return
	}
	reader, e := file.Open()
	if e != nil {
		c.Status(400)
		return
	}
	defer reader.Close()
	b, e := io.ReadAll(io.LimitReader(reader, (10<<20)+1))
	if e != nil || len(b) > 10<<20 {
		c.Status(413)
		return
	}
	if !imageWorker.TryLock() {
		c.JSON(429, gin.H{"error": "Другое изображение обрабатывается. Повторите через несколько секунд."})
		return
	}
	defer imageWorker.Unlock()
	img, e := decodeImage(b)
	if e != nil {
		fail(c, e)
		return
	}
	id := token()
	paths := []string{}
	variants := models.Document{}
	var full []byte
	for _, width := range []int{640, 960, 1280, img.Bounds().Dx()} {
		if width > img.Bounds().Dx() {
			continue
		}
		var encoded bytes.Buffer
		if e = webp.Encode(&encoded, resize(img, width), webp.Options{Quality: 86, Method: 4}); e != nil {
			break
		}
		name := fmt.Sprintf("%s-%d.webp", id, width)
		e = h.Store.Put(name, encoded.Bytes())
		if e != nil {
			break
		}
		paths = append(paths, name)
		variants[strconv.Itoa(width)] = "/media/" + name
		if width == img.Bounds().Dx() {
			full = encoded.Bytes()
		}
	}
	if e != nil {
		for _, p := range paths {
			h.Store.Delete(p)
		}
		c.JSON(500, gin.H{"error": "Не удалось обработать изображение"})
		return
	}
	name := fmt.Sprintf("%s-%d.webp", id, img.Bounds().Dx())
	asset := models.MediaAsset{Filename: name, OriginalFilename: filepath.Base(file.Filename), MimeType: "image/webp", Size: int64(len(full)), Width: img.Bounds().Dx(), Height: img.Bounds().Dy(), StoragePath: "/media/" + name, CreatedBy: uid(c), Variants: variants}
	e = h.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Create(&asset).Error; e != nil {
			return e
		}
		return audit(tx, uid(c), "upload", "media", strconv.Itoa(int(asset.ID)))
	})
	if e != nil {
		for _, p := range paths {
			h.Store.Delete(p)
		}
		c.Status(500)
		return
	}
	c.JSON(201, asset)
}
func (h Handler) mediaUses(db *gorm.DB, m models.MediaAsset) (int64, error) {
	var total, n int64
	for _, table := range []string{"services", "booking_options", "directions"} {
		if e := db.Table(table).Where("image=?", m.StoragePath).Count(&n).Error; e != nil {
			return 0, e
		}
		total += n
	}
	if e := db.Model(&models.PageSection{}).Where("settings->>'image'=?", m.StoragePath).Count(&n).Error; e != nil {
		return 0, e
	}
	total += n
	if e := db.Model(&models.CertificateAsset{}).Where("media_id=?", m.ID).Count(&n).Error; e != nil {
		return 0, e
	}
	return total + n, nil
}
func (h Handler) Media(c *gin.Context) {
	p, n := pageArgs(c)
	var rows []models.MediaAsset
	var total int64
	q := h.DB.Model(&models.MediaAsset{})
	if search := strings.TrimSpace(c.Query("q")); search != "" {
		if len(search) > 200 {
			c.Status(400)
			return
		}
		q = q.Where("original_filename ILIKE ?", "%"+search+"%")
	}
	if category := c.Query("category"); category == "certificates" {
		q = q.Where("EXISTS (SELECT 1 FROM certificate_assets ca WHERE ca.media_id=media_assets.id)")
	} else if category == "photos" {
		q = q.Where("NOT EXISTS (SELECT 1 FROM certificate_assets ca WHERE ca.media_id=media_assets.id)")
	}
	if e := q.Count(&total).Error; e != nil {
		c.Status(500)
		return
	}
	if e := q.Order("id DESC").Limit(n).Offset((p - 1) * n).Find(&rows).Error; e != nil {
		c.Status(500)
		return
	}
	items := []gin.H{}
	for _, r := range rows {
		uses, e := h.mediaUses(h.DB, r)
		if e != nil {
			c.Status(500)
			return
		}
		items = append(items, gin.H{"asset": r, "usage_count": uses, "uses": h.mediaLocations(r)})
	}
	c.JSON(200, gin.H{"items": items, "total": total})
}
func (h Handler) DeleteMedia(c *gin.Context) {
	var media models.MediaAsset
	e := h.DB.Transaction(func(tx *gorm.DB) error {
		if e := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&media, "id=?", c.Param("id")).Error; e != nil {
			return e
		}
		if !strings.HasPrefix(media.StoragePath, "/media/") {
			return errors.New("Исходные изображения проекта нельзя удалить через CMS")
		}
		uses, e := h.mediaUses(tx, media)
		if e != nil {
			return e
		}
		if uses > 0 {
			return errors.New("Изображение используется. Сначала замените его во всех блоках")
		}
		if e := tx.Delete(&media).Error; e != nil {
			return e
		}
		return audit(tx, uid(c), "delete media", "media", c.Param("id"))
	})
	if e != nil {
		fail(c, e)
		return
	}
	for _, path := range media.Variants {
		if p, ok := path.(string); ok {
			h.Store.Delete(strings.TrimPrefix(p, "/media/"))
		}
	}
	c.Status(204)
}
func (h Handler) ServeMedia(c *gin.Context) {
	path, e := h.Store.Path(c.Param("name"))
	if e != nil {
		c.Status(404)
		return
	}
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("Cache-Control", "public,max-age=31536000,immutable")
	c.Header("Content-Type", "image/webp")
	c.File(path)
}

// Names and links let an editor find the owning content before replacing a used image.
func (h Handler) mediaLocations(m models.MediaAsset) []gin.H {
	locations := []gin.H{}
	var sections []models.PageSection
	h.DB.Select("id,section_key").Where("settings->>'image'=?", m.StoragePath).Find(&sections)
	for _, s := range sections {
		locations = append(locations, gin.H{"kind": "sections", "key": s.SectionKey})
	}
	for _, kind := range []struct{ table, route, title string }{{"services", "services", "title_ru"}, {"booking_options", "options", "title_ru"}, {"directions", "directions", "(SELECT title FROM direction_translations WHERE direction_id=directions.id AND language='ru' LIMIT 1)"}} {
		var rows []map[string]any
		id := "id"
		if kind.table == "booking_options" {
			id = "code"
		}
		h.DB.Table(kind.table).Select(id+" AS id,"+kind.title+" AS title").Where("image=?", m.StoragePath).Find(&rows)
		for _, r := range rows {
			locations = append(locations, gin.H{"kind": kind.route, "id": r["id"], "title": r["title"]})
		}
	}
	var certificates []models.CertificateAsset
	h.DB.Where("media_id=?", m.ID).Find(&certificates)
	for _, a := range certificates {
		locations = append(locations, gin.H{"kind": "certificates", "id": a.CertificateID, "title": "Сертификат · " + strings.ToUpper(a.Language)})
	}
	return locations
}
