package utils

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var ipLocationCache sync.Map

func GetClientIP(c *gin.Context) string {
	if realIP := normalizeIP(c.GetHeader("X-Real-IP")); realIP != "" {
		return realIP
	}

	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		firstValid := ""
		for _, part := range parts {
			ip := normalizeIP(part)
			if ip == "" {
				continue
			}
			if firstValid == "" {
				firstValid = ip
			}
			if isPublicIP(ip) {
				return ip
			}
		}
		if firstValid != "" {
			return firstValid
		}
	}

	if clientIP := normalizeIP(c.ClientIP()); clientIP != "" {
		return clientIP
	}

	if host, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr)); err == nil {
		if remoteIP := normalizeIP(host); remoteIP != "" {
			return remoteIP
		}
	}

	return normalizeIP(c.Request.RemoteAddr)
}

func normalizeIP(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	ip := net.ParseIP(value)
	if ip == nil {
		return ""
	}
	return ip.String()
}

func isPublicIP(value string) bool {
	ip := net.ParseIP(value)
	if ip == nil {
		return false
	}
	return !ip.IsLoopback() && !ip.IsPrivate() && !ip.IsUnspecified() && !ip.IsMulticast() && !ip.IsLinkLocalUnicast() && !ip.IsLinkLocalMulticast()
}

func ResolveIPLocation(ipText string) string {
	ip := net.ParseIP(strings.TrimSpace(ipText))
	if ip == nil {
		return "未知地区"
	}

	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
		return "本地"
	}

	if cached, ok := ipLocationCache.Load(ip.String()); ok {
		return cached.(string)
	}

	location := lookupIPLocation(ip.String())
	ipLocationCache.Store(ip.String(), location)
	return location
}

type ipAPIResponse struct {
	Status     string `json:"status"`
	Country    string `json:"country"`
	RegionName string `json:"regionName"`
	City       string `json:"city"`
}

type ipWhoIsResponse struct {
	Success bool   `json:"success"`
	Country string `json:"country"`
	Region  string `json:"region"`
	City    string `json:"city"`
}

func lookupIPLocation(ip string) string {
	if location := lookupIPWhoIs(ip); location != "" {
		return location
	}
	if location := lookupIPAPI(ip); location != "" {
		return location
	}
	return "未知地区"
}

func lookupIPWhoIs(ip string) string {
	client := http.Client{Timeout: 2 * time.Second}
	endpoint := "https://ipwho.is/" + url.PathEscape(ip) + "?lang=zh-CN&fields=success,country,region,city"

	resp, err := client.Get(endpoint)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result ipWhoIsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}
	if !result.Success || strings.TrimSpace(result.Country) == "" {
		return ""
	}

	return compactLocation(result.Country, result.Region, result.City)
}

func lookupIPAPI(ip string) string {
	client := http.Client{Timeout: 2 * time.Second}
	endpoint := "http://ip-api.com/json/" + url.PathEscape(ip) + "?fields=status,country,regionName,city&lang=zh-CN"

	resp, err := client.Get(endpoint)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result ipAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}
	if result.Status != "success" || strings.TrimSpace(result.Country) == "" {
		return ""
	}

	return compactLocation(result.Country, result.RegionName, result.City)
}

func compactLocation(country string, regionName string, cityName string) string {
	country = strings.TrimSpace(country)
	regionName = strings.TrimSpace(regionName)
	cityName = strings.TrimSpace(cityName)

	parts := []string{country}
	if region := regionName; region != "" && region != country {
		parts = append(parts, region)
	} else if city := cityName; city != "" && city != country {
		parts = append(parts, city)
	}

	return strings.Join(parts, " ")
}

func ParseBrowser(userAgent string) string {
	value := strings.TrimSpace(userAgent)
	switch {
	case strings.Contains(value, "Edg/"):
		return "Edge"
	case strings.Contains(value, "Firefox/"):
		return "Firefox"
	case strings.Contains(value, "Chrome/"), strings.Contains(value, "CriOS/"):
		return "Chrome"
	case strings.Contains(value, "Safari/"):
		return "Safari"
	default:
		return ""
	}
}

func ParseOS(userAgent string) string {
	value := strings.TrimSpace(userAgent)
	switch {
	case strings.Contains(value, "Android"):
		return "Android"
	case strings.Contains(value, "iPhone"), strings.Contains(value, "iPad"), strings.Contains(value, "iPod"):
		return "iOS"
	case strings.Contains(value, "Mac OS X"), strings.Contains(value, "Macintosh"):
		return "macOS"
	case strings.Contains(value, "Windows NT"):
		return "Windows"
	case strings.Contains(value, "Linux"):
		return "Linux"
	default:
		return ""
	}
}
