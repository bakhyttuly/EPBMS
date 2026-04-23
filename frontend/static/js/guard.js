(function() {
    const token = localStorage.getItem("token");
    const userStr = localStorage.getItem("user");
    const currentPath = window.location.pathname;

    // Public pages that don't need a token
    const publicPages = ["/", "/login", "/register"];

    if (!token || !userStr) {
        if (!publicPages.includes(currentPath)) {
            window.location.href = "/";
        }
        return;
    }

    const user = JSON.parse(userStr);

    // Role-based access control on the client side
    if (currentPath === "/dashboard-page" && user.role !== "admin") {
        window.location.href = "/performers-page";
    }

    if (currentPath === "/my-schedule-page" && user.role !== "performer") {
        window.location.href = "/performers-page";
    }
})();

// Helper function to include token in fetch requests
async function authenticatedFetch(url, options = {}) {
    const token = localStorage.getItem("token");
    if (!options.headers) {
        options.headers = {};
    }
    options.headers["Authorization"] = `Bearer ${token}`;
    
    const response = await fetch(url, options);
    
    if (response.status === 401) {
        // Token expired or invalid
        localStorage.removeItem("token");
        localStorage.removeItem("user");
        window.location.href = "/";
        return null;
    }
    
    return response;
}
