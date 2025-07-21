package handler

// type AuthHandler struct {
// 	authService auth.AuthService
// }

// // NewAuthHandler creates a new authentication handler
// func NewAuthHandler(authService auth.AuthService) *AuthHandler {
// 	return &AuthHandler{
// 		authService: authService,
// 	}
// }

// // Register handles user registration
// func (h *AuthHandler) Register(c *fiber.Ctx) error {
// 	var req auth.RegisterRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	user, token, err := h.authService.Register(c.Context(), &req)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user":  user,
// 		"token": token,
// 	})
// }

// // Login handles user login
// func (h *AuthHandler) Login(c *fiber.Ctx) error {
// 	var req auth.LoginRequest
// 	if err := c.BodyParser(&req); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request body",
// 		})
// 	}

// 	user, token, err := h.authService.Login(c.Context(), &req)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user":  user,
// 		"token": token,
// 	})
// }

// // Profile returns the current user's profile
// func (h *AuthHandler) Profile(c *fiber.Ctx) error {
// 	user, ok := GetUserFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"user": user,
// 	})
// }

// // RefreshToken refreshes the access token
// func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
// 	userID, ok := GetUserIDFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	token, err := h.authService.RefreshToken(c.Context(), userID)
// 	if err != nil {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"token": token,
// 	})
// }

// // Logout logs out the current user
// func (h *AuthHandler) Logout(c *fiber.Ctx) error {
// 	userID, ok := GetUserIDFromContext(c)
// 	if !ok {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"error": "User not authenticated",
// 		})
// 	}

// 	if err := h.authService.Logout(c.Context(), userID); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"message": "Successfully logged out",
// 	})
// }
