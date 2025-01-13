package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
)

var db *sql.DB

// User struct for login credentials
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Server struct for endpoint data
type Server struct {
	Name       string
	Env        string
	URL        string
	StatusCode int
}

// Client struct for client data
type ClientResources struct {
    ClientName string `json:"client_name"`
    Resources string `json:"resources"`
}


// InfraResource struct for infra resources
type InfraResource struct {
    ServiceName string
    URL         string
    Owner       string
}
type ResourceDetails struct {
    ServiceName string `json:"serviceName"`
    URL         string `json:"url"`
    Owner       string `json:"owner"`
}

type Endpoint struct {
    NameOfTheServer string `json:"nameoftheserver"` // JSON key will match the form field name
    Env             string `json:"env"`
    URL             string `json:"url_endpoint"`
}



// Initialize the database
func initDB() {
	var err error
	dsn := "root:admin@tcp(localhost:3306)/mydatabase"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
}


func fetchServers() ([]Server, error) {
	rows, err := db.Query("SELECT nameoftheserver, env, url_endpoint FROM servers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var servers []Server
	var redServers []Server
	var greenServers []Server

	for rows.Next() {
		var server Server
		if err := rows.Scan(&server.Name, &server.Env, &server.URL); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Check the status of the endpoint
		server.StatusCode = getStatusCode(server.URL)

		// Categorize servers based on status code
		if server.StatusCode >= 200 && server.StatusCode < 400 {
			greenServers = append(greenServers, server)
		} else {
			redServers = append(redServers, server)
		}
	}

	// Combine red servers first, then green servers
	servers = append(servers, redServers...) // Unpack redServers slice
	servers = append(servers, greenServers...) // Unpack greenServers slice

	return servers, rows.Err()
}


func getStatusCode(url string) int {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Error fetching status for URL '%s': %v", url, err)
		return 0 // Return 0 for unreachable URLs
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// Fetch infra resources
func fetchInfraResources() ([]InfraResource, error) {
    rows, err := db.Query("SELECT service_name, url, owner FROM infra_resources")
    if err != nil {
        log.Printf("Error querying infra_resources: %v", err)
        return nil, err
    }
    defer rows.Close()

    var resources []InfraResource
    for rows.Next() {
        var resource InfraResource
        if err := rows.Scan(&resource.ServiceName, &resource.URL, &resource.Owner); err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        resources = append(resources, resource)
    }

    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }

    // Log the resources fetched from the database for debugging
    log.Printf("Fetched %d infra resources", len(resources))

    return resources, nil
}

func handleAddResource(c *gin.Context) {
    var resource ResourceDetails
    if err := c.ShouldBindJSON(&resource); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
    // Insert resource into the database
    _, err := db.Exec("INSERT INTO  infra_resources (service_name, url, owner) VALUES (?, ?, ?)", resource.ServiceName, resource.URL, resource.Owner)

    if err != nil {
        log.Printf("Error inserting resource: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add resource"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resource added successfully"})
}

// Handle Add Endpoint
func handleAddEndpoints(c *gin.Context) {
    var Endpoint Endpoint
    if err := c.ShouldBindJSON(&Endpoint); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }
	if Endpoint.NameOfTheServer == "" || Endpoint.Env == "" || Endpoint.URL == "" {
        c.JSON(http.StatusBadRequest, gin.H{"message": "All fields are required"})
        return
    }
    // Insert resource into the database
    _, err := db.Exec("INSERT INTO  servers (nameoftheserver, env, url_endpoint) VALUES (?, ?, ?)", Endpoint.NameOfTheServer, Endpoint.Env, Endpoint.URL)

    if err != nil {
        log.Printf("Error inserting resource: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add resource"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Resource added successfully"})
}
	

// Fetch clients and their resources
func fetchClientsWithResources() ([]ClientResources, error) {
    rows, err := db.Query("SELECT client_name, GROUP_CONCAT(DISTINCT resource_name ORDER BY resource_name ASC) AS resources FROM client_resources GROUP BY client_name")
    if err != nil {
        log.Printf("Error querying client_resources: %v", err)
        return nil, err
    }
    defer rows.Close()

    var clientResources []ClientResources
    for rows.Next() {
        var client ClientResources
        if err := rows.Scan(&client.ClientName, &client.Resources); err != nil {
            log.Printf("Error scanning row: %v", err)
            continue
        }
        clientResources = append(clientResources, client)
    }
    if err := rows.Err(); err != nil {
        log.Printf("Error iterating rows: %v", err)
        return nil, err
    }
    return clientResources, nil
}

// Middleware for authentication
func isAuthenticated(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		return
	}
	c.Next()
}

// Serve login page
func serveLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func handleSignup(c *gin.Context) {
    var credentials struct {
        Username     string `json:"username"`
        Password     string `json:"password"`
    }
    if err := c.ShouldBindJSON(&credentials); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input"})
        return
    }

    // Implement your signup logic here (e.g., store user in DB)
    // Ensure password is hashed and stored securely

    c.JSON(http.StatusOK, gin.H{"message": "Signup successful"})
}


// Handle login
func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE username = ?", user.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	if user.Password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid username or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// Serve home page
func serveHomePage(c *gin.Context) {
	Username := c.GetString("username")
	c.HTML(http.StatusOK, "home.html",  gin.H{"Username": Username})
}

// Serve endpoints page
func serveEndpointsPage(c *gin.Context) {
	servers, err := fetchServers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching servers"})
		return
	}
	c.HTML(http.StatusOK, "endpoints.html", gin.H{"servers": servers})
}

// Serve infra resources page
func serveInfraPage(c *gin.Context) {
    resources, err := fetchInfraResources()
    if err != nil {
        log.Printf("Error fetching infra resources: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching infra resources"})
        return
    }
    c.HTML(http.StatusOK, "infra.html", gin.H{"resources": resources})
}

// Cleints and its resources
func fetchClientsWithResourcesPage(c *gin.Context) {
    clients, err := fetchClientsWithResources()
    if err != nil {
        log.Printf("Error fetching Clients resources: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching Clients resources"})
        return
    }
    c.HTML(http.StatusOK, "clients.html", gin.H{"clients": clients})
}


func main() {
	initDB()
	defer db.Close()

	r := gin.Default()
	r.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("secret"))))

	r.SetFuncMap(template.FuncMap{})
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	r.GET("/", serveLoginPage)
	r.POST("/login", login)
	r.GET("/home", isAuthenticated, serveHomePage)
	r.GET("/endpoints", isAuthenticated, serveEndpointsPage)
	r.GET("/infra", isAuthenticated, serveInfraPage)
	r.GET("/clients", isAuthenticated, fetchClientsWithResourcesPage)
	r.POST("/signup", handleSignup)
	r.POST("/addResource", handleAddResource)
	r.POST("/addEndpoint", handleAddEndpoints)
	r.Run(":8080")
}
