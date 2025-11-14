# 🤝 Contributing to Sports Activities Platform

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## 📋 Table of Contents

- [🚀 Getting Started](#-getting-started)
- [🌿 Git Workflow](#-git-workflow)
- [💻 Development Standards](#-development-standards)
- [🧪 Testing Guidelines](#-testing-guidelines)
- [📝 Documentation](#-documentation)
- [🔍 Code Review Process](#-code-review-process)
- [🐛 Bug Reports](#-bug-reports)
- [💡 Feature Requests](#-feature-requests)

## 🚀 Getting Started

### Prerequisites

Before contributing, make sure you have:

1. **Docker & Docker Compose** installed and running
2. **Git** configured with your GitHub account
3. **Node.js & pnpm** for frontend development
4. **Go 1.21+** for backend development
5. Familiarity with **React/TypeScript** and **Go/Gin**

### Initial Setup

```bash
# 1. Fork the repository on GitHub

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/arq-soft-II.git
cd arq-soft-II

# 3. Add upstream remote
git remote add upstream https://github.com/agustindt/arq-soft-II.git

# 4. Verify Docker setup
docker-compose up --build
```

## 🌿 Git Workflow

### Branch Strategy

We use **Git Flow** with the following branches:

- `main` - Production releases (protected)
- `develop` - Integration branch (default target for PRs)
- `feature/*` - New features
- `hotfix/*` - Critical production fixes

### Working on Features

```bash
# 1. Ensure you're on develop and up to date
git checkout develop
git pull upstream develop

# 2. Create feature branch
git checkout -b feature/descriptive-name

# 3. Make your changes
# ... code, code, code ...

# 4. Commit with conventional messages
git add .
git commit -m "feat: add user profile avatar upload"

# 5. Push to your fork
git push origin feature/descriptive-name

# 6. Create Pull Request to develop branch
```

### Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### Types:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

#### Examples:
```bash
feat: add JWT authentication to users API
fix: resolve date parsing error in profile updates
docs: update API endpoint documentation
refactor: extract user service into separate module
perf: optimize database query performance
test: add unit tests for auth middleware
chore: update Docker compose configuration
```

### Keeping Your Fork Updated

```bash
# Fetch upstream changes
git fetch upstream

# Update develop branch
git checkout develop
git merge upstream/develop

# Push updates to your fork
git push origin develop
```

## 💻 Development Standards

### Go Backend Standards

#### Code Structure
```go
// Package organization
package handlers

import (
    "context"
    "net/http"
    
    "your-module/config"
    "your-module/models" 
    "your-module/utils"
    
    "github.com/gin-gonic/gin"
)

// Function documentation
// CreateUser creates a new user account
// @Summary Create new user
// @Description Create a new user with provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User details"
// @Success 201 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/users [post]
func CreateUser(c *gin.Context) {
    // Implementation...
}
```

#### Best Practices
- Use dependency injection
- Handle errors properly
- Use structured logging
- Implement proper validation
- Write unit tests
- Follow Go naming conventions
- Use interfaces for testability

#### Error Handling
```go
// Good error handling
if err := db.Create(&user).Error; err != nil {
    log.WithError(err).Error("Failed to create user")
    c.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to create user",
        "message": "Internal server error",
    })
    return
}
```

### TypeScript/React Frontend Standards

#### Component Structure
```typescript
// ComponentName.tsx
import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography,
  Button 
} from '@mui/material';

interface ComponentNameProps {
  id: string;
  onUpdate?: (data: UpdateData) => void;
}

const ComponentName: React.FC<ComponentNameProps> = ({ 
  id, 
  onUpdate 
}) => {
  const [loading, setLoading] = useState<boolean>(false);
  
  useEffect(() => {
    // Side effects
  }, [id]);

  const handleSubmit = async (): Promise<void> => {
    // Event handlers
  };

  return (
    <Box sx={{ padding: 2 }}>
      <Typography variant="h6">
        Component Title
      </Typography>
      {/* Component JSX */}
    </Box>
  );
};

export default ComponentName;
```

#### Best Practices
- Use TypeScript strict mode
- Implement proper error boundaries
- Follow Material-UI design system
- Write component tests
- Use proper state management
- Implement accessibility features
- Optimize performance with React.memo when needed

#### Type Definitions
```typescript
// types/index.ts
export interface User {
  id: number;
  email: string;
  username: string;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  role: 'user' | 'moderator' | 'admin' | 'root';
  created_at: string;
  updated_at: string;
}

export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}
```

### Docker Standards

#### Dockerfile Best Practices
```dockerfile
# Use specific versions
FROM node:18-alpine AS builder

# Create non-root user
RUN addgroup -g 1001 -S nodejs
RUN adduser -S nextjs -u 1001

# Set working directory
WORKDIR /app

# Copy package files first (for better caching)
COPY package*.json ./
RUN npm ci --only=production

# Copy source code
COPY . .

# Build application
RUN npm run build

# Production stage
FROM node:18-alpine AS runner
WORKDIR /app

# Copy built application
COPY --from=builder /app/build ./build

# Expose port and define health check
EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000/health || exit 1

CMD ["npm", "start"]
```

## 🧪 Testing Guidelines

### Backend Testing (Go)

#### Unit Tests
```go
// handlers/auth_test.go
package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.New()
    router.POST("/login", Login)

    // Test case
    req, _ := http.NewRequest("POST", "/login", strings.NewReader(`{
        "email": "test@example.com",
        "password": "password123"
    }`))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // Assertions
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "token")
}
```

#### Integration Tests
```go
// Test with database
func TestCreateUserIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    // Test logic
    // ...
}
```

### Frontend Testing (React)

#### Component Tests
```typescript
// __tests__/Login.test.tsx
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Login from '../Login';

describe('Login Component', () => {
  test('renders login form', () => {
    render(<Login />);
    
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  test('submits form with valid data', async () => {
    const mockLogin = jest.fn();
    render(<Login onLogin={mockLogin} />);
    
    fireEvent.change(screen.getByLabelText(/email/i), {
      target: { value: 'test@example.com' }
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'password123' }
    });
    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(mockLogin).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123'
      });
    });
  });
});
```

### Running Tests

```bash
# Backend tests
cd backend/users-api
go test ./... -v

# Frontend tests
cd frontend
pnpm test

# Integration tests with Docker
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

## 📝 Documentation

### Code Documentation

- **Go**: Use godoc comments for public functions and types
- **TypeScript**: Use JSDoc comments for complex functions
- **API**: Document endpoints with OpenAPI/Swagger annotations

### README Updates

When adding features, update relevant documentation:
- API endpoints in README.md
- Environment variables in environment.example
- Docker configuration changes
- Installation/setup procedures

### API Documentation

Use Swagger annotations for Go APIs:

```go
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/profile [get]
func GetProfile(c *gin.Context) {
    // Implementation
}
```

## 🔍 Code Review Process

### Creating Pull Requests

1. **Clear Title**: Use descriptive titles
   - ✅ "feat: implement user avatar upload functionality"
   - ❌ "fix stuff"

2. **Detailed Description**: Include:
   - What changes were made
   - Why the changes were necessary
   - How to test the changes
   - Screenshots for UI changes

3. **Checklist**: Ensure your PR includes:
   - [ ] Tests for new functionality
   - [ ] Updated documentation
   - [ ] No breaking changes (or documented)
   - [ ] Docker containers build successfully
   - [ ] Code follows project standards

### Review Criteria

Reviewers will check for:

- **Functionality**: Does the code work as intended?
- **Security**: Are there any security vulnerabilities?
- **Performance**: Are there performance implications?
- **Maintainability**: Is the code readable and maintainable?
- **Testing**: Are there adequate tests?
- **Documentation**: Is the code properly documented?

### Addressing Review Feedback

```bash
# Make changes based on feedback
git add .
git commit -m "fix: address review feedback - improve error handling"

# Push updates (no need for new PR)
git push origin feature/your-feature-name
```

## 🐛 Bug Reports

### Before Reporting

1. Check existing issues to avoid duplicates
2. Ensure you're using the latest version
3. Try reproducing with minimal configuration

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '...'
3. Scroll down to '...'
4. See error

**Expected behavior**
What you expected to happen.

**Screenshots**
If applicable, add screenshots.

**Environment:**
- OS: [e.g. macOS, Ubuntu]
- Docker version: [e.g. 20.10.7]
- Browser: [e.g. chrome, safari]
- API version: [e.g. v2.1.0]

**Additional context**
Any other context about the problem.
```

## 💡 Feature Requests

### Before Requesting

1. Check if feature already exists or is planned
2. Consider if it fits the project scope
3. Think about implementation complexity

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
A clear description of the problem.

**Describe the solution you'd like**
A clear description of what you want to happen.

**Describe alternatives you've considered**
Other solutions you've considered.

**Additional context**
Screenshots, mockups, or examples.

**Implementation notes**
Technical considerations or suggestions.
```

---

## 🚀 Ready to Contribute?

1. **Read this guide thoroughly**
2. **Set up your development environment**
3. **Pick an issue or create a feature request**
4. **Follow the git workflow**
5. **Write good code and tests**
6. **Submit a pull request**
7. **Respond to review feedback**

### Questions?

- 📧 Open an issue for technical questions
- 💬 Start a discussion for general questions
- 📖 Check the main README.md for setup help

**Thank you for contributing! 🎉**
