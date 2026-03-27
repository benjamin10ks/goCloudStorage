# Template Architecture Guide

## Overview

Your template system is now structured with HTMX for the hypermedia ecosystem, providing a modern, SPA-like UX without JavaScript frameworks.

## Template Structure

### Base Templates (Layouts)

#### 1. **base.tmpl** - Foundation
- Defines the HTML skeleton
- Includes HTMX library
- Contains global CSS variables and styles
- Provides reusable components (buttons, forms, modals, spinners)

#### 2. **auth_layout.tmpl** - Authentication Pages
- Extends `base.tmpl`
- Beautiful gradient background
- Centered card design
- Used for: login, register, passkey flows

#### 3. **app_layout.tmpl** - Application Pages  
- Extends `base.tmpl`
- Sidebar navigation
- Header with actions
- Main content area
- Used for: dashboard, file management

### Page Templates

#### Authentication Pages
- **login.tmpl** - Login page with passkey and GitHub OAuth
- **register.tmpl** - Registration with username and passkey setup
- **passkey_begin.tmpl** - Passkey registration flow (DO NOT CHANGE)

#### Application Pages
- **dashboard.tmpl** - Main file dashboard with grid layout

### Components (Reusable Partials)
- **file_card.tmpl** - Individual file display card
- **modal_share.tmpl** - Share file modal dialog
- **toast.tmpl** - Notification toasts
- **passkey_login.tmpl** - Passkey login flow

## HTMX Integration

### Key HTMX Patterns Used

#### 1. **Form Submissions**
```html
<form hx-post="/auth/register/username" 
      hx-target="#auth-form" 
      hx-swap="innerHTML">
```
- Posts form via AJAX
- Swaps result into target element
- No page reload needed

#### 2. **Button Actions**
```html
<button hx-delete="/api/delete/{{.ID}}"
        hx-confirm="Are you sure?"
        hx-target="#file-{{.ID}}"
        hx-swap="outerHTML">
```
- DELETE request on click
- Confirms before deleting
- Removes the element from DOM

#### 3. **File Upload**
```html
<input type="file" 
       hx-post="/api/upload" 
       hx-encoding="multipart/form-data"
       hx-target="#files-container"
       hx-swap="afterbegin">
```
- Uploads file without page reload
- Prepends new file card to container

#### 4. **Loading Indicators**
```html
<button class="btn">
    <span>Submit</span>
    <span class="htmx-indicator spinner"></span>
</button>
```
- Spinner shows during HTMX request
- Automatically hidden when complete

#### 5. **Modal Loading**
```html
<button hx-get="/api/files/{{.ID}}/share-modal"
        hx-target="#modal-container"
        hx-swap="innerHTML">
```
- Loads modal content on demand
- Reduces initial page size

## Backend Handler Patterns

### Full Page Render
```go
func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
    err := a.tmpl["login"].ExecuteTemplate(w, "base", nil)
    // ...
}
```

### Partial/Fragment Render (HTMX)
```go
func (a *App) handleRegisterUsername(w http.ResponseWriter, r *http.Request) {
    // Process username...
    
    // Return just the fragment, not full page
    a.tmpl["passkey_begin"].ExecuteTemplate(w, "passkey_begin", map[string]any{
        "UserID": userID,
    })
}
```

### Conditional HTMX Response
```go
func (a *App) handlePasskeyBeginLogin(w http.ResponseWriter, r *http.Request) {
    // HTMX request - return HTML fragment
    if r.Header.Get("HX-Request") == "true" {
        a.tmpl["login"].ExecuteTemplate(w, "passkey_login", nil)
        return
    }
    
    // Regular request - return JSON for API
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(options)
}
```

### Error Handling with HTMX
```go
func (a *App) handleRegisterUsername(w http.ResponseWriter, r *http.Request) {
    username := r.FormValue("username")
    userID, err := createUser(a.db, username)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        // Return HTML error to swap into page
        w.Write([]byte(`<div class="alert alert-error">Username already exists</div>`))
        return
    }
    // ...
}
```

