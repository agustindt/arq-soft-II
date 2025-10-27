package controllers

import (
	"net/http"
	"net/url"
	"search-api/services"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SearchController struct{ svc *services.SearchService }

func NewSearchController(s *services.SearchService) *SearchController {
	return &SearchController{svc: s}
}

func (s *SearchController) Health(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) }

// GET /search?q=texto&page=1&size=10&sort=title asc,created_at desc&category=boxing&location=Cordoba
func (s *SearchController) Search(c *gin.Context) {
	q := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	sortRaw := c.DefaultQuery("sort", "")
	// Convertir "field1 asc,field2 desc" -> "field1 asc,field2 desc" (Solr ya usa este formato)
	sort := strings.TrimSpace(sortRaw)

	filters := url.Values{}
	for k, vs := range c.Request.URL.Query() {
		if k == "q" || k == "page" || k == "size" || k == "sort" {
			continue
		}
		for _, v := range vs {
			filters.Add(k, v)
		}
	}

	res, err := s.svc.Search(q, page, size, sort, filters)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}
