// Package template provides utilities for rendering environment variable
// placeholders within string templates. It supports both ${VAR_NAME} and
// $VAR_NAME syntax, and can resolve variables from a provided map or
// directly from the current process environment.
//
// Example usage:
//
//	res := template.Render("host=${DB_HOST} port=${DB_PORT}", vars)
//	if len(res.Missing) > 0 {
//	    log.Printf("unresolved vars: %v", res.Missing)
//	}
//	fmt.Println(res.Rendered)
package template