## Best Practices

### 1. **Template Organization**
- ✅ Separate layouts (base, auth, app)
- ✅ Reusable components in `/components`
- ✅ One template per page in root `/templates`

### 2. **HTMX Usage**
- ✅ Use `hx-post`, `hx-get`, `hx-delete` for actions
- ✅ Use `hx-target` to specify where response goes
- ✅ Use `hx-swap` to control how content is inserted
- ✅ Add `hx-indicator` for loading states
- ✅ Use `hx-confirm` for destructive actions

### 3. **CSS Architecture**
- ✅ CSS variables in `:root` for theming
- ✅ Utility classes (btn, form-input, alert, etc.)
- ✅ Component-specific styles in layouts
- ✅ Responsive design with flexbox/grid

### 4. **Backend Patterns**
- ✅ Check `HX-Request` header for HTMX requests
- ✅ Return fragments for HTMX, full pages for direct access
- ✅ Use appropriate HTTP status codes
- ✅ Return HTML errors for HTMX to display inline

### 5. **User Experience**
- ✅ Loading spinners during requests
- ✅ Smooth transitions with CSS
- ✅ Toast notifications for success/error
- ✅ Confirm dialogs for destructive actions
- ✅ Empty states when no data

## File Upload Flow Example

1. User clicks "Upload File" button
2. Hidden file input is triggered
3. User selects file
4. HTMX sends multipart form data to `/api/upload`
5. Server processes upload, returns new file card HTML
6. HTMX prepends card to files container
7. User sees new file instantly

## Authentication Flow Example

### Registration
1. User visits `/auth/register`
2. Full page loads with registration form
3. User enters username, submits
4. HTMX posts to `/auth/register/username`
5. Server creates user, returns passkey setup HTML
6. Passkey HTML swaps into auth-form div
7. JavaScript initiates WebAuthn flow
8. On success, redirects to `/home`

### Login
1. User visits `/auth/login`
2. Clicks "Sign in with Passkey"
3. HTMX requests `/auth/login/passkey/begin/0`
4. Server returns passkey login HTML fragment
5. JavaScript initiates WebAuthn flow
6. On success, redirects to `/home`

## Adding New Features

### To Add a New Page:
1. Create template in `web/templates/yourpage.tmpl`
2. Choose layout (`auth_layout` or `app_layout`)
3. Define required blocks (`title`, `auth_content` or `main_content`)
4. Add to `loadTemplates()` in `utils.go`
5. Create handler in controller
6. Add route in `main.go`

### To Add a New Component:
1. Create in `web/templates/components/yourcomponent.tmpl`
2. Define template block: `{{define "component_name"}}`
3. Add to componentFiles in `utils.go`
4. Use in pages: `{{template "component_name" .}}`

### To Add HTMX Interaction:
1. Add `hx-*` attributes to HTML element
2. Create backend handler
3. Return appropriate HTML fragment
4. Set target and swap strategy

## Common Gotchas

❌ **Don't** return full page HTML for HTMX requests
✅ **Do** return just the fragment to swap

❌ **Don't** use onClick handlers for everything  
✅ **Do** use HTMX attributes declaratively

❌ **Don't** forget loading indicators
✅ **Do** add spinners to async buttons

❌ **Don't** mix template names with block names
✅ **Do** use consistent naming (template="login", block="base")

## Resources

- HTMX Docs: https://htmx.org/docs/
- HTMX Examples: https://htmx.org/examples/
- WebAuthn Guide: https://webauthn.guide/

## Next Steps

### TODO Items in Code:
1. Implement `getUserFiles()` function
2. Add file upload handler that returns file_card fragment
3. Add share modal endpoint (`/api/files/{id}/share-modal`)
4. Implement file sharing logic
5. Add "Recent" and "Shared" pages
6. Add profile settings page

### Recommended Enhancements:
- Add drag-and-drop file upload
- Add file preview modals
- Add progress bars for uploads
- Add search/filter functionality
- Add keyboard shortcuts
- Add dark mode toggle
