package permissions

// The Negroni middleware handler

import (
	"fmt"
	"net/http"
	"strings"
)

var perm Permissions

type Permissions struct {
	state              *UserState
	adminPathPrefixes  []string
	userPathPrefixes   []string
	publicPathPrefixes []string
	rootIsPublic       bool
	denied             http.HandlerFunc
}

func New() *Permissions {
	return NewPermissions(NewUserStateSimple())
}

func NewPermissions(state *UserState) *Permissions {
	// default permissions
	return &Permissions{state,
		[]string{"/admin"},                                                            // admin path prefixes
		[]string{"/repo", "/data"},                                                    // user path prefixes
		[]string{"/", "/login", "/register", "/favicon.ico", "/style", "/img", "/js"}, // public
		true,
		PermissionDenied}
}

func (perm *Permissions) SetDenyFunction(f http.HandlerFunc) {
	perm.denied = f
}

func (perm *Permissions) UserState() *UserState {
	return perm.state
}

// Add an url path prefix that is a page for the logged in administrators
func (perm *Permissions) AddAdminPath(prefix string) {
	perm.adminPathPrefixes = append(perm.adminPathPrefixes, prefix)
}

// Add an url path prefix that is a page for the logged in users
func (perm *Permissions) AddUserPath(prefix string) {
	perm.userPathPrefixes = append(perm.userPathPrefixes, prefix)
}

// Add an url path prefix that is a public page
func (perm *Permissions) AddPublicPath(prefix string) {
	perm.publicPathPrefixes = append(perm.publicPathPrefixes, prefix)
}

// Set all url path prefixes that are for the logged in administrator pages
func (perm *Permissions) SetAdminPath(pathPrefixes []string) {
	perm.adminPathPrefixes = pathPrefixes
}

// Set all url path prefixes that are for the logged in user pages
func (perm *Permissions) SetUserPath(pathPrefixes []string) {
	perm.userPathPrefixes = pathPrefixes
}

// Set all url path prefixes that are for the public pages
func (perm *Permissions) SetPublicPath(pathPrefixes []string) {
	perm.publicPathPrefixes = pathPrefixes
}

func PermissionDenied(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Permission denied.")
}

// Check if the user has the right admin/user rights
func (perm *Permissions) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	path := req.URL.Path // the path of the url that the user wish to visit
	reject := false

	// If it's not "/" and set to be public regardless of permissions
	if !(perm.rootIsPublic && path == "/") {

		// Reject if it is an admin page and user does not have admin permissions
		for _, prefix := range perm.adminPathPrefixes {
			if strings.HasPrefix(path, prefix) {
				if !perm.state.AdminRights(req) {
					reject = true
					break
				}
			}
		}

		if !reject {
			// Reject if it's a user page and the user does not have user rights
			for _, prefix := range perm.userPathPrefixes {
				if strings.HasPrefix(path, prefix) {
					if !perm.state.UserRights(req) {
						reject = true
						break
					}
				}
			}
		}

		if !reject {
			// Reject if it's not a public page
			found := false
			for _, prefix := range perm.publicPathPrefixes {
				if strings.HasPrefix(path, prefix) {
					found = true
					break
				}
			}
			if !found {
				reject = true
			}
		}

	}

	if reject {
		// Permission denied function
		perm.denied(rw, req)

		// Reject the request by not calling the next handler below
		return
	}

	// Call the next middleware handler
	next(rw, req)
}