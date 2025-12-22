# create a new repository on the command line

    echo "# States-" >> README.md
    git init
    git add README.md
    git commit -m "first commit"
    git branch -M main
    git remote add origin https://github.com/juanca0423/States-.git
    git push -u origin main
    
# push an existing repository from the command line


    git remote add origin https://github.com/juanca0423/States-.git
    git branch -M main
    git push -u origin main.z

    ssh -vT git@github.com
    # o si tu red bloquea el puerto 22:
    ssh -T -p 443 git@ssh.github.com

curl -X POST http://127.0.0.1:3000/register -H "Content-Type:aplication/json" -d '{"nombre":"Juan Carlos","apellido":"Perez Castro","email":"juanca0423@yahoo.com","pase":"654321","role":"admin"}'


    users
    ├─ id       (uint)
    ├─ email    (string)  único
    ├─ password (string)  bcrypt
    ├─ role     (string) "admin"|"cliente"|"usuario"
    └─ deleted_at

    go get github.com/gofiber/fiber/v2
    go get github.com/golang-jwt/jwt/v5
    go get golang.org/x/crypto/bcrypt


    // middleware/auth.go
    package middleware

    import (
        "errors"
            "fmt"
                "os"
                    "strings"
                        "time"

                            "github.com/gofiber/fiber/v2"
                                "github.com/golang-jwt/jwt/v5"
                                
    )

    // Claims personalizados
    type Claims struct {
        UserID uint   `json:"uid"`
            Role   string `json:"role"`
                jwt.RegisteredClaims
                
    }

    var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

    // GenerateToken crea un JWT firmado (úsalo en tu handler de login)
    func GenerateToken(userID uint, role string) (string, error) {
    claims := Claims{
            UserID: userID,
                    Role:   role,
                    RegisteredClaims: jwt.RegisteredClaims{
                                ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
                                            IssuedAt:  jwt.NewNumericDate(time.Now()),
                                                    
                    },
                        
    }
        return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
        
    }

    // Require verifica que el usuario tenga *exactamente* el rol indicado
    func Require(allowedRole string) fiber.Handler {
    return func(c *fiber.Ctx) error {
            role, err := extractRole(c)
            if err != nil {
                        return fiber.ErrUnauthorized
                                
            }
            if role != allowedRole {
                        return fiber.NewError(fiber.StatusForbidden,
                                        fmt.Sprintf("se requiere rol %s", allowedRole))
                                                
            }
                    return c.Next()
                        
    }
    
    }

    // Any verifica que el usuario tenga *alguno* de los roles listados
    func Any(roles ...string) fiber.Handler {
    return func(c *fiber.Ctx) error {
            userRole, err := extractRole(c)
            if err != nil {
                        return fiber.ErrUnauthorized
                                
            }
            for _, r := range roles {
            if userRole == r {
                            return c.Next()
                                        
            }
                    
            }
                    return fiber.NewError(fiber.StatusForbidden,
                                "sin permisos suficientes")
                                    
    }
    
    }

    // extractRole valida el JWT y devuelve el rol
    func extractRole(c *fiber.Ctx) (string, error) {
        auth := c.Get("Authorization")
        if auth == "" {
                return "", errors.New("falta header Authorization")
                    
        }
            parts := strings.Split(auth, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                    return "", errors.New("formato Authorization inválido")
                        
            }
                tokenStr := parts[1]

                token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
                if t.Method.Alg() != jwt.SigningMethodHS256.Name {
                            return nil, errors.New("método inesperado")
                                    
                }
                        return jwtSecret, nil
                            
                })
                if err != nil || !token.Valid {
                        return "", errors.New("token inválido")
                            
                }
                    claims, ok := token.Claims.(*Claims)
                    if !ok {
                            return "", errors.New("claims inválidos")
                                
                    }
                        // Guardamos datos útiles para handlers posteriores
                            c.Locals("uid", claims.UserID)
                                c.Locals("role", claims.Role)
                                    return claims.Role, nil
                                    
    }


package routes

import (
    "github.com/gofiber/fiber/v2"
        "tuproyecto/middleware"
        
)

func Setup(app *fiber.App) {
    api := app.Group("/api")

        // público
            api.Post("/login", handleLogin)

                // solo administradores
                    admin := api.Group("/admin", middleware.Require("admin"))
                        admin.Get("/dashboard", handleAdminDashboard)

                            // cliente o admin
                                cliente := api.Group("/cliente", middleware.Any("cliente", "admin"))
                                    cliente.Get("/pedidos", handlePedidosCliente)

                                        // cualquier usuario autenticado
                                            user := api.Group("/user", middleware.Any("usuario", "cliente", "admin"))
                                                user.Get("/perfil", handlePerfil)
                                                
}


    func handleLogin(c *fiber.Ctx) error {
	type input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var in input
	if err := c.BodyParser(&in); err != nil {
		return fiber.ErrBadRequest
	}

	var u models.User
	if err := database.DB.Where("email = ?", in.Email).First(&u).Error; err != nil {
		return fiber.ErrUnauthorized
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return fiber.ErrUnauthorized
	}

	tok, _ := middleware.GenerateToken(u.ID, u.Role)
	return c.JSON(fiber.Map{"token": tok})
}

app.Use(func(c *fiber.Ctx) error {
    fmt.Printf("Method=%s Path=%s Content-Type=%s\n", c.Method(), c.Path(), c.Get("Content-Type"))
    fmt.Println("Raw body --->", string(c.Body()))
    return c.Next()
})
