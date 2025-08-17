package auth

// func TestGenerateAndValidateJWT(t *testing.T) {
// 	os.Setenv("JWT_SECRET", "testsecret")
// 	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// 	testUser := &user.User{
// 		gorm.Model.ID : 123,
// 	}

// 	tokenString, err := GenerateJWT(testUser)
// 	if err != nil {
// 		t.Fatalf("GenerateJWT failed: %v", err)
// 	}
// 	if tokenString == "" {
// 		t.Fatal("GenerateJWT returned empty token")
// 	}

// 	claims, err := ValidateJWT(tokenString)
// 	if err != nil {
// 		t.Fatalf("ValidateJWT failed: %v", err)
// 	}

// 	if claims.UserId != testUser.ID {
// 		t.Errorf("UserId mismatch: got %d, want %d", claims.UserId, testUser.ID)
// 	}
// 	if claims.Issuer != "ITPL" {
// 		t.Errorf("Issuer mismatch: got %s, want ITPL", claims.Issuer)
// 	}
// 	t.Log("JWT Token:", tokenString)
// 	// 만료
// 	expiredClaims := &CustomClaims{
// 		UserId: testUser.ID,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // 이미 만료
// 			Issuer:    "ITPL",
// 		},
// 	}
// 	expiredToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims).SignedString(jwtSecret)
// 	_, err = ValidateJWT(expiredToken)
// 	if err == nil {
// 		t.Error("Expected error for expired token, got nil")
// 	}
// }
