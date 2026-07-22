package handlers

import "net/http"

func (r *Router) Dashboard(w http.ResponseWriter, req *http.Request) {
	r.render(w, "dashboard", map[string]interface{}{
		"Title":      "Dashboard",
		"ActivePage": "dashboard",
	})
}
